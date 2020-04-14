package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var logsNginxCommand = &cobra.Command{
	Use:   "nginx",
	Short: "Show nginx logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)

		nginxLogsAction, err := nitro.LogsNginx(name, flagNginxLogsKind)
		if err != nil {
			return err
		}

		return nitro.Run(nitro.NewMultipassRunner("multipass"), []nitro.Action{*nginxLogsAction})
	},
}

func init() {
	logsNginxCommand.Flags().StringVarP(&flagNginxLogsKind, "type", "t", "all", "Filter the logs by kind, access or error (defaults to all)")
}
