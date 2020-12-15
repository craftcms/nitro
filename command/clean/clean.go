package clean

import (
	"fmt"
	"strings"
	"time"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/datetime"
	"github.com/craftcms/nitro/labels"
	"github.com/craftcms/nitro/pkg/backup"
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

// NewCommand returns the command that is used to clean containers that do not exist in a specified
// environment. It will also perform the backup for database containers.
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
				// if this is a database container we need to back it up
				if c.Labels[labels.DatabaseEngine] != "" {
					output.Pending("backup up container", name)

					// get all of the databases
					databases, err := backup.Databases(cmd.Context(), docker, c.ID, c.Labels[labels.DatabaseCompatability])
					if err != nil {
						output.Warning()
						output.Info("Unable to get the databases from", name, err.Error())

						break
					}

					// backup each database
					for _, d := range databases {
						// create the database specific backup options
						opts := &backup.Options{
							BackupName:    fmt.Sprintf("%s-%s.sql", d, datetime.Parse(time.Now())),
							ContainerID:   c.ID,
							ContainerName: name,
							Database:      d,
							Home:          home,
						}

						// create the backup command based on the compatability type
						switch c.Labels[labels.DatabaseCompatability] {
						case "postgres":
							opts.Commands = []string{"pg_dump", "--username=nitro", d, "-f", "/tmp/" + opts.BackupName}
						default:
							opts.Commands = []string{"/usr/bin/mysqldump", "-h", "127.0.0.1", "-unitro", "--password=nitro", d, "--result-file=" + "/tmp/" + opts.BackupName}
						}

						// backup the container
						if err := backup.Perform(cmd.Context(), docker, opts); err != nil {
							output.Warning()
							output.Info("Unable to backup database", d, err.Error())

							break
						}
					}

					output.Done()
				}

				output.Pending("removing", name)

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
