package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var (
	flagCPUs            int
	flagMemory          string
	flagDisk            string
	flagPhpVersion      string
	flagDatabase        string
	flagDatabaseVersion string

	initCommand = &cobra.Command{
		Use:     "init",
		Aliases: []string{"bootstrap", "boot"},
		Short:   "Create a new machine",
		Run: func(cmd *cobra.Command, args []string) {
			name := config.GetString("machine", flagMachineName)
			cpus := config.GetInt("cpus", flagCPUs)
			memory := config.GetString("memory", flagMemory)
			disk := config.GetString("disk", flagDisk)
			phpVersion := config.GetString("phpVersion-version", flagPhpVersion)
			dbEngine := config.GetString("database.engine", flagDatabase)
			dbVersion := config.GetString("database.version", flagDatabaseVersion)

			if flagDyRun {
				fmt.Println("--- DEBUG ---")
				fmt.Println("machine:", name)
				fmt.Println("cpus:", cpus)
				fmt.Println("memory:", memory)
				fmt.Println("disk:", disk)
				fmt.Println("phpVersion-version:", phpVersion)
				fmt.Println("database-engine:", dbEngine)
				fmt.Println("database-version:", dbVersion)
				fmt.Println("--- DEBUG ---")
				return
			}

			if err := nitro.Run(
				nitro.NewMultipassRunner("multipass"),
				nitro.Init(name, strconv.Itoa(cpus), memory, disk, phpVersion, dbEngine, dbVersion),
			); err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	// attach local flags
	initCommand.Flags().IntVar(&flagCPUs, "cpus", 0, "number of cpus to allocate")
	initCommand.Flags().StringVar(&flagMemory, "memory", "", "amount of memory to allocate")
	initCommand.Flags().StringVar(&flagDisk, "disk", "", "amount of disk space to allocate")
	initCommand.Flags().StringVar(&flagPhpVersion, "php-version", "", "which version of PHP to make default")
	initCommand.Flags().StringVar(&flagDatabase, "database", "", "which database engine to make default")
	initCommand.Flags().StringVar(&flagDatabaseVersion, "database-version", "", "which version of the database to install")
}
