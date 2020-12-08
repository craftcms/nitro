package database

import (
	"archive/tar"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/craftcms/nitro/internal/datetime"
	"github.com/craftcms/nitro/internal/helpers"
	"github.com/craftcms/nitro/labels"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
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

			// TODO(jasonmccallister) prompt the user for the container to import
			var containerID, containerCompatability, containerName string

			var containerList []string
			for _, c := range containers {
				containerList = append(containerList, strings.TrimLeft(c.Names[0], "/"))
			}

			// prompt the user for which database to backup
			containerSelection, err := promptForOption(os.Stdin, containerList, "Which database engine? ", output)
			if err != nil {
				return err
			}

			// set the selected container
			containerName = containers[containerSelection].Names[0]
			containerID = containers[containerSelection].ID
			containerCompatability = containers[containerSelection].Labels[labels.DatabaseCompatability]

			// get a list of the databases
			var dbs []string
			switch strings.Contains(containerName, "mysql") || strings.Contains(containerName, "mariadb") {
			case true:
				// TODO(jasonmccallister) get a list of the mysql databases
				commands := []string{"mysql", "-unitro", "-pnitro", "-e", `SHOW DATABASES;`}

				// create the command and pass to exec
				exec, err := docker.ContainerExecCreate(ctx, containerID, types.ExecConfig{
					AttachStdout: true,
					AttachStderr: true,
					Tty:          false,
					Cmd:          commands,
				})
				if err != nil {
					return err
				}

				// attach to the container
				resp, err := docker.ContainerExecAttach(ctx, exec.ID, types.ExecConfig{
					AttachStdout: true,
					AttachStderr: true,
					Tty:          false,
					Cmd:          commands,
				})
				if err != nil {
					return err
				}
				defer resp.Close()

				// start the exec
				if err := docker.ContainerExecStart(ctx, exec.ID, types.ExecStartCheck{}); err != nil {
					return fmt.Errorf("unable to start the container, %w", err)
				}

				// get the output
				buf := new(bytes.Buffer)
				if _, err := io.Copy(buf, resp.Reader); err != nil {
					return err
				}

				// TODO(jasonmccallister) remove this to a helper?
				// get all the databases
				for _, db := range strings.Split(buf.String(), "\n") {
					// ignore the system defaults
					if db == "Database" || strings.Contains(db, `"Database`) || db == "information_schema" || db == "performance_schema" || db == "sys" || strings.Contains(db, "password on the command line") || db == "mysql" || db == "" {
						continue
					}

					dbs = append(dbs, db)
				}
			default:
				// get a list of the postgres databases
				commands := []string{"psql", "--username=nitro", "--command", `SELECT datname FROM pg_database WHERE datistemplate = false;`}

				// create the command and pass to exec
				exec, err := docker.ContainerExecCreate(ctx, containerID, types.ExecConfig{
					AttachStdout: true,
					AttachStderr: true,
					Tty:          false,
					Cmd:          commands,
				})
				if err != nil {
					return err
				}

				// attach to the container
				resp, err := docker.ContainerExecAttach(ctx, exec.ID, types.ExecConfig{
					AttachStdout: true,
					AttachStderr: true,
					Tty:          false,
					Cmd:          commands,
				})
				if err != nil {
					return err
				}
				defer resp.Close()

				// start the exec
				if err := docker.ContainerExecStart(ctx, exec.ID, types.ExecStartCheck{}); err != nil {
					return fmt.Errorf("unable to start the container, %w", err)
				}

				// get the output
				buf := new(bytes.Buffer)
				if _, err := io.Copy(buf, resp.Reader); err != nil {
					return err
				}

				// TODO(jasonmccallister) remove this to a helper?
				// split the lines
				sp := strings.Split(buf.String(), "\n")
				for i, d := range sp {
					// remove the first, second, last, rows, and empty lines
					if i == 0 || i == 1 || i == len(sp) || strings.Contains(d, "rows)") || d == "" {
						continue
					}

					dbs = append(dbs, strings.TrimSpace(d))
				}
			}

			// prompt the user for the specific database to backup
			var db string
			switch len(dbs) {
			case 1:
				output.Info("There is only one database to backup...")

				db = dbs[0]
			case 0:
				return fmt.Errorf("no databases found")
			default:
				dbSelection, err := promptForOption(os.Stdin, dbs, "Which database should we backup? ", output)
				if err != nil {
					return err
				}

				db = dbs[dbSelection]
			}

			output.Info("Preparing backup...")

			// create a backup with the current timestamp
			backup := fmt.Sprintf("%s-%s.sql", db, datetime.Parse(time.Now()))

			// create the backup command based on the compatability type
			var backupCmd []string
			switch containerCompatability {
			case "postgres":
				backupCmd = []string{"pg_dump", "--username=nitro", db, "-f", "/tmp/" + backup}
			default:
				backupCmd = []string{"/usr/bin/mysqldump", "-h", "127.0.0.1", "-unitro", "--password=nitro", db, "--result-file=" + "/tmp/" + backup}
			}

			output.Pending("creating backup", backup)

			// create the command and pass to exec
			exec, err := docker.ContainerExecCreate(ctx, containerID, types.ExecConfig{
				AttachStdout: true,
				AttachStderr: true,
				Tty:          false,
				Cmd:          backupCmd,
			})
			if err != nil {
				return err
			}

			// attach to the container
			stream, err := docker.ContainerExecAttach(ctx, exec.ID, types.ExecConfig{
				AttachStdout: true,
				AttachStderr: true,
				Tty:          false,
				Cmd:          backupCmd,
			})
			if err != nil {
				return err
			}
			defer stream.Close()

			// start the exec
			if err := docker.ContainerExecStart(ctx, exec.ID, types.ExecStartCheck{}); err != nil {
				return fmt.Errorf("unable to start the container, %w", err)
			}

			// wait for the container exec to complete
			waiting := true
			for waiting {
				resp, err := docker.ContainerExecInspect(ctx, exec.ID)
				if err != nil {
					return err
				}

				waiting = resp.Running
			}

			// copy the backup file from the container
			rdr, stat, err := docker.CopyFromContainer(ctx, containerID, "/tmp/"+backup)
			if err != nil || stat.Mode.IsRegular() == false {
				return err
			}
			defer rdr.Close()

			// read the content of the file, the file is in a tar format
			buf := new(bytes.Buffer)
			tr := tar.NewReader(rdr)

			for {
				_, err := tr.Next()
				// if end of tar archive
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}

				buf.ReadFrom(tr)
			}

			// make the backup directory if it does not exist
			backupDir := filepath.Join(home, ".nitro", "backups", env, containerName)
			if err := helpers.MkdirIfNotExists(backupDir); err != nil {
				return err
			}

			// write the file to the backups dir
			if err := ioutil.WriteFile(filepath.Join(backupDir, backup), buf.Bytes(), 0644); err != nil {
				return err
			}

			output.Done()

			output.Info("Backup saved in", filepath.Join(backupDir), "💾")

			return nil
		},
	}

	return cmd
}

func promptForOption(reader io.Reader, options []string, prompt string, output terminal.Outputer) (int, error) {
	for k, v := range options {
		output.Info(fmt.Sprintf("  %d. %s", k+1, v))
	}

	fmt.Print(prompt)
	fmt.Scan("")
	rdr := bufio.NewReader(reader)
	char, err := rdr.ReadString('\n')
	if err != nil {
		return 0, err
	}

	// remove the new line from string
	char = strings.TrimSpace(char)

	// convert the selection to an integer
	selection, err := strconv.Atoi(char)
	if err != nil {
		return 0, err
	}

	// make sure its there
	if len(options) < selection {
		return 0, err
	}

	// take away one from the selection
	selection = selection - 1

	return selection, nil
}
