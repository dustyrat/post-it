package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/DustyRat/post-it/cmd/method"
	"github.com/DustyRat/post-it/pkg/options"
	"github.com/spf13/cobra"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cmd := &cobra.Command{
		Use:   "post-it",
		Short: "post-it is a HTTP(S) CLI library for calling a variaty of urls from an input file.",
		Long: `post-it is a HTTP(S) CLI library for calling a variaty of urls from an input file.
TODO Long Description`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//Run: func(cmd *cobra.Command, args []string) { },
	}

	options := options.Options{}
	cmd.PersistentFlags().StringVarP(&options.Input, "input", "i", "input.csv", "Input File")
	cmd.PersistentFlags().StringVarP(&options.Output, "output", "o", "output.csv", "Output File")

	cmd.PersistentFlags().StringArrayVar(&options.Headers, "header", []string{}, "HTTP headers to use (\"K: V\")")
	//cmd.PersistentFlags().StringVarP(&options.RawUrl, "url", "u", "", "Url. Should be in the format 'http://localhost:3000/path/{column_name}' if input file is specified")
	//cmd.MarkPersistentFlagRequired("url")

	cmd.PersistentFlags().StringVar(&options.Flags.Type, "response-type", "none", "Response type to output. eg: all, error, status, none")
	cmd.PersistentFlags().StringVar(&options.Flags.Status, "response-status", "any", "Response status to output. eg: any, 2xx, -2xx (non 2xx statuses), 4xx, 5xx, 200, 301, 404, 503...")
	cmd.PersistentFlags().BoolVar(&options.Flags.Body, "record-body", false, "Output body")
	cmd.PersistentFlags().BoolVar(&options.Flags.Headers, "record-headers", false, "Output headers")

	cmd.PersistentFlags().IntVarP(&options.Connections, "connections", "c", 10, "connections")
	cmd.PersistentFlags().DurationVarP(&options.Client.Timeout, "timeout", "t", 3000*time.Millisecond, "Connection timeout")
	cmd.PersistentFlags().DurationVar(&options.Client.IdleConnTimeout, "idle-timeout", 500*time.Millisecond, "Idle Connection timeout")
	cmd.PersistentFlags().BoolVar(&options.Client.InsecureSkipVerify, "insecure-skip-verify", true, "Insecure Skip Verify")

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
