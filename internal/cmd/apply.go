package cmd

import (
	"bufio"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/jasonmccallister/hosts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/find"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/runas"
	"github.com/craftcms/nitro/internal/scripts"
	"github.com/craftcms/nitro/internal/task"
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

		// find the machines IP
		ip, err := find.IP(machine, output)
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

		var sites []config.Site

		// find sites that are enabled
		var confs []string
		if output, err := script.Run(false, `ls /etc/nginx/sites-enabled`); err == nil {
			sc := bufio.NewScanner(strings.NewReader(output))
			for sc.Scan() {
				if sc.Text() == "default" {
					continue
				}

				confs = append(confs, strings.TrimSpace(sc.Text()))
			}
		}

		// generate a list of sites that
		for _, conf := range confs {
			s := config.Site{}
			// get the webroot
			if output, err := script.Run(false, fmt.Sprintf(scripts.FmtNginxSiteWebroot, conf)); err == nil {
				sp := strings.Fields(output)
				if len(sp) >= 2 {
					s.Webroot = strings.TrimRight(sp[1], ";")
				}
			}

			// get the server_name
			if output, err := script.Run(false, fmt.Sprintf(`grep "server_name " /etc/nginx/sites-available/%s | while read -r line; do echo "$line"; done`, conf)); err == nil {
				sp := strings.Fields(output)
				if len(sp) >= 2 {
					s.Hostname = strings.TrimRight(sp[1], ";")
				}
			}

			// get the hostname
			if s.Webroot != "" && s.Hostname != "" {
				sites = append(sites, s)
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
			fmt.Println("Skipping editing the hosts file.")
			return nil
		}

		// find all records by IP
		var hostRecords []string
		if records, err := hosts.FindIP(ip); err == nil {
			for _, hr := range records {
				if hr.IP == ip {
					hostRecords = hr.Hosts
				}
			}
		}

		switch runtime.GOOS {
		case "windows":
			if len(hostRecords) > 0 {
				if err := hostsRemoveCommand.RunE(cmd, hostRecords); err != nil {
					return err
				}
			}

			return hostsCommand.RunE(cmd, args)
		default:
			if len(hostRecords) > 0 {
				if err := runas.Elevated(machine, append([]string{"hosts", "remove"}, hostRecords...)); err != nil {
					return err
				}
			}

			if err := runas.Elevated(machine, []string{"hosts"}); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	applyCommand.Flags().BoolVar(&flagSkipHosts, "skip-hosts", false, "Skip editing the hosts file.")
}
