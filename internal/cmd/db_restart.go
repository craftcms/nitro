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

var dbRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "Restart databases",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		mp, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}
		p := prompt.NewPrompt()

		// get all of the docker containers by name
		script := scripts.New(mp, machine)

		// make a list
		output, err := script.Run(false, scripts.DockerListContainerNames)
		if err != nil {
			return err
		}

		// check it twice
		containers := strings.Split(output, "\n")
		// find out if nice
		if len(containers) == 0 {
			return errors.New("there are no containers to perform actions on")
		}

		container, _, err := p.Select("Which database should we restart", containers, &prompt.SelectOptions{
			Default: 1,
		})
		if err != nil {
			return err
		}

		_, err = script.Run(false, fmt.Sprintf(scripts.FmtDockerRestartContainer, container))
		if err != nil {
			return err
		}

		fmt.Println("Restarted database", container)

		return nil
	},
}
