package keys

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/keys"
	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	site *config.Site
)

const exampleText = `  # keys command
  nitro keys`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys",
		Short:   "Adds SSH keys to a site container.",
		Example: exampleText,
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			// is there a site as the first arg?
			if len(args) > 0 {
				site, err = cfg.FindSiteByHostName(args[0])
				if err != nil {
					return err
				}

				output.Info("Preparing key import to", site.Hostname)

				return nil
			}

			// get a context aware list of sites
			sites := cfg.ListOfSitesByDirectory(home, wd)

			// create the options for the sites
			var options []string
			for _, s := range sites {
				options = append(options, s.Hostname)
			}

			// prompt for the site to ssh into
			selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
			if err != nil {
				return err
			}

			site = &sites[selected]

			output.Info("Preparing key import to", site.Hostname)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			path := filepath.Join(home, ".ssh")
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return errors.New("unable to find directory " + path)
			}

			// find all of the keys
			keys, err := keys.Find(path)
			if err != nil {
				return err
			}

			// if there are no keys
			if len(keys) == 0 {
				fmt.Println("Unable to find keys to add")
				return nil
			}

			var opts []string
			for k := range keys {
				opts = append(opts, k)
			}

			_, err = output.Select(os.Stdin, "Which key should we import?", opts)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
