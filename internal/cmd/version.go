package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "View Nitro version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Nitro version %s\n", cmd.Version)

		return nil
	},
}
