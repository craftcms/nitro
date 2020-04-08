package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var (
	flagNginxLogsKind string

	logsNginxCommand = &cobra.Command{
		Use:   "nginx",
		Short: "Show nginx logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := config.GetString("name", flagMachineName)

			return nitro.Run(nitro.NewMultipassRunner("multipass"), nitro.NginxLogs(name, flagNginxLogsKind))
		},
	}
)

func init() {
	logsNginxCommand.Flags().StringVarP(&flagNginxLogsKind, "type", "t", "all", "filter the logs by kind, access or error (defaults to all)")
}
