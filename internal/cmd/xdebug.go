package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/action"
	"github.com/craftcms/nitro/validate"
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
			name := config.GetString("name", flagMachineName)
			php := config.GetString("php", flagPhpVersion)
			if err := validate.PHPVersion(php); err != nil {
				return err
			}

			enableXdebugAction, err := action.EnableXdebug(name, php)
			if err != nil {
				return err
			}

			return action.Run(action.NewMultipassRunner("multipass"), []action.Action{*enableXdebugAction})
		},
	}

	xdebugOffCommand = &cobra.Command{
		Use:   "off",
		Short: "Disable Xdebug on machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := config.GetString("name", flagMachineName)
			php := config.GetString("php", flagPhpVersion)
			if err := validate.PHPVersion(php); err != nil {
				return err
			}

			disableXdebugAction, err := action.DisableXdebug(name, php)
			if err != nil {
				return err
			}

			return action.Run(action.NewMultipassRunner("multipass"), []action.Action{*disableXdebugAction})
		},
	}

	xdebugConfigureCommand = &cobra.Command{
		Use:   "configure",
		Short: "Configure xdebug on machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := config.GetString("name", flagMachineName)
			php := config.GetString("php", flagPhpVersion)

			var actions []action.Action
			xdebugConfigureAction, err := action.ConfigureXdebug(name, php)
			if err != nil {
				return err
			}
			actions = append(actions, *xdebugConfigureAction)

			restartPhpFpmAction, err := action.RestartPhpFpm(name, php)
			if err != nil {
				return err
			}
			actions = append(actions, *restartPhpFpmAction)

			return action.Run(action.NewMultipassRunner("multipass"), actions)
		},
	}
)

func init() {
	xdebugCommand.Flags().StringVarP(&flagPhpVersion, "php-version", "v", "", "version of PHP to enable/disable xdebug")
}
