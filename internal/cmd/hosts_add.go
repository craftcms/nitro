package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/txn2/txeh"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var hostsAddCommand = &cobra.Command{
	Use:   "add",
	Short: "Add an entry to your hosts file",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)

		uid := os.Geteuid()
		if uid != 0 {
			return errors.New("you do not appear to be running this command as root, so we cannot modify your hosts file")
		}

		// get the requested machines ip
		ip := nitro.IP(name, nitro.NewMultipassRunner("multipass"))

		// get all of the sites from the config file
		if !viper.IsSet("sites") {
			return errors.New("unable to read sites from " + viper.ConfigFileUsed())
		}

		var sites []config.Site
		if err := viper.UnmarshalKey("sites", &sites); err != nil {
			return err
		}

		hosts, err := txeh.NewHostsDefault()
		if err != nil {
			return err
		}

		var domains []string
		for _, site := range sites {
			domains = append(domains, site.Hostname)
		}

		hosts.AddHosts(ip, domains)

		return hosts.Save()
	},
}
