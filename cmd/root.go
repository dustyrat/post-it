package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/DustyRat/post-it/cmd/method"

	"github.com/DustyRat/post-it/internal/options"
	"github.com/spf13/cobra"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cmd := &cobra.Command{
		Use:   "post-it",
		Short: "post-it is a HTTP(S) CLI library for calling a variaty of urls from an input file.",
		Long: `post-it is a HTTP(S) CLI library for calling a variaty of urls from an input file.
		
All methods use the request_body column for requests.
		`,
	}

	opts := options.Options{}
	cmd.PersistentFlags().IntVarP(&opts.Connections, "connections", "c", 10, "Concurrent connections")
	cmd.PersistentFlags().BoolVarP(&opts.Flags.Errors, "errors", "e", false, "Record erorrs to output file")
	cmd.PersistentFlags().StringArrayVarP(&opts.Headers, "header", "H", []string{}, "HTTP headers to use (\"K: V\")")
	cmd.PersistentFlags().BoolVarP(&opts.Histogram, "histogram", "g", false, "Print histogram statistics")
	cmd.PersistentFlags().StringVarP(&opts.Input, "input", "i", "input.csv", "Input File")
	cmd.PersistentFlags().BoolVar(&opts.Client.InsecureSkipVerify, "insecure", true, "Insecure Skip Verify")
	cmd.PersistentFlags().BoolVarP(&opts.Latency, "latencies", "l", false, "Print latency statistics")
	cmd.PersistentFlags().StringVarP(&opts.Output, "output", "o", "output.csv", "Output File")
	cmd.PersistentFlags().BoolVarP(&opts.Flags.Body, "record-body", "b", false, "Record body to output file under the response_body column.")
	cmd.PersistentFlags().BoolVar(&opts.Flags.Headers, "record-headers", false, "Record headers to output file under the headers column.")
	cmd.PersistentFlags().StringVarP(&opts.Flags.Status, "response-status", "s", "-2xx", "Record response status to output file under the headers status. eg: any, none, 2xx, -2xx (non 2xx statuses), 4xx, 5xx, 200, 301, 404, 503...")
	cmd.PersistentFlags().DurationVarP(&opts.Client.Timeout, "timeout", "t", 3000*time.Millisecond, "Connection timeout")

	cmd.AddCommand(method.NewCmdDelete(&opts))
	cmd.AddCommand(method.NewCmdGet(&opts))
	cmd.AddCommand(method.NewCmdHead(&opts))
	cmd.AddCommand(method.NewCmdPatch(&opts))
	cmd.AddCommand(method.NewCmdPost(&opts))
	cmd.AddCommand(method.NewCmdPut(&opts))
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
