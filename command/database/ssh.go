package database

import (
	"os"
	"os/exec"
	"strings"

	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var sshExampleText = `  # ssh into a database container
  nitro db ssh`

func sshCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ssh",
		Short:   "SSH into a db container",
		Example: sshExampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// add filters to show only the environment and database containers
			filter := filters.NewArgs()
			filter.Add("label", labels.Nitro)
			filter.Add("label", labels.Type+"=database")

			// get a list of all the databases
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter})
			if err != nil {
				return err
			}

			// generate a list of engines for the prompt
			var containerList []string
			for _, c := range containers {
				containerList = append(containerList, strings.TrimLeft(c.Names[0], "/"))
			}

			// prompt for the container to ssh into
			selected, err := output.Select(cmd.InOrStdin(), "Select a database to connect to: ", containerList)
			if err != nil {
				return err
			}

			container := containerList[selected]

			return containerConnect(output, container)
		},
	}

	return cmd
}

func containerConnect(output terminal.Outputer, containerName string) error {
	// find the docker executable
	cli, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	c := exec.Command(cli, "exec", "-u", "root", "-it", containerName, "bash")

	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout

	return c.Run()
}
