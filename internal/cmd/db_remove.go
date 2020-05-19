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

var dbRemoveCommand = &cobra.Command{
	Use:   "remove",
	Short: "Remove databases",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		mp, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}
		p := prompt.NewPrompt()

		// get all of the docker containers by name
		script := scripts.New(mp, machine)
		output, err := script.Run(false, scripts.DockerListContainerNames)
		if err != nil {
			return err
		}

		// create a list
		containers := strings.Split(output, "\n")
		if len(containers) == 0 {
			return errors.New("there are no containers to perform actions on")
		}

		container, _, err := p.Select("Which database should we remove", containers, &prompt.SelectOptions{
			Default: 1,
		})
		if err != nil {
			return err
		}

		remove, err := p.Confirm("Are you sure you want to permanently remove the database "+container, &prompt.InputOptions{
			Default:   "no",
			Validator: nil,
		})
		if err != nil {
			return err
		}

		if remove {
			_, err = script.Run(false, fmt.Sprintf(scripts.FmtDockerRemoveContainer, container))
			if err != nil {
				return err
			}

			_, err = script.Run(false, fmt.Sprintf(scripts.FmtDockerRemoveVolume, container))
			if err != nil {
				return err
			}

			fmt.Println("Removed database", container)
			return nil
		}

		fmt.Println("We did not remove the database ", container)

		return nil
	},
}
