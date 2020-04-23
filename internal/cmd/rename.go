package cmd

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/find"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/prompt"
)

var renameCommand = &cobra.Command{
	Use:   "rename",
	Short: "Rename a site",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}

		configSites := configFile.GetSites()

		if len(configSites) == 0 {
			return errors.New("there are no sites to rename")
		}

		i, _ := prompt.Select("Select site to rename", configFile.SitesAsList())

		siteToRename := configSites[i]

		path, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}

		_, err = find.SitesEnabled(
			exec.Command(path, []string{"exec", machine, "--", "find", "/etc/nginx/sites-enabled/", "-maxdepth", "1", "-type", "l"}...),
		)
		if err != nil {
			return err
		}

		var actions []nitro.Action

		// remove the symlink
		removeSymlinkAction, err := nitro.RemoveSymlink(machine, siteToRename.Hostname)
		if err != nil {
			return err
		}
		actions = append(actions, *removeSymlinkAction)

		copyTemplateAction, err := nitro.CopyNginxTemplate(machine, siteToRename.Hostname)
		if err != nil {
			return err
		}
		actions = append(actions, *copyTemplateAction)

		// TODO add PHP back to the config file and add to apply
		changeNginxVariablesAction, err := nitro.ChangeTemplateVariables(machine, siteToRename.Webroot, siteToRename.Hostname, "7.4", siteToRename.Aliases)
		if err != nil {
			return err
		}
		actions = append(actions, *changeNginxVariablesAction...)

		// restart nginx
		restartNginxAction, err := nitro.NginxReload(machine)
		if err != nil {
			return err
		}
		actions = append(actions, *restartNginxAction)

		for _, action := range actions {
			fmt.Println(action.Args)
		}

		return nil
	},
}
