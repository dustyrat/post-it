package method

import (
	"log"
	"net/http"

	"github.com/DustyRat/post-it/internal/controller"
	"github.com/DustyRat/post-it/internal/file/csv"
	internal "github.com/DustyRat/post-it/internal/http"
	"github.com/DustyRat/post-it/internal/options"
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
			options.Client.Headers = internal.ParseHeaders(options.Headers)
			client, err := internal.New(options.Client)
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
				Client:   client,
				Routines: options.Connections,
				Writer:   writer,
			}

			err = ctrl.Run(options.Input, http.MethodPost, options.RawUrl, options.RequestBody)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	return cmd
}
