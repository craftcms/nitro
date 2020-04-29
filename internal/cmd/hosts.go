package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/txn2/txeh"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/hosts"
	"github.com/craftcms/nitro/internal/nitro"
)

var hostsCommand = &cobra.Command{
	Use:    "hosts",
	Short:  "Add sites to your hosts file",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

		if !flagDebug {
			uid := os.Geteuid()
			if uid != 0 {
				return errors.New("you do not appear to be running this command as root, so we cannot modify your hosts file")
			}
		}

		// get the requested machines ip
		ip := nitro.IP(machine, nitro.NewMultipassRunner("multipass"))

		// get all of the sites from the config file
		var sites []config.Site
		if err := viper.UnmarshalKey("sites", &sites); err != nil {
			return err
		}

		if sites == nil {
			fmt.Println("There are no sites in the config file to remove")
			return nil
		}

		var domains []string
		for _, site := range sites {
			domains = append(domains, site.Hostname)
		}

		he, err := txeh.NewHostsDefault()
		if err != nil {
			return err
		}

		if flagDebug {
			for _, domain := range domains {
				fmt.Println("adding", domain, "to hosts file")
			}

			return nil
		}

		return hosts.Add(he, ip, domains)
	},
}

func init() {
	hostsCommand.AddCommand(hostsRemoveCommand)
}
