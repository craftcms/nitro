package sudo

import (
	"fmt"
	"os"
	"os/exec"
)

// Run takes an executable and a list of arguments and will run the
// command as the sudo user
func Run(e string, args ...string) error {
	p, err := exec.LookPath(e)
	if err != nil {
		return fmt.Errorf("unable to find executable %q, %w", e, err)
	}

	b := []string{p}
	for _, a := range args {
		b = append(b, a)
	}

	cmd := exec.Command("sudo", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
