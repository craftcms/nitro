package sudo

import (
	"os"
	"os/exec"
)

func RunCommand(nitro, machine, command string) error {
	// TODO update this for windows support
	c := exec.Command("sudo", nitro, "-m", machine, command)

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return c.Run()
}
