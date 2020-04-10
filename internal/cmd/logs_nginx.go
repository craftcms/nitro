package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/action"
)

var logsNginxCommand = &cobra.Command{
	Use:   "nginx",
	Short: "Show nginx logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)

		nginxLogsAction, err := action.LogsNginx(name, flagNginxLogsKind)
		if err != nil {
			return err
		}

		return action.Run(action.NewMultipassRunner("multipass"), []action.Action{*nginxLogsAction})
	},
}

func init() {
	logsNginxCommand.Flags().StringVarP(&flagNginxLogsKind, "type", "t", "all", "filter the logs by kind, access or error (defaults to all)")
}
