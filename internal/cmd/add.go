package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/helpers"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/validate"
)

var addCommand = &cobra.Command{
	Use:   "add",
	Short: "Add site to machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)
		var wd string
		if len(args) > 0 {
			wd = args[0]
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			wd = cwd
		}

		// if the hostname flag is not set
		var hostname string
		if flagHostname != "" {
			hostname = flagHostname
		} else {
			pathName, err := helpers.PathName(wd)
			if err != nil {
				return err
			}

			hostnamePrompt := promptui.Prompt{
				Label:    fmt.Sprintf("what should the hostname be? [%s]", pathName),
				Validate: validate.Hostname,
			}

			hostnameEntered, err := hostnamePrompt.Run()
			if err != nil {
				return err
			}
			if hostnameEntered == "" {
				hostname = pathName
			}
		}

		// if the flag for webroot is not set, prompt
		var webroot string
		if flagWebroot != "" {
			webroot = flagWebroot
		} else {
			foundDir, err := helpers.FindWebRoot(wd)
			if err != nil {
				return err
			}
			webRootPrompt := promptui.Prompt{
				Label: fmt.Sprintf("where is the webroot? [%s]", foundDir),
			}

			webrootEntered, err := webRootPrompt.Run()
			if err != nil {
				return err
			}
			if webrootEntered == "" {
				webroot = foundDir
			}
		}

		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}

		mount := config.Mount{Source: wd}
		if err := configFile.AddMount(mount); err != nil {
			return err
		}

		site := config.Site{Hostname: hostname, Webroot: "/nitro/sites/" + wd + "/" + webroot}
		if err := configFile.AddSite(site); err != nil {
			return err
		}

		if !flagDebug {
			if err := configFile.Save(viper.ConfigFileUsed()); err != nil {
				return err
			}
		}

		fmt.Printf("%s has been added to nitro.yaml", hostname)

		applyPrompt := promptui.Prompt{
			Label: "apply nitro.yaml changes now? [y]",
		}

		apply, err := applyPrompt.Run()
		if err != nil {
			return err
		}
		if apply == "" {
			apply = "y"
		}

		if apply != "y" {
			fmt.Println("ok, you can apply new nitro.yaml changes later by running `nitro apply`.")

			return nil
		}

		var actions []nitro.Action
		// mount the directory
		m := configFile.Mounts[len(configFile.Mounts)-1]
		mountAction, err := nitro.MountDir(name, m.AbsSourcePath(), m.Dest)
		if err != nil {
			return err
		}
		actions = append(actions, *mountAction)

		// copy the nginx template
		copyTemplateAction, err := nitro.CopyNginxTemplate(name, site.Hostname)
		if err != nil {
			return err
		}
		actions = append(actions, *copyTemplateAction)

		// copy the nginx template
		changeNginxVariablesAction, err := nitro.ChangeTemplateVariables(name, site.Webroot, site.Hostname, configFile.PHP, []string{"test", "three", "four"})
		if err != nil {
			return err
		}
		actions = append(actions, *changeNginxVariablesAction...)

		if flagDebug {
			for _, action := range actions {
				fmt.Println(action.Args)
			}

			return nil
		}

		return errors.New("need to run the actions")
	},
}

func init() {
	addCommand.Flags().StringVar(&flagHostname, "hostname", "", "hostname of the site (e.g client.test)")
	addCommand.Flags().StringVar(&flagWebroot, "webroot", "", "webroot of the site (e.g. web)")
}
