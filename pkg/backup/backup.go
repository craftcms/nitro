package backup

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/craftcms/nitro/internal/helpers"
	"github.com/craftcms/nitro/labels"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Options are used to pass options to a database backup func.
// The options contain information such as the environment,
// container, home directory, and database to backup.
type Options struct {
	Environment   string
	Home          string
	ContainerID   string
	ContainerName string
	Database      string
	BackupName    string
	Commands      []string
}

// Prompt is used to ask a user for input and walk them
// through selecting a database engine (container) and
// a database to backup. It will return the container ID
// as the first string, the database name to backup, and
// the last return is an error.
func Prompt(ctx context.Context, reader io.Reader, docker client.ContainerAPIClient, output terminal.Outputer, containers []types.Container, containerList []string) (string, string, string, string, error) {
	// prompt the user for which database to backup
	selected, err := output.Select(reader, "Which database engine? ", containerList)
	if err != nil {
		return "", "", "", "", err
	}

	// get the selected container details
	name := containers[selected].Names[0]
	id := containers[selected].ID
	compatability := containers[selected].Labels[labels.DatabaseCompatability]

	// TODO(jasonmccallister) abstract to databases

	// get a list of the databases from the container
	var commands []string
	if strings.Contains(name, "mysql") || strings.Contains(name, "mariadb") {
		// get a list of the mysql databases
		commands = []string{"mysql", "-unitro", "-pnitro", "-e", `SHOW DATABASES;`}
	} else {
		commands = []string{"psql", "--username=nitro", "--command", `SELECT datname FROM pg_database WHERE datistemplate = false;`}
	}

	// create the command and pass to exec
	exec, err := docker.ContainerExecCreate(ctx, id, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Cmd:          commands,
	})
	if err != nil {
		return "", "", "", "", err
	}

	// attach to the container
	resp, err := docker.ContainerExecAttach(ctx, exec.ID, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Cmd:          commands,
	})
	if err != nil {
		return "", "", "", "", err
	}
	defer resp.Close()

	// start the exec
	if err := docker.ContainerExecStart(ctx, exec.ID, types.ExecStartCheck{}); err != nil {
		return "", "", "", "", fmt.Errorf("unable to start the container, %w", err)
	}

	// get all of the databases based on the engine
	databases, err := Databases(resp.Reader, compatability)
	if err != nil {
		return "", "", "", "", err
	}

	// TODO(jasonmccallister) abstract to databases

	// prompt the user for the specific database to backup
	var db string
	switch len(databases) {
	case 1:
		output.Info("There is only one database to backup...")

		db = databases[0]
	case 0:
		return "", "", "", "", fmt.Errorf("no databases found")
	default:
		selected, err := output.Select(os.Stdin, "Which database should we backup? ", databases)
		if err != nil {
			return "", "", "", "", err
		}

		db = databases[selected]
	}

	return id, name, compatability, db, nil
}

func Databases(reader io.Reader, compatability string) ([]string, error) {
	// get the output
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, reader); err != nil {
		return nil, err
	}

	var databases []string
	switch compatability {
	case "mysql":
		// get all the databases from the mysql engine
		for _, d := range strings.Split(buf.String(), "\n") {
			// ignore the system defaults
			if d == "Database" || strings.Contains(d, `"Database`) || d == "information_schema" || d == "performance_schema" || d == "sys" || strings.Contains(d, "password on the command line") || d == "mysql" || d == "" {
				continue
			}

			databases = append(databases, strings.TrimSpace(d))
		}
	default:
		// get all the databases from the postgres engine
		sp := strings.Split(buf.String(), "\n")
		for i, d := range sp {
			// remove the first, second, last, rows, and empty lines
			if i == 0 || i == 1 || i == len(sp) || strings.Contains(d, "rows)") || d == "" {
				continue
			}

			databases = append(databases, strings.TrimSpace(d))
		}
	}

	return databases, nil
}

// Perform is used to perform a backup for a database container, it does not prompt the user as it assumed the Prompt func above
// is used to determine the engine (container) and the specific database to backup. Perform accepts the backup commands and is
// agnositic to the database engine for the requested backup.
func Perform(ctx context.Context, docker client.ContainerAPIClient, opts *Options) error {
	// create the backup in the container
	exec, err := docker.ContainerExecCreate(ctx, opts.ContainerID, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Cmd:          opts.Commands,
	})
	if err != nil {
		return err
	}

	// attach to the container
	stream, err := docker.ContainerExecAttach(ctx, exec.ID, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Cmd:          opts.Commands,
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

	// copy the backup from the container into the host machine
	rdr, stat, err := docker.CopyFromContainer(ctx, opts.ContainerID, "/tmp/"+opts.BackupName)
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
	dir := filepath.Join(opts.Home, ".nitro", "backups", opts.Environment, opts.ContainerID)
	if err := helpers.MkdirIfNotExists(dir); err != nil {
		return err
	}

	// write the file to the backups dir
	if err := ioutil.WriteFile(filepath.Join(dir, opts.BackupName), buf.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}
