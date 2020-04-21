package sudo

import (
	"os"
	"os/exec"
)

func RunCommand(nitro, configFile, command string) error {
	hostsCmd := exec.Command("sudo", nitro, "-f", configFile, command)

	hostsCmd.Stdout = os.Stdout
	hostsCmd.Stderr = os.Stderr

	return hostsCmd.Run()
}
