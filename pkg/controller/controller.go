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
	"time"

	"github.com/DustyRat/post-it/pkg/stats"

	"github.com/DustyRat/post-it/pkg/worker"

	"github.com/goinggo/work"

	"github.com/vbauerster/mpb"

	"github.com/DustyRat/post-it/pkg/client"
)

type Controller struct {
	Client *client.Client
	//Method string
	//Url     string
	//headers []string
	chunks [][]client.Request

	//BatchSize int
	Routines int
	Stats    *stats.Stats

	//RequestBodyColumn string
	//RecordHeaders     bool
	//RecordBody        bool
	//ResponseType      string
	//Status            string
	//
	//Input  *os.File
	//Output *csv.Writer
}

func (c *Controller) Run(requests []*client.Request) error {
	c.Stats = stats.NewStats()
	wp, err := work.New(c.Routines, time.Hour*24, func(message string) {})
	if err != nil {
		return errors.New("error creating worker pools")
	}

	progress := mpb.New(mpb.WithContext(context.Background()))
	pool := worker.NewPool(wp, c.Client, c.Stats, progress, int64(len(requests)))
	for i := range requests {
		pool.NewWorker(requests[i])
	}
	pool.Run()
	c.Stats.Latencies.Print()
	return nil
}
