/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	"github.com/DustyRat/post-it/cmd/method"
	"github.com/spf13/cobra"
)

type BatchOptions struct {
	Input  string
	Output string

	BatchSize int
	//Connections int
	//
	//Headers []string
	//RawUrl  string
	//
	//Timeout            time.Duration
	//IdleTimeout        time.Duration
	//InsecureSkipVerify bool
}

func NewCmdBatch() *cobra.Command {
	options := BatchOptions{}
	cmd := &cobra.Command{
		Use:     "BATCH",
		Aliases: []string{"batch"},
	}

	cmd.PersistentFlags().StringVarP(&options.Input, "input", "i", "input.csv", "Input File")
	cmd.PersistentFlags().StringVarP(&options.Output, "output", "o", "output.csv", "Output File")
	cmd.PersistentFlags().IntVarP(&options.BatchSize, "batch", "b", 100, "Batch Size")

	////cmd.Flags().StringVar(&responseType, "response-type", "", "Response type to output. eg: all, error, status")
	////cmd.Flags().StringVar(&responseStatus, "response-status", "any", "Response status to output. eg: any, 2xx, -2xx (non 2xx statuses), 4xx, 5xx, 200, 301, 404, 503...")
	////cmd.Flags().BoolVar(&recordBody, "record-body", false, "Output body")
	////cmd.Flags().BoolVar(&recordHeaders, "record-headers", false, "Output headers")
	//
	//
	//cmd.Flags().StringArrayVar(&options.Headers, "header", []string{}, "HTTP headers to use (\"K: V\")")
	//cmd.Flags().StringVarP(&options.RawUrl, "url", "u", "", "Url. Should be in the format 'http://localhost:3000/path/{column_name}' if input file is specified")
	//cmd.MarkFlagRequired("url")
	//
	//cmd.Flags().IntVarP(&options.Connections, "connections", "c", 10, "connections")
	//cmd.Flags().DurationVarP(&options.Timeout, "timeout", "t", 3000*time.Millisecond, "Connection timeout")
	//cmd.Flags().DurationVar(&options.IdleTimeout, "idle-timeout", 500*time.Millisecond, "Idle Connection timeout")
	//cmd.Flags().BoolVar(&options.InsecureSkipVerify, "insecure-skip-verify", true, "Insecure Skip Verify")

	cmd.AddCommand(method.NewCmdGet())
	return cmd
}
