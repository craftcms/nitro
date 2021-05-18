package clean

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # remove unused containers
  nitro clean`

// NewCommand returns the command that is used to clean containers that do not exist in a specified
// environment. It will also perform the backup for database containers.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:        "clean",
		Deprecated: "Use the `apply` command instead",
		Short:      "Removes unused containers.",
		Example:    exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			output.Info("Cleaning up…")

			output.Pending("gathering details")

			// get all of the containers for the environment
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro+"=true")
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{All: true, Filters: filter})
			if err != nil {
				return err
			}

			// check if each container exists
			toRemove := []types.Container{}
			for _, c := range containers {
				// we should remove the container if it is a composer or npm container
				if c.Labels[containerlabels.Type] == "composer" || c.Labels[containerlabels.Type] == "npm" {
					toRemove = append(toRemove, c)
				}
			}

			output.Done()

			// if there is nothing to remove don't remove it
			if len(toRemove) == 0 {
				output.Info("Nothing to remove 😅")

				return nil
			}

			// remove each of the containers
			for _, c := range toRemove {
				// stop the container
				if err := docker.ContainerStop(cmd.Context(), c.ID, nil); err != nil {
					output.Warning()
					output.Info(err.Error())
					break
				}

				// remove the container
				if err := docker.ContainerRemove(cmd.Context(), c.ID, types.ContainerRemoveOptions{RemoveVolumes: true}); err != nil {
					output.Warning()
					output.Info(err.Error())
					break
				}

				output.Done()
			}

			output.Info("Cleanup completed 🛁")

			return nil
		},
	}

	return cmd
}
