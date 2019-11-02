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
package method

import (
	"log"
	"net/http"
	"time"

	"github.com/DustyRat/post-it/pkg/file"

	"github.com/DustyRat/post-it/pkg/client"
	"github.com/DustyRat/post-it/pkg/controller"

	"github.com/spf13/cobra"
)

type GetOptions struct {
	input string
	//output string

	batchSize   int
	connections int

	client client.Config

	headers []string
	rawUrl  string

	timeout            time.Duration
	idleTimeout        time.Duration
	insecureSkipVerify bool
}

func NewCmdGet() *cobra.Command {
	options := GetOptions{}
	cmd := &cobra.Command{
		Use:     "GET",
		Aliases: []string{"get"},
		Short:   "The GET method requests a representation of the specified resource.",
		Long:    `The HTTP GET method requests a representation of the specified resource.`,
		Example: "post-it GET -u http://localhost:3000/path/{column_name}",
		Run: func(cmd *cobra.Command, args []string) {
			options.client.Headers = client.ParseHeaders(options.headers)
			clt, err := client.NewClient(options.client)
			if err != nil {
				log.Fatal(err)
			}

			ctrl := controller.Controller{
				Client:   clt,
				Routines: options.connections,
			}

			chunks := file.ParseFile(options.input, http.MethodGet, options.rawUrl, options.batchSize)
			err = ctrl.Run(chunks)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVarP(&options.input, "input", "i", "input.csv", "Input File")

	//cmd.Flags().StringVarP(&options.output, "output", "o", "output.csv", "Output File")
	//cmd.Flags().StringVar(&responseType, "response-type", "", "Response type to output. eg: all, error, status")
	//cmd.Flags().StringVar(&responseStatus, "response-status", "any", "Response status to output. eg: any, 2xx, -2xx (non 2xx statuses), 4xx, 5xx, 200, 301, 404, 503...")
	//cmd.Flags().BoolVar(&recordBody, "record-body", false, "Output body")
	//cmd.Flags().BoolVar(&recordHeaders, "record-headers", false, "Output headers")

	cmd.Flags().IntVarP(&options.batchSize, "batch", "b", 100, "Batch Size")

	cmd.Flags().StringArrayVar(&options.headers, "header", []string{}, "HTTP headers to use (\"K: V\")")
	cmd.Flags().StringVarP(&options.rawUrl, "url", "u", "", "Url. Should be in the format 'http://localhost:3000/path/{column_name}' if input file is specified")
	cmd.MarkFlagRequired("url")

	cmd.Flags().IntVarP(&options.connections, "connections", "c", 10, "connections")
	cmd.Flags().DurationVarP(&options.client.Timeout, "timeout", "t", 3000*time.Millisecond, "Connection timeout")
	cmd.Flags().DurationVar(&options.client.IdleConnTimeout, "idle-timeout", 500*time.Millisecond, "Idle Connection timeout")
	cmd.Flags().BoolVar(&options.client.InsecureSkipVerify, "insecure-skip-verify", true, "Insecure Skip Verify")

	return cmd
}
