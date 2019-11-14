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

// NewCmdDelete ...
func NewCmdDelete(options *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "DELETE",
		Aliases: []string{"delete"},
		Short:   "The DELETE method deletes the specified resource.",
		Long:    `The DELETE method deletes the specified resource.`,
		Example: "post-it DELETE -u http://localhost:3000/path/{column_name}",
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

			headers, requests := file.ParseFile(options.Input, http.MethodDelete, options.RawUrl, options.RequestBody)
			err = ctrl.Run(headers, requests)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	return cmd
}
