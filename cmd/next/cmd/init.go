package cmd

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"

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
)

func init() {
	// attach the command to root
	rootCmd.AddCommand(initCommand)

	// attach local flags
	initCommand.Flags().StringVar(&flagMachineName, "machine", "", "name of machine")
	initCommand.Flags().IntVar(&flagCPUs, "cpus", 2, "number of cpus to allocate")
	initCommand.Flags().StringVar(&flagMemory, "memory", "4G", "amount of memory to allocate")
	initCommand.Flags().StringVar(&flagDisk, "disk", "40G", "amount of disk space to allocate")
	initCommand.Flags().StringVar(&flagPhpVersion, "php-version", "7.4", "which version of PHP to make default")
	initCommand.Flags().StringVar(&flagDatabase, "database", "mysql", "which database engine to make default")
	initCommand.Flags().StringVar(&flagPhpVersion, "database-version", "5.7", "which version of the database to install")
}

var initCommand = &cobra.Command{
	Use:     "init",
	Aliases: []string{"bootstrap", "boot"},
	Short:   "Create a new machine",
	PreRun: func(cmd *cobra.Command, args []string) {
		// set the defaults and load the yaml
		// TODO validate options for php and etc
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := nitro.Run(
			command.NewMultipassRunner("multipass"),
			nitro.Init(flagMachineName, strconv.Itoa(flagCPUs), flagMemory, flagDisk, flagPhpVersion, flagDatabase, flagDatabaseVersion),
		); err != nil {
			log.Fatal(err)
		}
	},
}
