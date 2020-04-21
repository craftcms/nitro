package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/nitro"
)

var restartCommand = &cobra.Command{
	Use:   "restart",
	Short: "Restart a machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := "nitro-dev"
		if flagMachineName != "" {
			machine = flagMachineName
		}

		restartAction, err := nitro.Restart(machine)
		if err != nil {
			return err
		}

		return nitro.Run(nitro.NewMultipassRunner("multipass"), []nitro.Action{*restartAction})
	},
}
