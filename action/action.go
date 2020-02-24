package action

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/phpdev/install"
)

// Initialize is used to create a new machine and setup any dependencies
func Initialize(c *cli.Context) error {
	machine := c.String("machine")
	php := c.String("php")
	// TODO remove the hardcoding
	database := "mariadb"

	fmt.Println("Creating a new machine:", machine)

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	// create the machine
	cmd := exec.Command(multipass, "launch", "--name", machine, "--cloud-init", "./scripts/cloud-init.yaml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// update the machine
	if err := c.App.RunContext(c.Context, []string{c.App.Name, "--machine", machine, "update"}); err != nil {
		return err
	}

	// install the PHP version request
	if php == "" {
		fmt.Println("ERROR")
		fmt.Println("php is empty")
		fmt.Println("ERROR")
		php = "7.4"
	}
	if err := c.App.RunContext(c.Context, []string{c.App.Name, "--machine", machine, "install", "php", "--version", php}); err != nil {
		return err
	}

	// install the database
	if err := c.App.RunContext(c.Context, []string{c.App.Name, "--machine", machine, "install", database}); err != nil {
		return err
	}

	// install redis
	if err := c.App.RunContext(c.Context, []string{c.App.Name, "--machine", machine, "install", "redis"}); err != nil {
		return err
	}

	// login to the machine
	if err := c.App.RunContext(c.Context, []string{c.App.Name, "--machine", machine, "ssh"}); err != nil {
		return err
	}

	return nil
}

// Update will perform system updates on a given machine
func Update(c *cli.Context) error {
	machine := c.String("machine")
	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	fmt.Println("Updating machine:", machine)

	cmd := exec.Command(multipass, "exec", machine, "--", "sudo", "apt-get", "update", "-y")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// SSH will perform system updates on a given machine
func SSH(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Connecting to machine:", machine)

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))
	cmd := exec.Command(multipass, "shell", machine)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func InstallPHP(c *cli.Context) error {
	machine := c.String("machine")
	version := c.String("version")

	phpCmds, err := install.PHP(version)
	if err != nil {
		return err
	}

	fmt.Println("Installing PHP on machine:", machine)

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))
	cmd := exec.Command(multipass, "exec", machine, "--", "sudo", "apt-get", "install", "-y", phpCmds)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func InstallNginx(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Installing Nginx on machine:", machine)

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))
	cmd := exec.Command(multipass, "exec", machine, "--", "sudo", "apt-get", "install", "-y", "nginx")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func InstallMariaDB(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Installing MariaDB on machine:", machine)

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))
	cmd := exec.Command(multipass, "exec", machine, "--", "sudo", "apt-get", "install", "-y", "mariadb-server")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func InstallRedis(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Installing redis on machine:", machine)

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))
	cmd := exec.Command(multipass, "exec", machine, "--", "sudo", "apt-get", "install", "-y", "redis")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func InstallPostgres(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Installing PostgreSQL on machine:", machine)

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))
	cmd := exec.Command(multipass, "exec", machine, "--", "sudo", "apt-get", "install", "-y", "postgresql", "postgresql-contrib")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func Delete(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Deleting machine:", machine)

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))
	cmd := exec.Command(multipass, "delete", machine)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func Stop(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Stopping machine:", machine)

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))
	cmd := exec.Command(multipass, "stop", machine)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
