package destroy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/backup"
	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/datetime"
	"github.com/craftcms/nitro/pkg/sudo"
	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	// ErrNoNetworks is returned when no networks are running for an environment
	ErrNoNetworks = fmt.Errorf("there are no networks")

	// ErrNoContainers is returned when no containers are running for an environment
	ErrNoContainers = fmt.Errorf("there are no running containers")

	// ErrNoVolumes is returned when no volumes are running for an environment
	ErrNoVolumes = fmt.Errorf("there are no volumes")
)

const exampleText = `  # remove all resources (networks, containers, and volumes)
  nitro destroy`

// NewCommand is used to destroy all resources for an environment. It will prompt for
// user verification and defaults to no. Part of the destroy process is to
// perform a backup for all databases in each container database.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "destroy",
		Short:   "Destroys Nitro’s Docker resources.",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// prompt the user for confirmation
			confirm, err := output.Confirm("Are you sure? (This will remove all containers, volumes, and networks.)", false, "")
			if err != nil {
				return err
			}

			if !confirm {
				output.Info("skipping destroy, all resources will remain 😅")

				return nil
			}

			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

			// get all related containers
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{
				All:     true,
				Filters: filter,
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

				output.Info("Removing Containers…")

				for _, c := range containers {
					name := strings.TrimLeft(c.Names[0], "/")

					// only perform a backup if the container is for databases
					if c.Labels[containerlabels.DatabaseEngine] != "" {
						// this container needs to be running before we can backup the system
						if c.State != "running" {
							if err := docker.ContainerStart(ctx, c.ID, types.ContainerStartOptions{}); err != nil {
								output.Info("unable to start the container to begin backups", name)
								break
							}
						}

						// get all of the databases
						databases, err := backup.Databases(ctx, docker, c.ID, c.Labels[containerlabels.DatabaseCompatibility])
						if err != nil {
							output.Info("unable to get the databases from", name, err.Error())

							break
						}

						// backup each database
						for _, db := range databases {
							// create the database specific backup options
							opts := &backup.Options{
								BackupName:    fmt.Sprintf("%s-%s.sql", db, datetime.Parse(time.Now())),
								ContainerID:   c.ID,
								ContainerName: name,
								Database:      db,
								Home:          home,
							}

							// create the backup command based on the compatibility type
							switch c.Labels[containerlabels.DatabaseCompatibility] {
							case "postgres":
								opts.Commands = []string{"pg_dump", "--username=nitro", db, "-f", "/tmp/" + opts.BackupName}
							default:
								opts.Commands = []string{"/usr/bin/mysqldump", "-h", "127.0.0.1", "-unitro", "--password=nitro", db, "--result-file=" + "/tmp/" + opts.BackupName}
							}

							output.Pending("creating backup", opts.BackupName)

							// backup the container
							if err := backup.Perform(ctx, docker, opts); err != nil {
								output.Warning()
								output.Info("Unable to backup database", db, err.Error())

								break
							}

							output.Done()
						}

						// show where all backups are saved for this container
						output.Info("Backups saved in", filepath.Join(home, config.DirectoryName, name), "💾")
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
				output.Info("Removing volumes…")

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
				output.Info("Removing Networks…")

				for _, n := range networks {
					output.Pending("removing", n.Name)

					if err := docker.NetworkRemove(ctx, n.ID); err != nil {
						output.Info("unable to remove network", n.Name, "you may need to manually remove network")
					}

					output.Done()
				}
			}

			// remove the config file when --clean is true
			if cmd.Flag("clean").Value.String() == "true" {
				if err := os.Remove(cfg.GetFile()); err != nil {
					output.Info("Unable to remove configuration file")
				}
			}

			// remove nitro hosts entries

			// get the executable
			nitro, err := os.Executable()
			if err != nil {
				return fmt.Errorf("unable to locate the nitro path, %w", err)
			}

			// run the hosts command
			switch runtime.GOOS {
			case "windows":
				// windows users should be running as admin, so just execute the hosts command as is
				c := exec.Command(nitro, "hosts", "remove")

				c.Stdout = os.Stdout
				c.Stderr = os.Stderr

				if c.Run() != nil {
					return err
				}
			default:
				output.Info("Updating hosts file (you might be prompted for your password)")

				// add the hosts
				if err := sudo.Run(nitro, "nitro", "hosts", "remove"); err != nil {
					return err
				}
			}

			output.Info("Nitro destroyed ✨")

			return nil
		},
	}

	// add flags to the command
	cmd.Flags().Bool("clean", false, "remove configuration file")

	return cmd
}
