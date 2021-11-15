package database

import (
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/craftcms/nitro/pkg/terminal"
)

func destroyCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroys a database engine.",
		Example: `  # remove a database engine from the config
  nitro db destroy`,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			cfg, err := config.Load(home)
			if err != nil {
				return nil, cobra.ShellCompDirectiveDefault
			}

			var options []string
			for _, d := range cfg.Databases {
				h, _ := d.GetHostname()
				options = append(options, h)
			}

			return options, cobra.ShellCompDirectiveDefault
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

			// get all of the databases
			dbs := cfg.Databases

			// create the options for the sites
			var options []string
			for _, d := range dbs {
				h, _ := d.GetHostname()
				options = append(options, h)
			}

			// prompt for the database
			selected, err := output.Select(cmd.InOrStdin(), "Select database to destroy: ", options)
			if err != nil {
				return err
			}

			db := dbs[selected]
			hostname, _ := db.GetHostname()

			output.Info("Removing", hostname)

			// remove the engine
			if err := cfg.RemoveDatabase(db); err != nil {
				return err
			}

			// save the config
			if err := cfg.Save(); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
