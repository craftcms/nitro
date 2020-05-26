package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/scripts"
)

var composerCommand = &cobra.Command{
	Use:   "composer",
	Short: "Install composer",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		mp, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}
		script := scripts.New(mp, machine)

		// create the local directory for mailhog
		if output, err := script.Run(false, fmt.Sprintf(scripts.FmtCreateDirectory, "/home/ubuntu/.composer")); err != nil {
			fmt.Println(output)
			return err
		}

		// download the installer
		if output, err := script.Run(false, `curl -sS https://getcomposer.org/installer -o composer-setup.php`); err != nil {
			fmt.Println(output)
			return err
		}

		// run the installer
		if output, err := script.Run(true, `php composer-setup.php --install-dir=/usr/local/bin --filename=composer`); err != nil {
			fmt.Println(output)
			return err
		}

		fmt.Println(fmt.Sprintf("Composer is now installed on %q.", machine))

		return nil
	},
}
