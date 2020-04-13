package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/action"
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

		var actions []action.Action

		// check for mounts
		var mounts []config.Mount
		if viper.IsSet("mounts") {
			if err := viper.UnmarshalKey("mounts", &mounts); err != nil {
				return err
			}

			for _, m := range mounts {
				// check for the tilde
				sourceDir := m.Source
				if strings.HasPrefix(m.Source, "~/") {
					home, err := homedir.Dir()
					if err != nil {
						return err
					}

					sourceDir = strings.Replace(m.Source, "~", home, 1)
				}
				if err := validate.Path(sourceDir); err != nil {
					return err
				}

				mountDirectoryAction, err := action.MountDirectory(name, sourceDir, m.Destination)
				if err != nil {
					return err
				}
				actions = append(actions, *mountDirectoryAction)
			}
		}

		launchAction, err := action.Launch(name, cpus, memory, disk, CloudConfig)
		if err != nil {
			return err
		}
		actions = append(actions, *launchAction)

		installAction, err := action.InstallPackages(name, phpVersion)
		if err != nil {
			return err
		}
		actions = append(actions, *installAction)

		// configure php settings that are specific to Craft
		configurePhpMemoryAction, err := action.ConfigurePHPMemoryLimit(name, phpVersion, "256M")
		if err != nil {
			return err
		}
		actions = append(actions, *configurePhpMemoryAction)

		configureExecutionTimeAction, err := action.ConfigurePHPExecutionTimeLimit(name, phpVersion, "240")
		if err != nil {
			return err
		}
		actions = append(actions, *configureExecutionTimeAction)

		xdebugConfigureAction, err := action.ConfigureXdebug(name, phpVersion)
		if err != nil {
			return err
		}
		actions = append(actions, *xdebugConfigureAction)

		restartPhpFpmAction, err := action.RestartPhpFpm(name, phpVersion)
		if err != nil {
			return err
		}
		actions = append(actions, *restartPhpFpmAction)

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

		var siteErrs []error
		if viper.IsSet("sites") {
			var sites []config.Site
			if err := viper.UnmarshalKey("sites", &sites); err != nil {
				return err
			}

			for _, site := range sites {
				// check if the site.Path is a mount
				for _, mount := range mounts {
					if strings.Contains(site.Path, mount.Source) && strings.Contains("/app/sites", mount.Destination) {
						fmt.Println("skipping: " + site.Path + " because the mount has already been set " + mount.Source)
						continue
					}

					// this is not already mounted, so we need to mount it
					if strings.HasPrefix(site.Path, "~/") {
						home, _ := homedir.Dir()
						site.Path = strings.Replace(site.Path, "~", home, 1)
					}

					mountAction, err := action.Mount(name, site.Path, site.Domain)
					if err != nil {
						siteErrs = append(siteErrs, err)
						continue
					}
					actions = append(actions, *mountAction)

					createDirectoryAction, err := action.CreateNginxSiteDirectory(name, site.Domain)
					if err != nil {
						siteErrs = append(siteErrs, err)
						continue
					}
					actions = append(actions, *createDirectoryAction)
				}

				copyTemplateAction, err := action.CopyNginxTemplate(name, site.Domain)
				if err != nil {
					siteErrs = append(siteErrs, err)
					continue
				}
				actions = append(actions, *copyTemplateAction)

				if site.Docroot == "" {
					site.Docroot = "web"
				}
				changeVarsActions, err := action.ChangeTemplateVariables(name, site.Domain, site.Docroot, phpVersion)
				if err != nil {
					siteErrs = append(siteErrs, err)
					continue
				}
				for _, a := range *changeVarsActions {
					actions = append(actions, a)
				}

				createSymlinkAction, err := action.CreateSiteSymllink(name, site.Domain)
				if err != nil {
					siteErrs = append(siteErrs, err)
					continue
				}
				actions = append(actions, *createSymlinkAction)

				reloadNginxAction, err := action.NginxReload(name)
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

			fmt.Printf("---- %d TOTAL COMMANDS ----", len(actions))

			return nil
		}

		if len(siteErrs) > 0 {
			for _, siteErr := range siteErrs {
				fmt.Println(siteErr)
			}
		}

		return action.Run(action.NewMultipassRunner("multipass"), actions)
	},
}

func init() {
	createCommand.Flags().IntVar(&flagCPUs, "cpus", 0, "Number of CPUs to allocate")
	createCommand.Flags().StringVar(&flagMemory, "memory", "", "Amount of memory to allocate")
	createCommand.Flags().StringVar(&flagDisk, "disk", "", "Amount of disk space to allocate")
	createCommand.Flags().StringVar(&flagPhpVersion, "php-version", "", "Which version of PHP to make default")
}
