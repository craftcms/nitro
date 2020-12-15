package clean

import (
	"fmt"
	"strings"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/labels"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var (
	// ErrExample is used when we want to share an error
	ErrExample = fmt.Errorf("some example error")
)

const exampleText = `  # remove unused containers
  nitro clean`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "clean",
		Short:   "Remove unused containers",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			env := cmd.Flag("environment").Value.String()

			// load the config file
			cfg, err := config.Load(home, env)
			if err != nil {
				return fmt.Errorf("unable to load config, %w", err)
			}

			output.Info("Cleaning up...")

			output.Pending("gathering details")

			// get a list of the known containers
			known := make(map[string]bool)

			// get all current sites
			for _, s := range cfg.Sites {
				known[s.Hostname] = true
			}

			// get all current databases
			for _, d := range cfg.Databases {
				hostname, err := d.GetHostname()
				if err != nil {
					return err
				}

				known[hostname] = true
			}

			// get all of the containers for the environment
			filter := filters.NewArgs()
			filter.Add("label", labels.Environment+"="+env)
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{All: true, Filters: filter})
			if err != nil {
				return err
			}

			// check if each container exists
			remove := make(map[string]types.Container)
			for _, c := range containers {
				name := strings.TrimLeft(c.Names[0], "/")

				// check if the container is a known container
				if _, ok := known[name]; ok {
					continue
				}

				// if this is a proxy container
				if c.Labels[labels.Proxy] == env {
					continue
				}

				// we should remove the container
				remove[name] = c
			}

			output.Done()

			// if there is nothing to remove don't remove it
			if len(remove) == 0 {
				output.Info("Nothing to remove 😅")

				return nil
			}

			// remove each of the containers
			for name, c := range remove {
				output.Pending("removing", name)

				// TODO(jasonmccallister) if this is a database container

				// stop the container
				if err := docker.ContainerStop(cmd.Context(), c.ID, nil); err != nil {
					output.Warning()
					output.Info(err.Error())
					continue
				}

				// remove the container
				if err := docker.ContainerRemove(cmd.Context(), c.ID, types.ContainerRemoveOptions{RemoveVolumes: true}); err != nil {
					output.Warning()
					output.Info(err.Error())
					continue
				}

				output.Done()
			}

			output.Info("Cleanup completed 🧹")

			return nil
		},
	}

	return cmd
}
