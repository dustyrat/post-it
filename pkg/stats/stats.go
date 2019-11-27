package stats

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"sync"
	"text/tabwriter"
)

// Stats ...
type Stats struct {
	Latencies *Histogram
	Codes     *shard
	Errors    *shard
}

type stats struct {
	Mean        float64
	Stddev      float64
	Max         float64
	Percentiles map[float64]uint64
}

type shard struct {
	m   map[interface{}]uint
	mux sync.RWMutex
}

// Add ...
func (s *shard) Add(key interface{}, amount uint) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.m[key] += amount
}

// Increment ...
func (s *shard) Increment(key interface{}) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.m[key]++
}

// Decrement ...
func (s *shard) Decrement(key interface{}) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.m[key]--
}

// Count ...
func (s *shard) Count() uint {
	s.mux.Lock()
	defer s.mux.Unlock()

	var count uint
	for _, value := range s.m {
		count += value
	}
	return count
}

func (s *shard) read() map[interface{}]uint {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.m
}

// New ...
func New() *Stats {
	return &Stats{
		Latencies: NewHistogram(),
		Codes: &shard{
			m:   make(map[interface{}]uint),
			mux: sync.RWMutex{},
		},
		Errors: &shard{
			m:   make(map[interface{}]uint),
			mux: sync.RWMutex{},
		},
	}
}

func (s *Stats) PrintCodes() {
	codes := make(sort.IntSlice, 0)
	values := make([]interface{}, 0)
	headers := ""
	line := ""

	m := s.Codes.read()
	for code := range m {
		codes = append(codes, code.(int))
	}
	codes.Sort()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	for _, code := range codes {
		if code != 0 {
			headers += fmt.Sprintf("%s: %d \t", http.StatusText(code), code)
			line += "%d \t"
			values = append(values, m[code])
		} else {
			headers += "Errors \t"
			line += "%d \t"
			values = append(values, m[code])
		}
	}

	fmt.Fprintln(w, headers)
	fmt.Fprintf(w, fmt.Sprintf("%s\n", line), values...)
	w.Flush()
}
