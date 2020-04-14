package cmd

import (
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/action"
)

var siteRemoveCommand = &cobra.Command{
	Use:   "remove",
	Short: "Remove a site from a machine",
	Long: `Removing a site will perform the following steps:

1. Remove the virtual host from NGINX sites enabled
2. Delete the directory in "app/sites/xmydomain.test"
3. Unmount the local directory from the machine
4. Restart the NGINX web server
5. Remove the site from your nitro.yaml sites configuration
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)
		// make sure we have sites in the configFile
		if !viper.IsSet("sites") {
			return errors.New("no sites found in " + viper.ConfigFileUsed())
		}

		// create the prompt
		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}

		var sites []string
		for _, site := range configFile.Sites {
			sites = append(sites, site.Domain)
		}

		prompt := promptui.Select{
			Label: "Select site to remove",
			Items: sites,
		}
		_, site, err := prompt.Run()
		if err != nil {
			return err
		}

		var actions []action.Action

		// remove symlink
		removeSymlinkAction, err := action.RemoveSymlink(name, site)
		if err != nil {
			return err
		}
		actions = append(actions, *removeSymlinkAction)

		// remove mount
		unmountAction, err := action.Unmount(name, site)
		if err != nil {
			return err
		}
		actions = append(actions, *unmountAction)

		// remove the directory
		removeNginxSiteDirectoryAction, err := action.RemoveNginxSiteDirectory(name, site)
		if err != nil {
			return err
		}
		actions = append(actions, *removeNginxSiteDirectoryAction)

		// restart nginx
		restartNginxAction, err := action.NginxReload(name)
		if err != nil {
			return err
		}
		actions = append(actions, *restartNginxAction)

		// remove from configFile file
		if err := configFile.RemoveSite(site); err != nil {
			return err
		}

		if err := action.Run(action.NewMultipassRunner("multipass"), actions); err != nil {
			return nil
		}

		if err := configFile.Save(viper.ConfigFileUsed()); err != nil {
			return err
		}

		fmt.Printf("removed %q from %q\n", site, name)

		return nil
	},
}
