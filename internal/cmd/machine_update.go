package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/action"
)

var updateCommand = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upgrade"},
	Short:   "Update machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)

		var actions []action.Action
		updateAction, err := action.Update(name)
		if err != nil {
			return err
		}
		actions = append(actions, *updateAction)

		upgradeAction, err := action.Upgrade(name)
		if err != nil {
			return err
		}
		actions = append(actions, *upgradeAction)

		return action.Run(action.NewMultipassRunner("multipass"), actions)
	},
}
