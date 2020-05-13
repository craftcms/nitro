package scripts

import (
	"os/exec"
	"strings"
)

const (
	FmtNginxSiteAvailable = `if test -f '/etc/nginx/sites-available/%s'; then echo 'exists'; fi`
	FmtNginxSiteEnabled   = `if test -f '/etc/nginx/sites-enabled/%s'; then echo 'exists'; fi`
	FmtNginxSiteWebroot   = `grep "root " /etc/nginx/sites-available/%s | while read -r line; do echo "$line"; done`
)

type Script struct {
	path    string
	machine string
}

// New will return a new Script struct that
// contains the path to multipass and
// the name of the machine
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
