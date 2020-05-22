package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/validate"
)

var xdebugOnCommand = &cobra.Command{
	Use:     "on",
	Aliases: []string{"xon"},
	Short:   "Enable Xdebug on a machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

		php := config.GetString("php", flagPhpVersion)
		if err := validate.PHPVersion(php); err != nil {
			return err
		}

		enableXdebugAction, err := nitro.EnableXdebug(machine, php)
		if err != nil {
			return err
		}

		actions := []nitro.Action{*enableXdebugAction}

		restartPhpFpmAction, err := nitro.RestartPhpFpm(machine, php)
		if err != nil {
			return err
		}
		actions = append(actions, *restartPhpFpmAction)

		fmt.Println("Enabling xdebug for", php, "on", machine)

		if err := nitro.Run(nitro.NewMultipassRunner("multipass"), actions); err != nil {
			return err
		}

		fmt.Println("Xdebug was enabled for PHP", php, "on", machine)

		return nil
	},
}
