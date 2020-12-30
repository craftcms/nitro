package database

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/h2non/filetype"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/database"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
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
			show, err := strconv.ParseBool(cmd.Flag("show-output").Value.String())
			if err != nil {
				// set to false
				show = false
			}

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

			output.Info("Preparing import...")

			// dectect the type of backup if not compressed
			detected := ""
			if compressed == false {
				output.Pending("detecting backup type")

				detected, err = database.DetermineEngine(file.Name())
				if err != nil {
					return err
				}

				output.Done()

				output.Info("Detected", detected, "backup")
			}

			// add filters to show only the envrionment and database containers
			filter := filters.NewArgs()
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

			// get all of the containers as a list
			var engineOpts []string
			for _, c := range containers {
				engineOpts = append(engineOpts, strings.TrimLeft(c.Names[0], "/"))
			}

			// prompt the user for the engine to import the backup into
			var containerID string
			selected, err := output.Select(os.Stdin, "Select a database engine: ", engineOpts)

			// set the container id
			containerID = containers[selected].ID
			if containerID == "" {
				return fmt.Errorf("unable to get the container")
			}

			// ask the user for the database to create
			msg := "Enter the database name: "

			fmt.Print(msg)

			var db string
			for {
				rdr := bufio.NewReader(os.Stdin)
				input, err := rdr.ReadString('\n')
				if err != nil {
					return err
				}

				if strings.ContainsAny(input, " -") == false {
					db = strings.TrimSpace(input)
					break
				}

				fmt.Println("  no spaces or hypens allowed 🙄...")
				fmt.Print(msg)
			}

			output.Pending("uploading backup")

			var rdr io.Reader
			switch compressed {
			case false:
				// generate the reader
				rdr, err = archive.Generate(file.Name())
				if err != nil {
					return err
				}
			default:
				// read the file and create a reader
				content, err := ioutil.ReadFile(file.Name())
				if err != nil {
					return err
				}

				rdr = bytes.NewReader(content)
			}

			// copy the file into the container
			if err := docker.CopyToContainer(cmd.Context(), containerID, "/tmp", rdr, types.CopyToContainerOptions{}); err != nil {
				return err
			}

			output.Done()

			// determine if the backup is to mysql or postgres and run the import file command
			var createCmd, importCmd []string
			switch detected {
			case "postgres":
				createCmd = []string{"psql", "--username=nitro", "--host=127.0.0.1", fmt.Sprintf(`-c CREATE DATABASE %s;`, db)}
				importCmd = []string{"psql", "--username=nitro", "--host=127.0.0.1", db, "--file", "/tmp/" + file.Name()}
			default:
				return fmt.Errorf("mysql imports have not been implemented")
			}

			output.Pending("importing database to", db)

			// create the database
			if _, err := execCreate(cmd.Context(), docker, containerID, createCmd, show); err != nil {
				return fmt.Errorf("unable to create the database, %w", err)
			}

			// import the database
			if _, err := execCreate(cmd.Context(), docker, containerID, importCmd, show); err != nil {
				return fmt.Errorf("unable to import the database, %w", err)
			}

			output.Done()

			output.Info("Import successful 💪")

			return nil
		},
	}

	cmd.Flags().Bool("show-output", false, "show debug from import")

	return cmd
}
