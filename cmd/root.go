package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"post-it/pkg/client"
	"post-it/pkg/controller"
	"post-it/pkg/csv"
	"strings"
	"time"

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
	responseTypes      string
	responseStatus     string
	timeout            time.Duration
	idleTimeout        time.Duration
	insecureSkipVerify bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "post-it",
	Short: "post-it is a HTTP(S) CLI library for calling a variaty of urls from an input file.",
	Long: `post-it is a HTTP(S) CLI library for calling a variaty of urls from an input file.
TODO Long Description`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&inputFile, "input", "i", "input.csv", "Input File")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "output.csv", "Output File")
	//rootCmd.PersistentFlags().StringVar(&responseTypes, "response-types", "", "Response types to output. eg: all, error, status")
	//rootCmd.PersistentFlags().StringVar(&responseStatus, "response-status", "any", "Response status to output. eg: any, 2xx, 4xx, 5xx, 200, 301, 404, 503...")

	rootCmd.PersistentFlags().IntVarP(&batchSize, "batch", "b", 100, "Batch Size")

	rootCmd.PersistentFlags().StringVarP(&rawUrl, "url", "u", "", "Url. Should be in the format 'http://localhost:3000/path/{column_name}' if input file is specified")
	rootCmd.PersistentFlags().StringArrayVar(&headers, "header", []string{}, "Header")
	rootCmd.PersistentFlags().IntVarP(&connections, "connections", "c", 10, "connections")
	rootCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 3000*time.Millisecond, "Connection timeout")
	rootCmd.PersistentFlags().DurationVar(&idleTimeout, "idle-timeout", 500*time.Millisecond, "Idle Connection timeout")
	rootCmd.PersistentFlags().BoolVar(&insecureSkipVerify, "insecure-skip-verify", true, "Insecure Skip Verify")

	rootCmd.MarkPersistentFlagRequired("url")
}

func getController(method string) controller.Controller {
	ctrl := controller.Controller{Method: method, Url: rawUrl, Client: getClient(), BatchSize: batchSize, Routines: connections}
	input, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	ctrl.Input = input

	output, err := csv.NewWriter(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	ctrl.Output = output
	return ctrl
}

func getClient() *client.Client {
	clt, err := client.NewClient(client.Config{
		Headers:            parseHeaders(headers),
		Timeout:            timeout,
		InsecureSkipVerify: insecureSkipVerify,
		MaxConnsPerHost:    connections,
		IdleConnTimeout:    idleTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	return clt
}

func parseHeaders(headers []string) http.Header {
	header := http.Header{}
	for _, h := range headers {
		head := strings.Split(h, ":")
		for _, v := range strings.Split(strings.TrimSpace(head[1]), ",") {
			header.Add(head[0], strings.TrimSpace(v))
		}
	}
	return header
}
