package dockerexec

import (
	"io"
	"os/exec"
)

var File = "docker"

// Connect is a helper package to specifically look up the path to the docker binary and execute
// an exec command into a container.
func Connect(r io.Reader, w io.Writer, user, container, shell string) error {
	p, err := exec.LookPath(File)
	if err != nil {
		return err
	}

	c := exec.Command(p, "exec", "--env", "TERM=xterm-256color", "--user", user, "--interactive", "--tty", container, shell)

	c.Stdin = r
	c.Stderr = w
	c.Stdout = w

	return c.Run()
}

// Exec is a helper package to specifically look up the path to the docker binary and execute
// an exec command into a container.
func Exec(r io.Reader, w io.Writer, user, container string, shell ...string) error {
	p, err := exec.LookPath(File)
	if err != nil {
		return err
	}

	args := []string{"exec", "--env", "TERM=xterm-256color", "--user", user, "--interactive", "--tty", container}
	args = append(args, shell...)

	c := exec.Command(p, args...)

	c.Stdin = r
	c.Stderr = w
	c.Stdout = w

	return c.Run()
}
