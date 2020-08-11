package cmd

import (
	"fmt"

	"github.com/pixelandtonic/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/internal/config"
	"github.com/craftcms/nitro/internal/helpers"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/validate"
	"github.com/craftcms/nitro/internal/webroot"
)

var addCommand = &cobra.Command{
	Use:   "add",
	Short: "Add site",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveFilterDirs
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		p := prompt.NewPrompt()
		runner := nitro.NewMultipassRunner("multipass")

		// check if the machine exists
		if ip := nitro.IP(machine, runner); ip == "" {
			create, err := p.Confirm(fmt.Sprintf("Unable to find machine %q, want to create it", machine), &prompt.InputOptions{Default: "yes"})
			if err != nil {
				return err
			}
			if create {
				return initCommand.RunE(cmd, args)
			}
		}

		// load the config
		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}

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
			hostname, err = p.Ask("Enter the hostname", &prompt.InputOptions{
				Default:   directoryName,
				Validator: validate.Hostname,
			})
			if err != nil {
				return err
			}
		default:
			hostname = helpers.RemoveTrailingSlash(flagHostname)
		}

		// set the webrootName var (e.g. web)
		var webrootDir string
		switch flagWebroot {
		case "":
			// look for the www,public,public_html,www using the absolutePath variable
			foundDir, err := webroot.Find(absolutePath)
			if err != nil {
				fmt.Println("Unable to locate a webroot, setting to web.")
				foundDir = "web"
			}

			webrootDir, err = p.Ask("Enter the webroot", &prompt.InputOptions{
				Default:   foundDir,
				Validator: nil,
			})
			if err != nil {
				return err
			}
		default:
			webrootDir = flagWebroot
		}

		webRootPath := fmt.Sprintf("/home/ubuntu/sites/%s/%s", directoryName, webrootDir)
		// create a new mount
		skipMount := true
		mount := config.Mount{Source: absolutePath}
		exists, found := configFile.AlreadyMounted(mount)
		if exists {
			fmt.Println(mount.Source, "is already mounted at", found.Dest, ". Using existing instead of creating new mount.")

			webRootPath = webroot.ForExistingMount(found, absolutePath, webrootDir)

			fmt.Println("Setting webroot to", webRootPath)
		} else {
			mount.Dest = "/home/ubuntu/sites/" + directoryName
			// add the mount to configfile
			if err := configFile.AddMount(mount); err != nil {
				return err
			}
			skipMount = false
		}

		// create a new site
		// add site to config file
		skipSite := true
		site := config.Site{Hostname: hostname, Webroot: webRootPath}
		if configFile.SiteExists(site) {
			fmt.Println(site.Hostname, "has already been set.")
		} else {
			if err := configFile.AddSite(site); err != nil {
				return err
			}
			skipSite = false
		}

		if skipMount && skipSite {
			fmt.Println("There are no changes to apply, skipping...")
			return nil
		}

		if !flagDebug {
			if err := configFile.Save(viper.ConfigFileUsed()); err != nil {
				return err
			}
		}

		fmt.Printf("Added %s to config file\n", hostname)

		apply, err := p.Confirm("Apply changes from config", &prompt.InputOptions{
			Default:            "yes",
			AppendQuestionMark: true,
		})
		if err != nil {
			return err
		}

		if !apply {
			fmt.Println("You can apply new config file changes later by running `nitro apply`.")

			return nil
		}

		return applyCommand.RunE(cmd, args)
	},
}

func init() {
	addCommand.Flags().StringVar(&flagHostname, "hostname", "", "Hostname of the site (e.g client.test)")
	addCommand.Flags().StringVar(&flagWebroot, "webroot", "", "webroot of the site (e.g. web)")
	addCommand.Flags().BoolVar(&flagSkipHosts, "skip-hosts", false, "Skip editing the hosts file.")
}
