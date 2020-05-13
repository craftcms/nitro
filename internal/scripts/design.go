package scripts

import (
	"os"
	"os/exec"
)

type Script struct {
	cmd *exec.Cmd
}

func (s *Script) Run(args []string) (string, error) {
	s.cmd.Args = args

	return "", nil
}

type Outputer interface {
	Run(args []string) (string, error)
}

func New(multipass string) Script {
	return Script{
		cmd: &exec.Cmd{
			Path:   multipass,
			Stdout: os.Stdout,
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
		},
	}
}
