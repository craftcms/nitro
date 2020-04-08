package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/action"
	"github.com/craftcms/nitro/internal/nitro"
)

var updateCommand = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u", "upgrade"},
	Short:   "Update a machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)

		updateAction, err := action.Update(name)
		if err != nil {
			return err
		}

		return nitro.RunAction(nitro.NewMultipassRunner("multipass"), []action.Action{*updateAction})
	},
}
