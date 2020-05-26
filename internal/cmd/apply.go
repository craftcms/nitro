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
	"github.com/craftcms/nitro/internal/runas"
	"github.com/craftcms/nitro/internal/scripts"
	"github.com/craftcms/nitro/internal/sudo"

	"github.com/craftcms/nitro/internal/task"
	"github.com/craftcms/nitro/internal/webroot"
)

var applyCommand = &cobra.Command{
	Use:   "apply",
	Short: "Apply changes",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

		// always read the config file so its updated from any previous commands
		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		// load the config file
		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}

		// ABSTRACT
		multipass, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}

		c := exec.Command(multipass, []string{"info", machine, "--format=csv"}...)
		output, err := c.Output()
		if err != nil {
			return err
		}

		// find mounts that already exist
		mounts, err := find.Mounts(machine, output)
		if err != nil {
			return err
		}
		// END ABSTRACT

		script := scripts.New(multipass, machine)

		// find sites that are created
		var sites []config.Site
		for _, site := range configFile.Sites {
			shouldAppend := false
			// check if the nginx config exists
			existsOutput, err := script.Run(false, fmt.Sprintf(scripts.FmtNginxSiteAvailable, site.Hostname))
			if err != nil {
				return err
			}
			if strings.Contains(existsOutput, "exists") {
				shouldAppend = true
			}

			// see if the webroot is the same
			webrootOutput, err := script.Run(false, fmt.Sprintf(scripts.FmtNginxSiteWebroot, site.Hostname))
			if err != nil {
				return err
			}

			matches, found := webroot.Matches(webrootOutput, site.Webroot)
			switch matches {
			case true:
				fmt.Println(fmt.Sprintf("Webroot for %q matches", site.Hostname))
			default:
				fmt.Println(fmt.Sprintf("Webroot for %q does not match, found %q", site.Hostname, found))
				site.Webroot = found
			}

			if shouldAppend {
				sites = append(sites, site)
			}
		}

		// find all existing databases
		databases, err := find.AllDatabases(exec.Command(multipass, []string{"exec", machine, "--", "docker", "container", "ls", "--format", `'{{ .Names }}'`}...))
		if err != nil {
			return err
		}

		// find the current version of php installed
		php, err := find.PHPVersion(exec.Command(multipass, "exec", machine, "--", "php", "--version"))
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

		if flagSkipHosts || len(configFile.Sites) == 0 {
			fmt.Println("Skipping editing the hosts file")
			return nil
		}

		fmt.Println("Editing your hosts file")

		return runas.Elevated(machine, []string{"hosts"})
	},
}

func init() {
	applyCommand.Flags().BoolVar(&flagSkipHosts, "skip-hosts", false, "skip editing the hosts file")
}
