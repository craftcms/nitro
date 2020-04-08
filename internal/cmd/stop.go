package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/action"
)

var stopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stop machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("machine", flagMachineName)

		stopAction, err := action.Stop(name)
		if err != nil {
			return err
		}

		return Run(NewMultipassRunner("multipass"), []action.Action{*stopAction})
	},
}
