package craft

import (
	"context"
	"fmt"
	"os"

	"github.com/craftcms/nitro/pkg/appaware"
	"github.com/craftcms/nitro/pkg/flags"
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

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

			var hostname string
			switch flags.AppName == "" {
			case true:
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

				hostname, err = appaware.Detect(*cfg, wd)
				if err != nil {
					return err
				}
			default:
				hostname = flags.AppName
			}

			output.Info("connecting to", hostname)

			// add the label to get the site
			filter.Add("label", containerlabels.Host+"="+hostname)

			// find the containers but limited to the app label
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter, All: true})
			if err != nil {
				return err
			}

			// are there any containers??
			if len(containers) == 0 {
				return fmt.Errorf("unable to find a matching app")
			}

			// start the container if it's not running
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

			// TODO(jasonmccallister) get the container path
			cmds = append(cmds, fmt.Sprintf("%s/%s", "/app", "craft"))

			switch len(args) == 0 {
			case true:
				// no args were provided, use the help command
				cmds = append(cmds, "help")
			default:
				// append the provided args to the command
				cmds = append(cmds, args...)
			}

			_, err = execCreate(cmd.Context(), docker, containers[0].ID, cmds, true)
			if err != nil {
				return err
			}


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
