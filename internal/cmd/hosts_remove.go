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
)

var hostsRemoveCommand = &cobra.Command{
	Use:    "remove",
	Short:  "Remove an entry from your hosts file",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !flagDebug {
			uid := os.Geteuid()
			if uid != 0 {
				return errors.New("you do not appear to be running this command as root, so we cannot modify your hosts file")
			}
		}

		// get all of the sites from the config file
		if !viper.IsSet("sites") {
			return errors.New("unable to read sites from " + viper.ConfigFileUsed())
		}

		var sites []config.Site
		if err := viper.UnmarshalKey("sites", &sites); err != nil {
			return err
		}

		he, err := txeh.NewHostsDefault()
		if err != nil {
			return err
		}

		var domains []string
		for _, site := range sites {
			domains = append(domains, site.Hostname)
		}

		if flagDebug {
			for _, domain := range domains {
				fmt.Println("removing", domain, "from hosts file")
			}

			return nil
		}

		return hosts.Remove(he, domains)
	},
}
