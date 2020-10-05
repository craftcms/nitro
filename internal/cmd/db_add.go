package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pixelandtonic/prompt"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/config"
	"github.com/craftcms/nitro/internal/scripts"
	"github.com/craftcms/nitro/internal/slug"
	"github.com/craftcms/nitro/internal/validate"
)

var dbAddCommand = &cobra.Command{
	Use:   "add",
	Short: "Add new database",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		mp, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}
		script := scripts.New(mp, machine)
		cfg, err := config.Read()
		if err != nil {
			return err
		}

		if len(cfg.Databases) == 0 {
			return errors.New("there are no databases engines we can add to")
		}

		// get all of the docker containers by name
		var containers []string
		for _, db := range cfg.Databases {
			containers = append(containers, db.Name())
		}

		p := prompt.NewPrompt()

		// if there is only one
		var container string
		switch len(containers) {
		case 1:
			container = containers[0]
		default:
			container, _, err = p.Select("Select database engine", containers, &prompt.SelectOptions{
				Default: 1,
			})
			if err != nil {
				return err
			}
		}

		// get the name
		database, err := p.Ask("Enter the name of the database", &prompt.InputOptions{Default: "", Validator: validate.DatabaseName})
		if err != nil {
			return err
		}

		// clean the database name
		database = slug.Generate(database)

		// run the scripts
		if strings.Contains(container, "mysql") {
			_, err = script.Run(false, fmt.Sprintf(scripts.FmtDockerMysqlCreateDatabaseIfNotExists, container, database))
			if err != nil {
				return err
			}
		} else {
			_, err = script.Run(false, fmt.Sprintf(scripts.FmtDockerPostgresCreateDatabase, container, database))
			if err != nil {
				return err
			}
		}

		fmt.Println(fmt.Sprintf("Added database %q to %q.", database, container))

		return nil
	},
}
