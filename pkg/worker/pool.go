package worker

import (
	"fmt"
	"sync"
	"time"

	"github.com/DustyRat/post-it/pkg/stats"

	"github.com/DustyRat/post-it/pkg/client"
	"github.com/goinggo/work"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

type Pool struct {
	client *client.Client

	progress *mpb.Progress
	bar      *mpb.Bar
	stats    *stats.Stats

	workerPool *work.Pool
	workers    []*worker
	total      int64

	mutex *sync.Mutex
}

func NewPool(workerPool *work.Pool, client *client.Client, stats *stats.Stats, progress *mpb.Progress) *Pool {
	return &Pool{
		client:   client,
		stats:    stats,
		progress: progress,
		bar: progress.AddBar(0,
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
		mutex:      &sync.Mutex{},
	}
}

func (p *Pool) NewWorker(requests []*client.Request) *worker {
	p.mutex.Lock()
	p.total += int64(len(requests))
	p.bar.SetTotal(p.total, false)
	w := &worker{pool: p, progress: p.bar, requests: requests}
	p.workers = append(p.workers, w)
	p.mutex.Unlock()
	return w
}

func (p *Pool) Run() {
	start := time.Now()
	defer func() {
		p.stats.Print()
		elapsed := time.Now().Sub(start)
		fmt.Printf("%d / %s | %.2f/sec\n", p.total, elapsed.Round(10*time.Millisecond), float64(p.total)/elapsed.Seconds())
		p.stats.Latencies.Print()
	}()

	for _, worker := range p.workers {
		p.workerPool.Run(worker)
	}
	p.workerPool.Shutdown()
	p.progress.Wait()
}
