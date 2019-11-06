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
	"errors"
	"log"
	"net/http"

	"github.com/DustyRat/post-it/pkg/client"
	"github.com/DustyRat/post-it/pkg/controller"
	"github.com/DustyRat/post-it/pkg/file"
	"github.com/DustyRat/post-it/pkg/file/csv"
	"github.com/DustyRat/post-it/pkg/options"
	"github.com/spf13/cobra"
)

func NewCmdHead(options *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "HEAD",
		Aliases: []string{"head"},
		Short:   "The HEAD method asks for a response identical to that of a GET request, but without the response body.",
		Long:    `The HEAD method asks for a response identical to that of a GET request, but without the response body.`,
		Example: "post-it HEAD -u http://localhost:3000/path/{column_name}",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing url")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			options.RawUrl = args[0]
			options.Client.Headers = client.ParseHeaders(options.Headers)
			clt, err := client.NewClient(options.Client)
			if err != nil {
				log.Fatal(err)
			}

			var writer *csv.Writer
			if options.Output != "" {
				writer, err = csv.NewWriter(options.Output)
				if err != nil {
					log.Fatal(err)
				}
			}

			ctrl := controller.Controller{
				Options:  options,
				Client:   clt,
				Routines: options.Connections,
				Writer:   writer,
			}

			headers, requests := file.ParseFile(options.Input, http.MethodHead, options.RawUrl, options.RequestBody)
			err = ctrl.Run(headers, requests)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	return cmd
}
