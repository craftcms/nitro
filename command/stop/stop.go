package stop

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

const exampleText = `  # stop all containers
  nitro stop

  # stop an individual site
  nitro stop tutorial.nitro`

// New is used to stop all running containers for an environment. The process
// of stopping to reduce usage and "finish" your work effort at the end of your session.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stop",
		Short:   "Stops containers.",
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

			return options, cobra.ShellCompDirectiveNoFileComp
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

				// if the user wants a single site only, skip all of the other sites
				if site != "" && hostname != site {
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
