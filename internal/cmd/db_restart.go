package cmd

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/pixelandtonic/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/scripts"
)

var dbRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "Restart databases",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		mp, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}
		p := prompt.NewPrompt()

		// get all of the docker containers by name
		script := scripts.New(mp, machine)

		var cfg config.Config
		if err := viper.Unmarshal(&cfg); err != nil {
			return err
		}

		if len(cfg.Databases) == 0 {
			return errors.New("there are no databases we can add to")
		}

		// get all of the docker containers by name
		var containers []string
		for _, db := range cfg.Databases {
			containers = append(containers, db.Name())
		}

		// if there is only one
		var container string
		switch len(containers) {
		case 1:
			container = containers[0]
		default:
			container, _, err = p.Select("Which database should we restart", containers, &prompt.SelectOptions{
				Default: 1,
			})
			if err != nil {
				return err
			}
		}

		_, err = script.Run(false, fmt.Sprintf(scripts.FmtDockerRestartContainer, container))
		if err != nil {
			return err
		}

		fmt.Println("Restarted database", container)

		return nil
	},
}
