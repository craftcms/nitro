package action

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/dev/install"
)

// Initialize is used to create a new machine and setup any dependencies
func Initialize(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Creating a new machine:", machine)
	multipass, err := exec.LookPath("multipass")
	if err != nil {
		fmt.Println(err)
		return err
	}

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

	multipass, err := exec.LookPath("multipass")
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Updating machine:", machine)

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

	multipass, err := exec.LookPath("multipass")
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Connecting to machine:", machine)

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
	version := c.String("php")

	multipass, err := exec.LookPath("multipass")
	if err != nil {
		fmt.Println(err)
		return err
	}

	phpArgs, cmdErr := install.PHP(version)
	if cmdErr != nil {
		fmt.Println(err)
		return err
	}

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
