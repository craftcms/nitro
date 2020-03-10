package internal

import (
	"strings"
)

// Runner is responsible for running command, it can defer to syscall.Exec or
// exec.Command where required.
type Runner interface {
	// Run is used when the command does not need to be interactive
	Run(args []string) error
	UseSyscall(t bool)
	SetReader(rdr *strings.Reader)
}
