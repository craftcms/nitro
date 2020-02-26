package action

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/nitro/install"
)

// Initialize is used to create a new machine and setup any dependencies
func Initialize(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Creating a new machine:", machine)

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	// create the machine
	cmd := exec.Command(multipass, "launch", "--name", machine, "--cloud-init", "./scripts/cloud-init.yaml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// if we are bootstrapping, call that command
	if c.Bool("bootstrap") {
		return c.App.RunContext(c.Context, []string{c.App.Name, "--machine", c.String("machine"), "bootstrap"})
	}

	return nil
}

func Bootstrap(c *cli.Context, e CommandLineExecutor) error {
	machine := c.String("machine")

	// TODO make this pass a PHP version and database
	args := []string{"multipass", "exec", machine, "--", "sudo", "bash", "/etc/nitro/bootstrap.sh"}
	err := e.Exec(e.Path(), args, os.Environ())
	if err != nil {
		return err
	}

	return nil
}

// Update will perform system updates on a given machine
func Update(c *cli.Context) error {
	machine := c.String("machine")
	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	fmt.Println("Updating machine:", machine)

	cmd := exec.Command(multipass, "exec", machine, "--", "sudo", "apt", "update", "-y")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// SSH will login a user to a specific machine
func SSH(m string, e CommandLineExecutor) error {
	fmt.Println("Connecting to machine:", m)

	args := []string{"multipass", "shell", m}
	err := e.Exec(e.Path(), args, os.Environ())
	if err != nil {
		return err
	}

	return nil
}

func Attach(path, machine string, exec CommandLineExecutor) error {
	// verify the path
	dir, err := os.Stat(path)
	if err != nil {
		return err
	}

	if dir.IsDir() == false {
		return errors.New("path must be a directory")
	}

	// mount the path
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

	cmd := exec.Command(multipass, "exec", machine, "--", "sudo", "apt", "install", "-y", phpCmds)
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
	cmd := exec.Command(multipass, "exec", machine, "--", "sudo", "apt", "install", "-y", "nginx")
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
	args := []string{"multipass", "exec", machine, "--", "sudo", "apt", "install", "-y", "mariadb-server"}
	err := syscall.Exec(multipass, args, os.Environ())
	if err != nil {
		return err
	}

	return nil
}

func InstallRedis(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Installing redis on machine:", machine)

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	args := []string{"multipass", "exec", machine, "--", "sudo", "apt", "install", "-y", "redis"}
	err := syscall.Exec(multipass, args, os.Environ())
	if err != nil {
		return err
	}

	return nil
}

func InstallPostgres(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Installing PostgreSQL on machine:", machine)

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))
	cmd := exec.Command(multipass, "exec", machine, "--", "sudo", "apt", "install", "-y", "postgresql", "postgresql-contrib")
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
