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
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Nitro version %s\n", Version)

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCommand)
}
