package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/action"
)

var restartCommand = &cobra.Command{
	Use:   "restart",
	Short: "Restart a machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)

		restartAction, err := action.Restart(name)
		if err != nil {
			return err
		}

		return action.Run(action.NewMultipassRunner("multipass"), []action.Action{*restartAction})
	},
}
