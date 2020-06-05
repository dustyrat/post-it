package controller

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/DustyRat/post-it/internal/file/csv"
	"github.com/DustyRat/post-it/internal/http"
	"github.com/DustyRat/post-it/internal/options"
	"github.com/DustyRat/post-it/internal/stats"
	"github.com/DustyRat/post-it/internal/worker"

	"github.com/goinggo/work"
	"github.com/vbauerster/mpb/v5"
)

// Controller ...
type Controller struct {
	Options  *options.Options
	Client   *http.Client
	Routines int
	Writer   *csv.Writer
}

// Run ...
func (c *Controller) Run(file, method, rawURL string) error {
	input, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	reader := csv.NewReader(input, method, rawURL, "request_body")
	wp, err := work.New(c.Routines, time.Hour*24, func(message string) {})
	if err != nil {
		return errors.New("error creating worker pools")
	}

	progress := mpb.New()
	pool := worker.NewPool(c.Options, wp, c.Client, progress, reader.Count(), reader, c.Writer)

	headers := reader.Headers()
	if c.Writer != nil {
		headers = append(headers, "status")
		if c.Options.Flags.Headers {
			headers = append(headers, "headers")
		}
		if c.Options.Flags.Body {
			headers = append(headers, "response_body")
		}
		headers = append(headers, "error")
		c.Writer.Write(headers)
	}

	for i := 0; i < reader.Count(); i++ {
		pool.NewWorker()
	}
	elapsed := pool.Run()
	stats.Print(*c.Options, elapsed)
	return nil
}
