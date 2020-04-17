package cmd

import (
	"fmt"

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

		// if there is no arg, get the current working dir
		// else get the first arg
		// set the directoryName variable
		directoryName, absolutePath, err := helpers.GetDirectoryArg(args)
		if err != nil {
			return err
		}

		// prompt for the hostname if --hostname == ""
		// else get the name of the current directory (e.g. nitro)
		var hostname string
		switch flagHostname {
		case "":
			hostnamePrompt := promptui.Prompt{
				Label:    fmt.Sprintf("what should the hostname be? [%s]", directoryName),
				Validate: validate.Hostname,
			}

			hostnameEntered, err := hostnamePrompt.Run()
			if err != nil {
				return err
			}

			switch hostnameEntered {
			case "":
				hostname = directoryName
			default:
				hostname = hostnameEntered
			}
		default:
			hostname = helpers.RemoveTrailingSlash(flagHostname)
		}

		// look for the www,public,public_html,www using the absolutePath variable
		// set the webrootName var (e.g. web)
		var webroot string
		switch flagWebroot {
		case "":
			foundDir, err := helpers.FindWebRoot(absolutePath)
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
			switch webrootEntered {
			case "":
				webroot = foundDir
			default:
				webroot = webrootEntered
			}
		default:
			webroot = flagWebroot
		}

		// create the vmWebRootPath (e.g. "/nitro/sites/"+ directoryName + "/" | webrootName
		webRootPath := fmt.Sprintf("/nitro/sites/%s/%s", directoryName, webroot)

		// load the config
		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}

		// create a new mount
		// add the mount to configfile
		mount := config.Mount{Source: absolutePath, Dest: webRootPath}
		if err := configFile.AddMount(mount); err != nil {
			return err
		}

		// create a new site
		// add site to config file
		site := config.Site{Hostname: hostname, Webroot: webRootPath}
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
		changeNginxVariablesAction, err := nitro.ChangeTemplateVariables(name, site.Webroot, site.Hostname, configFile.PHP, site.Aliases)
		if err != nil {
			return err
		}
		actions = append(actions, *changeNginxVariablesAction...)

		createSymlinkAction, err := nitro.CreateSiteSymllink(name, site.Hostname)
		if err != nil {
			return err
		}
		actions = append(actions, *createSymlinkAction)

		restartNginxAction, err := nitro.NginxReload(name)
		if err != nil {
			return err
		}
		actions = append(actions, *restartNginxAction)

		if flagDebug {
			for _, action := range actions {
				fmt.Println(action.Args)
			}

			return nil
		}

		if err = nitro.Run(nitro.NewMultipassRunner("multipass"), actions); err != nil {
			return err
		}

		fmt.Println("ok, we applied the changes and added", hostname, "to", name)

		return nil
	},
}

func init() {
	addCommand.Flags().StringVar(&flagHostname, "hostname", "", "hostname of site (e.g client.test)")
	addCommand.Flags().StringVar(&flagWebroot, "webroot", "", "webroot of site (e.g. web)")
}
