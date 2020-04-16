package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version = "0.0.0"

	versionCommand = &cobra.Command{
		Use:   "version",
		Short: "View Nitro version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("nitro version %s\n", Version)
			fmt.Println("")
			fmt.Println("visit https://github.com/craftcms/nitro for more information")

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCommand)
}
