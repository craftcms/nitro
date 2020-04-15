package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "View nitro version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("nitro version %s\n", cmd.Version)

		return nil
	},
}
