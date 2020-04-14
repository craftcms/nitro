package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var updateCommand = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upgrade"},
	Short:   "Update a machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)

		var actions []nitro.Action
		updateAction, err := nitro.Update(name)
		if err != nil {
			return err
		}
		actions = append(actions, *updateAction)

		upgradeAction, err := nitro.Upgrade(name)
		if err != nil {
			return err
		}
		actions = append(actions, *upgradeAction)

		return nitro.Run(nitro.NewMultipassRunner("multipass"), actions)
	},
}
