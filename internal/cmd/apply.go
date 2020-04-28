package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/find"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/sudo"
	"github.com/craftcms/nitro/internal/task"
)

var applyCommand = &cobra.Command{
	Use:   "apply",
	Short: "Apply changes from config",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

		// load the config file
		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}

		// ABSTRACT
		path, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}

		c := exec.Command(path, []string{"info", machine, "--format=csv"}...)
		output, err := c.Output()
		if err != nil {
			return err
		}

		mounts, err := find.Mounts(machine, output)
		if err != nil {
			return err
		}
		// END ABSTRACT

		// find sites not created
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

		actions, err := task.Apply(machine, configFile, mounts, sites, databases, php)
		if err != nil {
			return err
		}

		if flagDebug {
			for _, a := range actions {
				fmt.Println(a.Args)
			}

			return nil
		}

		if err := nitro.Run(nitro.NewMultipassRunner("multipass"), actions); err != nil {
			return err
		}

		fmt.Println("Applied changes from", viper.ConfigFileUsed())

		nitro, err := exec.LookPath("nitro")
		if err != nil {
			return err
		}

		fmt.Println("Editing your hosts file")

		// TODO check the current OS and call commands for windows
		return sudo.RunCommand(nitro, machine, "hosts")
	},
}
