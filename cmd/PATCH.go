package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// patchCmd represents the PATCH command
var patchCmd = &cobra.Command{
	Use:   "PATCH",
	Short: "The PATCH method is used to apply partial modifications to a resource.",
	Long:  `The PATCH method is used to apply partial modifications to a resource.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("PATCH called")
	},
}

func init() {
	rootCmd.AddCommand(patchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// patchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// patchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
