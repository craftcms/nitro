package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-input"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/prompt"
	"github.com/craftcms/nitro/validate"
)

var initCommand = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

		existingConfig := false
		if viper.ConfigFileUsed() != "" {
			fmt.Println("Using an existing config:", viper.ConfigFileUsed())
			existingConfig = true
		}

		ui := &input.UI{
			Writer: os.Stdout,
			Reader: os.Stdin,
		}

		// we don't have a config file
		// set the config file
		var cfg config.Config

		// TODO validate with https://golang.org/pkg/runtime/#NumCPU
		// ask how many cores
		cpuCores, err := prompt.Ask(ui, "How many CPU cores?", "2", true)
		if err != nil {
			return err
		}
		cfg.CPUs = cpuCores

		// ask how much memory
		memory, err := prompt.Ask(ui, "How much memory?", "4G", true)
		if err != nil {
			return err
		}
		cfg.Memory = memory

		// how much disk space
		disk, err := prompt.Ask(ui, "How much disk space??", "40G", true)
		if err != nil {
			return err
		}
		cfg.Disk = disk

		// which version of PHP
		if !existingConfig {
			php, _, err := prompt.Select(ui, "Which version of PHP?", "7.4", nitro.PHPVersions)
			if err != nil {
				return err
			}
			cfg.PHP = php
		} else {
			cfg.PHP = config.GetString("php", flagPhpVersion)

			// double check from the major update
			if cfg.PHP == "" {
				cfg.PHP = "7.4"
			}
		}

		if !existingConfig {
			// what database engine?
			engine, _, err := prompt.Select(ui, "Which database engine?", "mysql", nitro.DBEngines)
			if err != nil {
				return err
			}

			// which version should we use?
			versions := nitro.DBVersions[engine]
			defaultVersion := versions[0]
			version, _, err := prompt.Select(ui, "Select a version of "+engine+"?", defaultVersion, versions)
			if err != nil {
				return err
			}

			// get the port for the engine
			port := "3306"
			if strings.Contains(engine, "postgres") {
				port = "5432"
			}

			cfg.Databases = []config.Database{
				{
					Engine:  engine,
					Version: version,
					Port:    port,
				},
			}
		} else {
			var databases []config.Database
			if err := viper.UnmarshalKey("databases", &databases); err != nil {
				return err
			}

			if databases != nil {
				cfg.Databases = databases
			}
		}

		if len(cfg.Databases) > 0 {
			if err := validate.DatabaseConfig(cfg.Databases); err != nil {
				return err
			}
		}

		var mounts []config.Mount
		var sites []config.Site
		if existingConfig {
			if err := viper.UnmarshalKey("mounts", &mounts); err != nil {
				return err
			}
			if err := viper.UnmarshalKey("sites", &sites); err != nil {
				return err
			}
		}

		// save the config file if it does not exist
		if !existingConfig {
			home, err := homedir.Dir()
			if err != nil {
				return err
			}
			if err := cfg.SaveAs(home, machine); err != nil {
				return err
			}
		}

		cpuCoresInt := 0
		cpuCoresInt, err = strconv.Atoi(cpuCores)
		if err != nil {
			return err
		}

		actions, err := createActions(machine, memory, disk, cpuCoresInt, cfg.PHP, cfg.Databases, mounts, sites)
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

		fmt.Println("Applying the changes now...")

		if err := nitro.Run(nitro.NewMultipassRunner("multipass"), actions); err != nil {
			return err
		}

		return infoCommand.RunE(cmd, args)
	},
}

func init() {
	initCommand.Flags().IntVar(&flagCPUs, "cpus", 0, "Number of CPU cores for machine")
	initCommand.Flags().StringVar(&flagMemory, "memory", "", "Amount of memory for machine")
	initCommand.Flags().StringVar(&flagDisk, "disk", "", "Amount of disk space for machine")
	initCommand.Flags().StringVar(&flagPhpVersion, "php-version", "", "Version of PHP to make default")
}

func createActions(machine, memory, disk string, cpus int, phpVersion string, databases []config.Database, mounts []config.Mount, sites []config.Site) ([]nitro.Action, error) {
	var actions []nitro.Action
	launchAction, err := nitro.Launch(machine, cpus, memory, disk, CloudConfig)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *launchAction)

	installAction, err := nitro.InstallPackages(machine, phpVersion)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *installAction)

	// configure php settings that are specific to Craft
	configurePhpMemoryAction, err := nitro.ConfigurePHPMemoryLimit(machine, phpVersion, "256M")
	if err != nil {
		return nil, err
	}
	actions = append(actions, *configurePhpMemoryAction)

	configureExecutionTimeAction, err := nitro.ConfigurePHPExecutionTimeLimit(machine, phpVersion, "240")
	if err != nil {
		return nil, err
	}
	actions = append(actions, *configureExecutionTimeAction)

	xdebugConfigureAction, err := nitro.ConfigureXdebug(machine, phpVersion)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *xdebugConfigureAction)

	restartPhpFpmAction, err := nitro.RestartPhpFpm(machine, phpVersion)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *restartPhpFpmAction)

	// if there are mounts, set them
	for _, mount := range mounts {
		mountDirAction, err := nitro.MountDir(machine, mount.AbsSourcePath(), mount.Dest)
		if err != nil {
			return nil, err
		}
		actions = append(actions, *mountDirAction)
	}

	for _, database := range databases {
		volumeAction, err := nitro.CreateDatabaseVolume(machine, database.Engine, database.Version, database.Port)
		if err != nil {
			return nil, err
		}
		actions = append(actions, *volumeAction)

		createDatabaseAction, err := nitro.CreateDatabaseContainer(machine, database.Engine, database.Version, database.Port)
		if err != nil {
			return nil, err
		}
		actions = append(actions, *createDatabaseAction)

		setUserPermissions, err := nitro.SetDatabaseUserPermissions(machine, database)
		if err != nil {
			return nil, err
		}
		actions = append(actions, *setUserPermissions)
	}

	var siteErrs []error

	for _, site := range sites {
		copyTemplateAction, err := nitro.CopyNginxTemplate(machine, site.Hostname)
		if err != nil {
			siteErrs = append(siteErrs, err)
			continue
		}
		actions = append(actions, *copyTemplateAction)

		if site.Webroot == "" {
			site.Webroot = "web"
		}

		changeVarsActions, err := nitro.ChangeTemplateVariables(machine, site.Webroot, site.Hostname, phpVersion, site.Aliases)
		if err != nil {
			siteErrs = append(siteErrs, err)
			continue
		}
		for _, a := range *changeVarsActions {
			actions = append(actions, a)
		}

		createSymlinkAction, err := nitro.CreateSiteSymllink(machine, site.Hostname)
		if err != nil {
			siteErrs = append(siteErrs, err)
			continue
		}
		actions = append(actions, *createSymlinkAction)

		reloadNginxAction, err := nitro.NginxReload(machine)
		if err != nil {
			siteErrs = append(siteErrs, err)
			continue
		}
		actions = append(actions, *reloadNginxAction)
	}

	return actions, nil
}
