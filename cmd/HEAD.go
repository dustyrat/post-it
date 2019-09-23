package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// headCmd represents the HEAD command
var headCmd = &cobra.Command{
	Use:   "HEAD",
	Short: "The HEAD method asks for a response identical to that of a GET request, but without the response body.",
	Long: `The HEAD method asks for a response identical to that of a GET request, but without the response body.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("HEAD called")
	},
}

func init() {
	rootCmd.AddCommand(headCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// headCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// headCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
