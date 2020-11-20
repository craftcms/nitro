package initcmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/client"
	"github.com/craftcms/nitro/pkg/config"
)

// InitCommand is the command for creating new development environments
var InitCommand = &cobra.Command{
	Use:   "init",
	Short: "Create new environment",
	RunE:  initMain,
	Example: `  # create a new environment with the default name
  nitro init

  # create a new environment overriding the default name
  nitro init --environment my-new-env`,
}

func initMain(cmd *cobra.Command, args []string) error {
	env := cmd.Flag("environment").Value.String()

	// TODO(jasonmccallister) check for the env.yaml file in the home directory
	// TODO(jasonmccallister) ask for the default PHP version

	cfg, err := config.Umarshal()
	if err != nil {
		return fmt.Errorf("unable to read config, %w", err)
	}

	// create the new client
	nitro, err := client.NewClient()
	if err != nil {
		return err
	}

	// TODO (jasonmccallister) call the apply command

	if err := nitro.Init(cmd.Context(), env, args); err != nil {
		return err
	}

	if (len(cfg.Sites) > 0 || len(cfg.Databases) > 0) && cmd.Flag("skip-apply").Value.String() != "true" {
		return nitro.Apply(cmd.Context(), env, cfg)
	}

	return nil
}

func init() {
	flags := InitCommand.Flags()

	flags.BoolP("skip-apply", "s", false, "skip applying changes")
}
