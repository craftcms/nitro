package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var (
	restartDatabaseCommand = &cobra.Command{
		Use:   "database",
		Short: "Restart database",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := config.GetString("machine", flagMachineName)
			engine := config.GetString("database.engine", flagDatabase)
			version := config.GetString("database.version", flagDatabaseVersion)

			return nitro.Run(nitro.NewMultipassRunner("multipass"), nitro.RestartDatabase(name, engine, version))
		},
	}
	servicesDatabaseRestartCommand = &cobra.Command{
		Use:   "restart",
		Short: "Restart the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
)
