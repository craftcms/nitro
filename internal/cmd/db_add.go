package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pixelandtonic/prompt"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/scripts"
)

var dbAddCommand = &cobra.Command{
	Use:   "add",
	Short: "Add a new databases",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		mp, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}
		p := prompt.NewPrompt()

		// get all of the docker containers by name
		script := scripts.New(mp, machine)
		output, err := script.Run(scripts.DockerListContainerNames)
		if err != nil {
			return err
		}

		// create a list
		containers := strings.Split(output, "\n")
		if len(containers) == 0 {
			return errors.New("there are no databases we can add to")
		}

		// which database
		container, _, err := p.Select("Which database should we restart", containers, &prompt.SelectOptions{
			Default: 1,
		})
		if err != nil {
			return err
		}

		// get the name
		database, err := p.Ask("What is the name of the database to add", &prompt.InputOptions{Default: "", Validator: nil})
		if err != nil {
			return err
		}

		// run the scripts
		if strings.Contains(container, "mysql") {
			_, err = script.Run(fmt.Sprintf(scripts.FmtDockerMysqlCreateDatabaseIfNotExists, container, database))
			if err != nil {
				return err
			}
		} else {
			_, err = script.Run(fmt.Sprintf(scripts.FmtDockerPostgresCreateDatabase, container, database))
			if err != nil {
				return err
			}
		}

		fmt.Println(fmt.Sprintf("Added database %q to %q", database, container))

		return nil
	},
}
