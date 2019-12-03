package worker

import (
	"sync"
	"time"

	"github.com/DustyRat/post-it/pkg/options"

	"github.com/DustyRat/post-it/pkg/file"
	"github.com/DustyRat/post-it/pkg/file/csv"

	"github.com/DustyRat/post-it/pkg/stats"

	"github.com/DustyRat/post-it/pkg/client"
	"github.com/goinggo/work"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

// Pool ...
type Pool struct {
	start    time.Time
	time     time.Time
	requests int64

	options *options.Options
	client  *client.Client

	progress *mpb.Progress
	bar      *mpb.Bar
	stats    *stats.Stats

	workerPool *work.Pool
	workers    []*worker

	writer *csv.Writer
	mux    *sync.Mutex

	done chan interface{}
}

// NewPool ...
func NewPool(options *options.Options, workerPool *work.Pool, client *client.Client, stats *stats.Stats, progress *mpb.Progress, total int64, writer *csv.Writer) *Pool {
	return &Pool{
		options:  options,
		client:   client,
		stats:    stats,
		progress: progress,
		bar: progress.AddBar(total,
			mpb.BarID(0),
			mpb.PrependDecorators(
				decor.Counters(0, "%d / %d", decor.WCSyncSpaceR),
			),
			mpb.AppendDecorators(
				decor.OnComplete(decor.Percentage(decor.WCSyncSpaceR), "complete"),
				decor.AverageSpeed(0, "% .1f/s", decor.WCSyncSpaceR),
				decor.Name("Elapsed:", decor.WCSyncSpaceR),
				decor.Elapsed(decor.ET_STYLE_GO, decor.WCSyncSpaceR),
				decor.OnComplete(decor.Name("ETA:", decor.WCSyncSpaceR), ""),
				decor.OnComplete(decor.AverageETA(decor.ET_STYLE_GO, decor.WCSyncSpaceR), ""),
			),
		),
		workerPool: workerPool,
		writer:     writer,
		mux:        &sync.Mutex{},
		done:       make(chan interface{}),
	}
}

// NewWorker ...
func (p *Pool) NewWorker(data *file.Data) *worker {
	p.mux.Lock()
	w := &worker{pool: p, progress: p.bar, record: data.Record, request: data.Request}
	p.workers = append(p.workers, w)
	p.mux.Unlock()
	return w
}

// Run ...
func (p *Pool) Run() {
	defer close(p.done)
	p.start = time.Now()
	p.time = p.start
	go p.rateMeter()

	for _, worker := range p.workers {
		p.workerPool.Run(worker)
	}
	p.workerPool.Shutdown()
	p.progress.Wait()
	p.done <- struct{}{}
}

func (p *Pool) increment() {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.requests++
}

func (p *Pool) rateMeter() {
	interval := 10 * time.Millisecond
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	tick := ticker.C
	done := p.done
	for {
		select {
		case <-tick:
			rate := p.rps()
			p.stats.Rate.Increment(rate)
			continue
		case <-done:
			rate := p.rps()
			p.stats.Rate.Increment(rate)
			return
		}
	}
}

func (p *Pool) rps() float64 {
	p.mux.Lock()
	defer p.mux.Unlock()
	duration, requests := time.Since(p.time), p.requests
	p.requests = 0
	p.time = time.Now()
	return float64(requests) / duration.Seconds()
}
