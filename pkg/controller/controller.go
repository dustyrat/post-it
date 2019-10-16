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
	"context"
	"errors"
	"os"
	"time"

	"github.com/goinggo/work"

	"github.com/vbauerster/mpb"

	"github.com/DustyRat/post-it/pkg/client"
	"github.com/DustyRat/post-it/pkg/csv"
)

type Controller struct {
	Client  *client.Client
	Method  string
	Url     string
	headers []string
	chunks  [][]csv.Record

	BatchSize int
	Routines  int
	Stats     *Stats

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

	input := csv.Parse(c.Input)
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

	wp, err := work.New(c.Routines, time.Hour*24, func(message string) {})
	if err != nil {
		return errors.New("error creating worker pools")
	}

	progress := mpb.New(mpb.WithContext(context.Background()))

	pool := NewPool(wp, c.Client, c.Method, c.Url, c.Stats, c.Output, progress)

	var chunks [][]csv.Record
	for c.BatchSize < len(input.Records) {
		input.Records, chunks = input.Records[c.BatchSize:], append(chunks, input.Records[0:c.BatchSize:c.BatchSize])
	}
	chunks = append(chunks, input.Records)
	for i := range chunks {
		pool.NewWorker(
			chunks[i],
			Flags{
				Body:    c.RecordBody,
				Headers: c.RecordHeaders,
			},
			Keys{
				Body:         c.RequestBodyColumn,
				Status:       c.Status,
				ResponseType: c.ResponseType,
			},
		)
	}
	pool.Run()
	return nil
}
