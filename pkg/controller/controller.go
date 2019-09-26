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
package controller

import (
	"bufio"
	"errors"
	"os"
	"time"

	"github.com/DustyRat/post-it/pkg/client"
	"github.com/DustyRat/post-it/pkg/csv"

	"github.com/cheggaaa/pb/v3"
	"github.com/goinggo/work"
)

type Controller struct {
	Client  *client.Client
	Method  string
	Url     string
	headers []string
	chunks  [][]csv.Record

	WorkerPool *work.Pool
	BatchSize  int
	Routines   int
	Stats      *Stats

	RequestBodyColumn string
	RecordHeaders     bool
	RecordBody        bool
	ResponseType      string
	Status            string

	Input  *os.File
	Output *csv.Writer
}

func (c *Controller) reset() {
	c.Stats = &Stats{
		Responses: make(map[int]int, 0),
		Entries:   make([]Entry, 0),
	}
}

func (c *Controller) Run() error {
	c.reset()
	defer func() {
		c.Stats.Print()
	}()

	from := 1
	to := c.BatchSize
	batch := 1

	input := csv.Parse(bufio.NewReader(c.Input))
	input.Headers = append(input.Headers, "status")
	if c.RecordHeaders {
		input.Headers = append(input.Headers, "headers")
	}
	if c.RecordBody {
		input.Headers = append(input.Headers, "response_body")
	}
	input.Headers = append(input.Headers, "error")

	c.Output.Write(input.Headers)
	c.Output.Flush()

	total := len(input.Records)
	progress := pb.Full.Start(total)
	defer progress.Finish()

	wp, err := work.New(c.Routines, time.Hour*24, func(message string) {})
	if err != nil {
		return errors.New("error creating worker pools")
	}

	var chunks [][]csv.Record
	for c.BatchSize < len(input.Records) {
		input.Records, chunks = input.Records[c.BatchSize:], append(chunks, input.Records[0:c.BatchSize:c.BatchSize])
	}
	chunks = append(chunks, input.Records)
	for i := range chunks {
		w := worker{
			writer:            c.Output,
			client:            c.Client,
			method:            c.Method,
			url:               c.Url,
			chunk:             chunks[i],
			batch:             batch,
			from:              from,
			to:                to,
			stats:             c.Stats,
			progress:          progress,
			status:            c.Status,
			responseType:      c.ResponseType,
			recordBody:        c.RecordBody,
			recordHeaders:     c.RecordHeaders,
			requestBodyColumn: c.RequestBodyColumn,
		}
		wp.Run(&w)

		from = from + c.BatchSize
		to = to + c.BatchSize
		batch++
	}

	// wait for all worker Routines to finish doing their work
	wp.Shutdown()
	return nil
}
