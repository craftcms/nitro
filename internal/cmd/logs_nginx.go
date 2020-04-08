package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
)

var logsNginxCommand = &cobra.Command{
	Use:   "nginx",
	Short: "Show nginx logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = config.GetString("name", flagMachineName)

		return errors.New("TODO")
	},
}

func init() {
	logsNginxCommand.Flags().StringVarP(&flagNginxLogsKind, "type", "t", "all", "filter the logs by kind, access or error (defaults to all)")
}
