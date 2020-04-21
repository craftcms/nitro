package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var logsCommand = &cobra.Command{
	Use:   "logs",
	Short: "Show machine logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := "nitro-dev"
		if flagMachineName != "" {
			machine = flagMachineName
		}

		// define the flags
		opts := []string{"nginx", "database", "docker"}

		logType := promptui.Select{
			Label: "Select the type of logs to view",
			Items: opts,
			Size:  len(opts),
		}

		_, kind, err := logType.Run()
		if err != nil {
			return err
		}

		var actions []nitro.Action
		switch kind {
		case "docker":
			validate := func(input string) error {
				if input == "" {
					return errors.New("container machine cannot be empty")
				}
				if strings.Contains(input, " ") {
					return errors.New("container names cannot contain spaces")
				}
				return nil
			}

			containerNamePrompt := promptui.Prompt{
				Label:    "Enter container machine",
				Validate: validate,
			}

			containerName, err := containerNamePrompt.Run()
			if err != nil {
				return err
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
				dbs = append(dbs, fmt.Sprintf("%s_%s_%s", db.Engine, db.Version, db.Port))
			}
			databaseContainerName := promptui.Select{
				Label: "Select database",
				Items: dbs,
			}

			_, containerName, err := databaseContainerName.Run()
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
