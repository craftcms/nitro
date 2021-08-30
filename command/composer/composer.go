package composer

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/composer"
	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerargs"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/contextor"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/craftcms/nitro/pkg/volumename"
	volumetypes "github.com/docker/docker/api/types/volume"
)

var (
	// ErrNoComposerFile is returned when there is no composer.json file in a directory
	ErrNoComposerFile = fmt.Errorf("no composer.json or composer.lock was found")
)

const exampleText = `  # run composer install in a sites container
  nitro composer sitename.nitro -- install

  # require a plugin for a site
  nitro composer sitename.nitro -- require craftcms/contact-form

  # run composer install in the current directory using a new container
  nitro composer install

  # use composer (without local installation) to create a new project
  nitro composer create-project craftcms/craft my-project`

var showHelp bool

// NewCommand returns a new command that runs composer install or update for a directory.
// This command allows users to skip installing composer on the host machine and will run
// all the commands in a disposable docker container.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "composer",
		Short:              "Runs a Composer command.",
		Example:            exampleText,
		DisableFlagParsing: true,
		Args:               cobra.MinimumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			cfg, err := config.Load(home)
			if err != nil {
				return nil, cobra.ShellCompDirectiveDefault
			}

			var options []string
			for _, s := range cfg.Sites {
				options = append(options, s.Hostname)
			}

			return options, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// check if the first arg is the help command/flag
			if args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
				return cmd.Help()
			}

			// get the context
			ctx := contextor.New(cmd.Context())

			// load the configuration
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// parse the args for the container
			command, err := containerargs.Parse(args)
			if err != nil {
				return err
			}

			// if the container was defined
			if command.Container != "" {
				// make sure its a valid site
				if _, err := cfg.FindSiteByHostName(command.Container); err != nil {
					return err
				}

				// create a filter for the environment
				filter := filters.NewArgs()
				filter.Add("label", containerlabels.Nitro)
				filter.Add("label", containerlabels.Host+"="+command.Container)

				// find the containers but limited to the site label
				containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter, All: true})
				if err != nil {
					return err
				}

				// are there any containers??
				if len(containers) == 0 {
					return fmt.Errorf("unable to find an matching site")
				}

				// start the container if its not running
				if containers[0].State != "running" {
					if err := docker.ContainerStart(ctx, command.Container, types.ContainerStartOptions{}); err != nil {
						return err
					}
				}

				// create the command for running the craft console
				cmds := []string{"exec", "-it", command.Container, "composer"}
				cmds = append(cmds, command.Args...)

				// find the docker executable
				cli, err := exec.LookPath("docker")
				if err != nil {
					return err
				}

				// create the command
				c := exec.Command(cli, cmds...)
				c.Stdin = cmd.InOrStdin()
				c.Stderr = cmd.ErrOrStderr()
				c.Stdout = cmd.OutOrStdout()

				if err := c.Run(); err != nil {
					return err
				}

				output.Info("composer", command.Args[0], "completed 🤘")

				return nil
			}

			// fallback to the stand alone container
			image := "docker.io/composer"

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

			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("unable to get the current directory, %w", err)
			}

			path, err := filepath.Abs(wd)
			if err != nil {
				return fmt.Errorf("unable to find the absolute path, %w", err)
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
			volumeName := volumename.FromPath(path)

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

			output.Info("composer command completed 🤘")

			// remove the container
			if err := docker.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{}); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
