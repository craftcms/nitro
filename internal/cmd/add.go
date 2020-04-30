package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-input"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/helpers"
	"github.com/craftcms/nitro/internal/prompt"
	"github.com/craftcms/nitro/internal/webroot"
)

var addCommand = &cobra.Command{
	Use:   "add",
	Short: "Add a site to a machine",
	RunE: func(cmd *cobra.Command, args []string) error {
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

		ui := &input.UI{
			Writer: os.Stdout,
			Reader: os.Stdin,
		}

		// prompt for the hostname if --hostname == ""
		// else get the name of the current directory (e.g. nitro)
		var hostname string
		switch flagHostname {
		case "":
			hostname, err = prompt.Ask(ui, "What should the hostname be?", directoryName, true)
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
				return err
			}

			webrootDir, err = prompt.Ask(ui, "Where is the webroot?", foundDir, true)
			if err != nil {
				return err
			}
		default:
			webrootDir = flagWebroot
		}

		// create the vmWebRootPath (e.g. "/nitro/sites/"+ hostName + "/" | webrootName
		webRootPath := fmt.Sprintf("/nitro/sites/%s/%s", hostname, webrootDir)

		// create a new mount
		skipMount := true
		mount := config.Mount{Source: absolutePath, Dest: "/nitro/sites/" + hostname}
		if configFile.MountExists(mount.Dest) {
			fmt.Println(mount.Source, "is already mounted at", mount.Dest, ". Using that instead of creating a new mount.")
		} else {
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

		fmt.Printf("%s has been added to nitro.yaml\n", hostname)

		applyChanges, err := prompt.Verify(ui, "Apply changes from config now?", "y")
		if err != nil {
			return err
		}

		if !applyChanges {
			fmt.Println("You can apply new config file changes later by running `nitro apply`.")

			return nil
		}

		return applyCommand.RunE(cmd, args)
	},
}

func init() {
	addCommand.Flags().StringVar(&flagHostname, "hostname", "", "Hostname of the site (e.g client.test)")
	addCommand.Flags().StringVar(&flagWebroot, "webroot", "", "webroot of the site (e.g. web)")
}
