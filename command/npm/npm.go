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
	"github.com/craftcms/nitro/pkg/containerargs"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/contextor"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # run npm install in a sites container
  nitro npm sitename.nitro -- install

  # run build steps in a sites container
  nitro npm sitename.nitro -- run dev

  # run npm install in the current directory using a new container
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

			// if there was a container passed into the args
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
				cmds := []string{"exec", "-it", command.Container, "npm"}
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

				output.Info("npm", command.Args[0], "completed 🤘")

				return nil
			}

			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			// determine the command
			action := args

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
