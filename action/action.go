package action

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/urfave/cli/v2"
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
	args := []string{"multipass", "launch", "--name", machine}
	launchErr := syscall.Exec(multipass, args, os.Environ())
	if launchErr != nil {
		fmt.Println(launchErr)
		return launchErr
	}

	return nil
}

// Prepare will prepare a machine for development work
func Prepare(c *cli.Context) error {
	machine := c.String("machine")
	php := c.String("php")

	installScript := "./scripts/php" + php + "/install.sh"

	_, err := os.Stat(installScript)
	if os.IsNotExist(err) {
		return errors.New("unable to find the file " + installScript)
	}

	multipass, err := exec.LookPath("multipass")
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Preparing the machine:", machine)

	args := []string{"multipass", "transfer", installScript, machine + ":/tmp/install.sh"}
	execErr := syscall.Exec(multipass, args, os.Environ())
	if execErr != nil {
		fmt.Println(execErr)
		return execErr
	}

	return nil
}

// Build will prepare a machine for development work
func Build(c *cli.Context) error {
	machine := c.String("machine")

	multipass, err := exec.LookPath("multipass")
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Building the machine:", machine)

	args := []string{"multipass", "exec", machine, "bash", "/tmp/install.sh"}
	execErr := syscall.Exec(multipass, args, os.Environ())
	if execErr != nil {
		fmt.Println(execErr)
		return execErr
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
