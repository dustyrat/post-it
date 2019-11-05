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
	"github.com/DustyRat/post-it/pkg/options"
	"github.com/spf13/cobra"
)

func NewCmdPost(options *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "POST",
		Aliases: []string{"post"},
		Short:   "The POST method is used to submit an entity to the specified resource, often causing a change in state or side effects on the server.",
		Long:    `The POST method is used to submit an entity to the specified resource, often causing a change in state or side effects on the server.`,
		Example: "post-it POST -u http://localhost:3000/path/{column_name}",
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

			ctrl := controller.Controller{
				Client:   clt,
				Routines: options.Connections,
			}

			requests := file.ParseFile(options.Input, http.MethodPost, options.RawUrl, options.RequestBody)
			err = ctrl.Run(requests)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	return cmd
}
