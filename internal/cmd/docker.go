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

var dockerCommand = &cobra.Command{
	Use:   "docker",
	Short: "Perform Docker actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		mp, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}
		p := prompt.NewPrompt()

		action, _, err := p.Select("What action would you like to perform", []string{"restart", "stop", "start", "remove"}, &prompt.SelectOptions{
			Default: 1,
		})
		if err != nil {
			return err
		}

		// get all of the docker containers by name
		script := scripts.New(mp, machine)
		output, err := script.Run(scripts.DockerListContainerNames)
		if err != nil {
			return err
		}

		// create a list
		containers := strings.Split(output, "\n")
		if len(containers) == 0 {
			return errors.New("there are no containers to perform actions on")
		}

		switch action {
		case "stop":
			container, _, err := p.Select("Which container should we stop", containers, &prompt.SelectOptions{
				Default: 1,
			})
			if err != nil {
				return err
			}

			_, err = script.Run(fmt.Sprintf(scripts.FmtDockerStopContainer, container))
			if err != nil {
				return err
			}

			fmt.Println("Stopped container", container)
		case "start":
			container, _, err := p.Select("Which container should we start", containers, &prompt.SelectOptions{
				Default: 1,
			})
			if err != nil {
				return err
			}

			_, err = script.Run(fmt.Sprintf(scripts.FmtDockerStartContainer, container))
			if err != nil {
				return err
			}

			fmt.Println("Started container", container)
		case "remove":
			container, _, err := p.Select("Which container should we remove", containers, &prompt.SelectOptions{
				Default: 1,
			})
			if err != nil {
				return err
			}

			remove, err := p.Confirm("Are you sure you want to permanently remove the container "+container, &prompt.InputOptions{
				Default:   "no",
				Validator: nil,
			})
			if err != nil {
				return err
			}

			if remove {
				_, err = script.Run(fmt.Sprintf(scripts.FmtDockerRemoveContainer, container))
				if err != nil {
					return err
				}

				fmt.Println("Removed container", container)
				return nil
			}

			fmt.Print("We did nto remove the container ", container)
		default:
			container, _, err := p.Select("Which container should we restart", containers, &prompt.SelectOptions{
				Default: 1,
			})
			if err != nil {
				return err
			}

			_, err = script.Run(fmt.Sprintf(scripts.FmtDockerRestartContainer, container))
			if err != nil {
				return err
			}

			fmt.Println("Restarted container", container)
		}

		return nil
	},
}
