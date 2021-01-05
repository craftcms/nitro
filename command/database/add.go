package database

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
)

var addExampleTest = `  # add a new database
  nitro db add`

func addCommand(docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add",
		Short:   "Add a new database",
		Example: addExampleTest,
		RunE: func(cmd *cobra.Command, args []string) error {
			show, err := strconv.ParseBool(cmd.Flag("show-output").Value.String())
			if err != nil {
				// set to false
				show = false
			}

			// add filters to show only the environment and database containers
			filter := filters.NewArgs()
			filter.Add("label", labels.Nitro)
			filter.Add("label", labels.Type+"=database")

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

			// prompt the user for the engine to add the database
			var containerID, databaseEngine string
			selected, err := output.Select(os.Stdin, "Select the database engine: ", engineOpts)

			// set the container id and db engine
			containerID = containers[selected].ID
			databaseEngine = containers[selected].Labels[labels.DatabaseCompatibility]
			if containerID == "" {
				return fmt.Errorf("unable to get the container")
			}

			// ask the user for the database to create
			msg := "Enter the new database name: "

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

				fmt.Println("  no spaces or hyphens are allowed…")
				fmt.Print(msg)
			}

			output.Pending("creating database", db)

			// set the commands based on the engine type
			var cmds, privileges []string
			switch databaseEngine {
			case "mysql":
				cmds = []string{"mysql", "-uroot", "-pnitro", fmt.Sprintf(`-e CREATE DATABASE IF NOT EXISTS %s;`, db)}
				privileges = []string{"mysql", "-uroot", "-pnitro", fmt.Sprintf(`-e GRANT ALL PRIVILEGES ON * TO '%s'@'%s';`, "nitro", "%")}
			default:
				cmds = []string{"psql", "--username=nitro", "--host=127.0.0.1", fmt.Sprintf(`-c CREATE DATABASE %s;`, db)}
			}

			// execute the command to create the database
			if _, err := execCreate(cmd.Context(), docker, containerID, cmds, show); err != nil {
				return err
			}

			// check if we should grant privileges
			if privileges != nil {
				if _, err := execCreate(cmd.Context(), docker, containerID, privileges, show); err != nil {
					return err
				}
			}

			output.Done()

			output.Info("Database added 💪")

			return nil
		},
	}

	cmd.Flags().Bool("show-output", false, "show debug from import")

	return cmd
}
