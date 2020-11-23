package stop

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

const exampleText = `  # stop containers for the default environment
  nitro stop`

// New is used for scaffolding new commands
func New(docker client.CommonAPIClient, terminal terminal.Terminal) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stop",
		Short:   "Stop environment",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			env := cmd.Flag("environment").Value.String()
			ctx := cmd.Context()

			// get all the containers using a filter, we only want to stop containers which
			// have the environment label
			filter := filters.NewArgs()
			filter.Add("label", labels.Environment+"="+env)

			// get all of the container
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter})
			if err != nil {
				return fmt.Errorf("unable to get a list of the containers, %w", err)
			}

			// if there are no containers, were done
			if len(containers) == 0 {
				return ErrNoContainers
			}

			// stop each environment container
			for _, c := range containers {
				n := strings.TrimLeft(c.Names[0], "/")

				terminal.Pending("stopping", n)

				// stop the container
				if err := docker.ContainerStop(ctx, c.ID, nil); err != nil {
					return fmt.Errorf("unable to stop container %s: %w", n, err)
				}

				terminal.Done()
			}

			terminal.Info(env, "shutdown 😴")

			return nil
		},
	}

	return cmd
}
