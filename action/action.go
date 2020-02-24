package action

import (
	"fmt"
	"os"
	"syscall"

	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/phpdev/install"
)

// Initialize is used to create a new machine and setup any dependencies
func Initialize(c *cli.Context) error {
	machine := c.String("machine")

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	fmt.Println("Creating a new machine:", machine)

	// create the machine
	args := []string{"multipass", "launch", "--name", machine, "--cloud-init", "./scripts/cloud-init.yaml"}
	launchErr := syscall.Exec(multipass, args, os.Environ())
	if launchErr != nil {
		fmt.Println(launchErr)
		return launchErr
	}

	return nil
}

// Update will perform system updates on a given machine
func Update(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Updating machine:", machine)

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))
	args := []string{"multipass", "exec", machine, "--", "sudo", "apt-get", "update", "-y"}
	execErr := syscall.Exec(multipass, args, os.Environ())
	if execErr != nil {
		fmt.Println(execErr)
		return execErr
	}

	return nil
}

// SSH will perform system updates on a given machine
func SSH(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Connecting to machine:", machine)

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))
	args := []string{"multipass", "shell", machine}
	execErr := syscall.Exec(multipass, args, os.Environ())
	if execErr != nil {
		fmt.Println(execErr)
		return execErr
	}

	return nil
}

func InstallPHP(c *cli.Context) error {
	machine := c.String("machine")
	version := c.String("version")

	phpArgs, err := install.PHP(version)
	if err != nil {
		return err
	}

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))
	args := []string{"multipass", "exec", machine, "--", "sudo", "apt-get", "install", "-y"}
	for _, v := range phpArgs {
		args = append(args, v)
	}

	fmt.Println("Installing PHP on machine:", machine)

	execErr := syscall.Exec(multipass, args, os.Environ())
	if execErr != nil {
		fmt.Println(execErr)
		return execErr
	}

	return nil
}

func InstallNginx(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Installing Nginx on machine:", machine)

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))
	args := []string{"multipass", "exec", machine, "--", "sudo", "apt-get", "install", "-y", "nginx"}
	execErr := syscall.Exec(multipass, args, os.Environ())
	if execErr != nil {
		fmt.Println(execErr)
		return execErr
	}

	return nil
}

func InstallMariaDB(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Installing MariaDB on machine:", machine)

	args := []string{"multipass", "exec", machine, "--", "sudo", "apt-get", "install", "-y", "mariadb-server"}
	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))
	execErr := syscall.Exec(multipass, args, os.Environ())
	if execErr != nil {
		fmt.Println(execErr)
		return execErr
	}

	return nil
}

func InstallRedis(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Installing redis on machine:", machine)

	args := []string{"multipass", "exec", machine, "--", "sudo", "apt-get", "install", "-y", "redis"}
	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))
	execErr := syscall.Exec(multipass, args, os.Environ())
	if execErr != nil {
		fmt.Println(execErr)
		return execErr
	}

	return nil
}

func InstallPostgres(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Installing postgres on machine:", machine)

	args := []string{"multipass", "exec", machine, "--", "sudo", "apt-get", "install", "-y", "postgresql", "postgresql-contrib"}
	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))
	execErr := syscall.Exec(multipass, args, os.Environ())
	if execErr != nil {
		fmt.Println(execErr)
		return execErr
	}

	return nil
}

func Delete(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Deleting machine:", machine)

	args := []string{"multipass", "delete", machine}
	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))
	execErr := syscall.Exec(multipass, args, os.Environ())
	if execErr != nil {
		fmt.Println(execErr)
		return execErr
	}

	return nil
}

func Stop(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Stopping machine:", machine)

	args := []string{"multipass", "stop", machine}
	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))
	execErr := syscall.Exec(multipass, args, os.Environ())
	if execErr != nil {
		fmt.Println(execErr)
		return execErr
	}

	return nil
}
