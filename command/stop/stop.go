package stop

import (
	"fmt"
	"os"
	"strings"

	"github.com/craftcms/nitro/pkg/appaware"
	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/flags"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # stop all containers
  nitro stop

  # stop an individual app
  nitro stop --app tutorial.nitro`

// NewCommand is used to stop all running containers for an environment. The process
// of stopping to reduce usage and "finish" your work effort at the end of your session.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stop",
		Short:   "Stops containers.",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// get all the containers using a filter, we only want to stop containers which
			// have the environment label
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

			var appName string
			if flags.AppName != "" {
				// add the label to get the app
				filter.Add("label", containerlabels.Host+"="+flags.AppName)

				appName = flags.AppName
			} else {
				wd, err := os.Getwd()
				if err != nil {
					return err
				}

				cfg, err := config.Load(home)
				if err != nil {
					return err
				}

				// don't return an error because we should start all
				app, err := appaware.Detect(*cfg, wd)
				if err != nil {
					return err
				}

				appName = app
			}

			// get all containers
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter})
			if err != nil {
				return fmt.Errorf("unable to get a list of the containers, %w", err)
			}

			// if there are no containers, were done
			if len(containers) == 0 {
				output.Info("there are no running containers")
				return nil
			}

			output.Info("Stopping Nitro…")

			// stop each environment container
			for _, c := range containers {
				hostname := strings.TrimLeft(c.Names[0], "/")

				// if the user wants a single app only, skip other apps
				if appName != "" && hostname != appName {
					continue
				}

				output.Pending("stopping", hostname)

				// stop the container
				if err := docker.ContainerStop(ctx, c.ID, nil); err != nil {
					return fmt.Errorf("unable to stop container %s: %w", hostname, err)
				}

				output.Done()
			}

			output.Info("Nitro shutdown 😴")

			return nil
		},
	}

	return cmd
}
