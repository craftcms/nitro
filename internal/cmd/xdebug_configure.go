package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var xdebugConfigureCommand = &cobra.Command{
	Use:   "configure",
	Short: "Configure Xdebug on a machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)
		php := config.GetString("php", flagPhpVersion)

		var actions []nitro.Action
		xdebugConfigureAction, err := nitro.ConfigureXdebug(name, php)
		if err != nil {
			return err
		}
		actions = append(actions, *xdebugConfigureAction)

		restartPhpFpmAction, err := nitro.RestartPhpFpm(name, php)
		if err != nil {
			return err
		}
		actions = append(actions, *restartPhpFpmAction)

		return nitro.Run(nitro.NewMultipassRunner("multipass"), actions)
	},
}
