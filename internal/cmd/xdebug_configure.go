package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/action"
)

var xdebugConfigureCommand = &cobra.Command{
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
