package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/action"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/validate"
)

var (
	flagCPUs            int
	flagMemory          string
	flagDisk            string
	flagPhpVersion      string
	flagDatabase        string
	flagDatabaseVersion string

	createCommand = &cobra.Command{
		Use:     "create",
		Aliases: []string{"bootstrap", "boot"},
		Short:   "Create machine",
		Example: "nitro machine create --name example-machine --cpus 4 --memory 4G --disk 60G --php-version 7.4",
		PreRun: func(cmd *cobra.Command, args []string) {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name := config.GetString("name", flagMachineName)
			cpus := config.GetInt("cpus", flagCPUs)
			memory := config.GetString("memory", flagMemory)
			disk := config.GetString("disk", flagDisk)
			phpVersion := config.GetString("php", flagPhpVersion)

			// validate
			if err := validate.DiskSize(disk); err != nil {
				return err
			}
			if err := validate.Memory(memory); err != nil {
				return err
			}
			if err := validate.PHPVersion(phpVersion); err != nil {
				return err
			}
			if !viper.IsSet("databases") {
				return errors.New("no databases defined in " + viper.ConfigFileUsed())
			}

			var databases []config.Database
			_ = viper.UnmarshalKey("databases", &databases)

			if err := validate.DatabaseConfig(databases); err != nil {
				return err
			}

			var actions []action.Action
			launchAction, err := action.Launch(name, cpus, memory, disk, nitro.CloudConfig)
			if err != nil {
				return err
			}
			actions = append(actions, *launchAction)

			installAction, err := action.InstallPackages(name, phpVersion)
			if err != nil {
				return err
			}
			actions = append(actions, *installAction)

			for _, database := range databases {
				volumeAction, err := action.CreateDatabaseVolume(name, database.Engine, database.Version, database.Port)
				if err != nil {
					return err
				}
				actions = append(actions, *volumeAction)

				createDatabaseAction, err := action.CreateDatabaseContainer(name, database.Engine, database.Version, database.Port)
				if err != nil {
					return err
				}
				actions = append(actions, *createDatabaseAction)
			}

			if flagDebug {
				for _, a := range actions {
					fmt.Println(a.Args)
				}

				return nil
			}

			return nitro.RunAction(nitro.NewMultipassRunner("multipass"), actions)
		},
	}
)

func init() {
	createCommand.Flags().IntVar(&flagCPUs, "cpus", 0, "number of cpus to allocate")
	createCommand.Flags().StringVar(&flagMemory, "memory", "", "amount of memory to allocate")
	createCommand.Flags().StringVar(&flagDisk, "disk", "", "amount of disk space to allocate")
	createCommand.Flags().StringVar(&flagPhpVersion, "php-version", "", "which version of PHP to make default")
}
