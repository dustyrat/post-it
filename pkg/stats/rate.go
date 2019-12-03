package stats

import (
	"math"
	"sort"
	"sync"
)

// Rate ...
type Rate struct {
	m   map[float64]uint
	mux sync.RWMutex
}

// New ...
func NewRate() *Rate {
	return &Rate{
		m:   make(map[float64]uint),
		mux: sync.RWMutex{},
	}
}

// Gather ...
func (r *Rate) Gather(percentiles []float64) Request {
	stat := Request{}
	count := uint(0)
	sum := float64(0)
	max := float64(0)
	pairs := make([]struct {
		k float64
		v uint
	}, 0, r.Count())

	for key, value := range r.read() {
		if key > max {
			max = key
		}
		sum += key
		count += value
		pairs = append(pairs, struct {
			k float64
			v uint
		}{key, value})
	}
	stat.Max = max

	if count < 1 {
		return stat
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].k < pairs[j].k
	})

	m := make(map[float64]float64)
	for _, pc := range percentiles {
		if _, calculated := m[pc]; calculated {
			continue
		}
		if pc < 0 || pc > 1 {
			// Drop percentiles outside of [0, 1] range
			continue
		}
		rank := uint(pc*float64(count) + 0.5)
		total := uint(0)
		for _, p := range pairs {
			total += p.v
			if total >= rank {
				m[pc] = p.k
				break
			}
		}
	}

	mean := sum / float64(count)
	sumOfSquares := float64(0)
	for key := range m {
		sumOfSquares += math.Pow(key-mean, 2)
	}
	stat.Mean = mean

	stddev := 0.0
	if count > 2 {
		stddev = math.Sqrt(sumOfSquares / float64(count))
	}
	stat.Stddev = stddev
	return stat
}

// Add ...
func (r *Rate) Add(key float64, amount uint) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.m[key] += amount
}

// Increment ...
func (r *Rate) Increment(key float64) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.m[key]++
}

// Decrement ...
func (r *Rate) Decrement(key float64) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.m[key]--
}

// Count ...
func (r *Rate) Count() uint {
	r.mux.Lock()
	defer r.mux.Unlock()

	var count uint
	for _, value := range r.m {
		count += value
	}
	return count
}

// Sum ...
func (r *Rate) Sum() float64 {
	r.mux.RLock()
	defer r.mux.RUnlock()

	var total float64
	for key := range r.m {
		total += key
	}
	return total
}

func (r *Rate) read() map[float64]uint {
	r.mux.RLock()
	defer r.mux.RUnlock()
	return r.m
}

func (r *Rate) Read() map[float64]uint {
	return r.read()
}
