package composer

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/craftcms/nitro/labels"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/cobra"
)

var (
	// ErrNoComposerFile is returned when there is no composer.json file in a directory
	ErrNoComposerFile = fmt.Errorf("No composer.json or composer.lock was found")
)

const exampleText = `  # run composer install in a current directory using a container
  nitro composer install

  # update a composer project using verison 1
  nitro composer update --version 1`

// NewCommand returns a new command that runs composer install or update for a directory.
// This command allows users to skip installing composer on the host machine and will run
// all the commands in a disposable docker container.
func NewCommand(docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:       "composer",
		Short:     "Run composer install or update",
		Example:   exampleText,
		ValidArgs: []string{"install", "update"},
		RunE: func(cmd *cobra.Command, args []string) error {
			env := cmd.Flag("environment").Value.String()
			version := cmd.Flag("version").Value.String()
			ctx := cmd.Context()
			if ctx == nil {
				// when we call commands from other commands (e.g. create)
				// the context could be nil, so we set it to the parent
				// context just in case.
				ctx = context.Background()
			}

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

			// get the full file path
			composerPath := filepath.Join(path, "composer.json")

			// set the container name to keep the ephemeral
			name := containerName(path, version, action)

			output.Pending("checking", composerPath)

			// make sure the file exists
			if _, err = os.Stat(composerPath); os.IsNotExist(err) {
				return ErrNoComposerFile
			}

			output.Done()

			image := fmt.Sprintf("docker.io/library/%s:%s", "composer", version)

			// filter for the image ref
			imageFilter := filters.NewArgs()
			imageFilter.Add("reference", image)

			// look for the image
			images, err := docker.ImageList(ctx, types.ImageListOptions{Filters: imageFilter})
			if err != nil {
				return fmt.Errorf("unable to get a list of images, %w", err)
			}

			// if we don't have the image, pull it
			if len(images) == 0 {
				output.Pending("pulling image")

				rdr, err := docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
				if err != nil {
					return fmt.Errorf("unable to pull the docker image, %w", err)
				}

				buf := &bytes.Buffer{}
				if _, err := buf.ReadFrom(rdr); err != nil {
					return fmt.Errorf("unable to read the output from pulling the image, %w", err)
				}

				output.Done()
			}

			// build the args
			commands := []string{"composer"}
			commands = append(commands, args...)

			// set filters for the container name and environment
			containerFilter := filters.NewArgs()
			containerFilter.Add("label", labels.Environment+"="+env)
			containerFilter.Add("name", name)

			// check if there is an existing container
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: containerFilter})
			if err != nil {
				return err
			}

			// check the length of the container
			var containerID string
			switch len(containers) {
			case 1:
				containerID = containers[0].ID
			default:
				// create the container
				resp, err := docker.ContainerCreate(ctx,
					&container.Config{
						Image: image,
						Cmd:   commands,
						Tty:   false,
						Labels: map[string]string{
							labels.Environment: env,
							labels.Type:        "composer",
							// TODO abstract this?
							"com.craftcms.nitro.path": path,
						},
					},
					&container.HostConfig{
						Mounts: []mount.Mount{
							{
								Type:   mount.TypeBind,
								Source: path,
								Target: "/app",
							},
						},
					},
					nil,
					name)
				if err != nil {
					return fmt.Errorf("unable to create the composer container\n%w", err)
				}

				containerID = resp.ID
			}

			// attach to the container
			stream, err := docker.ContainerAttach(ctx, containerID, types.ContainerAttachOptions{
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
			if err := docker.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
				return fmt.Errorf("unable to start the container, %w", err)
			}

			// show the output to stdout and stderr
			if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, stream.Reader); err != nil {
				return fmt.Errorf("unable to copy the output of the container logs, %w", err)
			}

			output.Info("composer", action, "completed 🤘")

			// should we remove the container
			if cmd.Flag("keep").Value.String() == "false" {
				if err := docker.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{}); err != nil {
					return err
				}
			}

			return nil
		},
	}

	// set flags for the command
	cmd.Flags().String("version", "2", "which composer version to use")
	cmd.Flags().Bool("keep", true, "keep the container (faster since it will cache dependencies)")

	return cmd
}

func containerName(path, version, action string) string {
	// combine the path and version
	n := fmt.Sprintf("%s_%s_%s_%s", path, "composer", version, action)

	// make it lower case
	n = strings.ToLower(n)

	// replace path separators with underscores
	n = strings.Replace(n, string(os.PathSeparator), "_", -1)

	// remove the first underscore
	return strings.TrimLeft(n, "_")
}
