package cmd

import (
	"github.com/spf13/cobra"
)

var xdebugCommand = &cobra.Command{
	Use:   "xdebug",
	Short: "Manage Xdebug on machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	xdebugCommand.Flags().StringVarP(&flagPhpVersion, "php-version", "v", "", "Version of PHP to enable/disable Xdebug")
}
