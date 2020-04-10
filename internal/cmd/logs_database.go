package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var logsDatabaseCommand = &cobra.Command{
	Use:    "database",
	Short:  "Show database logs",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not implemented yet")
	},
}

func init() {
	logsDatabaseCommand.Flags().StringVar(&flagDatabase, "database", "", "which database engine")
	logsDatabaseCommand.Flags().StringVar(&flagDatabaseVersion, "database-version", "", "which version of the database")
}
