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

// NewCmdPatch ...
func NewCmdPatch(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "PATCH",
		Aliases: []string{"patch"},
		Short:   "The PATCH method is used to apply partial modifications to a resource.",
		Example: "post-it PATCH -u http://localhost:3000/path/{column_name}",
		Run: func(cmd *cobra.Command, args []string) {
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

			err = ctrl.Run(opts.Input, http.MethodPatch, opts.RawUrl)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	return cmd
}
