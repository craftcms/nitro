package cmd

import (
	"github.com/spf13/cobra"
)

var machineCommand = &cobra.Command{
	Use:    "machine",
	Short:  "Manage multipass machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}
