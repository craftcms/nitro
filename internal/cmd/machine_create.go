package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/validate"
)

var createCommand = &cobra.Command{
	Use:     "create",
	Aliases: []string{"bootstrap", "boot"},
	Short:   "Create a machine",
	Example: "nitro machine create --name example-machine --cpus 4 --memory 4G --disk 60G --php-version 7.4",
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
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

		var actions []nitro.Action
		launchAction, err := nitro.Launch(name, cpus, memory, disk, CloudConfig)
		if err != nil {
			return err
		}
		actions = append(actions, *launchAction)

		installAction, err := nitro.InstallPackages(name, phpVersion)
		if err != nil {
			return err
		}
		actions = append(actions, *installAction)

		// configure php settings that are specific to Craft
		configurePhpMemoryAction, err := nitro.ConfigurePHPMemoryLimit(name, phpVersion, "256M")
		if err != nil {
			return err
		}
		actions = append(actions, *configurePhpMemoryAction)

		configureExecutionTimeAction, err := nitro.ConfigurePHPExecutionTimeLimit(name, phpVersion, "240")
		if err != nil {
			return err
		}
		actions = append(actions, *configureExecutionTimeAction)

		xdebugConfigureAction, err := nitro.ConfigureXdebug(name, phpVersion)
		if err != nil {
			return err
		}
		actions = append(actions, *xdebugConfigureAction)

		restartPhpFpmAction, err := nitro.RestartPhpFpm(name, phpVersion)
		if err != nil {
			return err
		}
		actions = append(actions, *restartPhpFpmAction)

		// if there are mounts, set them
		if viper.IsSet("mounts") {
			var mounts []config.Mount
			if err := viper.UnmarshalKey("mounts", &mounts); err != nil {
				return err
			}

			for _, mount := range mounts {
				mountDirAction, err := nitro.MountDir(name, mount.AbsSourcePath(), mount.Dest)
				if err != nil {
					return err
				}
				actions = append(actions, *mountDirAction)
			}
		}

		for _, database := range databases {
			volumeAction, err := nitro.CreateDatabaseVolume(name, database.Engine, database.Version, database.Port)
			if err != nil {
				return err
			}
			actions = append(actions, *volumeAction)

			createDatabaseAction, err := nitro.CreateDatabaseContainer(name, database.Engine, database.Version, database.Port)
			if err != nil {
				return err
			}
			actions = append(actions, *createDatabaseAction)
		}

		var siteErrs []error
		if viper.IsSet("sites") {
			var sites []config.Site
			if err := viper.UnmarshalKey("sites", &sites); err != nil {
				return err
			}
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
		}

		if flagDebug {
			fmt.Println("---- COMMANDS ----")
			for _, a := range actions {
				fmt.Println(a.Args)
			}

			fmt.Println("---- CONFIG FILE ----")

			var configFile config.Config
			if err := viper.Unmarshal(&configFile); err != nil {
				return err
			}

			configData, err := yaml.Marshal(configFile)
			if err != nil {
				return err
			}

			fmt.Println(string(configData))

			return nil
		}

		if len(siteErrs) > 0 {
			for _, siteErr := range siteErrs {
				fmt.Println(siteErr)
			}
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
