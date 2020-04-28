package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-input"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/prompt"
)

var renameCommand = &cobra.Command{
	Use:   "rename",
	Short: "Rename a site",
	RunE: func(cmd *cobra.Command, args []string) error {
		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}

		sites := configFile.GetSites()

		if len(sites) == 0 {
			return errors.New("there are no sites to rename")
		}

		ui := &input.UI{
			Writer: os.Stdout,
			Reader: os.Stdin,
		}

		// ask to select a site
		var site config.Site
		_, i, err := prompt.Select(ui, "Select a site to rename:", sites[0].Hostname, configFile.SitesAsList())
		if err != nil {
			return err
		}
		site = sites[i]

		// ask for the new newHostname
		var newHostname string
		newHostname, err = prompt.Ask(ui, "What should the new hostname be?", site.Hostname, true)
		if err != nil {
			return err
		}
		if site.Hostname == newHostname {
			return errors.New("the new and original hostnames match, nothing to do")
		}

		// update the config
		if err := configFile.RenameSite(site, newHostname); err != nil {
			return err
		}

		// save the file
		if !flagDebug {
			if err := configFile.Save(viper.ConfigFileUsed()); err != nil {
				return err
			}
		}

		applyChanges, err := prompt.Verify(ui, "Apply changes from config now?", "y")
		if err != nil {
			return err
		}

		if applyChanges {
			fmt.Println("Ok, applying changes from the config file...")
			return applyCommand.RunE(cmd, args)
		}

		return nil
	},
}
