package database

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/database"
	"github.com/craftcms/nitro/pkg/filetype"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/pathexists"
	"github.com/craftcms/nitro/pkg/terminal"
)

var importExampleText = `  # import a sql file into a database
  nitro db import filename.sql

  # use a relative path
  nitro db import ~/Desktop/backup.sql

  # use an absolute path
  nitro db import /Users/oli/Desktop/backup.sql`

// importCommand is the command for creating new development environments
func importCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import a database",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				fmt.Println(cmd.UsageString())

				return fmt.Errorf("database backup file path param missing")
			}

			return nil
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{"sql", "gz", "zip", "dump"}, cobra.ShellCompDirectiveFilterFileExt
		},
		Example: importExampleText,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// make sure the file exists
			if exists := pathexists.IsFile(args[0]); !exists {
				return fmt.Errorf("unable to find file %s", args[0])
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			show, err := strconv.ParseBool(cmd.Flag("show-output").Value.String())
			if err != nil {
				// set to false
				show = false
			}

			// replace the relative path with the full directory
			path := args[0]
			if strings.HasPrefix(path, "~") {
				path = strings.Replace(path, "~", home, 1)
			}

			// check if this is a known archive type for docker
			supportedArchive := archive.IsArchivePath(path)

			// check if this is a zip file
			var isZip bool
			kind, err := filetype.Determine(path)
			if err != nil {
				return err
			}

			switch kind {
			case "zip", "tar":
				isZip = true
			}

			// detect the type of backup if not compressed
			detected := ""
			if !supportedArchive && !isZip {
				output.Pending("detecting backup type")

				// determine the database engine
				detected, err = database.DetermineEngine(path)
				if err != nil {
					output.Warning()

					return err
				}

				output.Done()

				output.Info("Detected", detected, "backup")
			}

			// add filters to show only the environment and database containers
			filter := filters.NewArgs()
			filter.Add("label", labels.Nitro)
			filter.Add("label", labels.Type+"=database")

			// if we detected the engine type, add the compatibility label to the filter
			switch detected {
			case "mysql":
				filter.Add("label", labels.DatabaseCompatibility+"=mysql")
			case "postgres":
				filter.Add("label", labels.DatabaseCompatibility+"=postgres")
			}

			// get a list of all the databases
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter})
			if err != nil {
				return err
			}

			// sort containers by the name
			sort.SliceStable(containers, func(i, j int) bool {
				return containers[i].Names[0] < containers[j].Names[0]
			})

			// get all of the containers as a list
			var options []string
			for _, c := range containers {
				options = append(options, strings.TrimLeft(c.Names[0], "/"))
			}

			// prompt the user for the engine to import the backup into
			var containerID string
			selected, err := output.Select(os.Stdin, "Select a database engine: ", options)
			if err != nil {
				return err
			}

			// set the container id
			containerID = containers[selected].ID
			if containerID == "" {
				return fmt.Errorf("unable to get the container")
			}

			// ask the user for the database to create
			db, err := output.Ask("Enter the database name", "", ":", nil)
			if err != nil {
				return err
			}

			output.Info("Preparing import…")

			// get the reader, along with the filename to use to set the container path
			rdr, filename, err := database.PrepareArchiveFromPath(path)
			if err != nil {
				return err
			}

			output.Pending("uploading backup", filename)

			// copy the file into the container
			if err := docker.CopyToContainer(cmd.Context(), containerID, "/tmp", rdr, types.CopyToContainerOptions{}); err != nil {
				output.Warning()
				return err
			}

			// set the path to the file in the container
			containerPath := "/tmp/" + filename

			// wait for the file to exist
			waiting := true
			for waiting {
				_, err := docker.ContainerStatPath(cmd.Context(), containerID, containerPath)
				if err == nil {
					waiting = false
				}

				if !waiting {
					break
				}
			}

			// determine if the backup is to mysql or postgres and run the import file command
			var createCmd, importCmd []string
			switch detected {
			case "postgres":
				createCmd = []string{"psql", "--username=nitro", "--host=127.0.0.1", fmt.Sprintf(`-c CREATE DATABASE %s;`, db)}
				importCmd = []string{"psql", "--username=nitro", "--host=127.0.0.1", db, "--file", containerPath}
			default:
				createCmd = []string{"mysql", "-uroot", "-pnitro", fmt.Sprintf(`-e CREATE DATABASE IF NOT EXISTS %s;`, db)}
				// https: //dev.mysql.com/doc/refman/8.0/en/mysql-command-options.html
				importCmd = []string{"mysql", "-unitro", "-pnitro", db, fmt.Sprintf(`-e source %s`, containerPath)}
			}

			// create the database
			if _, err := execCreate(cmd.Context(), docker, containerID, createCmd, show); err != nil {
				output.Warning()
				return fmt.Errorf("unable to create the database, %w", err)
			}

			// create the exec for create
			createExec, err := docker.ContainerExecCreate(cmd.Context(), containerID, types.ExecConfig{
				AttachStdout: true,
				AttachStderr: true,
				AttachStdin:  true,
				Tty:          false,
				Cmd:          createCmd,
			})
			if err != nil {
				return err
			}

			// attach to the container
			createResp, err := docker.ContainerExecAttach(cmd.Context(), createExec.ID, types.ExecStartCheck{
				Tty: false,
			})
			if err != nil {
				return err
			}
			defer createResp.Close()

			// should we display output?
			if show {
				// show the output to stdout and stderr
				if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, createResp.Reader); err != nil {
					return fmt.Errorf("unable to copy the output of container, %w", err)
				}
			}

			// start the exec
			if err := docker.ContainerExecStart(cmd.Context(), createExec.ID, types.ExecStartCheck{}); err != nil {
				return fmt.Errorf("unable to start the container, %w", err)
			}

			// wait for the container exec to complete
			createWaiting := true
			for createWaiting {
				resp, err := docker.ContainerExecInspect(cmd.Context(), createExec.ID)
				if err != nil {
					return err
				}

				createWaiting = resp.Running
			}

			// create the exec for import
			importExec, err := docker.ContainerExecCreate(cmd.Context(), containerID, types.ExecConfig{
				AttachStdout: true,
				AttachStderr: true,
				AttachStdin:  true,
				Tty:          true,
				Cmd:          importCmd,
			})
			if err != nil {
				return err
			}

			// attach to the container
			importResp, err := docker.ContainerExecAttach(cmd.Context(), importExec.ID, types.ExecStartCheck{
				Tty: false,
			})
			if err != nil {
				return err
			}
			defer importResp.Close()

			// should we display output?
			if show {
				// show the output to stdout and stderr
				if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, importResp.Reader); err != nil {
					return fmt.Errorf("unable to copy the output of container, %w", err)
				}
			}

			// start the exec
			if err := docker.ContainerExecStart(cmd.Context(), importExec.ID, types.ExecStartCheck{}); err != nil {
				return fmt.Errorf("unable to start the container, %w", err)
			}

			// wait for the container exec to complete
			importWaiting := true
			for importWaiting {
				resp, err := docker.ContainerExecInspect(cmd.Context(), importExec.ID)
				if err != nil {
					return err
				}

				importWaiting = resp.Running
			}

			output.Done()

			output.Info("Import successful 💪")

			return nil
		},
	}

	cmd.Flags().Bool("show-output", false, "show debug from import")

	return cmd
}

func getReader(supportedArchive bool, kind, path string) (io.Reader, string, error) {
	_, filename := filepath.Split(path)

	// if the file is a zip or gzip file
	if !supportedArchive {
		// get the file content based on the kind
		switch kind {
		case "zip":
			// create a new zip reader
			r, err := zip.OpenReader(path)
			if err != nil {
				return nil, "", err
			}
			defer r.Close()

			// read each of the files
			for _, file := range r.File {
				if strings.HasSuffix(file.Name, ".sql") {
					// create the temp file
					temp, err := ioutil.TempFile(os.TempDir(), "nitro-import-zip-"+filename)
					if err != nil {
						return nil, "", err
					}
					defer temp.Close()

					// read the content of the zip file
					rc, err := file.Open()
					if err != nil {
						return nil, "", err
					}

					buf := new(bytes.Buffer)
					if _, err := buf.ReadFrom(rc); err != nil && !errors.Is(err, io.EOF) {
						return nil, "", err
					}

					// write the content to the temp file
					if _, err := temp.Write(buf.Bytes()); err != nil {
						return nil, "", err
					}

					// read the file content from the temp file
					content, err := ioutil.ReadFile(temp.Name())
					if err != nil {
						return nil, "", err
					}

					reader, err := archive.Generate(file.Name, string(content))
					return reader, file.Name, err
				}
			}

			// we did not find a sql file, so we need to return an error
			return nil, "", fmt.Errorf("unable to find a .sql file in the zip")
		default:
			return nil, "", fmt.Errorf("unsupported kind %q provided", kind)
		}
	}

	// read the file content
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, "", err
	}

	// generate the reader
	reader, err := archive.Generate(filename, string(content))
	return reader, filename, err
}
