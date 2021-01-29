package ssh

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # ssh into a container - assuming its the current working directory
  nitro ssh

  # ssh into the container as root - changes may not persist after "nitro apply"
  nitro ssh --root`

// NewCommand returns the ssh command to get a shell in a container. The command is context aware and if
// it is not in a known project directory, it will provide a list of known sites to the user.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ssh",
		Short:   "SSH into a container",
		Example: exampleText,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", labels.Nitro)

			// get a context aware list of sites
			sites := cfg.ListOfSitesByDirectory(home, wd)

			// create the options for the sites
			var options []string
			for _, s := range sites {
				options = append(options, s.Hostname)
			}

			// if there are found sites we want to show or connect to the first one, otherwise prompt for
			// which site to connect to.
			switch len(sites) {
			case 0:
				// prompt for the site to ssh into
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
				if err != nil {
					return err
				}

				// add the label to get the site
				filter.Add("label", labels.Host+"="+sites[selected].Hostname)
			case 1:
				output.Info("connecting to", sites[0].Hostname)

				// add the label to get the site
				filter.Add("label", labels.Host+"="+sites[0].Hostname)
			default:
				// prompt for the site to ssh into
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
				if err != nil {
					return err
				}

				// add the label to get the site
				filter.Add("label", labels.Host+"="+sites[selected].Hostname)
			}

			// find the containers but limited to the site label
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter, All: true})
			if err != nil {
				return err
			}

			// are there any containers??
			if len(containers) == 0 {
				return fmt.Errorf("unable to find an matching site")
			}

			// start the container if its not running
			if containers[0].State != "running" {
				if err := docker.ContainerStart(cmd.Context(), containers[0].ID, types.ContainerStartOptions{}); err != nil {
					return err
				}
			}

			return containerConnect(cmd, output, containers[0].ID)
		},
	}

	cmd.Flags().Bool("root", false, "ssh as the root user")
	// cmd.Flags().Bool("proxy", false, "connect to the proxy container")

	return cmd
}

func containerConnect(cmd *cobra.Command, output terminal.Outputer, containerID string) error {
	// find the docker executable
	cli, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	// check if the root user should be used
	user := "www-data"
	if cmd.Flag("root").Value.String() == "true" {
		user = "root"
	}

	// show a notice about changes
	if user == "root" {
		output.Info("using root… system changes are ephemeral…")
	}

	c := exec.Command(cli, "exec", "-u", user, "-it", containerID, "sh")

	c.Stdin = cmd.InOrStdin()
	c.Stderr = cmd.ErrOrStderr()
	c.Stdout = cmd.OutOrStdout()

	return c.Run()
}
