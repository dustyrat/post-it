/*
Copyright © 2019 Dustin Ratcliffe <dustin.k.ratcliffe@gmail.com>

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
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// patchCmd represents the PATCH command
var patchCmd = &cobra.Command{
	Use:   "PATCH",
	Short: "The PATCH method is used to apply partial modifications to a resource.",
	Long:  `The PATCH method is used to apply partial modifications to a resource.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// CONTROLLER
		ctrl = getController(http.MethodPatch)
	},
	Run: func(cmd *cobra.Command, args []string) {
		defer ctrl.Input.Close()
		err := ctrl.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(patchCmd)
	patchCmd.PersistentFlags().StringVar(&body, "body-column", "request_body", "Column name to use for the request body.")
}
