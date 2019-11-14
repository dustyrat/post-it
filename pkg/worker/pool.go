package worker

import (
	"sync"

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
	options *options.Options
	client  *client.Client

	progress *mpb.Progress
	bar      *mpb.Bar
	stats    *stats.Stats

	workerPool *work.Pool
	workers    []*worker

	writer *csv.Writer
	mutex  *sync.Mutex
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
		mutex:      &sync.Mutex{},
	}
}

// NewWorker ...
func (p *Pool) NewWorker(data *file.Data) *worker {
	p.mutex.Lock()
	w := &worker{pool: p, progress: p.bar, record: data.Record, request: data.Request}
	p.workers = append(p.workers, w)
	p.mutex.Unlock()
	return w
}

// Run ...
func (p *Pool) Run() {
	for _, worker := range p.workers {
		p.workerPool.Run(worker)
	}
	p.workerPool.Shutdown()
	p.progress.Wait()
}
