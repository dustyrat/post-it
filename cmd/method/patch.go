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

// NewCmdPatch ...
func NewCmdPatch(options *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "PATCH",
		Aliases: []string{"patch"},
		Short:   "The PATCH method is used to apply partial modifications to a resource.",
		Example: "post-it PATCH -u http://localhost:3000/path/{column_name}",
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

			headers, requests := file.ParseFile(options.Input, http.MethodPatch, options.RawUrl, options.RequestBody)
			err = ctrl.Run(headers, requests)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	return cmd
}
