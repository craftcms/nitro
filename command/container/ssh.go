package container

import (
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var sshExampleText = `  # ssh into a custom container
  nitro container ssh`

func sshCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ssh",
		Short:   "SSH into custom container",
		Example: sshExampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// add filters to show only the environment and database containers
			filter := filters.NewArgs()
			filter.Add("label", labels.Nitro)
			filter.Add("label", labels.Type+"=custom")

			// get a list of all the containers
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
					if err := docker.ContainerStart(cmd.Context(), c.ID, types.ContainerStartOptions{}); err != nil {
						return err
					}
				}

				containerList = append(containerList, strings.TrimLeft(c.Names[0], "/"))
			}

			// prompt for the container to ssh into
			selected, err := output.Select(cmd.InOrStdin(), "Select a container to connect to: ", containerList)
			if err != nil {
				return err
			}

			container := containerList[selected]

			return containerConnect(container, output)
		},
	}

	return cmd
}

func containerConnect(name string, output terminal.Outputer) error {
	// find the docker executable
	cli, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	c := exec.Command(cli, "exec", "-u", "root", "-it", name, "bash")

	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout

	return c.Run()
}
