package composer

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/dockerexec"
	"github.com/craftcms/nitro/pkg/flags"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # run composer install in an app container
  nitro composer install

  # specify the app to run the command in
  nitro --app myapp.nitro composer install

  # run complex commands with flags (args will be ignored)
  nitro --app myapp.nitro composer --command 'install --ignore-platform-reqs'`

var (
	flagCommand string
)

// NewCommand returns a new command that runs composer install or update for a directory.
// This command allows users to skip installing composer on the host machine and will run
// all the commands in a docker container.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "composer",
		Short:   "Runs Composer in an app.",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// assign the flag for the app name to a local variable
			appName := flags.AppName

			// load the configuration
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			if appName != "" {
				// find the app by the hostname and only return an error if it's not found
				_, err := cfg.FindAppByHostname(appName)
				if err != nil {
					return err
				}
			}

			// did the user specify the flag for commands?
			switch flagCommand != "" {
			case true:
				// use the command flag as inputs and ignore the args
				if err := dockerexec.Connect(cmd.InOrStdin(), cmd.OutOrStdout(), "nitro", appName, fmt.Sprintf("composer %s", flagCommand)); err != nil {
					return err
				}
			default:
				return fmt.Errorf("not yet implemented")
			}

			output.Info("composer command", "completed 🤘")

			return nil
		},
	}

	cmd.Flags().StringVar(&flagCommand, "command", "", "the command to execute in the container")

	return cmd
}
