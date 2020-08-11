package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/nitro"
)

var restartCommand = &cobra.Command{
	Use:   "restart",
	Short: "Restart machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		runner := nitro.NewMultipassRunner("multipass")

		// check if the machine is running, if not start it
		if nitro.IP(machine, runner) == "" {
			fmt.Println("The " + machine + " machine is not running...")
			if err := startCommand.RunE(cmd, args); err != nil {
				return err
			}

			return nil
		}

		restartAction, err := nitro.Restart(machine)
		if err != nil {
			return err
		}

		if err := nitro.Run(nitro.NewMultipassRunner("multipass"), []nitro.Action{*restartAction}); err != nil {
			return err
		}

		fmt.Println("Restarted", machine)

		return nil
	},
}
