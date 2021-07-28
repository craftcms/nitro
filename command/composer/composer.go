package composer

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/composer"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/contextor"
	"github.com/craftcms/nitro/pkg/pathexists"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/craftcms/nitro/pkg/volumename"
)

var (
	// ErrNoComposerFile is returned when there is no composer.json file in a directory
	ErrNoComposerFile = fmt.Errorf("no composer.json or composer.lock was found")
)

const exampleText = `  # run composer install in a current directory using a container
  nitro composer install

  # use composer (without local installation) to create a new project
  nitro composer create-project craftcms/craft my-project`

// NewCommand returns a new command that runs composer install or update for a directory.
// This command allows users to skip installing composer on the host machine and will run
// all the commands in a disposable docker container.
func NewCommand(docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "composer",
		Short:              "Runs a Composer command.",
		Example:            exampleText,
		DisableFlagParsing: true,
		Args:               cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := contextor.New(cmd.Context())
			var version string
			version, args = versionFromArgs(args)

			// get the path from args or current directory
			var path string
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("unable to get the current directory, %w", err)
			}

			path, err = filepath.Abs(wd)
			if err != nil {
				return fmt.Errorf("unable to find the absolute path, %w", err)
			}

			// determine the default action
			action := args[0]
			// if this is not a create project request, check for a composer.json
			if action != "create-project" {
				// get the full file path
				composerPath := filepath.Join(path, "composer.json")

				output.Pending("checking", composerPath)

				// see if the file exists
				if exists := pathexists.IsFile(composerPath); !exists {
					output.Warning()
					return fmt.Errorf("unable to find file %s", composerPath)
				}

				output.Done()
			}

			image := fmt.Sprintf("docker.io/craftcms/%s:%s-dev", "cli", version)

			// filter for the image ref
			filter := filters.NewArgs()
			filter.Add("reference", image)

			// look for the image
			images, err := docker.ImageList(ctx, types.ImageListOptions{Filters: filter})
			if err != nil {
				return fmt.Errorf("unable to get a list of images, %w", err)
			}

			// if we don't have the image, pull it
			if len(images) == 0 {
				rdr, err := docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
				if err != nil {
					return fmt.Errorf("unable to pull the docker image, %w", err)
				}

				buf := &bytes.Buffer{}
				if _, err := buf.ReadFrom(rdr); err != nil {
					return fmt.Errorf("unable to read the output from pulling the image, %w", err)
				}
			}

			// remove the image ref filter
			filter.Del("reference", image)

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

			// add filters for the volume
			filter.Add("label", containerlabels.Type+"=composer")
			filter.Add("label", containerlabels.Path+"="+path)

			// check if there is an existing volume
			volumes, err := docker.VolumeList(ctx, filter)
			if err != nil {
				return err
			}

			// set the volume name
			volumeName := volumename.FromPath(strings.Join([]string{path, version}, string(os.PathSeparator)))

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
						containerlabels.Type: "composer",
						containerlabels.Path: path,
					},
				})
				if err != nil {
					return fmt.Errorf("unable to create the volume, %w", err)
				}

				pathVolume = volume
			}

			// build the container options
			opts := &composer.Options{
				Image:    image,
				Commands: args,
				Labels: map[string]string{
					containerlabels.Nitro: "true",
					containerlabels.Type:  "composer",
					containerlabels.Path:  path,
				},
				Volume: &pathVolume,
				Path:   path,
				NetworkConfig: &network.NetworkingConfig{
					EndpointsConfig: map[string]*network.EndpointSettings{
						"nitro-network": {
							NetworkID: networkID,
						},
					},
				},
			}

			// create the container
			container, err := composer.CreateContainer(ctx, docker, opts)
			if err != nil {
				return fmt.Errorf("unable to create the composer container\n%w", err)
			}

			// attach to the container
			stream, err := docker.ContainerAttach(ctx, container.ID, types.ContainerAttachOptions{
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
			if err := docker.ContainerStart(ctx, container.ID, types.ContainerStartOptions{}); err != nil {
				return fmt.Errorf("unable to start the container, %w", err)
			}

			// show the output to stdout and stderr
			if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, stream.Reader); err != nil {
				return fmt.Errorf("unable to copy the output of the container logs, %w", err)
			}

			output.Info("composer", action, "completed 🤘")

			// remove the container
			if err := docker.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{}); err != nil {
				return err
			}

			return nil
		},
	}

	// set flags for the command
	cmd.Flags().String("php-version", "7.4", "which php version to use")

	return cmd
}

func versionFromArgs(args []string) (string, []string) {
	var version string
	var newArgs []string
	for i, a := range args {
		// get the version if using =
		if strings.Contains(a, "--php-version=") {
			parts := strings.Split(a, "=")
			version = parts[len(parts)-1]
			continue
		}

		// get the version if using a space
		if a == "--php-version" {
			version = args[i+1]
			continue
		}

		// append the new args
		newArgs = append(newArgs, a)
	}

	// if the version is not set, use the default
	if version == "" {
		version = "7.4"
	}

	return version, newArgs
}
