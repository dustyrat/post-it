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
func NewCmdPost(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "POST <url>",
		Aliases: []string{"post"},
		Args:    cobra.ExactArgs(1),
		Short:   "The POST method is used to submit an entity to the specified resource, often causing a change in state or side effects on the server.",
		Example: "post-it POST http://localhost:3000/path/{column_name}",
		Run: func(cmd *cobra.Command, args []string) {
			opts.RawUrl = args[0]
			opts.Client.Headers = internal.ParseHeaders(opts.Headers)
			client, err := internal.New(opts.Client)
			if err != nil {
				log.Fatal(err)
			}

			var writer *csv.Writer
			if opts.Output != "" {
				writer, err = csv.NewWriter(opts.Output)
				if err != nil {
					log.Fatal(err)
				}
			}

			ctrl := controller.Controller{
				Options:  opts,
				Client:   client,
				Routines: opts.Connections,
				Writer:   writer,
			}

			err = ctrl.Run(opts.Input, http.MethodPost, opts.RawUrl)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	return cmd
}
