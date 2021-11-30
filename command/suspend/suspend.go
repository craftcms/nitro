package suspend

import (
	"os"

	"github.com/craftcms/nitro/pkg/appaware"
	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/flags"
	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # suspend the app in the current directory
  nitro suspend

  # suspend a specific app using the global flag
  nitro --app myapp.nitro suspend`

// NewCommand returns the command to suspend an app from automatically starting.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "suspend",
		Short:   "Suspends an app.",
		Example: exampleText,
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return prompt.RunApply(cmd, args, false, output)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// get the app
			name := flags.AppName
			if name == "" {
				// get the current working directory
				wd, err := os.Getwd()
				if err != nil {
					return err
				}

				name, err = appaware.Detect(*cfg, wd)
				if err != nil {
					return err
				}
			}

			// suspend the app
			if err := cfg.SuspendApp(name); err != nil {
				return err
			}

			// save the config
			if err := cfg.Save(); err != nil {
				return err
			}

			output.Info("Suspended", name)

			return nil
		},
	}

	return cmd
}
