package mount

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

const sshExampleText = `  # ssh into a mounts container - assuming its the current working directory
  nitro mount ssh`

func sshCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ssh",
		Short:   "SSH into a mount",
		Example: sshExampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", labels.Nitro)

			// get all of the mounts
			var mounts, found []string
			for _, m := range cfg.Mounts {
				p, _ := m.GetAbsPath(home)

				// check if the path matches a sites path, then we are in a known site
				if strings.Contains(wd, p) {
					found = append(found, m.Hostname())
				}

				// add the site to the list in case we cannot find the directory
				mounts = append(mounts, m.Hostname())
			}

			// if there are found sites we want to show or connect to the first one, otherwise prompt for
			// which site to connect to.
			switch len(found) {
			case 0:
				// prompt for the site to ssh into
				selected, err := output.Select(cmd.InOrStdin(), "Select a mount to connect: ", mounts)
				if err != nil {
					return err
				}

				// add the name to get the mount
				filter.Add("name", mounts[selected])
			case 1:
				// add the name to get the mount
				filter.Add("name", found[0])
			default:
				// prompt for the mount to ssh into
				selected, err := output.Select(cmd.InOrStdin(), "Select a mount to connect: ", found)
				if err != nil {
					return err
				}

				// add the name to get the mount
				filter.Add("name", found[selected])
			}

			// find the containers but limited to the mount name
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter})
			if err != nil {
				return err
			}

			// are there any containers??
			if len(containers) == 0 {
				return fmt.Errorf("unable to find an matching mount")
			}

			// start the container if its not running
			if containers[0].State != "running" {
				if err := docker.ContainerStart(cmd.Context(), containers[0].ID, types.ContainerStartOptions{}); err != nil {
					return err
				}
			}

			return containerConnect(cmd, output, containers[0].ID)
		},
	}

	cmd.Flags().Bool("root", false, "ssh as the root user")

	return cmd
}

func containerConnect(cmd *cobra.Command, output terminal.Outputer, containerID string) error {
	// find the docker executable
	cli, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	// check if the root user should be used
	user := "www-data"
	if cmd.Flag("root").Value.String() == "true" {
		user = "root"
	}

	// show a notice about changes
	if user == "root" {
		output.Info("using root… system changes are ephemeral…")
	}

	c := exec.Command(cli, "exec", "-u", user, "-it", containerID, "sh")

	c.Stdin = cmd.InOrStdin()
	c.Stderr = cmd.ErrOrStderr()
	c.Stdout = cmd.OutOrStdout()

	return c.Run()
}
