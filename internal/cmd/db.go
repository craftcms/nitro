package cmd

import (
	"github.com/spf13/cobra"
)

var dbCommand = &cobra.Command{
	Use:       "db",
	Short:     "Manage databases",
	ValidArgs: []string{"add", "backup", "import", "restart", "stop", "start"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	dbCommand.AddCommand(dbAddCommand, dbImportCommand, dbRestartCommand, dbStopCommand, dbStartCommand, dbRemoveCommand, dbBackupCommand)
}
