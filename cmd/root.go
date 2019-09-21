/*
Copyright © 2019 Dustin Ratcliffe

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Input TODO
	Input string

	// Output TODO
	Output string

	// BatchSize TODO
	BatchSize int

	// Routines TODO
	Routines int

	// URL TODO
	URL string

	// ResponseTypes TODO
	ResponseTypes string

	// ResponseStatus TODO
	ResponseStatus string

	// Timeout TODO
	Timeout int

	// IdleTimeout TODO
	IdleTimeout int

	// InsecureSkipVerify TODO
	InsecureSkipVerify bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "post-it",
	Short: "post-it is a HTTP(S) CLI library for calling a variaty of urls from an input file.",
	Long: `post-it is a HTTP(S) CLI library for calling a variaty of urls from an input file.
TODO Long Description`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.post-it.yaml)")

	rootCmd.PersistentFlags().StringVarP(&Input, "input", "i", "", "Input File")
	rootCmd.PersistentFlags().StringVarP(&Output, "output", "o", "", "Output File")
	rootCmd.PersistentFlags().StringVar(&Output, "response-types", "", "Response types to output. eg: all, error, status")
	rootCmd.PersistentFlags().StringVar(&Output, "response-status", "any", "Response status to output. eg: 1xx, 2xx, 3xx, 4xx, 5xx, 200, 201, 404, 503, any")

	rootCmd.PersistentFlags().IntVarP(&BatchSize, "batch", "b", 100, "Batch Size")
	rootCmd.PersistentFlags().IntVarP(&Routines, "routines", "r", 10, "Routines")

	rootCmd.PersistentFlags().StringVarP(&URL, "url", "u", "", "Url. Should be in the format 'http://localhost:3000/path' or 'http://localhost:3000/path/{{column_name}}' if input file is specified")
	rootCmd.PersistentFlags().IntVarP(&Timeout, "timeout", "t", 3000, "Connection Timeout")
	rootCmd.PersistentFlags().IntVar(&IdleTimeout, "idle-timeout", 5000, "Idle Connection Timeout")
	rootCmd.PersistentFlags().BoolVar(&InsecureSkipVerify, "insecure-skip-verify", true, "Insecure Skip Verify")

	rootCmd.MarkPersistentFlagRequired("url")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
