package destroy

import (
	"fmt"
	"strings"
	"time"

	"github.com/craftcms/nitro/labels"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var (
	// ErrNoNetworks is returned when no networks are running for an environment
	ErrNoNetworks = fmt.Errorf("there are no networks")

	// ErrNoContainers is returned when no containers are running for an environment
	ErrNoContainers = fmt.Errorf("there are no running containers")

	// ErrNoVolumes is returned when no volumes are running for an environment
	ErrNoVolumes = fmt.Errorf("there are no volumes")
)

const exampleText = `  # remove all resources (networks, containers, and volumes) for an environment
  nitro destroy

  # destroy resources for a specific environment
  nitro destroy --environment my-testing-environment`

// New is used to destroy all resources for an environment. It will prompt for
// user verification and defaults to no. Part of the destroy process is to
// perform a backup for all databases in each container database.
func New(docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "destroy",
		Short:   "Destroy environment",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			env := cmd.Flag("environment").Value.String()
			ctx := cmd.Context()

			fmt.Print("Are you sure (this will remove all containers, volumes, and networks) [Y/n] ")

			// prompt the user for confirmation
			var response string
			_, err := fmt.Scanln(&response)
			if err != nil {
				return fmt.Errorf("unable to provide a prompt, %w", err)
			}

			var confirm bool
			resp := strings.TrimSpace(response)
			for _, answer := range []string{"y", "Y", "yes", "Yes", "YES"} {
				if resp == answer {
					confirm = true
				}
			}

			if !confirm {
				output.Info("skipping destroy, all resources related to", env, "will remain 😅")

				return nil
			}

			filter := filters.NewArgs()
			filter.Add("label", labels.Environment+"="+env)

			// get all related containers
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{
				All: true,
			})
			if err != nil {
				return fmt.Errorf("unable to list the containers, %w", err)
			}

			// make sure there are containers
			if len(containers) == 0 {
				output.Info(ErrNoContainers.Error())
			}

			// get all related volumes
			volumes, err := docker.VolumeList(ctx, filter)
			if err != nil {
				return err
			}

			// make sure there are volumes
			if len(volumes.Volumes) == 0 {
				output.Info(ErrNoVolumes.Error())
			}

			// get all related networks
			networks, err := docker.NetworkList(ctx, types.NetworkListOptions{Filters: filter})
			if err != nil {
				return err
			}

			// make sure there are networks
			if len(networks) == 0 {
				output.Info(ErrNoNetworks.Error())
			}

			// stop all of the container
			if len(containers) > 0 {
				timeout := time.Duration(5000) * time.Millisecond

				output.Info("Removing Containers...")

				for _, c := range containers {
					name := strings.TrimLeft(c.Names[0], "/")

					// only perform a backup if the container is for databases
					if c.Labels[labels.DatabaseEngine] != "" {
						output.Info("backing up databases is not yet supported")

						// TODO(jasonmccallister) implement backups of the databases
						// fmt.Println("Backing up database")
						// time.Sleep(time.Second * 2)
					}

					// stop the container
					output.Pending("removing", name)

					// stop the container
					if err := docker.ContainerStop(ctx, c.ID, &timeout); err != nil {
						return fmt.Errorf("unable to stop the container, %w", err)
					}

					// remove the container
					if err := docker.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
						return fmt.Errorf("unable to remove the container, %w", err)
					}

					output.Done()
				}
			}

			// get all the volumes
			if len(volumes.Volumes) > 0 {
				output.Info("Removing Volumes...")

				for _, v := range volumes.Volumes {
					output.Pending("removing", v.Name)

					// remove the volume
					if err := docker.VolumeRemove(ctx, v.Name, true); err != nil {
						output.Info("unable to remove volume", v.Name)
						break
					}

					output.Done()
				}
			}

			// get all the networks
			if len(networks) > 0 {
				output.Info("Removing Networks...")

				for _, n := range networks {
					output.Pending("removing", n.Name)

					if err := docker.NetworkRemove(ctx, n.ID); err != nil {
						output.Info("unable to remove network", n.Name, "you may need to manually remove network")
					}

					output.Done()
				}
			}

			output.Info(env, "destroyed ✨")

			return nil
		},
	}

	// set flags for the command
	// cmd.Flags().StringP("example", "e", "example", "an example flag")

	return cmd
}
