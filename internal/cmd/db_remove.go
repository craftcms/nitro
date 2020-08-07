package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pixelandtonic/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/internal/config"
	"github.com/craftcms/nitro/internal/scripts"
)

var dbRemoveCommand = &cobra.Command{
	Use:   "remove",
	Short: "Remove database engine",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		mp, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}
		p := prompt.NewPrompt()

		// get all of the docker containers by name
		script := scripts.New(mp, machine)

		var cfg config.Config
		if err := viper.Unmarshal(&cfg); err != nil {
			return err
		}

		if len(cfg.Databases) == 0 {
			return errors.New("there are no databases we can add to")
		}

		// get all of the docker containers by name
		var containers []string
		for _, db := range cfg.Databases {
			containers = append(containers, db.Name())
		}

		// if there is only one
		var container string
		switch len(containers) {
		case 1:
			container = containers[0]
		default:
			container, _, err = p.Select("Select the database engine", containers, &prompt.SelectOptions{
				Default: 1,
			})
			if err != nil {
				return err
			}
		}

		// get all of the databases in the engine
		var dbs []string
		switch strings.Contains(container, "mysql") {
		case false:
			if output, err := script.Run(false, fmt.Sprintf(scripts.FmtDockerPostgresShowAllDatabases, container)); err == nil {
				sp := strings.Split(output, "\n")
				for i, d := range sp {
					if i == 0 || i == 1 || i == len(sp) || strings.Contains(d, "rows)") {
						continue
					}

					dbs = append(dbs, strings.TrimSpace(d))
				}
			}
		default:
			if output, err := script.Run(false, fmt.Sprintf(scripts.FmtDockerMysqlShowAllDatabases, container)); err == nil {
				for _, db := range strings.Split(output, "\n") {
					// ignore the system defaults
					if db == "Database" || db == "information_schema" || db == "performance_schema" || db == "sys" || strings.Contains(db, "password on the command line") || db == "mysql" {
						continue
					}
					dbs = append(dbs, db)
				}
			}
		}

		// ask the user which database to remove
		database, _, err := p.Select("Select database to remove", dbs, &prompt.SelectOptions{Default: 1})
		if err != nil {
			return err
		}

		// make sure the user wants to do this
		remove, err := p.Confirm(fmt.Sprintf("Are you sure you want to permanently remove the database %q", database), &prompt.InputOptions{
			Default:            "no",
			Validator:          nil,
			AppendQuestionMark: true,
		})
		if err != nil {
			return err
		}

		if remove {
			switch strings.Contains(container, "mysql") {
			case false:
				// its postgres so remove the db
				if output, err := script.Run(false, fmt.Sprintf(`docker exec -i %s psql --username nitro -c "DROP DATABASE IF EXISTS %s;"`, container, database)); err != nil {
					fmt.Println(output)
					return err
				}
			default:
				// its mysql, drop the db
				if output, err := script.Run(false, fmt.Sprintf(`docker exec -i %s mysql -unitro -pnitro -e "DROP DATABASE IF EXISTS %s;"`, container, database)); err != nil {
					fmt.Println(output)
					return err
				}
			}

			fmt.Println("Removed database", database)
			return nil
		}

		fmt.Println("There was a problem removing the database", database)

		return nil
	},
}
