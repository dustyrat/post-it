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
	"regexp"
	"runtime/debug"

	"github.com/DustyRat/post-it/pkg/client"

	"github.com/vbauerster/mpb"

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

type worker struct {
	id       int
	requests []*client.Request

	pool     *Pool
	progress *mpb.Bar
}

func (w *worker) Work(id int) {
	w.id = id
	defer w.done()

	for i := range w.requests {
		w.work(w.requests[i])
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

func (w *worker) work(request *client.Request) {
	defer w.progress.Increment()
	response, err := w.pool.client.Do(request.Method, request.URL, request.Header, request.Body)
	if err != nil {
		w.pool.stats.Errors.Increment(err)
		return
	}
	request.Response = response
	w.pool.stats.Codes.Increment(response.StatusCode)
	w.pool.stats.Latencies.Increment(response.Duration)

	//var body io.Reader
	////if requestBody, ok := record.Fields[w.keys.Body]; ok {
	////	body = bytes.NewBuffer([]byte(requestBody))
	////}
	//
	//request, err := client.NewRequest(w.pool.method, w.pool.rawurl, http.Header{}, body, record.Fields)
	//if err != nil {
	//	w.pool.stats.Errors.Increment(err)
	//	return
	//}
}

//func (w *worker) write(entry *stats.Entry) {
//	if w.pool.writer == nil {
//		return
//	}
//
//	defer w.pool.writer.Flush()
//switch w.keys.ResponseType {
//case "all":
//	w.pool.writer.Write(entry.Strings())
//case "error":
//	if entry.Error != nil {
//		w.pool.writer.Write(entry.Strings())
//	}
//case "status":
//	if w.keys.Status == "any" {
//		w.pool.writer.Write(entry.Strings())
//	} else if status_dxx.MatchString(w.keys.Status) {
//		status, _ := strconv.Atoi(strings.Replace(w.keys.Status, "xx", "00", 1))
//		if status > 0 {
//			a := status
//			b := a + 100
//			if client.InRange(entry.Status, a, b) {
//				w.pool.writer.Write(entry.Strings())
//			}
//		} else {
//			a := -status
//			b := a + 100
//			if !client.InRange(entry.Status, a, b) {
//				w.pool.writer.Write(entry.Strings())
//			}
//		}
//	} else if status_ddd.MatchString(w.keys.Status) {
//		status, _ := strconv.Atoi(w.keys.Status)
//		if entry.Status == status {
//			w.pool.writer.Write(entry.Strings())
//		}
//	}
//}
//}

//func headerToString(headers http.Header) string {
//	keys := make([]string, 0)
//	for key := range headers {
//		keys = append(keys, key)
//	}
//	sort.Strings(keys)
//
//	str := ""
//	for i, key := range keys {
//		str += fmt.Sprintf("%s: ", key)
//		values := headers[key]
//		for j, value := range values {
//			str += value
//			if j < len(values)-1 {
//				str += ", "
//			}
//		}
//
//		if i < len(headers)-1 {
//			str += "; "
//		} else {
//			str += ";"
//		}
//	}
//	return str
//}
