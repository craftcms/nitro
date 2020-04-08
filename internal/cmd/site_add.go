package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/action"
	"github.com/craftcms/nitro/internal/nitro"
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

		if err := validate.Path(localDirectory); err != nil {
			return err
		}
		if err := validate.Domain(domainName); err != nil {
			return err
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

		if flagDebug {
			for _, a := range actions {
				fmt.Println(a.Args)
			}

			return nil
		}

		return nitro.RunAction(nitro.NewMultipassRunner("multipass"), actions)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println(
			fmt.Sprintf("added site %q to %q", args[1], config.GetString("name", flagMachineName)),
		)
	},
}

func init() {
	siteAddCommand.Flags().StringVarP(&flagPhpVersion, "php-version", "p", "", "version of PHP to use")
	siteAddCommand.Flags().StringVarP(&flagPublicDir, "public-dir", "r", "web", "name of the public directory")
}
