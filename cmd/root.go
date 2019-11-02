package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/DustyRat/post-it/pkg/controller"
	"github.com/spf13/cobra"
)

var (
	ctrl               controller.Controller
	body               string
	inputFile          string
	outputFile         string
	batchSize          int
	connections        int
	rawUrl             string
	headers            []string
	responseType       string
	responseStatus     string
	recordBody         bool
	recordHeaders      bool
	timeout            time.Duration
	idleTimeout        time.Duration
	insecureSkipVerify bool
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
	cmd.AddCommand(NewCmdBatch())
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//rootCmd.PersistentFlags().StringVarP(&inputFile, "input", "i", "input.csv", "Input File")
	//rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "output.csv", "Output File")
	//
	//rootCmd.PersistentFlags().StringVar(&responseType, "response-type", "", "Response type to output. eg: all, error, status")
	//rootCmd.PersistentFlags().StringVar(&responseStatus, "response-status", "any", "Response status to output. eg: any, 2xx, -2xx (non 2xx statuses), 4xx, 5xx, 200, 301, 404, 503...")
	//rootCmd.PersistentFlags().BoolVar(&recordBody, "record-body", false, "Output body")
	//rootCmd.PersistentFlags().BoolVar(&recordHeaders, "record-headers", false, "Output headers")
	//
	//rootCmd.PersistentFlags().IntVarP(&batchSize, "batch", "b", 100, "Batch Size")
	//
	//rootCmd.PersistentFlags().StringVarP(&rawUrl, "url", "u", "", "Url. Should be in the format 'http://localhost:3000/path/{column_name}' if input file is specified")
	//rootCmd.PersistentFlags().StringArrayVar(&headers, "header", []string{}, "HTTP headers to use (\"K: V\")")
	//rootCmd.PersistentFlags().IntVarP(&connections, "connections", "c", 10, "connections")
	//rootCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 3000*time.Millisecond, "Connection timeout")
	//rootCmd.PersistentFlags().DurationVar(&idleTimeout, "idle-timeout", 500*time.Millisecond, "Idle Connection timeout")
	//rootCmd.PersistentFlags().BoolVar(&insecureSkipVerify, "insecure-skip-verify", true, "Insecure Skip Verify")
	//
	//rootCmd.MarkPersistentFlagRequired("url")
}
