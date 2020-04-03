package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/command"
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
			// look for the defaults from the config
			var name string
			if viper.IsSet("machine") && flagMachineName == "" {
				name = viper.GetString("machine")
			} else {
				name = flagMachineName
			}

			var cpus string
			if viper.IsSet("cpus") && flagCPUs == 0 {
				cpus = strconv.Itoa(viper.GetInt("cpus"))
			} else {
				cpus = strconv.Itoa(flagCPUs)
			}

			var memory string
			if viper.IsSet("memory") && flagMemory == "" {
				memory = viper.GetString("memory")
			} else {
				memory = flagMemory
			}

			var disk string
			if viper.IsSet("disk") && flagDisk == "" {
				disk = viper.GetString("disk")
			} else {
				disk = flagDisk
			}

			var php string
			if viper.IsSet("php") && flagPhpVersion == "" {
				php = viper.GetString("php")
			} else {
				php = flagPhpVersion
			}

			var db string
			if viper.IsSet("database.engine") && flagDatabase == "" {
				db = viper.GetString("database.engine")
			} else {
				db = flagDatabase
			}

			var dbVersion string
			if viper.IsSet("database.version") && flagDatabaseVersion == "" {
				dbVersion = viper.GetString("database.version")
			} else {
				dbVersion = flagDatabaseVersion
			}

			if flagDyRun {
				fmt.Println("--- DEBUG ---")
				fmt.Println("machine:", name)
				fmt.Println("cpus:", cpus)
				fmt.Println("memory:", memory)
				fmt.Println("disk:", disk)
				fmt.Println("php-version:", php)
				fmt.Println("database-engine:", db)
				fmt.Println("database-version:", dbVersion)
				fmt.Println("--- DEBUG ---")
				return
			}

			if err := nitro.Run(
				command.NewMultipassRunner("multipass"),
				nitro.Init(name, cpus, memory, disk, php, db, dbVersion),
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
