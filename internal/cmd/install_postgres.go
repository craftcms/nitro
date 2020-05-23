package cmd

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/pixelandtonic/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/scripts"
	"github.com/craftcms/nitro/validate"
)

type newPostgresValidator struct {
	cfg *config.Config
}

func (v newPostgresValidator) ValidateVersion(version string) error {
	if err := validate.DatabaseEngineAndVersion("postgres", version); err != nil {
		return err
	}

	for _, db := range v.cfg.Databases {
		if version == db.Version {
			return errors.New(fmt.Sprintf("postgres version %q is already installed", version))
		}
	}

	return nil
}

func (v newPostgresValidator) ValidatePort(port string) error {
	for _, db := range v.cfg.Databases {
		if port == db.Port {
			return errors.New(fmt.Sprintf("postgres port %q is already in use", port))
		}
	}

	return nil
}

var postgresCommand = &cobra.Command{
	Use:     "postgres",
	Aliases: []string{"postgresql", "psql", "pg"},
	Short:   "Setup PostgreSQL",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		mp, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}
		p := prompt.NewPrompt()
		_ = scripts.New(mp, machine)

		// get the config
		cfg, err := config.Read()
		if err != nil {
			return err
		}

		validator := newPostgresValidator{cfg: cfg}

		// ask for the version
		version, err := p.Ask(fmt.Sprintf("Which version of postgres should we install"), &prompt.InputOptions{
			Validator: validator.ValidateVersion,
		})
		if err != nil {
			return err
		}

		// ask for the port assignment
		port, err := p.Ask(fmt.Sprintf("Which port should postgres listen on"), &prompt.InputOptions{
			Validator: validator.ValidatePort,
		})
		if err != nil {
			return err
		}

		// save to the config file
		cfg.Databases = append(cfg.Databases, config.Database{
			Engine:  "postgres",
			Version: version,
			Port:    port,
		})

		// save the file
		if err := cfg.Save(viper.ConfigFileUsed()); err != nil {
			fmt.Println("Error saving the config file")
			return err
		}

		fmt.Println(fmt.Sprintf("Adding PostgreSQL version %q on port %q", version, port))

		// prompt for the apply command
		apply, err := p.Confirm("Apply changes from config now", &prompt.InputOptions{
			Default:   "yes",
			Validator: nil,
		})
		if err != nil {
			return err
		}

		if apply {
			return applyCommand.RunE(cmd, args)
		}

		return nil
	},
}
