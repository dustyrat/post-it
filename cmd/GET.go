package cmd

import (
	"post-it/pkg/controller"

	"github.com/spf13/cobra"
)

var ctrl controller.Controller

// getCmd represents the GET command
var getCmd = &cobra.Command{
	Use:     "GET",
	Short:   "The GET method requests a representation of the specified resource.",
	Long:    `The HTTP GET method requests a representation of the specified resource.`,
	Example: "post-it GET -u http://localhost:3000/path",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

	},
	PreRun: func(cmd *cobra.Command, args []string) {
		// // API CLIENT
		// clt, err := client.NewClient(client.Config{
		// 	Timeout:             3000,
		// 	InsecureSkipVerify:  true,
		// 	MaxConnsPerHost:     100,
		// 	MaxIdleConns:        10,
		// 	MaxIdleConnsPerHost: 10,
		// 	IdleConnTimeout:     5000,
		// })
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// // CONTROLLER
		// ctrl = controller.Controller{Method: http.MethodGet, Url: URL, Client: clt, BatchSize: BatchSize, Routines: Routines}

		// ctrl.Input, err = os.Open(Input)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// ctrl.Output, err = controller.NewWriter(Output)
		// if err != nil {
		// 	log.Fatal(err)
		// }
	},
	Run: func(cmd *cobra.Command, args []string) {
		// defer ctrl.Input.Close()
		// err := ctrl.Run()
		// if err != nil {
		// 	log.Fatal(err)
		// }
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		// ctrl.WorkerPool.Shutdown()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
