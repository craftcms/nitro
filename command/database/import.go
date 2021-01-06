package database

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
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

			// get the clean path for the file
			path := filepath.Join(args[0])

			// open the file
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// read the file
			b, err := ioutil.ReadFile(file.Name())
			if err != nil {
				return err
			}

			// check if the file is an archive
			compressed := false
			if filetype.IsArchive(b) {
				compressed = true
			}

			output.Info("Preparing import…")

			// detect the type of backup if not compressed
			detected := ""
			if !compressed {
				output.Pending("detecting backup type")

				detected, err = database.DetermineEngine(file.Name())
				if err != nil {
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

			// get all of the containers as a list
			var engineOpts []string
			for _, c := range containers {
				engineOpts = append(engineOpts, strings.TrimLeft(c.Names[0], "/"))
			}

			// prompt the user for the engine to import the backup into
			var containerID string
			selected, err := output.Select(os.Stdin, "Select a database engine: ", engineOpts)
			if err != nil {
				return err
			}

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

				if !strings.ContainsAny(input, " -") {
					db = strings.TrimSpace(input)
					break
				}

				fmt.Println("  no spaces or hyphens allowed…")
				fmt.Print(msg)
			}

			output.Pending("importing backup")

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

			// should we steam or copy into the container
			stream := true
			if cmd.Flag("stream").Value.String() != "true" {
				stream = false
			}

			// should we stream the output
			switch stream {
			case true:
				if err := streamToContainer(cmd.Context(), file, containers[selected], db); err != nil {
					return err
				}
			default:
				// copy into the container if not streaming
				if err := copyToContainer(cmd.Context(), docker, output, show, file, rdr, detected, db, containerID); err != nil {
					output.Warning()
					return err
				}
			}

			output.Done()

			output.Info("Import successful 💪")

			return nil
		},
	}

	cmd.Flags().Bool("show-output", false, "show debug from import")
	cmd.Flags().Bool("stream", false, "stream the file contents into the database")

	return cmd
}

func streamToContainer(ctx context.Context, file *os.File, container types.Container, database string) error {
	// create the connection to the database
	if len(container.Ports) == 0 {
		return fmt.Errorf("expected the container to have ports")
	}

	// get the IP and port of the database
	ip := container.Ports[0].IP
	port := container.Ports[0].PublicPort

	var db *sql.DB
	var createStatement, privilegesStatement string
	switch container.Labels[labels.DatabaseCompatibility] {
	case "mysql":
		conn, err := sql.Open("mysql", fmt.Sprintf("root:nitro@tcp(%s:%d)/nitro", ip, port))
		if err != nil {
			return fmt.Errorf("error opening mysql connection: %w", err)
		}

		createStatement = fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", database)
		privilegesStatement = fmt.Sprintf(`GRANT ALL PRIVILEGES ON * TO '%s'@'%s';`, "nitro", "%")

		db = conn
	default:
		conn, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", ip, port, "nitro", "nitro"))
		if err != nil {
			return fmt.Errorf("error opening postgres connection: %w", err)
		}

		createStatement = fmt.Sprintf("CREATE DATABASE %s;", database)

		db = conn
	}
	defer db.Close()

	// make sure we can reach the database
	for {
		if err := db.PingContext(ctx); err == nil {
			break
		}
	}

	if createStatement == "" {
		return fmt.Errorf("expected the create statement but it was empty")
	}

	// create the database
	if _, err := db.ExecContext(ctx, createStatement); err != nil {
		return err
	}

	// do we need to set the privileges on the database?
	if privilegesStatement != "" {
		if _, err := db.ExecContext(ctx, privilegesStatement); err != nil {
			return err
		}
	}

	reader := bufio.NewReader(file)
	for {
		n, err := reader.ReadString(';')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// check for comments and empty lines
		if strings.HasPrefix(n, "--") || n == "" {
			continue
		}

		if _, err := db.ExecContext(ctx, string(n)); err != nil {
			return fmt.Errorf("\nerror executing import on line %s\nerror: %w", string(n), err)
		}
	}

	return nil
}

func copyToContainer(ctx context.Context, docker client.CommonAPIClient, output terminal.Outputer, show bool, file *os.File, rdr io.Reader, detected, db, containerID string) error {
	// copy the file into the container
	if err := docker.CopyToContainer(ctx, containerID, "/tmp", rdr, types.CopyToContainerOptions{}); err != nil {
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
	if _, err := execCreate(ctx, docker, containerID, createCmd, show); err != nil {
		return fmt.Errorf("unable to create the database, %w", err)
	}

	// import the database
	if _, err := execCreate(ctx, docker, containerID, importCmd, show); err != nil {
		return fmt.Errorf("unable to import the database, %w", err)
	}

	return nil
}
