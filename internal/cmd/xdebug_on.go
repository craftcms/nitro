package cmd

import (
	"github.com/spf13/cobra"
)

var xdebugOnCommand = &cobra.Command{
	Use:   "on",
	Short: "Enable Xdebug on a machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		return xonCommand.RunE(cmd, args)
	},
}
