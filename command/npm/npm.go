package npm

import (
	"bytes"
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
	// ErrNoPackageFile is returned when there is no package.json or package-lock.json file in a directory
	ErrNoPackageFile = fmt.Errorf("No package.json or package-lock.json was found")
)

const exampleText = `  # run npm install in a current directory
  nitro npm

  # run npm install in current directory
  nitro npm --install

  # updating a node project outside of the current directory
  nitro npm ./project-dir --version 10 --update`

// New is used for scaffolding new commands
func New(docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "npm",
		Short:   "Run npm install or update",
		Example: exampleText,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveFilterDirs
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			env := cmd.Flag("environment").Value.String()
			ctx := cmd.Context()
			if ctx == nil {
				// when we call commands from other commands (e.g. create)
				// the context could be nil, so we set it to the parent
				// context just in case.
				ctx = cmd.Parent().Context()
			}
			version := cmd.Flag("version").Value.String()

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
			nodePath := filepath.Join(path, "package.json")
			if action == "update" {
				nodePath = filepath.Join(path, "package-lock.json")
			}

			// set the container name to keep the ephemeral
			name := containerName(path, version, action)

			output.Pending("checking", nodePath)

			// make sure the file exists
			_, err := os.Stat(nodePath)
			if os.IsNotExist(err) {
				fmt.Println("")
				return ErrNoPackageFile
			}

			output.Done()

			image := fmt.Sprintf("docker.io/library/%s:%s", "node", version)

			imageFilter := filters.NewArgs()
			imageFilter.Add("reference", image)

			// look for the image
			images, err := docker.ImageList(ctx, types.ImageListOptions{Filters: imageFilter})
			if err != nil {
				return fmt.Errorf("unable to get a list of images, %w", err)
			}

			// if we don't have the image, pull it
			if len(images) == 0 {
				output.Pending("pulling")

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

			var commands []string
			switch action {
			case "install":
				commands = []string{"npm", "install"}
			default:
				commands = []string{"npm", "update"}
			}

			// set filters for the container name and environment
			containerFilter := filters.NewArgs()
			containerFilter.Add("label", labels.Environment+"="+env)
			containerFilter.Add("name", name)

			// check if there is an existing container
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{All: true, Filters: containerFilter})
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
							labels.Type:        "npm",
							// TODO abstract this?
							"com.craftcms.nitro.path": path,
						},
					},
					&container.HostConfig{
						Mounts: []mount.Mount{
							{
								Type:   "bind",
								Source: path,
								Target: "/app",
							},
						},
					},
					nil,
					name)
				if err != nil {
					return fmt.Errorf("unable to create container\n%w", err)
				}

				containerID = resp.ID
			}

			output.Info("Running npm", action)

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

			// copy the stream to stdout
			if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, stream.Reader); err != nil {
				return fmt.Errorf("unable to copy the output of the container logs, %w", err)
			}

			output.Info("npm", action, "complete 🤘")

			return nil
		},
	}

	// set flags for the command
	cmd.Flags().Bool("update", false, "run node update instead of install")
	cmd.Flags().String("version", "14", "which node version to use")

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
