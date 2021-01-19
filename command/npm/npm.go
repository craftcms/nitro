package npm

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	volumetypes "github.com/docker/docker/api/types/volume"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
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
		Short:   "Run npm commands",
		Example: exampleText,
		Hidden:  true,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				// when we call commands from other commands (e.g. create)
				// the context could be nil, so we set it to the parent
				// context just in case.
				ctx = cmd.Parent().Context()
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
			if _, err := os.Stat(nodePath); os.IsNotExist(err) {
				output.Warning()
				return ErrNoPackageFile
			}

			output.Done()

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
			filter.Add("label", labels.Type+"=npm")
			filter.Add("label", labels.Path+"="+path)

			// check if there is an existing volume
			volumes, err := docker.VolumeList(ctx, filter)
			if err != nil {
				return err
			}

			// set the volume name
			// TODO(jasonmccallister) remove this hardcoded version
			volumeName := name(path, "14")

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
						labels.Type: "npm",
						labels.Path: path,
					},
				})
				if err != nil {
					return fmt.Errorf("unable to create the volume, %w", err)
				}

				pathVolume = volume
			}

			commands := append([]string{"npm"}, args...)

			// create the container
			resp, err := docker.ContainerCreate(ctx,
				&container.Config{
					Image: image,
					Cmd:   commands,
					Tty:   false,
					Labels: map[string]string{
						labels.Type: "npm",
						labels.Path: path,
					},
					WorkingDir: "/home/node/app",
				},
				&container.HostConfig{
					Mounts: []mount.Mount{
						{
							Type:   mount.TypeVolume,
							Source: pathVolume.Name,
							// TODO(jasonmccallister) get the path where node sotres deps
							Target: "/root",
						},
						{
							Type:   "bind",
							Source: path,
							Target: "/home/node/app",
						},
					},
				},
				nil,
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

func name(path, version string) string {
	// combine the path and version
	n := fmt.Sprintf("%s_%s_%s", path, "npm", version)

	// make it lower case
	n = strings.ToLower(n)

	// replace path separators with underscores
	n = strings.Replace(n, string(os.PathSeparator), "_", -1)

	// remove : to prevent error on windows
	n = strings.Replace(n, ":", "_", -1)

	// remove the first underscore
	return strings.TrimLeft(n, "_")
}
