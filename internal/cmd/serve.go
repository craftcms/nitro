package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/validate"
)

var serveCommand = &cobra.Command{
	Use:   "serve",
	Short: "Serve a website",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("machine", flagMachineName)
		_ = config.GetString("php", flagPhpVersion)
		path := args[0]
		domain := args[1]

		if err := validate.Path(path); err != nil {
			return err
		}

		if err := validate.Domain(domain); err != nil {
			return err
		}

		var commands []nitro.Command

		// attach the provided path to /home/ubuntu/sites/domain.test
		commands = append(commands, nitro.Command{
			Machine:   name,
			Type:      "mount",
			Args:      []string{path, name + ":/home/ubuntu/sites/" + domain},
		})

		// todo run the nginx configuration script
		// todo run the nginx restart service script
		// todo edit the hosts file

		return nitro.Run(nitro.NewMultipassRunner("multipass"), commands)
	},
}

func init() {
	serveCommand.Flags().StringVarP(&flagPhpVersion, "php-version", "p", "", "version of PHP to use")
}
