package start

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
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
		Short:   "Start containers",
		Example: exampleText,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			cfg, err := config.Load(home)
			if err != nil {
				return nil, cobra.ShellCompDirectiveDefault
			}

			var options []string
			for _, s := range cfg.Sites {
				options = append(options, s.Hostname)
			}

			return options, cobra.ShellCompDirectiveDefault
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// is the docker api alive?
			if _, err := docker.Ping(cmd.Context()); err != nil {
				return fmt.Errorf("Couldn’t connect to Docker; please make sure Docker is running.")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			var site string
			if len(args) > 0 {
				site = args[0]
			}

			// get all the containers using a filter, we only want to stop containers which
			// have the environment label
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

			// get all of the container
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

				// if the user wants a single site only, skip all of the other sites
				if site != "" && hostname != site && containerType == "site" {
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
