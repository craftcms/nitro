package npm

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/pathexists"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/craftcms/nitro/pkg/volumename"
)

var (
	// ErrNoPackageFile is returned when there is no package.json or package-lock.json file in a directory
	ErrNoPackageFile = fmt.Errorf("no package.json or package-lock.json was found")
)

const exampleText = `  # run npm install in a current directory
  nitro npm install

  # run npm update
  nitro npm update

  # run a script
  nitro npm run dev`

// NewCommand is the command used to run npm commands in a container.
func NewCommand(docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "npm",
		Short:   "Runs an npm command.",
		Example: exampleText,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				// when we call commands from other commands (e.g. create)
				// the context could be nil, so we set it to the parent
				// context just in case.
				ctx = context.Background()
			}
			version := cmd.Flag("version").Value.String()

			var path string
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("unable to get the current directory, %w", err)
			}

			path, err = filepath.Abs(wd)
			if err != nil {
				return fmt.Errorf("unable to find the absolute path, %w", err)
			}

			// determine the command
			action := args[0]

			// get the full file path
			nodePath := filepath.Join(path, "package.json")

			output.Pending("checking", nodePath)

			// make sure the file exists
			// see if the file exists
			if exists := pathexists.IsFile(nodePath); !exists {
				output.Warning()
				return fmt.Errorf("unable to find file %s", nodePath)
			}

			output.Done()

			// find the network
			networkFilter := filters.NewArgs()
			networkFilter.Add("name", "nitro-network")

			// check if the network needs to be created
			networks, err := docker.NetworkList(ctx, types.NetworkListOptions{Filters: networkFilter})
			if err != nil {
				return fmt.Errorf("unable to list the docker networks, %w", err)
			}

			var networkID string
			for _, n := range networks {
				if n.Name == "nitro-network" || strings.TrimLeft(n.Name, "/") == "nitro-network" {
					networkID = n.ID
				}
			}

			image := fmt.Sprintf("docker.io/library/%s:%s-alpine", "node", version)

			filter := filters.NewArgs()
			filter.Add("reference", image)

			// look for the image
			images, err := docker.ImageList(ctx, types.ImageListOptions{Filters: filter})
			if err != nil {
				return fmt.Errorf("unable to get a list of images, %w", err)
			}

			// remove the image ref filter
			filter.Del("reference", image)

			// if we don't have the image, pull it
			if len(images) == 0 {
				output.Pending("pulling", image)

				rdr, err := docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
				if err != nil {
					return fmt.Errorf("unable to pull docker image, %w", err)
				}

				buf := &bytes.Buffer{}
				if _, err := buf.ReadFrom(rdr); err != nil {
					return fmt.Errorf("unable to read the output from pulling the image, %w", err)
				}

				output.Done()
			}

			// add filters for the volume
			filter.Add("label", containerlabels.Type+"=npm")
			filter.Add("label", containerlabels.Path+"="+path)

			// check if there is an existing volume
			volumes, err := docker.VolumeList(ctx, filter)
			if err != nil {
				return err
			}

			// set the volume name
			// TODO(jasonmccallister) remove this hardcoded version
			volumeName := volumename.FromPath(strings.Join([]string{path, "14"}, string(os.PathSeparator)))

			var pathVolume types.Volume
			switch len(volumes.Volumes) {
			case 1:
				pathVolume = *volumes.Volumes[0]
			case 0:
				// create the volume if it does not exist
				volume, err := docker.VolumeCreate(ctx, volumetypes.VolumeCreateBody{
					Driver: "local",
					Name:   volumeName,
					Labels: map[string]string{
						containerlabels.Type: "npm",
						containerlabels.Path: path,
					},
				})
				if err != nil {
					return fmt.Errorf("unable to create the volume, %w", err)
				}

				pathVolume = volume
			}

			commands := append([]string{"npm"}, args...)

			networkConfig := &network.NetworkingConfig{}
			if networkID != "" {
				networkConfig = &network.NetworkingConfig{
					EndpointsConfig: map[string]*network.EndpointSettings{
						"nitro-network": {
							NetworkID: networkID,
						},
					},
				}
			}

			// create the container
			resp, err := docker.ContainerCreate(ctx,
				&container.Config{
					Image: image,
					Cmd:   commands,
					Tty:   false,
					Labels: map[string]string{
						containerlabels.Nitro: "true",
						containerlabels.Type:  "npm",
						containerlabels.Path:  path,
					},
					WorkingDir: "/home/node/app",
				},

				&container.HostConfig{
					Mounts: []mount.Mount{
						{
							Type:   mount.TypeVolume,
							Source: pathVolume.Name,
							Target: "/root",
						},
						{
							Type:   "bind",
							Source: path,
							Target: "/home/node/app",
						},
					},
				},
				networkConfig,
				nil,
				"")
			if err != nil {
				return fmt.Errorf("unable to create container\n%w", err)
			}

			output.Info("Running npm", action)

			// attach to the container
			stream, err := docker.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{
				Stream: true,
				Stdout: true,
				Stderr: true,
				Logs:   true,
			})
			if err != nil {
				return fmt.Errorf("unable to attach to container, %w", err)
			}
			defer stream.Close()

			// run the container
			if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
				return fmt.Errorf("unable to start the container, %w", err)
			}

			// copy the stream to stdout
			if _, err := stdcopy.StdCopy(cmd.OutOrStdout(), cmd.ErrOrStderr(), stream.Reader); err != nil {
				return fmt.Errorf("unable to copy the output of the container logs, %w", err)
			}

			output.Info("npm", action, "complete 🤘")

			if err := docker.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{}); err != nil {
				return err
			}

			return nil
		},
	}

	// set flags for the command
	cmd.Flags().String("version", "14", "which node version to use")

	return cmd
}
