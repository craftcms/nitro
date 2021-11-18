package start

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
	"github.com/craftcms/nitro/pkg/contextor"
	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	// ErrNoContainers is returned when no containers are running for an environment
	ErrNoContainers = fmt.Errorf("there are no running containers")
)

const exampleText = `  # start all containers
  nitro start`

// NewCommand returns the command used to start all of the containers for an environment.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Short:   "Starts containers.",
		Example: exampleText,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// is the docker api alive?
			if _, err := docker.Ping(cmd.Context()); err != nil {
				return fmt.Errorf("couldn’t connect to Docker; please make sure Docker is running")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := contextor.New(cmd.Context())

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
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
			if err != nil {
				return fmt.Errorf("unable to get a list of the containers, %w", err)
			}

			// if there are no containers, were done
			if len(containers) == 0 {
				return ErrNoContainers
			}

			output.Info("Starting Nitro…")

			// start each environment container
			for _, c := range containers {
				// don't start composer or npm containers
				if c.Labels[containerlabels.Type] == "composer" || c.Labels[containerlabels.Type] == "npm" {
					continue
				}

				// identify the type of container
				containerType := containerlabels.Identify(c)

				hostname := strings.TrimLeft(c.Names[0], "/")

				// if the user wants a single app only, skip other apps
				if appName != "" && hostname != appName && containerType == "app" {
					continue
				}

				// if the container is already running
				if c.State == "running" {
					output.Success(hostname)
					continue
				}

				output.Pending("starting", hostname)

				// start the container
				if err := docker.ContainerStart(ctx, c.ID, types.ContainerStartOptions{}); err != nil {
					return fmt.Errorf("unable to start container %s: %w", hostname, err)
				}

				output.Done()
			}

			output.Info("Nitro started 👍")

			return nil
		},
	}

	return cmd
}
