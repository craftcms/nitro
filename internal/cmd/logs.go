package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/pixelandtonic/go-input"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/prompt"
)

var logsCommand = &cobra.Command{
	Use:   "logs",
	Short: "Show machine logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

		// define the flags
		opts := []string{"nginx", "database", "docker"}
		ui := &input.UI{
			Writer: os.Stdout,
			Reader: os.Stdin,
		}

		kind, _, err := prompt.Select(ui, "Select the type of logs to view", "nginx", opts)
		if err != nil {
			return err
		}

		var actions []nitro.Action
		switch kind {
		case "docker":
			containerName, err := prompt.Ask(ui, "Enter container name:", "", true)
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
			fmt.Println("Here are the docker logs for", containerName, "...")
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

			containerName, _, err := prompt.Select(ui, "Select database", dbs[0], dbs)
			if err != nil {
				return err
			}

			dockerLogsAction, err := nitro.LogsDocker(machine, containerName)
			if err != nil {
				return err
			}
			actions = append(actions, *dockerLogsAction)
			fmt.Println("Here are the database logs for", containerName, "...")
		default:
			fmt.Println("Here are the nginx logs...")
			nginxLogsAction, err := nitro.LogsNginx(machine, flagNginxLogsKind)
			if err != nil {
				return err
			}
			actions = append(actions, *nginxLogsAction)
			fmt.Println("Here are the nginx logs...")
		}

		return nitro.Run(nitro.NewMultipassRunner("multipass"), actions)
	},
}
