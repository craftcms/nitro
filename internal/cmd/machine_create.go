package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/helpers"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/prompt"
	"github.com/craftcms/nitro/validate"
)

var createCommand = &cobra.Command{
	Use:     "create",
	Aliases: []string{"bootstrap", "boot"},
	Short:   "Create a machine",
	Example: "nitro machine create --name example-machine --cpus 4 --memory 4G --disk 60G --php-version 7.4",
	PreRun: func(cmd *cobra.Command, args []string) {
		if viper.ConfigFileUsed() != "" {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// load the config file
		// if the config does not exist, prompt the user
		if viper.ConfigFileUsed() == "" {
			home, err := homedir.Dir()
			if err != nil {
				return err
			}

			// make the ~/.nitro/ directory
			nitroDir := home + "/.nitro/"
			if err := helpers.MkdirIfNotExists(nitroDir); err != nil {
				return err
			}

			if err := helpers.CreateFileIfNotExist(nitroDir + "nitro.yaml"); err != nil {
				fmt.Println(err)
			}

			// set the config file
			var configFile config.Config

			// name
			name, err := prompt.Ask("What should the machine be named?", "nitro-dev", validate.MachineName)
			if err != nil {
				return err
			}
			configFile.Name = name

			// number of cpus 1
			cpus, err := prompt.Ask("How many CPUs should the machine have?", "1", nil)
			if err != nil {
				return err
			}
			configFile.CPUs = cpus

			// how much memory 4G
			memory, err := prompt.Ask("How much memory should the machine have?", "4G", validate.Memory)
			if err != nil {
				return err
			}
			configFile.Memory = memory

			// how large should the disk size be? 40G
			disk, err := prompt.Ask("How much disk space should the machine have?", "40G", validate.Memory)
			if err != nil {
				return err
			}
			configFile.Disk = disk

			// which version of PHP would you like installed? 7.5
			phpPrompt := promptui.Select{
				Label:     "Which version of PHP should we install?",
				Items:     nitro.PHPVersions,
				CursorPos: 0,
			}
			_, phpVersion, err := phpPrompt.Run()
			if err != nil {
				return err
			}
			configFile.PHP = phpVersion

			// what database engine would you like to use? mysql
			dbEnginePrompt := promptui.Select{
				Label:     "Which database engine should the machine have?",
				Items:     nitro.DBEngines,
				CursorPos: 0,
			}
			_, dbEngine, err := dbEnginePrompt.Run()
			if err != nil {
				return err
			}

			_, dbVersion := prompt.Select("Select a version of "+dbEngine+" to use:", nitro.DBVersions[dbEngine])

			dbPort := "3306"
			if dbEngine == "postgres" {
				dbPort = "5432"
			}

			db := config.Database{
				Engine:  dbEngine,
				Version: dbVersion,
				Port:    dbPort,
			}
			configFile.Databases = []config.Database{db}

			if err := validate.DatabaseConfig(configFile.Databases); err != nil {
				return err
			}

			// save the config file
			if err := configFile.Save(nitroDir + "nitro.yaml"); err != nil {
				return err
			}

			cpu, err := strconv.Atoi(cpus)
			if err != nil {
				return err
			}

			actions, err := createActions(name, memory, disk, cpu, phpVersion, configFile.Databases, nil, nil)
			if err != nil {
				return err
			}

			if flagDebug {
				fmt.Println("---- COMMANDS ----")
				for _, a := range actions {
					fmt.Println(a.Args)
				}

				return nil
			}

			return nitro.Run(nitro.NewMultipassRunner("multipass"), actions)
		}

		// run the actions

		// grab the config/options for the command
		name := config.GetString("name", flagMachineName)
		cpus := config.GetInt("cpus", flagCPUs)
		memory := config.GetString("memory", flagMemory)
		disk := config.GetString("disk", flagDisk)
		phpVersion := config.GetString("php", flagPhpVersion)

		// validate options
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
		if err := viper.UnmarshalKey("databases", &databases); err != nil {
			return err
		}

		if err := validate.DatabaseConfig(databases); err != nil {
			return err
		}

		var mounts []config.Mount
		if viper.IsSet("mounts") {
			if err := viper.UnmarshalKey("mounts", &mounts); err != nil {
				return err
			}
		}

		var sites []config.Site
		if viper.IsSet("sites") {
			if err := viper.UnmarshalKey("sites", &sites); err != nil {
				return err
			}
		}

		actions, err := createActions(name, memory, disk, cpus, phpVersion, databases, mounts, sites)
		if err != nil {
			return err
		}

		if flagDebug {
			fmt.Println("---- COMMANDS ----")
			for _, a := range actions {
				fmt.Println(a.Args)
			}

			return nil
		}

		return nitro.Run(nitro.NewMultipassRunner("multipass"), actions)
	},
}

func init() {
	createCommand.Flags().IntVar(&flagCPUs, "cpus", 0, "Number of CPUs to allocate")
	createCommand.Flags().StringVar(&flagMemory, "memory", "", "Amount of memory to allocate")
	createCommand.Flags().StringVar(&flagDisk, "disk", "", "Amount of disk space to allocate")
	createCommand.Flags().StringVar(&flagPhpVersion, "php-version", "", "Which version of PHP to make default")
}

func createActions(name, memory, disk string, cpus int, phpVersion string, databases []config.Database, mounts []config.Mount, sites []config.Site) ([]nitro.Action, error) {
	var actions []nitro.Action
	launchAction, err := nitro.Launch(name, cpus, memory, disk, CloudConfig)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *launchAction)

	installAction, err := nitro.InstallPackages(name, phpVersion)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *installAction)

	// configure php settings that are specific to Craft
	configurePhpMemoryAction, err := nitro.ConfigurePHPMemoryLimit(name, phpVersion, "256M")
	if err != nil {
		return nil, err
	}
	actions = append(actions, *configurePhpMemoryAction)

	configureExecutionTimeAction, err := nitro.ConfigurePHPExecutionTimeLimit(name, phpVersion, "240")
	if err != nil {
		return nil, err
	}
	actions = append(actions, *configureExecutionTimeAction)

	xdebugConfigureAction, err := nitro.ConfigureXdebug(name, phpVersion)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *xdebugConfigureAction)

	restartPhpFpmAction, err := nitro.RestartPhpFpm(name, phpVersion)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *restartPhpFpmAction)

	// if there are mounts, set them
	for _, mount := range mounts {
		mountDirAction, err := nitro.MountDir(name, mount.AbsSourcePath(), mount.Dest)
		if err != nil {
			return nil, err
		}
		actions = append(actions, *mountDirAction)
	}

	for _, database := range databases {
		volumeAction, err := nitro.CreateDatabaseVolume(name, database.Engine, database.Version, database.Port)
		if err != nil {
			return nil, err
		}
		actions = append(actions, *volumeAction)

		createDatabaseAction, err := nitro.CreateDatabaseContainer(name, database.Engine, database.Version, database.Port)
		if err != nil {
			return nil, err
		}
		actions = append(actions, *createDatabaseAction)
	}

	var siteErrs []error

	for _, site := range sites {
		copyTemplateAction, err := nitro.CopyNginxTemplate(name, site.Hostname)
		if err != nil {
			siteErrs = append(siteErrs, err)
			continue
		}
		actions = append(actions, *copyTemplateAction)

		if site.Webroot == "" {
			site.Webroot = "web"
		}

		changeVarsActions, err := nitro.ChangeTemplateVariables(name, site.Webroot, site.Hostname, phpVersion, site.Aliases)
		if err != nil {
			siteErrs = append(siteErrs, err)
			continue
		}
		for _, a := range *changeVarsActions {
			actions = append(actions, a)
		}

		createSymlinkAction, err := nitro.CreateSiteSymllink(name, site.Hostname)
		if err != nil {
			siteErrs = append(siteErrs, err)
			continue
		}
		actions = append(actions, *createSymlinkAction)

		reloadNginxAction, err := nitro.NginxReload(name)
		if err != nil {
			siteErrs = append(siteErrs, err)
			continue
		}
		actions = append(actions, *reloadNginxAction)
	}

	return actions, nil
}
