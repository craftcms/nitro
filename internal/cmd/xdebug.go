package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var (
	xdebugCommand = &cobra.Command{
		Use:     "xdebug",
		Aliases: []string{"x"},
		Short:   "Perform Xdebug operations on machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	xdebugOnCommand = &cobra.Command{
		Use:   "on",
		Short: "Enable XDebug on machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := config.GetString("machine", flagMachineName)
			php := config.GetString("php", flagPhpVersion)

			if err := nitro.Run(
				nitro.NewMultipassRunner("multipass"),
				nitro.EnableXdebug(name, php),
			); err != nil {
				return err
			}

			return nil
		},
	}

	xdebugOffCommand = &cobra.Command{
		Use:   "off",
		Short: "Disable Xdebug on machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := config.GetString("machine", flagMachineName)
			php := config.GetString("php", flagPhpVersion)

			if err := nitro.Run(
				nitro.NewMultipassRunner("multipass"),
				nitro.DisableXdebug(name, php),
			); err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	xdebugCommand.Flags().StringVarP(&flagPhpVersion, "php-version", "v", "", "version of PHP to enable/disable xdebug")
}
