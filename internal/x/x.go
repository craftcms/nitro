package x

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/nitro/internal/scripts"
)

var phpScript = `
#!/bin/bash
apt install -y php7.4 php7.4-mbstring php7.4-cli php7.4-curl php7.4-fpm php7.4-gd php7.4-intl php7.4-json \
php7.4-mysql php7.4-opcache php7.4-pgsql php7.4-zip php7.4-xml
`

// MultipleCommands is used to test running multiple commands
func MultipleCommands(c *cli.Context) error {
	machine := c.String("machine")
	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	script := make(map[string][]string)
	script["php"] = scripts.InstallPHP("7.4")
	script["mariadb"] = scripts.InstallMariaDB()

	//scripts["update"] = []string{"sudo", "apt", "update", "-y"}
	//scripts["upgrade"] = []string{"sudo", "apt", "update", "-y"}
	//scripts["install"] = []string{"sudo", "apt", "install", "-y", "postgresql-10"}
	//scripts["list"] = []string{"ls", "-la"}
	//scripts["print"] = []string{"pwd"}
	//scripts["who"] = []string{"whoami"}

	for _, c := range script {
		args := []string{"exec", machine, "--"}

		for _, c := range c {
			args = append(args, c)
		}

		cmd := exec.Command(multipass, args...)
		if cmd.Stdin == nil {
			cmd.Stdin = os.Stdin
		}
		if cmd.Stdout == nil {
			cmd.Stdout = os.Stdout
		}

		if err := cmd.Run(); err != nil {
			return err
		}

		//time.Sleep(1 * time.Millisecond)
	}

	return nil
}
