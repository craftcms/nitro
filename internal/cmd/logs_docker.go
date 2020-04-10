package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/action"
)

var logsDockerCommand = &cobra.Command{
	Use:   "docker",
	Short: "Show docker logs",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)

		dockerLogsAction, err := action.LogsDocker(name, args[0])
		if err != nil {
			return err
		}

		return action.Run(action.NewMultipassRunner("multipass"), []action.Action{*dockerLogsAction})
	},
}
