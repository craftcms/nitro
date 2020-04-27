package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/find"
	"github.com/craftcms/nitro/internal/helpers"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/sudo"
	"github.com/craftcms/nitro/internal/task"
	"github.com/craftcms/nitro/internal/webroot"
	"github.com/craftcms/nitro/validate"
)

var addCommand = &cobra.Command{
	Use:   "add",
	Short: "Add site to machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

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
				Label:    fmt.Sprintf("What should the hostname be? [%s]", directoryName),
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
		var webrootDir string
		switch flagWebroot {
		case "":
			foundDir, err := webroot.Find(absolutePath)
			if err != nil {
				return err
			}
			webRootPrompt := promptui.Prompt{
				Label: fmt.Sprintf("Where is the webroot? [%s]", foundDir),
			}

			webrootEntered, err := webRootPrompt.Run()
			if err != nil {
				return err
			}
			switch webrootEntered {
			case "":
				webrootDir = foundDir
			default:
				webrootDir = webrootEntered
			}
		default:
			webrootDir = flagWebroot
		}

		// create the vmWebRootPath (e.g. "/nitro/sites/"+ hostName + "/" | webrootName
		webRootPath := fmt.Sprintf("/nitro/sites/%s/%s", hostname, webrootDir)

		// load the config
		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}

		// create a new mount
		// add the mount to configfile
		mount := config.Mount{Source: absolutePath, Dest: "/nitro/sites/" + hostname}
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
			Label: "Apply changes now? [y]",
		}

		apply, err := applyPrompt.Run()
		if err != nil {
			return err
		}
		if apply == "" {
			apply = "y"
		}

		if apply != "y" {
			fmt.Println("You can apply new nitro.yaml changes later by running `nitro apply`.")

			return nil
		}

		path, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}

		// find the existingMounts

		c := exec.Command(path, []string{"info", machine, "--format=csv"}...)
		output, err := c.Output()
		if err != nil {
			return err
		}
		existingMounts, err := find.Mounts(machine, output)
		if err != nil {
			return err
		}

		// find the sites
		var sites []config.Site
		for _, site := range configFile.Sites {
			output, err := exec.Command(path, "exec", machine, "--", "sudo", "bash", "/opt/nitro/scripts/site-exists.sh", site.Hostname).Output()
			if err != nil {
				return err
			}
			if strings.Contains(string(output), "exists") {
				sites = append(sites, site)
			}
		}

		// check if a database already exists
		var databases []config.Database
		for _, db := range configFile.Databases {
			database, err := find.ExistingContainer(exec.Command(path, []string{"exec", machine, "--", "sudo", "bash", "/opt/nitro/scripts/docker-container-exists.sh", db.Name()}...), db)
			if err != nil {
				return err
			}

			if database != nil {
				fmt.Println("Database", db.Name(), "exists, skipping...")
				databases = append(databases, *database)
			}
		}

		php, err := find.PHPVersion(exec.Command(path, "exec", machine, "--", "php", "--version"))
		if err != nil {
			return err
		}

		actions, err := task.Apply(machine, configFile, existingMounts, sites, databases, php)

		if flagDebug {
			for _, action := range actions {
				fmt.Println(action.Args)
			}

			return nil
		}

		if err = nitro.Run(nitro.NewMultipassRunner("multipass"), actions); err != nil {
			return err
		}

		fmt.Println("Applied the changes and added", hostname, "to", machine)

		// prompt to add hosts file
		nitro, err := exec.LookPath("nitro")
		if err != nil {
			return err
		}

		fmt.Println("Adding", site.Hostname, "to your hosts file")

		return sudo.RunCommand(nitro, machine, "hosts")
	},
}

func init() {
	addCommand.Flags().StringVar(&flagHostname, "hostname", "", "Hostname of the site (e.g client.test)")
	addCommand.Flags().StringVar(&flagWebroot, "webroot", "", "webroot of the site (e.g. web)")
}
