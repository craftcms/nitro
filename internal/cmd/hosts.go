package cmd

import (
	"errors"
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
		machine := "nitro-dev"
		if flagMachineName != "" {
			machine = flagMachineName
		}

		uid := os.Geteuid()
		if uid != 0 {
			return errors.New("you do not appear to be running this command as root, so we cannot modify your hosts file")
		}

		// get the requested machines ip
		ip := nitro.IP(machine, nitro.NewMultipassRunner("multipass"))

		// get all of the sites from the config file
		if !viper.IsSet("sites") {
			return errors.New("unable to read sites from " + viper.ConfigFileUsed())
		}

		var sites []config.Site
		if err := viper.UnmarshalKey("sites", &sites); err != nil {
			return err
		}

		var domains []string
		for _, site := range sites {
			domains = append(domains, site.Hostname)
		}

		he, err := txeh.NewHostsDefault()
		if err != nil {
			return err
		}

		return hosts.Add(he, ip, domains)
	},
}
