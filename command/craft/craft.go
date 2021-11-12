package craft

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # run craft console command in a sites container
  nitro craft update all

  # view craft console help command
  nitro craft

  # enter the craft shell
  nitro craft shell`

// NewCommand returns the craft command which allows users to pass craft specific commands to a sites
// container. Its context aware and will prompt the user for the site if its not in a directory.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "craft",
		Short:              "Runs a Craft console command.",
		DisableFlagParsing: true,
		Example:            exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			// load the config
			cfg, err := config.Load(home, false)
			if err != nil {
				return err
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

			// get a context aware list of sites
			sites := cfg.ListOfSitesByDirectory(home, wd)

			// create the options for the sites
			var options []string
			for _, s := range sites {
				options = append(options, s.Hostname)
			}

			var site config.Site
			switch len(sites) {
			case 0:
				// prompt for the site
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
				if err != nil {
					return err
				}

				// set the site we selected
				site = sites[selected]

				// add the label to get the site
				filter.Add("label", containerlabels.Host+"="+sites[selected].Hostname)
			case 1:
				output.Info("connecting to", sites[0].Hostname)

				// set the site we selected
				site = sites[0]

				// add the label to get the site
				filter.Add("label", containerlabels.Host+"="+sites[0].Hostname)
			default:
				// prompt for the site to ssh into
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
				if err != nil {
					return err
				}

				// set the site we selected
				site = sites[selected]

				// add the label to get the site
				filter.Add("label", containerlabels.Host+"="+sites[selected].Hostname)
			}

			// find the containers but limited to the site label
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter, All: true})
			if err != nil {
				return err
			}

			// are there any containers??
			if len(containers) == 0 {
				return fmt.Errorf("unable to find an matching site")
			}

			// start the container if its not running
			if containers[0].State != "running" {
				for _, command := range cmd.Root().Commands() {
					if command.Use == "start" {
						if err := command.RunE(cmd, []string{}); err != nil {
							return err
						}
					}
				}
			}

			// create the command for running the craft console
			cmds := []string{"exec", "-it", containers[0].ID, "php"}

			// get the container path
			path := site.GetContainerPath()
			if path != "" {
				cmds = append(cmds, fmt.Sprintf("%s/%s", path, "craft"))
			} else {
				cmds = append(cmds, "craft")
			}

			switch len(args) == 0 {
			case true:
				// no args were provided, use the help command
				cmds = append(cmds, "help")
			default:
				// append the provided args to the command
				cmds = append(cmds, args...)
			}

			// find the docker executable
			cli, err := exec.LookPath("docker")
			if err != nil {
				return err
			}

			// create the command
			c := exec.Command(cli, cmds...)

			c.Stdin = cmd.InOrStdin()
			c.Stderr = cmd.ErrOrStderr()
			c.Stdout = cmd.OutOrStdout()

			if err := c.Run(); err != nil {
				return err
			}

			// attach

			return nil
		},
	}

	return cmd
}

func execCreate(ctx context.Context, docker client.ContainerAPIClient, containerID string, cmds []string, show bool) (bool, error) {
	// create the exec
	e, err := docker.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Cmd:          cmds,
	})
	if err != nil {
		return false, err
	}

	// attach to the container
	resp, err := docker.ContainerExecAttach(ctx, e.ID, types.ExecStartCheck{
		Tty: false,
	})
	defer resp.Close()

	// should we display output?
	if show {
		// show the output to stdout and stderr
		if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, resp.Reader); err != nil {
			return false, fmt.Errorf("unable to copy the output of container, %w", err)
		}
	}

	// start the exec
	if err := docker.ContainerExecStart(ctx, e.ID, types.ExecStartCheck{}); err != nil {
		return false, fmt.Errorf("unable to start the container, %w", err)
	}

	// wait for the container exec to complete
	waiting := true
	for waiting {
		resp, err := docker.ContainerExecInspect(ctx, e.ID)
		if err != nil {
			return false, err
		}

		waiting = resp.Running
	}

	return true, nil
}
