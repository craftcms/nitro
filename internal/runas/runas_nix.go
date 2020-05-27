// +build linux, darwin, !windows

package runas

import (
	"os"
	"os/exec"
)

// Elevated allows the command to be run as sudo user. We
// explicit pass the name of the machine and args that
// we are going to pass to the nitro cli.
// (e.g sudo nitro -m machine-name hosts remove)
func Elevated(machine string, args []string) error {
	nitro, err := os.Executable()
	if err != nil {
		return err
	}

	b := []string{nitro, "-m", machine}
	for _, command := range args {
		b = append(b, command)
	}

	c := exec.Command("sudo", b...)

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return c.Run()
}
