package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "0.0.0"

func init() {
	rootCmd.AddCommand(versionCommand)
}

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "View Nitro version",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("nitro version %s\n", Version)
		fmt.Println("")
		fmt.Println("visit https://github.com/craftcms/nitro for more information")

		return nil
	},
}
