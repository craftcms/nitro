package cmd

import (
	"github.com/spf13/cobra"
)

var machineCommand = &cobra.Command{
	Use:    "machine",
	Short:  "Perform actions on the multipass machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}
