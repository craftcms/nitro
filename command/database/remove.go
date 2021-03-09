package database

import (
	"fmt"
	"sort"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/craftcms/nitro/pkg/backup"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/craftcms/nitro/protob"
)

var removeExampleText = `  # remove a database
  nitro db remove`

func removeCommand(docker client.CommonAPIClient, nitrod protob.NitroClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Short:   "Remove a database",
		Example: removeExampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// add filters to show only the environment and database containers
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)
			filter.Add("label", containerlabels.Type+"=database")

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
			selectedEngine, err := output.Select(cmd.InOrStdin(), "Which database engine? ", containerList)
			if err != nil {
				return err
			}

			// get the containers info
			info, err := docker.ContainerInspect(cmd.Context(), containers[selectedEngine].ID)
			if err != nil {
				return err
			}

			// get the containers details
			engine := info.Config.Labels[containerlabels.DatabaseCompatibility]
			hostname := strings.TrimLeft(info.Name, "/")
			version := info.Config.Labels[containerlabels.DatabaseVersion]
			var port string
			// get the port from the container info
			for p, bind := range info.HostConfig.PortBindings {
				for _, v := range bind {
					if v.HostPort != "" {
						port = p.Port()
					}
				}
			}

			// get all of the databases
			databases, err := backup.Databases(cmd.Context(), docker, info.ID, engine)
			if err != nil {
				return err
			}

			// ask the user which database
			selected, err := output.Select(cmd.InOrStdin(), "Which database should we remove? ", databases)
			if err != nil {
				return err
			}

			db := databases[selected]

			// wait for the api to be ready
			for {
				_, err := nitrod.Ping(cmd.Context(), &protob.PingRequest{})
				if err == nil {
					break
				}
			}

			output.Pending("removing", db)

			// remove the database
			resp, err := nitrod.RemoveDatabase(cmd.Context(), &protob.RemoveDatabaseRequest{
				Database: &protob.DatabaseInfo{
					Engine:   engine,
					Hostname: hostname,
					Version:  version,
					Port:     port,
					Database: db,
				},
			})
			// check if the error code is unimplemented
			if code := status.Code(err); code == codes.Unimplemented {
				output.Warning()

				// ask if the update command should run
				confirm, err := output.Confirm("The API does not appear to be updated, run `nitro update` now", true, "?")
				if err != nil {
					return err
				}

				if !confirm {
					output.Info("Skipping the update command, you need to update before using this command")

					return nil
				}

				// run the update command
				for _, c := range cmd.Parent().Commands() {
					// set the update command
					if c.Use == "update" {
						if err := c.RunE(c, args); err != nil {
							return err
						}
					}
				}
			}
			if err != nil {
				return err
			}

			output.Done()

			output.Info(fmt.Sprintf("%s 💪", resp.Message))

			return nil
		},
	}

	return cmd
}
