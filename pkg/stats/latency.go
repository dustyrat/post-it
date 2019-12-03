package stats

import (
	"math"
	"sort"
	"sync"
	"time"
)

// Latency ...
type Latency struct {
	m   map[time.Duration]uint
	mux sync.RWMutex
}

// New ...
func NewLatency() *Latency {
	return &Latency{
		m:   make(map[time.Duration]uint),
		mux: sync.RWMutex{},
	}
}

// Gather ...
func (l *Latency) Gather(percentiles []float64) Response {
	stat := Response{
		Percentiles: map[float64]time.Duration{},
	}
	count := uint(0)
	sum := time.Duration(0)
	max := time.Duration(0)
	pairs := make([]struct {
		k time.Duration
		v uint
	}, 0, l.Count())

	for key, value := range l.read() {
		if key > max {
			max = key
		}
		sum += key
		count += value
		pairs = append(pairs, struct {
			k time.Duration
			v uint
		}{key, value})
	}
	stat.Max = max.Truncate(10 * time.Microsecond)

	if count < 1 {
		return stat
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].k < pairs[j].k
	})

	m := make(map[float64]time.Duration)
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
	stat.Percentiles = m

	mean := float64(sum) / float64(count)
	sumOfSquares := float64(0)
	for key := range m {
		sumOfSquares += math.Pow(float64(key)-mean, 2)
	}
	stat.Mean = time.Duration(mean).Truncate(10 * time.Microsecond)

	stddev := 0.0
	if count > 2 {
		stddev = math.Sqrt(sumOfSquares / float64(count))
	}
	stat.Stddev = time.Duration(stddev).Truncate(10 * time.Microsecond)

	for key, value := range stat.Percentiles {
		stat.Percentiles[key] = value.Truncate(10 * time.Microsecond)
	}
	return stat
}

// Add ...
func (l *Latency) Add(key time.Duration, amount uint) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.m[key] += amount
}

// Increment ...
func (l *Latency) Increment(key time.Duration) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.m[key]++
}

// Decrement ...
func (l *Latency) Decrement(key time.Duration) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.m[key]--
}

// Count ...
func (l *Latency) Count() uint {
	l.mux.Lock()
	defer l.mux.Unlock()

	var count uint
	for _, value := range l.m {
		count += value
	}
	return count
}

// Sum ...
func (l *Latency) Sum() time.Duration {
	l.mux.RLock()
	defer l.mux.RUnlock()

	var total time.Duration
	for key := range l.m {
		total += key
	}
	return total
}

func (l *Latency) read() map[time.Duration]uint {
	l.mux.RLock()
	defer l.mux.RUnlock()
	return l.m
}
