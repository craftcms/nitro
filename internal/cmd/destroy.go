package cmd

import (
	"fmt"
	"github.com/pixelandtonic/prompt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/sudo"
)

var destroyCommand = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy a machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

		// get the sites
		var cfg config.Config
		if err := viper.Unmarshal(&cfg); err != nil {
			return err
		}

		var domains []string
		for _, site := range cfg.Sites {
			domains = append(domains, site.Hostname)
		}

		p := prompt.NewPrompt()

		reallyDestroy, err := p.Confirm("Are you sure you want to permanently destroy your " + machine + " machine", &prompt.InputOptions{
			Default:   "no",
			Validator: nil,
		})
		if err != nil {
			return err
		}

		if !reallyDestroy {
			return nil
		}

		destroyAction, err := nitro.Destroy(machine)
		if err != nil {
			return err
		}

		if err := nitro.Run(nitro.NewMultipassRunner("multipass"), []nitro.Action{*destroyAction}); err != nil {
			return err
		}

		if flagClean {
			if err := os.Remove(viper.ConfigFileUsed()); err != nil {
				fmt.Println("unable to remove the config:", viper.ConfigFileUsed())
			}
		}

		if len(domains) == 0 {
			fmt.Println("Permanently destroyed", machine)
			return nil
		}

		cmds := []string{"hosts", "remove"}
		for _, domain := range domains {
			cmds = append(cmds, domain)
		}

		// prompt to remove hosts file
		nitro, err := exec.LookPath("nitro")
		if err != nil {
			return err
		}

		fmt.Println("Removing sites from your hosts file")

		return sudo.RunCommand(nitro, machine, cmds...)
	},
}

func init() {
	destroyCommand.Flags().BoolVar(&flagClean, "clean", false, "remove the config file when destroying the machine")
}
