package worker

import (
	"sync"
	"time"

	"github.com/DustyRat/post-it/internal/file/csv"
	"github.com/DustyRat/post-it/internal/http"
	"github.com/DustyRat/post-it/internal/options"

	"github.com/goinggo/work"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

// Pool ...
type Pool struct {
	requests int64

	options *options.Options
	client  *http.Client

	progress *mpb.Progress
	bar      *mpb.Bar

	pool    *work.Pool
	workers []*worker

	reader *csv.Reader
	writer *csv.Writer
	mux    *sync.Mutex
}

// NewPool ...
func NewPool(opts *options.Options, pool *work.Pool, client *http.Client, progress *mpb.Progress, total int, reader *csv.Reader, writer *csv.Writer) *Pool {
	return &Pool{
		options:  opts,
		client:   client,
		progress: progress,
		bar: progress.AddBar(int64(total),
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
		pool:   pool,
		reader: reader,
		writer: writer,
		mux:    &sync.Mutex{},
	}
}

// NewWorker ...
func (p *Pool) NewWorker() *worker {
	p.mux.Lock()
	defer p.mux.Unlock()
	record := p.reader.Read()
	if record != nil {
		w := &worker{pool: p, progress: p.bar, record: record, request: record.Request}
		p.workers = append(p.workers, w)
		return w
	}
	return nil
}

// Run ...
func (p *Pool) Run() time.Duration {
	start := time.Now()

	for _, worker := range p.workers {
		p.pool.Run(worker)
	}
	p.pool.Shutdown()
	p.progress.Wait()
	return time.Since(start)
}

func (p *Pool) increment() {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.requests++
}
