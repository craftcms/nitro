package container

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/portavail"
	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/craftcms/nitro/pkg/validate"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

func newCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "Add a new custom container",
		Example: `  # add a new custom container
  nitro container new

  # expand the number of images from the search
  nitro container new --limit 50`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// verify the config exists
			_, err := config.Load(home)
			if errors.Is(err, config.ErrNoConfigFile) {
				output.Info("Warning:", err.Error())

				// ask if the init command should run
				init, err := output.Confirm("Run `nitro init` now to create the config", true, "?")
				if err != nil {
					return err
				}

				// if init is false return nil
				if !init {
					return fmt.Errorf("You must run `nitro init` in order to add a site...")
				}

				// run the init command
				for _, c := range cmd.Parent().Commands() {
					// set the init command
					if c.Use == "init" {
						if err := c.RunE(c, args); err != nil {
							return err
						}
					}
				}
			}

			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return prompt.RunApply(cmd, args, output)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse limit flag
			limitFlag := cmd.Flag("limit").Value.String()
			limit, err := strconv.Atoi(limitFlag)
			if err != nil {
				limit = 10
			}

			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// ask for the image
			resp, err := output.Ask("What image are you trying to add", "", "?", &validate.HostnameValidator{})
			if err != nil {
				return err
			}

			// search for an image with the name
			images, err := docker.ImageSearch(cmd.Context(), resp, types.ImageSearchOptions{Limit: limit})
			if err != nil {
				return err
			}

			// if there are no images return
			if len(images) == 0 {
				return fmt.Errorf("no images found matching %q", resp)
			}

			// show the found images as a selection
			options := []string{}
			for _, i := range images {
				options = append(options, i.Name)
			}

			// prompt for the image we found
			selection, err := output.Select(cmd.InOrStdin(), "Which image should we use?", options)
			if err != nil {
				return err
			}

			image := images[selection].Name

			// ask for the tag (default to latest)
			tag, err := output.Ask("What tag should we use", "latest", "?", nil)
			if err != nil {
				return err
			}

			// generate the image ref
			var ref string
			switch strings.Contains(image, "/") {
			case true:
				ref = fmt.Sprintf("%s:%s", image, tag)
			default:
				ref = fmt.Sprintf("docker.io/library/%s:%s", image, tag)
			}

			output.Pending("downloading", ref)

			// pull the image
			rc, err := docker.ImagePull(cmd.Context(), ref, types.ImagePullOptions{All: false})
			if err != nil {
				output.Warning()

				return err
			}
			defer rc.Close()

			buf := &bytes.Buffer{}
			if _, err := buf.ReadFrom(rc); err != nil {
				output.Warning()
				return fmt.Errorf("unable to read the output from pulling the image, %w", err)
			}
			output.Done()

			// inspect the recently pulled image
			imageSpecs, _, err := docker.ImageInspectWithRaw(cmd.Context(), ref)
			if err != nil {
				return err
			}

			// get the ports for the image
			ports := []string{}
			for port := range imageSpecs.ContainerConfig.ExposedPorts {
				// find the first available port
				p, err := portavail.FindNext("", port.Port())
				if err != nil {
					return err
				}

				// should we prompt for the port to be exposed?
				add, err := output.Confirm(fmt.Sprintf("Expose port `%s` on host", p), true, "?")
				if err != nil {
					return err
				}

				if add {
					ports = append(ports, fmt.Sprintf("%s:%s", p, port.Port()))
				}
			}

			exposesUI, err := output.Confirm("Does the image contain a web based UI", true, "?")
			if err != nil {
				return err
			}

			var uiPort int
			if exposesUI {
				// format the ports to grab only the right side (container port)
				opts := []string{}
				for _, p := range ports {
					p := strings.Split(p, ":")
					opts = append(opts, p[len(p)-1])
				}

				// prompt the user for the port
				selected, err := output.Select(cmd.InOrStdin(), "Which port should we use for the UI?", opts)
				if err != nil {
					return err
				}

				clean := strings.Split(ports[selected], ":")

				// get the container port not the host port
				p, err := strconv.Atoi(clean[1])
				if err != nil {
					return err
				}

				// assign the port
				uiPort = p
			}

			// inspect the images volumes
			var volumes []string
			for v := range imageSpecs.ContainerConfig.Volumes {
				// should we create a volume for the volume?
				add, err := output.Confirm(fmt.Sprintf("Create volume `%q` for container", v), true, "?")
				if err != nil {
					return err
				}

				if add {
					volumes = append(volumes, v)
				}
			}

			// get a suggest image name
			suggested := image
			if strings.Contains(image, "/") {
				suggested = strings.Split(image, "/")[0]
			}

			// prompt for the container name
			name, err := output.Ask("What is the name of the container", suggested, "?", &validate.HostnameValidator{})
			if err != nil {
				return err
			}

			// create the container and save it
			container := config.Container{
				Name:    name,
				Image:   image,
				Tag:     tag,
				Ports:   ports,
				WebGui:  uiPort,
				Volumes: volumes,
			}

			// setup a custom env file?
			createEnvfile, err := output.Confirm("Create a file to add environment variables", true, "?")
			if err != nil {
				return err
			}

			var envFile string
			if createEnvfile {
				// create the file
				file := filepath.Join(home, config.DirectoryName, "."+name)
				if _, err := os.Create(file); err != nil {
					output.Warning()

					return fmt.Errorf("unable to create environment file: %w", err)
				}

				_, envFile = filepath.Split(file)

				output.Info(fmt.Sprintf("Created environment variables file at %q...", file))
			}

			container.EnvFile = envFile

			// add the container to the config
			if err := cfg.AddContainer(container); err != nil {
				return err
			}

			// save the config
			if err := cfg.Save(); err != nil {
				return err
			}

			output.Info(fmt.Sprintf("New container %q added! 🐳", name+".containers.nitro"))

			return nil
		},
	}

	cmd.Flags().Int("limit", 10, "number of images to return from the registry")

	return cmd
}
