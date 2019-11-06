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

func (s *shard) Add(key interface{}, amount uint) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.m[key] += amount
}

func (s *shard) Increment(key interface{}) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.m[key]++
}

func (s *shard) Decrement(key interface{}) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.m[key]--
}

func (s *shard) Count() uint {
	s.mux.Lock()
	defer s.mux.Unlock()

	var count uint
	for _, value := range s.m {
		count += value
	}
	return count
}

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
