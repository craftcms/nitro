package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-input"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/prompt"
	"github.com/craftcms/nitro/validate"
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

		// ask to select a existingSite
		var existingSite config.Site
		_, i, err := prompt.Select(ui, "Select a site to rename:", "1", configFile.SitesAsList())
		if err != nil {
			return err
		}

		existingSite = sites[i]

		// ask for the new newHostname
		var newHostname string
		hostnamePrompt := promptui.Prompt{
			Label:    fmt.Sprintf("What should the new newHostname be? [current: %s]", existingSite.Hostname),
			Validate: validate.Hostname,
		}

		hostnameEntered, err := hostnamePrompt.Run()
		if err != nil {
			return err
		}

		switch hostnameEntered {
		case "":
			newHostname = existingSite.Hostname
		default:
			newHostname = hostnameEntered
		}

		if existingSite.Hostname == newHostname {
			return errors.New("the new and original hostnames match, nothing to do")
		}

		// update the config
		if err := configFile.RenameSite(existingSite, newHostname); err != nil {
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
