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
	"fmt"

	"github.com/spf13/cobra"
)

// GETCmd represents the GET command
var GETCmd = &cobra.Command{
	Use:     "GET",
	Short:   "The GET method requests a representation of the specified resource.",
	Long:    `The HTTP GET method requests a representation of the specified resource.`,
	Example: "post-it GET -u http://localhost:3000/path",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("GET called")
		// if URL == "" {
		// 	fmt.Println("No url provided")
		// 	os.Exit(1)
		// }

		// // API CLIENT
		// clt, err := client.NewClient(client.Config{
		// Timeout:             3000,
		// InsecureSkipVerify:  true,
		// MaxConnsPerHost:     100,
		// MaxIdleConns:        10,
		// MaxIdleConnsPerHost: 10,
		// IdleConnTimeout:     5000,
		// })
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// // CONTROLLER
		// ctrl := controller.Controller{Method: "get", Url: *url, Client: clt, BatchSize: *batchSize, Routines: *routines}
		// err = ctrl.Run(*input, *output)
		// if err != nil {
		// 	log.Fatal(err)
		// }
	},
}

func init() {
	rootCmd.AddCommand(GETCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// GETCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// GETCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
