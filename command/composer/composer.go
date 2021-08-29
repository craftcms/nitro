package composer

import (
	"fmt"
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

			// if there is a container name
			if command.Container != "" {
				// make sure its a valid site
				if _, err := cfg.FindSiteByHostName(command.Container); err != nil {
					return err
				}
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

			// find the containers but limited to the site label
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter, All: true})
			if err != nil {
				return err
			}

			// are there any containers??
			if len(containers) == 0 {
				return fmt.Errorf("unable to find an matching site")
			}

			// make sure the container is started
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
		},
	}

	return cmd
}
