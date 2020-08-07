package cmd

import (
	"errors"
	"fmt"

	"github.com/pixelandtonic/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/internal/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var logsCommand = &cobra.Command{
	Use:   "logs",
	Short: "Show logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

		p := prompt.NewPrompt()

		// define the flags
		opts := []string{"nginx", "database", "docker"}

		kind, _, err := p.Select("Select the type of logs", opts, &prompt.SelectOptions{
			Default: 1,
		})
		if err != nil {
			return err
		}

		var actions []nitro.Action
		switch kind {
		case "docker":
			containerName, err := p.Ask("Enter the name of the container", &prompt.InputOptions{
				Default:   "",
				Validator: nil,
			})
			if err != nil {
				return err
			}

			if containerName == "" {
				return errors.New("container name cannot be empty")
			}

			dockerLogsAction, err := nitro.LogsDocker(machine, containerName)
			if err != nil {
				return err
			}
			actions = append(actions, *dockerLogsAction)
			fmt.Println("Docker logs for", containerName, "...")
		case "database":
			var databases []config.Database
			if err := viper.UnmarshalKey("databases", &databases); err != nil {
				return err
			}
			var dbs []string
			for _, db := range databases {
				dbs = append(dbs, db.Name())
			}

			if len(dbs) == 0 {
				return errors.New("there are no databases to view logs from")
			}

			containerName, _, err := p.Select("Select database", dbs, &prompt.SelectOptions{
				Default:   1,
				Validator: nil,
			})
			if err != nil {
				return err
			}

			dockerLogsAction, err := nitro.LogsDocker(machine, containerName)
			if err != nil {
				return err
			}
			actions = append(actions, *dockerLogsAction)
			fmt.Println("Database logs for", containerName, "...")
		default:
			nginxLogsAction, err := nitro.LogsNginx(machine, flagNginxLogsKind)
			if err != nil {
				return err
			}
			actions = append(actions, *nginxLogsAction)
			fmt.Println("nginx logs...")
		}

		return nitro.Run(nitro.NewMultipassRunner("multipass"), actions)
	},
}
