package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// putCmd represents the PUT command
var putCmd = &cobra.Command{
	Use:   "PUT",
	Short: "The PUT method replaces all current representations of the target resource with the request payload.",
	Long: `The PUT method replaces all current representations of the target resource with the request payload.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("PUT called")
	},
}

func init() {
	rootCmd.AddCommand(putCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// putCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// putCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
