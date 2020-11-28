package database

import (
	"io/ioutil"
	"os"

	"github.com/craftcms/nitro/internal/database"
	"github.com/craftcms/nitro/labels"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/h2non/filetype"
	"github.com/spf13/cobra"
)

var importExampleText = `  # import a sql file into a database
  nitro db import filename.sql

  # use a relative path
  nitro db import ~/Desktop/backup.sql

  # use an absolute path
  nitro db import /Users/oli/Desktop/backup.sql`

// importCommand is the command for creating new development environments
func importCommand(docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import a database",
		Args:  cobra.MinimumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{"sql", "gz", "zip", "dump"}, cobra.ShellCompDirectiveFilterFileExt
		},
		Example: importExampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			env := cmd.Flag("environment").Value.String()

			output.Info("Preparing import...")

			// TODO(jasonmccallister) get the abs clean path for the file
			file, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer file.Close()

			b, err := ioutil.ReadFile(file.Name())
			if err != nil {
				return err
			}

			// check if the file is an archive
			compressed := false
			if filetype.IsArchive(b) {
				compressed = true
			}

			// dectect the type of backup if not compressed
			detected := ""
			if compressed == false {
				detected, err = database.DetermineEngine(file.Name())
				if err != nil {
					return err
				}
			}

			if detected != "" {
				output.Success("detected", detected, "backup")
			}

			// TODO(jasonmccallister) get a list of all the databases
			filter := filters.NewArgs()
			filter.Add("label", labels.Environment+"="+env)
			filter.Add("label", labels.Type+"=database")

			switch detected {
			case "mysql":
				filter.Add("label", labels.DatabaseCompatability+"=mysql")
			case "postgres":
				filter.Add("label", labels.DatabaseCompatability+"=postgres")
			}

			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter})
			if err != nil {
				return err
			}

			// TODO(jasonmccallister) prompt the user for the container to import
			var containerID string
			for _, c := range containers {
				containerID = c.ID
				output.Info(c.Names[0])
			}

			// TODO(jasonmccallister) if the file is not compressed, compress it to send to the api
			switch compressed {
			case false:
				rdr, err := archive.Generate(file.Name())
				if err != nil {
					return err
				}

				// https://stackoverflow.com/questions/47641799/unable-to-understand-docker-cp-command
				if err := docker.CopyToContainer(cmd.Context(), containerID, "/tmp", rdr, types.CopyToContainerOptions{AllowOverwriteDirWithFile: true}); err != nil {
					return err
				}
			}

			// TODO(jasonmccallister) copy the file, in tar format, to the container in the tmp directory

			// TODO(jasonmccallister) determine if the backup is to mysql or postgres and run the import file command
			return nil
		},
	}

	return cmd
}
