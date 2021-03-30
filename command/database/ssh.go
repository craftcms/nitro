package database

import (
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/craftcms/nitro/pkg/containerlabels"
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
		Short:   "Opens a shell in a database container.",
		Example: sshExampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// add filters to show only the environment and database containers
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)
			filter.Add("label", containerlabels.Type+"=database")

			// get a list of all the databases
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter})
			if err != nil {
				return err
			}

			// sort containers by the name
			sort.SliceStable(containers, func(i, j int) bool {
				return containers[i].Names[0] < containers[j].Names[0]
			})

			// generate a list of engines for the prompt
			var containerList []string
			for _, c := range containers {
				// start the container if not running
				if c.State != "running" {
					for _, command := range cmd.Root().Commands() {
						if command.Use == "start" {
							if err := command.RunE(cmd, []string{}); err != nil {
								return err
							}
						}
					}
				}

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
