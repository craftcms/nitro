package npm

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/contextor"
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
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "npm",
		Short:   "Runs an npm command.",
		Example: exampleText,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := contextor.New(cmd.Context())
			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			// determine the command
			action := args

			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

			// get a context aware list of sites
			sites := cfg.ListOfSitesByDirectory(home, wd)

			// create the options for the sites
			var options []string
			for _, s := range sites {
				options = append(options, s.Hostname)
			}

			var site config.Site
			switch len(sites) {
			case 0:
				// prompt for the site
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
				if err != nil {
					return err
				}

				// set the site we selected
				site = sites[selected]

				// add the label to get the site
				filter.Add("label", containerlabels.Host+"="+sites[selected].Hostname)
			case 1:
				output.Info("connecting to", sites[0].Hostname)

				// set the site we selected
				site = sites[0]

				// add the label to get the site
				filter.Add("label", containerlabels.Host+"="+sites[0].Hostname)
			default:
				// prompt for the site to ssh into
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
				if err != nil {
					return err
				}

				// set the site we selected
				site = sites[selected]

				// add the label to get the site
				filter.Add("label", containerlabels.Host+"="+sites[selected].Hostname)
			}

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
				for _, command := range cmd.Root().Commands() {
					if command.Use == "start" {
						if err := command.RunE(cmd, []string{}); err != nil {
							return err
						}
					}
				}
			}

			// create the command for running the craft console
			cmds := []string{"exec", "-it", containers[0].ID, "npm"}

			// get the container path
			path := site.GetContainerPath()
			if path != "" {
				cmds = append(cmds, fmt.Sprintf("%s/%s", path, action))
			} else {
				cmds = append(cmds, action...)
			}

			switch len(args) == 0 {
			case true:
				// no args were provided, use the help command
				cmds = append(cmds, "help")
			default:
				// append the provided args to the command
				cmds = append(cmds, args...)
			}

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

			output.Info(fmt.Sprintf("npm command %q complete 🤘", action))

			return nil
		},
	}

	return cmd
}
