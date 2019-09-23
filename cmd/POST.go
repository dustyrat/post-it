package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// postCmd represents the POST command
var postCmd = &cobra.Command{
	Use:   "POST",
	Short: "The POST method is used to submit an entity to the specified resource, often causing a change in state or side effects on the server.",
	Long:  `The POST method is used to submit an entity to the specified resource, often causing a change in state or side effects on the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("POST called")
	},
}

func init() {
	rootCmd.AddCommand(postCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// postCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// postCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
