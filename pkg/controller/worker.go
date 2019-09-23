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
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"post-it/pkg/client"
	"post-it/pkg/csv"
	"runtime/debug"
	"strings"

	"github.com/cheggaaa/pb/v3"
	log "github.com/sirupsen/logrus"
)

type worker struct {
	id     int
	client *client.Client
	method string
	url    string
	chunk  []csv.Record
	writer *csv.Writer

	to    int
	from  int
	batch int

	stats    *Stats
	progress *pb.ProgressBar
}

func (w *worker) Work(id int) {
	w.id = id
	defer w.done()

	for i := range w.chunk {
		w.call(w.chunk[i])
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

func (w *worker) call(record csv.Record) {
	defer w.progress.Increment()
	defer w.writer.Flush()

	entry := Entry{Record: record}
	defer func() {
		if !client.IsSuccessful(entry.Status) {
			w.writer.Write(entry.Strings())
		}
	}()

	rawurl := w.url
	for k, v := range record.Fields {
		rawurl = strings.Replace(rawurl, fmt.Sprintf("{%s}", k), v, 1)
	}

	_url, err := url.Parse(rawurl)
	if err != nil {
		entry.Error = err
		return
	}

	var response *client.Response
	if method, ok := record.Fields["method"]; ok {
		response, err = w.client.Do(method, _url, http.Header{}, nil)
	} else {
		response, err = w.client.Do(w.method, _url, http.Header{}, nil)
	}
	if err != nil {
		entry.Error = err
		return
	}

	entry.Status = response.StatusCode
	w.stats.Increment(response.StatusCode)

	entry.ResponseBody = string(response.Body)
	entry.Headers = headerToString(response.Header)
}

func headerToString(headers http.Header) string {
	str := ""
	i := 0
	for key, values := range headers {
		str += fmt.Sprintf("%s: ", key)
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
		i++
	}
	return str
}
