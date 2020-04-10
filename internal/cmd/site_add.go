package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/action"
	"github.com/craftcms/nitro/validate"
)

var siteAddCommand = &cobra.Command{
	Use:   "add",
	Short: "Add a site to machine",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)
		php := config.GetString("php", flagPhpVersion)
		localDirectory := args[0]
		domainName := args[1]
		fullDirectoryPath, err := filepath.Abs(localDirectory)
		if err != nil {
			return err
		}

		if err := validate.Path(fullDirectoryPath); err != nil {
			return err
		}
		if err := validate.Domain(domainName); err != nil {
			return err
		}

		// grab the config file and unmarshal
		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}

		site := config.Site{
			Domain:  domainName,
			Path:    fullDirectoryPath,
			Docroot: flagPublicDir,
		}

		var actions []action.Action

		mountAction, _ := action.Mount(name, localDirectory, domainName)
		actions = append(actions, *mountAction)

		createDirectoryAction, _ := action.CreateNginxSiteDirectory(name, domainName)
		actions = append(actions, *createDirectoryAction)

		copyTemplateAction, _ := action.CopyNginxTemplate(name, domainName)
		actions = append(actions, *copyTemplateAction)

		changeVarsActions, _ := action.ChangeTemplateVariables(name, domainName, flagPublicDir, php)
		for _, a := range *changeVarsActions {
			actions = append(actions, a)
		}

		createSymlinkAction, _ := action.CreateSiteSymllink(name, domainName)
		actions = append(actions, *createSymlinkAction)

		reloadNginxAction, _ := action.NginxReload(name)
		actions = append(actions, *reloadNginxAction)

		if err := configFile.AddSite(site); err != nil {
			return err
		}

		if flagDebug {
			fmt.Println("---- COMMANDS ----")
			for _, a := range actions {
				fmt.Println(a.Args)
			}

			fmt.Println("---- CONFIG FILE ----")

			configData, err := yaml.Marshal(configFile)
			if err != nil {
				return err
			}

			fmt.Println(string(configData))

			return nil
		}

		if err := configFile.Save(viper.ConfigFileUsed()); err != nil {
			return err
		}

		return action.Run(action.NewMultipassRunner("multipass"), actions)
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		if flagDebug {
			return nil
		}
		fmt.Println(
			fmt.Sprintf("added site %q to %q", args[1], config.GetString("name", flagMachineName)),
		)

		return nil
	},
}

func init() {
	siteAddCommand.Flags().StringVarP(&flagPhpVersion, "php-version", "p", "", "version of PHP to use")
	siteAddCommand.Flags().StringVarP(&flagPublicDir, "public-dir", "r", "web", "name of the public directory")
}
