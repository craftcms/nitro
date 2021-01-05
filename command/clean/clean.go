package clean

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/backup"
	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/datetime"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # remove unused containers
  nitro clean`

// NewCommand returns the command that is used to clean containers that do not exist in a specified
// environment. It will also perform the backup for database containers.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "clean",
		Short:   "Remove unused containers",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// load the config file
			cfg, err := config.Load(home)
			if err != nil {
				return fmt.Errorf("unable to load config, %w", err)
			}

			output.Info("Cleaning up…")

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
			filter.Add("label", labels.Nitro+"=true")
			filter.Add("label", labels.Type+"=composer")
			filter.Add("label", labels.Type+"=node")
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
				if c.Labels[labels.Proxy] != "true" {
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
				// if this is a database container we need to back it up
				if c.Labels[labels.DatabaseEngine] != "" {
					// get all of the databases
					databases, err := backup.Databases(cmd.Context(), docker, c.ID, c.Labels[labels.DatabaseCompatability])
					if err != nil {
						output.Info("Unable to get the databases from", name, err.Error())

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

						// create the backup command based on the compatability type
						switch c.Labels[labels.DatabaseCompatability] {
						case "postgres":
							opts.Commands = []string{"pg_dump", "--username=nitro", db, "-f", "/tmp/" + opts.BackupName}
						default:
							opts.Commands = []string{"/usr/bin/mysqldump", "-h", "127.0.0.1", "-unitro", "--password=nitro", db, "--result-file=" + "/tmp/" + opts.BackupName}
						}

						output.Pending("creating backup", opts.BackupName)

						// backup the container
						if err := backup.Perform(cmd.Context(), docker, opts); err != nil {
							output.Warning()
							output.Info("Unable to backup database", db, err.Error())

							break
						}

						output.Done()
					}

					// show where all backups are saved for this container
					output.Info("Backups saved in", filepath.Join(home, ".nitro", name), "💾")
				}

				output.Pending("removing", name)

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

				// remove the containers volume if its a database
				if c.Labels[labels.DatabaseEngine] != "" {
					// append the database labels to the filter
					filter.Add("label", labels.DatabaseEngine+"="+c.Labels[labels.DatabaseEngine])
					filter.Add("label", labels.DatabaseVersion+"="+c.Labels[labels.DatabaseVersion])
					filter.Add("label", labels.Type+"=database")

					// get the volume for the database
					resp, err := docker.VolumeList(cmd.Context(), filter)
					if err != nil {
						return err
					}

					// no volumes found
					if len(resp.Volumes) > 0 {
						// remove the found volume
						if err := docker.VolumeRemove(cmd.Context(), resp.Volumes[0].Name, false); err != nil {
							return err
						}
					}

					// remove the database labels from the filter
					filter.Del("label", labels.DatabaseEngine+"="+c.Labels[labels.DatabaseEngine])
					filter.Del("label", labels.DatabaseVersion+"="+c.Labels[labels.DatabaseVersion])
					filter.Del("label", labels.Type+"=database")
				}

				output.Done()
			}

			output.Info("Cleanup completed 🛁")

			return nil
		},
	}

	return cmd
}
