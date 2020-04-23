package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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

		if viper.ConfigFileUsed() != "" {
			// TODO prompt for the confirmation of re initing the machine
			return errors.New("using a config file already")
		}

		// we don't have a config file
		// set the config file
		var cfg config.Config

		// hardcode the CPUs until this issue is resolved
		// https://github.com/craftcms/nitro/issues/65
		hardCodedCpus := "2"
		cpuInt, err := strconv.Atoi(hardCodedCpus)
		if err != nil {
			return err
		}
		cfg.CPUs = hardCodedCpus

		// ask how much memory
		memory, err := prompt.AskWithDefault("How much memory should we assign?", "4G", nil)
		if err != nil {
			return err
		}
		cfg.Memory = memory

		// how much disk space
		disk, err := prompt.AskWithDefault("How much disk space should the machine have?", "40G", nil)
		if err != nil {
			return err
		}
		cfg.Disk = disk

		// which version of PHP
		_, php := prompt.SelectWithDefault("Which version of PHP should we install?", "7.4", nitro.PHPVersions)
		cfg.PHP = php

		// what database engine?
		_, engine := prompt.SelectWithDefault("Which database engine should we setup?", "mysql", nitro.DBEngines)

		// which version should we use?
		versions := nitro.DBVersions[engine]
		defaultVersion := versions[0]
		_, version := prompt.SelectWithDefault("Select a version of "+engine+" to use:", defaultVersion, versions)

		// get the port for the engine
		port := "3306"
		if strings.Contains(engine, "postgres") {
			port = "5432"
		}
		// TODO check if the port has already been used and +1 it

		cfg.Databases = []config.Database{
			{
				Engine:  engine,
				Version: version,
				Port:    port,
			},
		}

		if err := validate.DatabaseConfig(cfg.Databases); err != nil {
			return err
		}

		// save the config file
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		if err := cfg.SaveAs(home, machine); err != nil {
			return err
		}

		actions, err := createActions(machine, memory, disk, cpuInt, php, cfg.Databases, nil, nil)
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

		fmt.Println("Ok, applying the changes now")

		return nitro.Run(nitro.NewMultipassRunner("multipass"), actions)
	},
}

func init() {
	// initCommand.Flags().IntVar(&flagCPUs, "cpus", 0, "Number of CPUs to allocate")
	initCommand.Flags().StringVar(&flagMemory, "memory", "", "Amount of memory to allocate")
	initCommand.Flags().StringVar(&flagDisk, "disk", "", "Amount of disk space to allocate")
	initCommand.Flags().StringVar(&flagPhpVersion, "php-version", "", "Which version of PHP to make default")
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
