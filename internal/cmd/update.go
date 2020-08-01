package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/nitro"
)

// TODO move this to the nitrod API
var updateCommand = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upgrade"},
	Short:   "Update machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

		var actions []nitro.Action
		updateAction, err := nitro.Update(machine)
		if err != nil {
			return err
		}
		actions = append(actions, *updateAction)

		upgradeAction, err := nitro.Upgrade(machine)
		if err != nil {
			return err
		}
		actions = append(actions, *upgradeAction)

		if err := nitro.Run(nitro.NewMultipassRunner("multipass"), actions); err != nil {
			return err
		}

		fmt.Println("Updated", machine)

		return nil
	},
}
