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
