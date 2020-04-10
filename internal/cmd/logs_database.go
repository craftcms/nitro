package cmd

import (
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/action"
)

var logsDatabaseCommand = &cobra.Command{
	Use:    "database",
	Short:  "Show database logs",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)
		if !viper.IsSet("databases") {
			return errors.New("no databases found in " + viper.ConfigFileUsed())
		}

		var databases []config.Database
		if err := viper.UnmarshalKey("databases", &databases); err != nil {
			return err
		}

		// TODO abstract this
		var dbs []string
		for _, db := range databases {
			dbs = append(dbs, fmt.Sprintf("%s_%s_%s", db.Engine, db.Version, db.Port))
		}

		prompt := promptui.Select{
			Label: "Select database",
			Items: dbs,
		}

		_, container, err := prompt.Run()
		if err != nil {
			return err
		}

		dockerLogsAction, err := action.LogsDocker(name, container)
		if err != nil {
			return err
		}

		return action.Run(action.NewMultipassRunner("multipass"), []action.Action{*dockerLogsAction})
	},
}

func init() {
	logsDatabaseCommand.Flags().StringVar(&flagDatabase, "database", "", "which database engine")
	logsDatabaseCommand.Flags().StringVar(&flagDatabaseVersion, "database-version", "", "which version of the database")
}
