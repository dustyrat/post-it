
package stats

import (
	"sync"
)

// Stats ...
type Stats struct {
	Latencies *Histogram
	Codes     *shard
	Errors    *shard

	//Requests  Histogram

	//Responses Histogram
	//Errors    Histogram

	//template *template.Template
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

// NewStats ...
func NewStats() *Stats {
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
