package composer

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/moby/moby/pkg/stdcopy"
	"github.com/spf13/cobra"
)

var (
	// ErrExample is used when we want to share an error
	ErrExample = fmt.Errorf("some example error")
)

const exampleText = `  # run composer install in a current directory
  nitro composer

  # updating a composer project outside of the current directory
  nitro composer ./project-dir --version 2 --update`

// New is used for scaffolding new commands
func New(docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "composer",
		Short:   "Run composer install or update",
		Example: exampleText,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveFilterDirs
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			version := cmd.Flag("version").Value.String()

			// get the path from args or current directory
			var path string
			switch len(args) {
			case 0:
				wd, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("unable to get the current directory, %w", err)
				}

				path, err = filepath.Abs(wd)
				if err != nil {
					return fmt.Errorf("unable to find the absolute path, %w", err)
				}
			default:
				var err error
				path, err = filepath.Abs(args[0])
				if err != nil {
					return fmt.Errorf("unable to find the absolute path, %w", err)
				}
			}

			// determine the default action
			action := "install"
			if cmd.Flag("update").Value.String() == "true" {
				action = "update"
			}

			// get the full file path
			var composerPath string
			switch action {
			case "install":
				composerPath = fmt.Sprintf("%s%c%s", path, os.PathSeparator, "composer.json")
			default:
				composerPath = fmt.Sprintf("%s%c%s", path, os.PathSeparator, "composer.lock")
			}

			// make sure the file exists
			output.Pending("checking", composerPath)
			_, err := os.Stat(composerPath)
			if os.IsNotExist(err) {
				return fmt.Errorf("unable to locate a composer file at %s", path)
			}
			output.Done()

			image := fmt.Sprintf("docker.io/library/%s:%s", "composer", version)

			// filter for the image ref
			filters := filters.NewArgs()
			filters.Add("reference", image)

			// look for the image
			images, err := docker.ImageList(cmd.Context(), types.ImageListOptions{
				Filters: filters,
			})
			if err != nil {
				return fmt.Errorf("unable to get a list of images, %w", err)
			}

			// if we don't have the image, pull it
			if len(images) == 0 {
				output.Pending("pulling image")

				rdr, err := docker.ImagePull(cmd.Context(), image, types.ImagePullOptions{All: false})
				if err != nil {
					return fmt.Errorf("unable to pull the docker image, %w", err)
				}

				buf := &bytes.Buffer{}
				if _, err := buf.ReadFrom(rdr); err != nil {
					return fmt.Errorf("unable to read the output from pulling the image, %w", err)
				}

				output.Done()
			}

			var commands []string
			switch action {
			case "install":
				commands = []string{"composer", "install", "--ignore-platform-reqs", "--prefer-dist"}
			default:
				commands = []string{"composer", "update", "--ignore-platform-reqs", "--prefer-dist"}
			}

			// create the temp container
			resp, err := docker.ContainerCreate(cmd.Context(),
				&container.Config{
					Image: image,
					Cmd:   commands,
					Tty:   false,
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
				"")
			if err != nil {
				return fmt.Errorf("unable to create the composer container\n%w", err)
			}

			// attach to the container
			stream, err := docker.ContainerAttach(cmd.Context(), resp.ID, types.ContainerAttachOptions{
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
			if err := docker.ContainerStart(cmd.Context(), resp.ID, types.ContainerStartOptions{}); err != nil {
				return fmt.Errorf("unable to start the container, %w", err)
			}

			// show the output to stdout and stderr
			if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, stream.Reader); err != nil {
				return fmt.Errorf("unable to copy the output of the container logs, %w", err)
			}

			output.Pending("cleaning up")

			// remove the container
			if err := docker.ContainerRemove(cmd.Context(), resp.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
				return fmt.Errorf("unable to remove the temporary container %q, %w", resp.ID, err)
			}

			output.Done()

			output.Info("Composer", action, "completed 🤘")

			return ErrExample
		},
	}

	// set flags for the command
	cmd.Flags().BoolP("update", "u", false, "run composer update instead of install")
	cmd.Flags().StringP("version", "v", "2", "which composer version to use")

	return cmd
}
