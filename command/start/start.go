package start

import (
	"fmt"
	"strings"

	"github.com/craftcms/nitro/labels"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var (
	// ErrNoContainers is returned when no containers are running for an environment
	ErrNoContainers = fmt.Errorf("there are no running containers")
)

const exampleText = `  # start containers for the default environment
  nitro start`

// New is used for scaffolding new commands
func New(docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Short:   "Start environment",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			env := cmd.Flag("environment").Value.String()
			ctx := cmd.Context()

			// get all the containers using a filter, we only want to stop containers which
			// have the environment label
			filter := filters.NewArgs()
			filter.Add("label", labels.Environment+"="+env)

			// get all of the container
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
			if err != nil {
				return fmt.Errorf("unable to get a list of the containers, %w", err)
			}

			// if there are no containers, were done
			if len(containers) == 0 {
				return ErrNoContainers
			}

			output.Info(fmt.Sprintf("Starting %s...", env))

			// start each environment container
			for _, c := range containers {
				n := strings.TrimLeft(c.Names[0], "/")

				// if the container is already running
				if c.State == "running" {
					output.Success(n, "ready")
					continue
				}

				output.Pending("starting", n)

				// start the container
				if err := docker.ContainerStart(ctx, c.ID, types.ContainerStartOptions{}); err != nil {
					return fmt.Errorf("unable to start container %s: %w", n, err)
				}

				output.Done()
			}

			output.Info(env, "started 👍")

			return nil
		},
	}

	return cmd
}
