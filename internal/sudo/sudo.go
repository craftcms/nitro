package sudo

import (
	"os"
	"os/exec"
)

func RunCommand(nitro, machine string, commands ...string) error {
	// TODO update this for windows support

	b := []string{nitro, "-m", machine}
	for _, command := range commands {
		b = append(b, command)
	}

	c := exec.Command("sudo", b...)

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return c.Run()
}
