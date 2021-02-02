package database

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/craftcms/nitro/pkg/database"
	"github.com/craftcms/nitro/pkg/filetype"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/pathexists"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/craftcms/nitro/pkg/validate"
	"github.com/craftcms/nitro/protob"
)

var importExampleText = `  # import a sql file into a database
  nitro db import filename.sql

  # use a relative path
  nitro db import ~/Desktop/backup.sql

  # use an absolute path
  nitro db import /Users/oli/Desktop/backup.sql`

// importCommand is the command for creating new development environments
func importCommand(home string, docker client.CommonAPIClient, nitrod protob.NitroClient, output terminal.Outputer) *cobra.Command {
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
			return []string{"sql", "gz", "zip"}, cobra.ShellCompDirectiveFilterFileExt
		},
		Example: importExampleText,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// make sure the file exists
			if exists := pathexists.IsFile(args[0]); !exists {
				output.Info(cmd.UsageString())

				return fmt.Errorf("unable to find file %s", args[0])
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// replace the relative path with the full directory
			path := args[0]
			if strings.HasPrefix(path, "~") {
				path = strings.Replace(path, "~", home, 1)
			}

			// check if this is a zip file
			var compressed bool
			kind, err := filetype.Determine(path)
			if err != nil {
				return err
			}

			switch kind {
			case "zip", "tar":
				compressed = true
			}

			// detect the type of backup if not compressed
			detected := ""
			if !compressed {
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
			db, err := output.Ask("Enter the database name", "", ":", &validate.DatabaseName{})
			if err != nil {
				return err
			}

			output.Info("Preparing import…")

			stream, err := nitrod.ImportDatabase(cmd.Context())
			if err != nil {
				return err
			}

			// get the containers info
			info, err := docker.ContainerInspect(cmd.Context(), containers[selected].ID)
			if err != nil {
				return err
			}

			// get the database compatability from the container labelsmake l
			detected = info.Config.Labels[labels.DatabaseCompatibility]
			hostname := strings.TrimLeft(info.Name, "/")
			version := info.Config.Labels[labels.DatabaseVersion]

			var port string
			// get the port from the container info
			for p, bind := range info.HostConfig.PortBindings {
				for _, v := range bind {
					if v.HostPort != "" {
						port = p.Port()
					}
				}
			}

			// create a request with the database information to populate the database info for the import
			err = stream.Send(&protob.ImportDatabaseRequest{
				Payload: &protob.ImportDatabaseRequest_Database{
					Database: &protob.DatabaseInfo{
						Engine:     detected,
						Version:    version,
						Database:   db,
						Port:       port,
						Compressed: compressed,
						Hostname:   hostname,
					},
				},
			})
			// check if the error code is unimplemented
			if code := status.Code(err); code == codes.Unimplemented {
				output.Info("The API does not appear to be updated, please run `nitro update` to get the latest")

				return nil
			}
			if err != nil {
				return stream.RecvMsg(nil)
			}

			// create a timer
			start := time.Now()

			// open the file
			file, err := os.Open(path)
			if err != nil {
				return err
			}

			// create a buffer to handle large files more gracefully
			buffer := make([]byte, 1024*20)
			reader := bufio.NewReader(file)

			output.Pending(fmt.Sprintf("importing database %q into %q", db, hostname))

			// stream to backup file to the api
			for {
				n, err := reader.Read(buffer)
				if err == io.EOF {
					break
				}
				if err != nil {
					output.Warning()

					return stream.RecvMsg(nil)
				}

				// send the chunked file data in pieces
				if err := stream.Send(&protob.ImportDatabaseRequest{
					Payload: &protob.ImportDatabaseRequest_Data{
						Data: buffer[:n],
					},
				}); err != nil {
					output.Warning()

					return err
				}
			}

			// handle the response
			reply, err := stream.CloseAndRecv()
			if err != nil {
				output.Warning()

				return stream.RecvMsg(nil)
			}

			output.Done()

			output.Info(fmt.Sprintf("%s, took %.2f seconds 💪...", reply.Message, time.Since(start).Seconds()))

			return nil
		},
	}

	return cmd
}
