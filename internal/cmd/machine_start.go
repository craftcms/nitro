package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/action"
)

var startCommand = &cobra.Command{
	Use:     "start",
	Aliases: []string{"up"},
	Short:   "Start machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)

		startAction, err := action.Start(name)
		if err != nil {
			return err
		}

		return Run(NewMultipassRunner("multipass"), []action.Action{*startAction})
	},
}
