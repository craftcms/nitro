package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
)

var (
	xdebugCommand = &cobra.Command{
		Use:   "xdebug",
		Short: "Manage Xdebug on machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	xdebugOnCommand = &cobra.Command{
		Use:     "on",
		Aliases: []string{"xon"},
		Short:   "Enable XDebug on machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = config.GetString("machine", flagMachineName)
			_ = config.GetString("php", flagPhpVersion)

			return errors.New("not implemented")
		},
	}

	xdebugOffCommand = &cobra.Command{
		Use:   "off",
		Short: "Disable Xdebug on machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = config.GetString("name", flagMachineName)
			_ = config.GetString("php", flagPhpVersion)

			return errors.New("not implemented")
		},
	}
)

func init() {
	xdebugCommand.Flags().StringVarP(&flagPhpVersion, "php-version", "v", "", "version of PHP to enable/disable xdebug")
}
