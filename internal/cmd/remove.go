package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/prompt"
)

var removeCommand = &cobra.Command{
	Use:   "remove",
	Short: "Manage your nitro sites",
	RunE: func(cmd *cobra.Command, args []string) error {
		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}
		name := configFile.Name

		sites := configFile.GetSites()
		i, _ := prompt.Select("Select site to remove", configFile.SitesAsList())

		site := sites[i]

		if !prompt.Verify(fmt.Sprintf("this will remove %s from %s, do you want to apply the changes? [y]", site.Hostname, name)) {
			fmt.Println("ok, you can apply new nitro.yaml changes later by running `nitro apply`.")

			return nil
		}

		// remove site
		if err := configFile.RemoveSite(site.Hostname); err != nil {
			return err
		}

		// save the config
		if err := configFile.Save(viper.ConfigFileUsed()); err != nil {
			return err
		}

		var actions []nitro.Action

		// unmount
		unmountAction, err := nitro.UnmountDir(name, site.Webroot)
		if err != nil {
			return err
		}
		actions = append(actions, *unmountAction)

		// remove nginx symlink
		removeSymlinkAction, err := nitro.RemoveSymlink(name, site.Hostname)
		if err != nil {
			return err
		}
		actions = append(actions, *removeSymlinkAction)

		restartNginxAction, err := nitro.NginxReload(name)
		if err != nil {
			return err
		}
		actions = append(actions, *restartNginxAction)

		for _, a := range actions {
			fmt.Println(a.Args)
		}

		return nil
	},
}
