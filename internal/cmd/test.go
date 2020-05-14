package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pixelandtonic/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/scripts"
)

var testCommand = &cobra.Command{
	Use:   "test",
	Short: "Testing",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}

		mp, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}

		var databases []config.Database
		if err := viper.UnmarshalKey("databases", &databases); err != nil {
			return err
		}
		var dbs []string
		for _, db := range databases {
			dbs = append(dbs, db.Name())
		}

		if len(dbs) == 0 {
			return errors.New("there are no databases")
		}

		// PROMPT FOR INPUT
		p := prompt.NewPrompt()

		containerName, _, err := p.Select("Which database engine to import the backup", dbs, &prompt.SelectOptions{Default: 1})

		databaseName, err := p.Ask("What is the database name to create for the import", &prompt.InputOptions{Default: "", Validator: nil})
		if err != nil {
			return err
		}

		script := scripts.New(mp, machine)

		// check if the site it available
		if strings.Contains(containerName, "mysql") {
			_, err := script.Run(fmt.Sprintf(scripts.FmtDockerMysqlCreateDatabaseIfNotExists, containerName, databaseName))
			if err != nil {
				return err
			}
			fmt.Println("Created database", databaseName)

			_, err = script.Run(fmt.Sprintf(scripts.FmtDockerMysqlGrantPrivileges, containerName))
			if err != nil {
				return err
			}
			fmt.Println("Set permissions for the user nitro on", databaseName)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCommand)
}
