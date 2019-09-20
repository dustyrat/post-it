package controller

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"text/tabwriter"
)

type Stats struct {
	Responses map[int]int
	Entries []Entry
	mux       sync.Mutex
}

type Entry struct {
	Record     Record
	Status int
	//Body string
	Error      error

	//recordBody bool
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
	//if e.recordBody {
	//	out = append(out, e.Body)
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

func (s *Stats) append(entry Entry){
	s.mux.Lock()
	s.Entries = append(s.Entries, entry)
	s.mux.Unlock()
}