package stats

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

// Histogram ...
type Histogram struct {
	m   map[time.Duration]uint
	mux sync.RWMutex
}

// NewHistogram ...
func NewHistogram() *Histogram {
	return &Histogram{
		m:   make(map[time.Duration]uint),
		mux: sync.RWMutex{},
	}
}

// Print ...
func (h *Histogram) Print() {
	count := uint(0)
	sum := time.Duration(0)
	max := time.Duration(0)
	pairs := make([]struct {
		k time.Duration
		v uint
	}, 0, h.Count())

	m := h.read()
	for key, value := range m {
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

	if count < 1 {
		return
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].k < pairs[j].k
	})

	percentiles := []float64{0.5, 0.75, 0.9, 0.95, 0.99}
	percentilesMap := make(map[float64]time.Duration)
	for _, pc := range percentiles {
		if _, calculated := percentilesMap[pc]; calculated {
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
				percentilesMap[pc] = p.k
				break
			}
		}
	}

	mean := float64(sum) / float64(count)
	sumOfSquares := float64(0)
	for key := range m {
		sumOfSquares += math.Pow(float64(key)-mean, 2)
	}

	stddev := 0.0
	if count > 2 {
		stddev = math.Sqrt(sumOfSquares / float64(count))
	}

	fmt.Printf("%s | %s | %s | \n", time.Duration(mean), time.Duration(stddev), max)
	for _, key := range percentiles {
		value := percentilesMap[key]
		fmt.Printf("%.2f: %s\n", key, value.Truncate(time.Millisecond))
	}
}

// Add ...
func (h *Histogram) Add(key time.Duration, amount uint) {
	h.mux.Lock()
	defer h.mux.Unlock()
	h.m[key] += amount
}

// Increment ...
func (h *Histogram) Increment(key time.Duration) {
	h.mux.Lock()
	defer h.mux.Unlock()
	h.m[key]++
}

// Decrement ...
func (h *Histogram) Decrement(key time.Duration) {
	h.mux.Lock()
	defer h.mux.Unlock()
	h.m[key]--
}

// Count ...
func (h *Histogram) Count() uint {
	h.mux.Lock()
	defer h.mux.Unlock()

	var count uint
	for _, value := range h.m {
		count += value
	}
	return count
}

// Sum ...
func (h *Histogram) Sum() time.Duration {
	h.mux.RLock()
	defer h.mux.RUnlock()

	var total time.Duration
	for key := range h.m {
		total += key
	}
	return total
}

func (h *Histogram) read() map[time.Duration]uint {
	h.mux.RLock()
	defer h.mux.RUnlock()
	return h.m
}
