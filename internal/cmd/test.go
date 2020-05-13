package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pixelandtonic/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/scripts"
	"github.com/craftcms/nitro/internal/webroot"
)

var testCommand = &cobra.Command{
	Use:   "test",
	Short: "Testing",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}

		sites := configFile.GetSites()

		if len(sites) == 0 {
			return errors.New("there are no sites to select")
		}

		p := prompt.NewPrompt()

		mp, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}

		script := scripts.New(mp, machine)

		var site config.Site
		_, i, err := p.Select("Select a site to check", configFile.SitesAsList(), &prompt.SelectOptions{
			Default: 1,
		})
		if err != nil {
			return err
		}
		site = sites[i]

		// check if the site it available
		output, err := script.Run(fmt.Sprintf(scripts.FmtNginxSiteAvailable, site.Hostname))
		if err != nil {
			return err
		}
		if strings.Contains(output, "exists") {
			fmt.Println(site.Hostname, "exists")
		}

		// check if the site webroot matches
		webrootOutput, err := script.Run(fmt.Sprintf(scripts.FmtNginxSiteWebroot, site.Hostname))
		if err != nil {
			return err
		}
		if webroot.Matches(webrootOutput, site.Webroot) {
			fmt.Println(fmt.Sprintf("The webroot %q matches", site.Webroot))
		} else {
			fmt.Println(fmt.Sprintf("The webroot %q does NOT match", site.Webroot))
		}

		return nil

		//var actions []nitro.Action

		//complexAction := nitro.Action{
		//	Type:       "exec",
		//	UseSyscall: false,
		//	Args:       []string{"exec", machine, "--", `bash`, `-c`, `if test -f 'test'; then echo 'exists'; fi`},
		//}

		//actions = append(actions, complexAction)

		//if err := nitro.Run(nitro.NewMultipassRunner("multipass"), actions); err != nil {
		//	return err
		//}

		// return nil
	},
}

func init() {
	rootCmd.AddCommand(testCommand)
}
