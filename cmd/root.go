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
	}

	options := options.Options{}
	cmd.PersistentFlags().StringVarP(&options.Input, "input", "i", "input.csv", "Input File")
	cmd.PersistentFlags().StringVarP(&options.Output, "output", "o", "output.csv", "Output File")

	cmd.PersistentFlags().BoolVarP(&options.Latency, "latencies", "l", false, "Print latency statistics")
	cmd.PersistentFlags().BoolVarP(&options.Histogram, "histogram", "g", false, "Print histogram statistics")

	cmd.PersistentFlags().StringArrayVarP(&options.Headers, "header", "H", []string{}, "HTTP headers to use (\"K: V\")")
	cmd.PersistentFlags().StringVarP(&options.RawUrl, "url", "u", "", "Url. Should be in the format 'http://localhost:3000/path/{column_name}' if input file is specified")
	cmd.MarkPersistentFlagRequired("url")

	cmd.PersistentFlags().BoolVarP(&options.Flags.Errors, "errors", "e", false, "Record erorrs to output file")
	cmd.PersistentFlags().StringVarP(&options.Flags.Status, "response-status", "s", "-2xx", "Record response status to output file under the headers status. eg: any, none, 2xx, -2xx (non 2xx statuses), 4xx, 5xx, 200, 301, 404, 503...")
	cmd.PersistentFlags().BoolVarP(&options.Flags.Body, "record-body", "b", false, "Record body to output file under the response_body column.")
	// cmd.PersistentFlags().BoolVar(&options.Flags.Headers, "record-headers", false, "Record headers to output file under the headers column.")

	cmd.PersistentFlags().IntVarP(&options.Connections, "connections", "c", 10, "Concurrent connections")
	cmd.PersistentFlags().DurationVarP(&options.Client.Timeout, "timeout", "t", 3000*time.Millisecond, "Connection timeout")
	cmd.PersistentFlags().BoolVar(&options.Client.InsecureSkipVerify, "insecure", true, "Insecure Skip Verify")

	cmd.AddCommand(method.NewCmdDelete(&options))
	cmd.AddCommand(method.NewCmdGet(&options))
	cmd.AddCommand(method.NewCmdHead(&options))
	cmd.AddCommand(method.NewCmdPatch(&options))
	cmd.AddCommand(method.NewCmdPost(&options))
	cmd.AddCommand(method.NewCmdPut(&options))
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
