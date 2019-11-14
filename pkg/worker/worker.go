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
package worker

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"

	"github.com/DustyRat/post-it/pkg/options"

	"github.com/DustyRat/post-it/pkg/file/csv"

	"github.com/DustyRat/post-it/pkg/client"

	"github.com/vbauerster/mpb"

	log "github.com/sirupsen/logrus"
)

var (
	statusDxx *regexp.Regexp
	statusDdd *regexp.Regexp
)

func init() {
	statusDxx = regexp.MustCompile("^-?\\dxx")
	statusDdd = regexp.MustCompile("\\d{3}")
}

type worker struct {
	id      int
	record  csv.Record
	request *client.Request

	pool     *Pool
	progress *mpb.Bar
}

type entry struct {
	record  csv.Record
	request *client.Request
	err     error
}

// Strings ...
func (e *entry) Strings(flags options.Flags) []string {
	out := make([]string, 0)
	for _, header := range e.record.Headers {
		if value, ok := e.record.Fields[header]; ok {
			out = append(out, value)
		} else {
			out = append(out, "")
		}
	}
	out = append(out, strconv.Itoa(e.request.Response.StatusCode))

	if flags.Headers {
		out = append(out, toString(e.request.Response.Header))
	}

	if flags.Body {
		out = append(out, string(e.request.Response.Body))
	}

	if e.err != nil {
		out = append(out, e.err.Error())
	} else {
		out = append(out, "")
	}
	return out
}

// Work ...
func (w *worker) Work(id int) {
	w.id = id
	entry := &entry{record: w.record, request: w.request}
	defer func() {
		w.write(entry)
		w.progress.Increment()
		w.done()
	}()

	response, err := w.pool.client.Do(w.request.Method, w.request.URL, w.request.Header, w.request.Body)
	if err != nil {
		entry.err = err
		w.pool.stats.Errors.Increment(err)
		return
	}
	w.request.Response = response
	w.pool.stats.Codes.Increment(response.StatusCode)
	w.pool.stats.Latencies.Increment(response.Duration)
}

func (w *worker) write(entry *entry) {
	if w.pool.writer == nil {
		return
	}
	defer w.pool.writer.Flush()

	switch w.pool.options.Flags.Type {
	case "all":
		w.pool.writer.Write(entry.Strings(w.pool.options.Flags))
	case "error":
		if entry.err != nil {
			w.pool.writer.Write(entry.Strings(w.pool.options.Flags))
		}
	case "status":
		if w.pool.options.Flags.Status == "any" {
			w.pool.writer.Write(entry.Strings(w.pool.options.Flags))
		} else if statusDxx.MatchString(w.pool.options.Flags.Status) {
			status, _ := strconv.Atoi(strings.Replace(w.pool.options.Flags.Status, "xx", "00", 1))
			if status > 0 {
				a := status
				b := a + 100
				if client.InRange(entry.request.Response.StatusCode, a, b) {
					w.pool.writer.Write(entry.Strings(w.pool.options.Flags))
				}
			} else {
				a := -status
				b := a + 100
				if !client.InRange(entry.request.Response.StatusCode, a, b) {
					w.pool.writer.Write(entry.Strings(w.pool.options.Flags))
				}
			}
		} else if statusDdd.MatchString(w.pool.options.Flags.Status) {
			status, _ := strconv.Atoi(w.pool.options.Flags.Status)
			if entry.request.Response.StatusCode == status {
				w.pool.writer.Write(entry.Strings(w.pool.options.Flags))
			}
		}
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

func toString(headers http.Header) string {
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
