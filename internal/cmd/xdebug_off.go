package cmd

import (
	"github.com/spf13/cobra"
)

var xdebugOffCommand = &cobra.Command{
	Use:   "off",
	Short: "Disable Xdebug on a machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		return xoffCommand.RunE(cmd, args)
	},
}
