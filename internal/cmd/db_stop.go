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

var dbStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stop databases",
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

		container, _, err := p.Select("Which database should we stop", containers, &prompt.SelectOptions{
			Default: 1,
		})
		if err != nil {
			return err
		}

		_, err = script.Run(false, fmt.Sprintf(scripts.FmtDockerStopContainer, container))
		if err != nil {
			return err
		}

		fmt.Println("Stopped database", container)

		return nil
	},
}
