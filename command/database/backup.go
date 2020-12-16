package database

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/backup"
	"github.com/craftcms/nitro/datetime"
	"github.com/craftcms/nitro/labels"
	"github.com/craftcms/nitro/terminal"
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
			env := cmd.Flag("environment").Value.String()
			ctx := cmd.Context()

			// add filters to show only the envrionment and database containers
			filter := filters.NewArgs()
			filter.Add("label", labels.Environment+"="+env)
			filter.Add("label", labels.Type+"=database")

			// get a list of all the databases
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter})
			if err != nil {
				return err
			}

			// generate a list of engines for the prompt
			var containerList []string
			for _, c := range containers {
				containerList = append(containerList, strings.TrimLeft(c.Names[0], "/"))
			}

			output.Info("Getting ready to backup...")

			// get the container id, name, and database from the user
			containerID, containerName, compatability, db, err := backup.Prompt(ctx, os.Stdin, docker, output, containers, containerList)
			if err != nil {
				return err
			}

			output.Info("ID:", containerID)
			output.Info("Name:", containerName)
			output.Info("Compatability:", compatability)

			output.Info("Preparing backup...")

			// create the options for the backup
			opts := &backup.Options{
				BackupName:    fmt.Sprintf("%s-%s.sql", db, datetime.Parse(time.Now())),
				ContainerID:   containerID,
				ContainerName: containerName,
				Database:      db,
				Environment:   env,
				Home:          home,
			}

			// create the backup command based on the compatability type
			switch compatability {
			case "postgres":
				opts.Commands = []string{"pg_dump", "--username=nitro", db, "-f", "/tmp/" + opts.BackupName}
			default:
				opts.Commands = []string{"/usr/bin/mysqldump", "-h", "127.0.0.1", "-unitro", "--password=nitro", db, "--result-file=" + "/tmp/" + opts.BackupName}
			}

			output.Pending("creating backup", opts.BackupName)

			// perform the backup
			if err := backup.Perform(ctx, docker, opts); err != nil {
				output.Warning()

				return fmt.Errorf("unable to backup the database, %w", err)
			}

			output.Done()

			output.Info("Backup saved in", filepath.Join(opts.Home, ".nitro", opts.ContainerName), "💾")

			return nil
		},
	}

	return cmd
}
