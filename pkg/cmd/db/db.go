package db

import (
	"github.com/spf13/cobra"
)

// DBCommand is the command for creating new development environments
var DBCommand = &cobra.Command{
	Use:   "db",
	Short: "Perform database actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
	Example: `  # list all containers for the environment
  nitro db import`,
}

func init() {
	DBCommand.AddCommand(ImportCommand)
}
