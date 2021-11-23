package restart

import (
	"fmt"
	"strings"
	"time"

	"github.com/craftcms/nitro/pkg/flags"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	// ErrNoContainers is returned when no containers are running for an environment
	ErrNoContainers = fmt.Errorf("there are no running containers")
)

const exampleText = `  # restart all containers
  nitro restart

  # restart specific app
  nitro restart --app tutorial.nitro`

// New returns the command to restart all of an environments containers
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "restart",
		Short:   "Restarts all containers.",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// get all the containers using a filter, we only want to restart containers which
			// have the label com.craftcms.nitro.environment=name
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

			if flags.AppName != "" {
				// add the label to get the app
				filter.Add("label", containerlabels.Host+"="+flags.AppName)
			}

			// get the containers
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter})
			if err != nil {
				return fmt.Errorf("unable to get a list of the containers, %w", err)
			}

			// if there are no containers, were done
			if len(containers) == 0 {
				return ErrNoContainers
			}

			output.Info("Restarting Nitro…")

			// set a timeout, consider making this a flag
			timeout := time.Duration(5000) * time.Millisecond

			// restart each container for the environment
			for _, c := range containers {
				n := strings.TrimLeft(c.Names[0], "/")

				// make sure the app is not disabled
				if c.Labels[containerlabels.Disabled] == "true" {
					output.Info(fmt.Sprintf("%s is disabled, skipping", n))
					continue
				}

				output.Pending("restarting", n)

				// restart the container
				if err := docker.ContainerRestart(ctx, c.ID, &timeout); err != nil {
					return fmt.Errorf("unable to restart container %s: %w", n, err)
				}

				output.Done()
			}

			fmt.Println("Nitro restarted 🎉")

			return nil
		},
	}

	return cmd
}
