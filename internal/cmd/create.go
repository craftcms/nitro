package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/validate"
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
		Short:   "Create a new machine",
		PreRun: func(cmd *cobra.Command, args []string) {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name := config.GetString("machine", flagMachineName)
			cpus := config.GetInt("cpus", flagCPUs)
			memory := config.GetString("memory", flagMemory)
			disk := config.GetString("disk", flagDisk)
			phpVersion := config.GetString("php", flagPhpVersion)
			dbEngine := config.GetString("database.engine", flagDatabase)
			dbVersion := config.GetString("database.version", flagDatabaseVersion)

			if viper.IsSet("sites") {
				var sites []config.Site
				viper.UnmarshalKey("sites", &sites)
				fmt.Println(sites)
			}

			if err := validate.DiskSize(disk); err != nil {
				return err
			}
			if err := validate.Memory(memory); err != nil {
				return err
			}
			if err := validate.PHPVersion(phpVersion); err != nil {
				return err
			}
			if err := validate.DatabaseEngineAndVersion(dbEngine, dbVersion); err != nil {
				return err
			}

			if flagDebug {
				fmt.Println("--- DEBUG ---")
				fmt.Println("machine:", name)
				fmt.Println("cpus:", cpus)
				fmt.Println("memory:", memory)
				fmt.Println("disk:", disk)
				fmt.Println("php:", phpVersion)
				fmt.Println("database:", dbEngine, dbVersion)
				fmt.Println("--- DEBUG ---")
				return nil
			}

			if err := nitro.Run(
				nitro.NewMultipassRunner("multipass"),
				nitro.Create(name, strconv.Itoa(cpus), memory, disk, phpVersion, dbEngine, dbVersion),
			); err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	// attach local flags
	createCommand.Flags().IntVar(&flagCPUs, "cpus", 0, "number of cpus to allocate")
	createCommand.Flags().StringVar(&flagMemory, "memory", "", "amount of memory to allocate")
	createCommand.Flags().StringVar(&flagDisk, "disk", "", "amount of disk space to allocate")
	createCommand.Flags().StringVar(&flagPhpVersion, "php-version", "", "which version of PHP to make default")
	createCommand.Flags().StringVar(&flagDatabase, "database", "", "which database engine to make default")
	createCommand.Flags().StringVar(&flagDatabaseVersion, "database-version", "", "which version of the database to install")
}
