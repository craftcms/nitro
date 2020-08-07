package cmd

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/pixelandtonic/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/runas"
	"github.com/craftcms/nitro/suggest"
	"github.com/craftcms/nitro/validate"
)

var initCommand = &cobra.Command{
	Use:   "init",
	Short: "Create new machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		existingConfig := false
		runner := nitro.NewMultipassRunner("multipass")
		p := prompt.NewPrompt()

		// check if the machine exists
		if ip := nitro.IP(machine, runner); ip != "" {
			fmt.Println(fmt.Sprintf("Machine %q already exists, skipping the init process", machine))
			return nil
		}

		if viper.ConfigFileUsed() != "" {
			fmt.Println("Using an existing config:", viper.ConfigFileUsed())
			existingConfig = true
		}

		if existingConfig == false {
			initMachine, err := p.Confirm("Initialize the primary machine now", &prompt.InputOptions{Default: "yes", AppendQuestionMark: true})
			if err != nil {
				return err
			}
			if initMachine == false {
				fmt.Println("You can create a new machine later by running `nitro init`")
				return nil
			}
		}

		// we don't have a config file
		// set the config file
		var cfg config.Config
		actual := runtime.NumCPU()
		cpuCores, err := p.Ask("How many CPU cores", &prompt.InputOptions{
			Default:            suggest.NumberOfCPUs(actual),
			Validator:          validate.NewCPUValidator(actual).Validate,
			AppendQuestionMark: true,
		})
		if err != nil {
			return err
		}

		// ask how much memory
		mem, err := p.Ask("How much memory", &prompt.InputOptions{
			Default:            "4G",
			Validator:          validate.Memory,
			AppendQuestionMark: true,
		})
		if err != nil {
			return err
		}
		memory := mem

		// how much disk space
		di, err := p.Ask("How much disk space", &prompt.InputOptions{
			Default:            "40G",
			Validator:          validate.DiskSize,
			AppendQuestionMark: true,
		})
		if err != nil {
			return err
		}
		disk := di

		// which version of PHP
		if !existingConfig {
			var loop bool
			for ok := true; ok; ok = !loop {
				php, err := p.Ask("Which version of PHP", &prompt.InputOptions{
					Default:            "7.4",
					Validator:          validate.PHPVersion,
					AppendQuestionMark: true,
				})

				if err == nil {
					loop = true
					cfg.PHP = php
				} else {
					loop = false
					fmt.Println("Invalid input. Possible PHP versions are:", strings.Join(nitro.PHPVersions, ", "))
				}
			}
		} else {
			// double check from the major update
			if cfg.PHP == "" {
				cfg.PHP = "7.4"
			}
		}

		if !existingConfig {
			// what database engine?
			var dbEngineLoop bool
			var engine string
			for ok := true; ok; ok = !dbEngineLoop {
				engine, err = p.Ask("Which database engine", &prompt.InputOptions{
					Default:            "mysql",
					Validator:          validate.DatabaseEngine,
					AppendQuestionMark: true,
				})

				if err == nil {
					dbEngineLoop = true
				} else {
					fmt.Println("Invalid input. Possible database engines are:", strings.Join(nitro.DBEngines, ", "))
					dbEngineLoop = false
				}
			}

			// get the database version
			var dbVersionLoop bool
			var version string
			for ok := true; ok; ok = !dbVersionLoop {
				versions := nitro.DBVersions[engine]
				defaultVersion := versions[0]
				version, _ = p.Ask("Which version of "+engine, &prompt.InputOptions{
					Default:            defaultVersion,
					AppendQuestionMark: true,
				})

				err := validate.DatabaseEngineAndVersion(engine, version)

				if err == nil {
					dbVersionLoop = true
				} else {
					fmt.Println("Invalid input. Possible database versions are:", strings.Join(nitro.DBVersions[engine], ", "))
					dbVersionLoop = false
				}
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

		if err := nitro.Run(runner, actions); err != nil {
			return err
		}

		// if there are sites, edit the hosts file
		if len(sites) > 0 {
			switch runtime.GOOS {
			case "windows":
				if err := hostsCommand.RunE(cmd, args); err != nil {
					return err
				}
			default:
				if err := runas.Elevated(machine, []string{"hosts"}); err != nil {
					return err
				}
			}
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
