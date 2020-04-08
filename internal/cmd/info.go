package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/action"
)

var infoCommand = &cobra.Command{
	Use:   "info",
	Short: "Show machine info",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)

		infoAction, err := action.Info(name)
		if err != nil {
			return err
		}

		return Run(NewMultipassRunner("multipass"), []action.Action{*infoAction})
	},
}
