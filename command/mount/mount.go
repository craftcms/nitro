package mount

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/phpversions"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # mount the current directory into a container
  nitro mount

  # mount a directory
  nitro mount path`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mount",
		Short:   "Mount directory into a container",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			output.Info("Mounting directory…")

			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			// get working directory or provided arg
			var absDir string
			switch len(args) {
			case 1:
				absDir = filepath.Join(wd, args[0])
				absDir = filepath.Clean(absDir)
			default:
				absDir = filepath.Clean(wd)
			}

			displayName := strings.Replace(absDir, home, "~", 1)
			fmt.Println(displayName)

			mount := config.Mount{
				Path: displayName,
			}

			// prompt for the php version
			versions := phpversions.Versions
			selected, err := output.Select(cmd.InOrStdin(), "Choose a PHP version: ", versions)
			if err != nil {
				return err
			}

			mount.Version = versions[selected]

			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// add the site to the config
			if err := cfg.AddMount(mount); err != nil {
				return err
			}

			// save the config file
			if err := cfg.Save(); err != nil {
				output.Warning()

				return err
			}

			output.Info("Mount added")

			// ask if the apply command should run
			var response string
			fmt.Print("Apply changes now [Y/n]? ")
			if _, err := fmt.Scanln(&response); err != nil {
				return fmt.Errorf("unable to provide a prompt, %w", err)
			}

			// get the response
			resp := strings.TrimSpace(response)
			var confirm bool
			for _, answer := range []string{"y", "Y", "yes", "Yes", "YES"} {
				if resp == answer {
					confirm = true
				}
			}

			// we are skipping the apply step or there is no parent command
			if !confirm || cmd.Parent() == nil {
				return nil
			}

			// get the apply command and run it
			for _, c := range cmd.Parent().Commands() {
				if c.Use == "apply" {
					return c.RunE(c, args)
				}
			}

			return nil
		},
	}

	// set sub commands
	cmd.AddCommand(sshCommand(home, docker, output))

	return cmd
}
