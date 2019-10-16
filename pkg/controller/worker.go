/*
Copyright © 2019 Dustin Ratcliffe <dustin.k.ratcliffe@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package controller

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/goinggo/work"

	"github.com/vbauerster/mpb/decor"

	"github.com/vbauerster/mpb"

	"github.com/DustyRat/post-it/pkg/client"
	"github.com/DustyRat/post-it/pkg/csv"

	log "github.com/sirupsen/logrus"
)

var (
	status_dxx *regexp.Regexp
	status_ddd *regexp.Regexp
)

func init() {
	status_dxx = regexp.MustCompile("^-?\\dxx")
	status_ddd = regexp.MustCompile("\\d{3}")
}

type Pool struct {
	client *client.Client
	method string
	rawurl string

	writer   *csv.Writer
	progress *mpb.Progress
	bar      *mpb.Bar
	bars     map[int]*mpb.Bar
	stats    *Stats

	workerPool *work.Pool
	workers    []*worker
	total      int64

	mutex *sync.Mutex
}

func NewPool(workerPool *work.Pool, client *client.Client, method, rawurl string, stats *Stats, writer *csv.Writer, progress *mpb.Progress) *Pool {
	return &Pool{
		client:   client,
		method:   method,
		rawurl:   rawurl,
		stats:    stats,
		writer:   writer,
		progress: progress,
		bar: progress.AddBar(0,
			mpb.BarID(0),
			mpb.PrependDecorators(
				decor.Name("Total", decor.WCSyncSpaceR),
				decor.Counters(0, "%d / %d", decor.WCSyncSpaceR),
			),
			mpb.AppendDecorators(
				decor.OnComplete(decor.Percentage(decor.WCSyncSpaceR), "complete"),
				decor.AverageSpeed(0, "% .1f/s", decor.WCSyncSpaceR),
				decor.Name("Elapsed:", decor.WCSyncSpaceR),
				decor.Elapsed(decor.ET_STYLE_GO, decor.WCSyncSpaceR),
				decor.Name("ETA:", decor.WCSyncSpaceR),
				decor.AverageETA(decor.ET_STYLE_GO, decor.WCSyncSpaceR),
			),
		),
		bars:       map[int]*mpb.Bar{},
		workerPool: workerPool,
		mutex:      &sync.Mutex{},
	}
}

func (p *Pool) NewWorker(chunk []csv.Record, flags Flags, keys Keys) *worker {
	p.mutex.Lock()
	p.total += int64(len(chunk))
	p.bar.SetTotal(p.total, false)
	w := &worker{pool: p, progress: p.bar, chunk: chunk, flags: flags, keys: keys}
	p.workers = append(p.workers, w)
	p.mutex.Unlock()
	return w
}

func (p *Pool) Run() {
	for _, worker := range p.workers {
		p.workerPool.Run(worker)
	}
	p.workerPool.Shutdown()
	p.progress.Wait()
}

func (p *Pool) addBar(bar *mpb.Bar) {
	p.mutex.Lock()
	p.bars[bar.ID()] = bar
	p.mutex.Unlock()
}

func (p *Pool) getBar(id int, total int64) *mpb.Bar {
	bar := p.progress.AddBar(total,
		mpb.BarRemoveOnComplete(),
		mpb.BarID(id),
		mpb.BarPriority(id),
		mpb.PrependDecorators(
			decor.Name(fmt.Sprintf("Worker#%d", id), decor.WCSyncSpaceR),
			decor.Counters(0, "%d / %d", decor.WCSyncSpaceR),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WCSyncSpaceR),
			decor.AverageSpeed(0, "% .1f/s", decor.WCSyncSpaceR),
		),
	)
	p.addBar(bar)
	return bar
}

type worker struct {
	id    int
	chunk []csv.Record

	flags Flags
	keys  Keys

	pool     *Pool
	progress *mpb.Bar
}

type Flags struct {
	Headers bool
	Body    bool
}

type Keys struct {
	Body         string
	Status       string
	ResponseType string
}

func (w *worker) Work(id int) {
	w.id = id
	defer w.done()

	var bar *mpb.Bar
	//bar := w.pool.getBar(id, int64(len(w.chunk)))

	for i := range w.chunk {
		w.call(bar, w.chunk[i])
	}
}

func (w *worker) done() {
	if r := recover(); r != nil {
		log.Debug("recovered from ", r)
		stack := debug.Stack()
		var err error
		switch t := r.(type) {
		case string:
			err = errors.New(t)
		case error:
			err = t
		default:
			err = errors.New("unknown error")
		}
		log.WithFields(log.Fields{"stacktrace": string(stack)}).Errorf("[%d] %s", w.id, err)
	}
}

func (w *worker) call(bar *mpb.Bar, record csv.Record) {
	defer func() {
		if bar != nil {
			bar.Increment()
		}
		if w.progress != nil {
			w.progress.Increment()
		}
	}()

	entry := Entry{Record: record}
	defer w.write(&entry)

	rawurl := w.pool.rawurl
	for k, v := range record.Fields {
		rawurl = strings.Replace(rawurl, fmt.Sprintf("{%s}", k), v, 1)
	}

	_url, err := url.Parse(rawurl)
	if err != nil {
		entry.Error = err
		return
	}

	var body io.Reader
	if requestBody, ok := record.Fields[w.keys.Body]; ok {
		body = bytes.NewBuffer([]byte(requestBody))
	}

	response, err := w.pool.client.Do(w.pool.method, _url, http.Header{}, body)
	if err != nil {
		w.pool.stats.Increment(0)
		entry.Error = err
		return
	}

	w.pool.stats.Increment(response.StatusCode)
	entry.Status = response.StatusCode

	if w.flags.Headers {
		entry.Headers = headerToString(response.Header)
	}

	if w.flags.Body {
		entry.ResponseBody = string(response.Body)
	}
}

func (w *worker) write(entry *Entry) {
	if w.pool.writer == nil {
		return
	}

	defer w.pool.writer.Flush()
	switch w.keys.ResponseType {
	case "all":
		w.pool.writer.Write(entry.Strings())
	case "error":
		if entry.Error != nil {
			w.pool.writer.Write(entry.Strings())
		}
	case "status":
		if w.keys.Status == "any" {
			w.pool.writer.Write(entry.Strings())
		} else if status_dxx.MatchString(w.keys.Status) {
			status, _ := strconv.Atoi(strings.Replace(w.keys.Status, "xx", "00", 1))
			if status > 0 {
				a := status
				b := a + 100
				if client.InRange(entry.Status, a, b) {
					w.pool.writer.Write(entry.Strings())
				}
			} else {
				a := -status
				b := a + 100
				if !client.InRange(entry.Status, a, b) {
					w.pool.writer.Write(entry.Strings())
				}
			}
		} else if status_ddd.MatchString(w.keys.Status) {
			status, _ := strconv.Atoi(w.keys.Status)
			if entry.Status == status {
				w.pool.writer.Write(entry.Strings())
			}
		}
	}
}

func headerToString(headers http.Header) string {
	keys := make([]string, 0)
	for key := range headers {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	str := ""
	for i, key := range keys {
		str += fmt.Sprintf("%s: ", key)
		values := headers[key]
		for j, value := range values {
			str += value
			if j < len(values)-1 {
				str += ", "
			}
		}

		if i < len(headers)-1 {
			str += "; "
		} else {
			str += ";"
		}
	}
	return str
}
