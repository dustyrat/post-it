package method

import (
	"log"
	"net/http"

	"github.com/DustyRat/post-it/pkg/client"
	"github.com/DustyRat/post-it/pkg/controller"
	"github.com/DustyRat/post-it/pkg/file"
	"github.com/DustyRat/post-it/pkg/file/csv"
	"github.com/DustyRat/post-it/pkg/options"
	"github.com/spf13/cobra"
)

// NewCmdPost ...
func NewCmdPost(options *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "POST",
		Aliases: []string{"post"},
		Short:   "The POST method is used to submit an entity to the specified resource, often causing a change in state or side effects on the server.",
		Example: "post-it POST -u http://localhost:3000/path/{column_name}",
		Run: func(cmd *cobra.Command, args []string) {
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

			headers, requests := file.ParseFile(options.Input, http.MethodPost, options.RawUrl, options.RequestBody)
			err = ctrl.Run(headers, requests)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	return cmd
}
