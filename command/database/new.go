package database

import (
	"runtime"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/portavail"
	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/craftcms/nitro/pkg/terminal"
)

var newExampleTest = `  # add a new database engine
  nitro db new`

func newCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "new",
		Short:   "Add a database engine",
		Example: newExampleTest,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prompt.VerifyInit(cmd, args, home, output)
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return prompt.RunApply(cmd, args, false, output)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// define the options
			var options []string
			switch runtime.GOOS {
			case "arm64", "arm":
				options = []string{"mariadb", "postgres"}
			default:
				options = []string{"mariadb", "mysql", "postgres"}
			}

			// prompt for the engine
			selection, err := output.Select(cmd.InOrStdin(), "Which database engine should we use", options)
			if err != nil {
				return err
			}

			// get the engine
			engine := options[selection]

			// ask for the version
			version, err := output.Ask("Which version should we use", "", "?", nil)
			if err != nil {
				return err
			}

			// set the default port
			var defaultPort string
			switch engine {
			case "postgres":
				defaultPort = "5432"
			default:
				defaultPort = "3306"
			}

			// find the first available port
			p, err := portavail.FindNext("", defaultPort)
			if err != nil {
				return err
			}

			// confirm the port to use
			port, err := output.Ask("Which port should we use for "+engine, p, "?", nil)
			if err != nil {
				return err
			}

			// add the database to the config
			cfg.Databases = append(cfg.Databases, config.Database{
				Engine:  engine,
				Version: version,
				Port:    port,
			})

			// save the config
			if err := cfg.Save(); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
