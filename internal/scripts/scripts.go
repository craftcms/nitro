package scripts

import (
	"os/exec"
	"strings"
)

// Run is a process that takes the path to multipass and args that
// should be sent to the multipass CLI and return the cmd output
// or an error. Output will auto trim spaces (e.g. new lines)
func Run(multipass string, args []string) (string, error) {
	cmd := exec.Command(multipass, args...)

	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	output := strings.TrimSpace(string(bytes))

	return output, nil
}
