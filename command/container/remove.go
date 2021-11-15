package container

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/craftcms/nitro/pkg/terminal"
)

func removeCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Removes a custom container.",
		Example: `  # remove a custom container from the config
  nitro container remove`,
		Aliases: []string{"rm"},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return prompt.RunApply(cmd, args, false, output)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// if there are no container
			if len(cfg.Containers) == 0 {
				return fmt.Errorf("there are no custom containers to remove")
			}

			// get all of the sites
			var options []string
			for _, c := range cfg.Containers {
				// add the site to the list
				options = append(options, c.Name)
			}

			// prompt for the container to remove
			selected, err := output.Select(cmd.InOrStdin(), "Select the custom container to remove: ", options)
			if err != nil {
				return err
			}

			container, err := cfg.FindContainerByName(options[selected])
			if err != nil {
				return err
			}

			// remove the container
			if err := cfg.RemoveContainer(container); err != nil {
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
