package database

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/backup"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
)

var removeExampleText = `  # remove a database
  nitro db remove`

func removeCommand(docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Short:   "Remove a database",
		Example: removeExampleText,
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

			// sort containers by the name
			sort.SliceStable(containers, func(i, j int) bool {
				return containers[i].Names[0] < containers[j].Names[0]
			})

			// generate a list of engines for the prompt
			var containerList []string
			for _, c := range containers {
				containerList = append(containerList, strings.TrimLeft(c.Names[0], "/"))
			}

			// prompt the user for which database to backup
			engine, err := output.Select(cmd.InOrStdin(), "Which database engine? ", containerList)
			if err != nil {
				return err
			}

			id := containers[engine].ID
			compatibility := containers[engine].Labels[labels.DatabaseCompatibility]

			// get all of the databases
			databases, err := backup.Databases(cmd.Context(), docker, id, compatibility)
			if err != nil {
				return err
			}

			// ask the user which database
			selected, err := output.Select(cmd.InOrStdin(), "Which database should we remove? ", databases)
			if err != nil {
				return err
			}

			db := databases[selected]

			// set the commands based on the engine type
			var cmds []string
			switch compatibility {
			case "mysql":
				cmds = []string{"mysqladmin", "-user=root", "-pnitro", "drop", db}
			default:
				cmds = []string{"psql", "--username=nitro", "--host=127.0.0.1", fmt.Sprintf(`-c DROP DATABASE IF EXISTS %s;`, db)}
			}

			output.Pending("removing", db)

			// execute the command to create the database
			if _, err := execCreate(cmd.Context(), docker, id, cmds, show); err != nil {
				return err
			}

			output.Done()

			output.Info("Database removed 💪")

			return nil
		},
	}

	cmd.Flags().Bool("show-output", false, "show debug from import")

	return cmd
}
