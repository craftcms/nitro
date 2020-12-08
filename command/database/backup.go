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
			for k, c := range containers {
				output.Info(fmt.Sprintf("  %d. %s", k+1, strings.TrimLeft(c.Names[0], "/")))
			}

			// ask for the engine
			fmt.Print("Select database engine: ")
			fmt.Scan("")
			reader := bufio.NewReader(os.Stdin)
			char, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			// remove the new line from string
			char = strings.TrimSpace(char)

			// convert the selection to an integer
			selection, err := strconv.Atoi(char)
			if err != nil {
				return err
			}

			// make sure its there
			if len(containers) < selection {
				return err
			}

			// take away one from the selection
			selection = selection - 1

			// set the selected container
			containerName = containers[selection].Names[0]
			containerID = containers[selection].ID
			containerCompatability = containers[selection].Labels[labels.DatabaseCompatability]

			output.Info("Preparing backup...")

			// create a backup with the current timestamp
			// TODO(jasonmccallister) replace with the database to backup from the prompt
			backup := fmt.Sprintf("nitro-%s.sql", datetime.Parse(time.Now()))

			// create the backup command based on the compatability type
			var backupCmd []string
			switch containerCompatability {
			case "postgres":
				backupCmd = []string{"pg_dump", "-Unitro", "-f", "/tmp/" + backup}
			default:
				backupCmd = []string{"/usr/bin/mysqldump", "-h", "127.0.0.1", "-unitro", "--password=nitro", "nitro", "--result-file=" + "/tmp/" + backup}
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
