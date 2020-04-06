package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var sqlCommand = &cobra.Command{
	Use:   "sql",
	Short: "Enter a SQL shell",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("machine", flagMachineName)
		engine := config.GetString("database.engine", flagDatabase)
		version := config.GetString("database.version", flagDatabaseVersion)

		return nitro.Run(nitro.NewMultipassRunner("multipass"), nitro.SQL(name, engine, version, false))
	},
}
