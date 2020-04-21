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
	Use: "init",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := "nitro-dev"
		if flagMachineName != "" {
			machine = flagMachineName
		}

		if viper.ConfigFileUsed() != "" {
			return errors.New("using a config file already")
		}

		// we don't have a config file
		// set the config file
		var cfg config.Config

		// hardcode the CPUs until this issue is resolved
		// https://github.com/craftcms/nitro/issues/65
		hardCodedCpus := "1"
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

		return errors.New("would create the file " + machine + ".yaml")
	},
}

func getMachine(name string) string {
	if name == "" {
		return "nitro-dev"
	}

	return name
}
