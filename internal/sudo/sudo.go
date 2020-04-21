package sudo

import (
	"os"
	"os/exec"
)

func RunCommand(nitro, configFile, command string) error {
	// TODO update this for windows support
	c := exec.Command("sudo", nitro, "-f", configFile, command)

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return c.Run()
}
