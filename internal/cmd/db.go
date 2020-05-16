package cmd

import (
	"github.com/spf13/cobra"
)

var dbCommand = &cobra.Command{
	Use:   "db",
	Short: "Manage databases",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	dbCommand.AddCommand(dbAddCommand, dbImportCommand, dbRestartCommand, dbStopCommand, dbStartCommand, dbRemoveCommand)
}
