package database

import (
	"archive/tar"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/craftcms/nitro/internal/database"
	"github.com/craftcms/nitro/labels"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
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

			// add filters to show only the envrionment and database containers
			filter := filters.NewArgs()
			filter.Add("label", labels.Environment+"="+env)
			filter.Add("label", labels.Type+"=database")

			// if we detected the engine type, add the compatability label to the filter
			switch detected {
			case "mysql":
				filter.Add("label", labels.DatabaseCompatability+"=mysql")
			case "postgres":
				filter.Add("label", labels.DatabaseCompatability+"=postgres")
			}

			// get a list of all the databases
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter})
			if err != nil {
				return err
			}

			// TODO(jasonmccallister) prompt the user for the container to import
			var containerID string
			for _, c := range containers {
				containerID = c.ID
			}

			// if the file is not compressed, compress it to send to the api
			switch compressed {
			case false:
				rdr, err := newTarArchiveFromFile(file)
				if err != nil {
					return err
				}

				// copy the file into the container
				if err := docker.CopyToContainer(cmd.Context(), containerID, "/tmp", rdr, types.CopyToContainerOptions{}); err != nil {
					return err
				}
			}

			// TODO(jasonmccallister) determine if the backup is to mysql or postgres and run the import file command

			output.Info("Import successful")

			return nil
		},
	}

	return cmd
}

func newTarArchiveFromFile(file *os.File) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	info, err := os.Stat(file.Name())
	if err != nil {
		return nil, err
	}

	header, err := tar.FileInfoHeader(info, file.Name())
	if err != nil {
		return nil, err
	}

	// header.Name = strings.TrimPrefix(strings.Replace(file, path, "", -1), string(filepath.Separator))
	err = tw.WriteHeader(header)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, fmt.Errorf("is directory")
	}

	_, err = io.Copy(tw, file)
	if err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}

	return bufio.NewReader(&buf), nil
}
