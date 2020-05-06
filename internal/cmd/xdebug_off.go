package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/validate"
)

var xdebugOffCommand = &cobra.Command{
	Use:   "off",
	Short: "Disable Xdebug on a machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

		php := config.GetString("php", flagPhpVersion)
		if err := validate.PHPVersion(php); err != nil {
			return err
		}

		var actions []nitro.Action
		disableXdebugAction, err := nitro.DisableXdebug(machine, php)
		if err != nil {
			return err
		}
		actions = append(actions, *disableXdebugAction)

		restartPhpFpmAction, err := nitro.RestartPhpFpm(machine, php)
		if err != nil {
			return err
		}
		actions = append(actions, *restartPhpFpmAction)

		for _, action := range actions {
			if action.Output != "" {
				fmt.Println(action.Output)
			}
		}

		return nitro.Run(nitro.NewMultipassRunner("multipass"), actions)
	},
}
