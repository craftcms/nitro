package database

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/backup"
	"github.com/craftcms/nitro/pkg/datetime"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
)

var backupExampleText = `  # backup a database
  nitro db backup`

// backupCommand is the command for backing up an individual database or
func backupCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "backup",
		Short:   "Backup a database",
		Example: backupExampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// add filters to show only the environment and database containers
			filter := filters.NewArgs()
			filter.Add("label", labels.Nitro)
			filter.Add("label", labels.Type+"=database")

			// get a list of all the databases
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter})
			if err != nil {
				return err
			}

			// sort containers by the name
			sort.SliceStable(containers, func(i, j int) bool {
				return containers[i].Names[0] < containers[j].Names[0]
			})

			// generate a list of engines for the prompt
			var containerList []string
			for _, c := range containers {
				containerList = append(containerList, strings.TrimLeft(c.Names[0], "/"))
			}

			output.Info("Getting ready to backup…")

			// get the container id, name, and database from the user
			containerID, containerName, compatibility, db, err := backup.Prompt(ctx, os.Stdin, docker, output, containers, containerList)
			if err != nil {
				return err
			}

			output.Info("Preparing backup…")

			// create the options for the backup
			opts := &backup.Options{
				BackupName:    fmt.Sprintf("%s-%s.sql", db, datetime.Parse(time.Now())),
				ContainerID:   containerID,
				ContainerName: containerName,
				Database:      db,
				Home:          home,
			}

			// create the backup command based on the compatibility type
			switch compatibility {
			case "postgres":
				opts.Commands = []string{"pg_dump", "--username=nitro", db, "-f", "/tmp/" + opts.BackupName}
			default:
				opts.Commands = []string{"/usr/bin/mysqldump", "-h", "127.0.0.1", "--username=nitro", "--password=nitro", db, "--result-file=" + "/tmp/" + opts.BackupName}
			}

			output.Pending("creating backup", opts.BackupName)

			// perform the backup
			if err := backup.Perform(ctx, docker, opts); err != nil {
				output.Warning()

				return fmt.Errorf("unable to backup the database, %w", err)
			}

			output.Done()

			output.Info("Backup saved in", filepath.Join(opts.Home, ".nitro", "backups", opts.ContainerName), "💾")

			return nil
		},
	}

	return cmd
}
