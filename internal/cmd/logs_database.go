package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var (
	logsDatabaseCommand = &cobra.Command{
		Use:   "database",
		Short: "Show database logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := config.GetString("machine", flagMachineName)
			engine := config.GetString("database.engine", flagDatabase)
			version := config.GetString("database.version", flagDatabaseVersion)

			return nitro.Run(nitro.NewMultipassRunner("multipass"), nitro.DatabaseLogs(name, engine, version))
		},
	}
)

func init() {
	logsDatabaseCommand.Flags().StringVar(&flagDatabase, "database", "", "which database engine")
	logsDatabaseCommand.Flags().StringVar(&flagDatabaseVersion, "database-version", "", "which version of the database")
}
