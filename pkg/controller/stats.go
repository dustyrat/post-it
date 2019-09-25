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
	"fmt"
	"net/http"
	"os"
	"post-it/pkg/csv"
	"sort"
	"strconv"
	"sync"
	"text/tabwriter"
)

type Stats struct {
	Responses map[int]int
	Entries   []Entry
	mux       sync.Mutex
}

type Entry struct {
	Record csv.Record

	Status       int
	Headers      string
	ResponseBody string
	Error        error

	RecordHeaders      bool
	RecordResponseBody bool
}

func (e Entry) Strings() []string {
	out := make([]string, 0)
	for _, header := range e.Record.Headers {
		if value, ok := e.Record.Fields[header]; ok {
			out = append(out, value)
		} else {
			out = append(out, "")
		}
	}
	out = append(out, strconv.Itoa(e.Status))

	//if e.RecordHeaders {
	out = append(out, e.Headers)
	//}

	//if e.RecordResponseBody {
	out = append(out, e.ResponseBody)
	//}

	if e.Error != nil {
		out = append(out, e.Error.Error())
	} else {
		out = append(out, "")
	}

	return out
}

func (s *Stats) Print() {
	codes := make(sort.IntSlice, 0)
	values := make([]interface{}, 0)
	headers := ""
	line := ""

	for code := range s.Responses {
		codes = append(codes, code)
	}
	codes.Sort()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	for _, code := range codes {
		headers += fmt.Sprintf("%s: %d \t", http.StatusText(code), code)
		line += "%d \t"
		values = append(values, s.Responses[code])
	}

	fmt.Fprintln(w, headers)
	fmt.Fprintf(w, fmt.Sprintf("%s\n", line), values...)
	w.Flush()
}

func (s *Stats) Increment(code int) {
	s.mux.Lock()
	s.Responses[code]++
	s.mux.Unlock()
}

func (s *Stats) Add(f *int, v int) {
	s.mux.Lock()
	*f += v
	s.mux.Unlock()
}

func (s *Stats) append(entry Entry) {
	s.mux.Lock()
	s.Entries = append(s.Entries, entry)
	s.mux.Unlock()
}
