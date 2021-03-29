package database

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/craftcms/nitro/pkg/validate"
	"github.com/craftcms/nitro/protob"
)

var addExampleTest = `  # add a new database
  nitro db add`

func addCommand(docker client.CommonAPIClient, nitrod protob.NitroClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add",
		Short:   "Adds a new database.",
		Example: addExampleTest,
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

			// get all of the containers as a list
			var engineOpts []string
			for _, c := range containers {
				engineOpts = append(engineOpts, strings.TrimLeft(c.Names[0], "/"))
			}

			// prompt the user for the engine to add the database
			selected, err := output.Select(os.Stdin, "Select the database engine: ", engineOpts)
			if err != nil {
				return err
			}

			// get the containers info
			info, err := docker.ContainerInspect(cmd.Context(), containers[selected].ID)
			if err != nil {
				return err
			}

			// ask the user for the database to create
			db, err := output.Ask("Enter the new database name", "", ":", &validate.DatabaseName{})
			if err != nil {
				return err
			}

			output.Pending("creating database", db)

			// wait for the api to be ready
			for {
				_, err := nitrod.Ping(cmd.Context(), &protob.PingRequest{})
				if err == nil {
					break
				}
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

			// create the database
			resp, err := nitrod.AddDatabase(cmd.Context(), &protob.AddDatabaseRequest{
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
