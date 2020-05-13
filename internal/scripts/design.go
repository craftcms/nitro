package scripts

import (
	"os/exec"
	"strings"
)

type Script struct {
	path    string
	machine string
}

func New(multipass, machine string) *Script {
	return &Script{
		path:    multipass,
		machine: machine,
	}
}

// Run is used to make running scripts on a nitro machine
// a lot easier, using New will store the path to the
// nitro path and machine name. Run will then run
// the script on the machine and
func (s Script) Run(arg ...string) (string, error) {
	args := []string{"exec", s.machine, "--", "bash", "-c"}
	args = append(args, arg...)

	cmd := exec.Command(s.path, args...)

	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	output := strings.TrimSpace(string(bytes))

	return output, nil
}
