package terminal

import (
	"fmt"
	"strings"
)

type Outputer interface {
	Info(s ...string)
	Success(s ...string)
	Pending(s ...string)
	Done()
}

type terminal struct{}

func New() Outputer {
	return terminal{}
}

func (t terminal) Info(s ...string) {
	fmt.Printf("%s\n", strings.Join(s, " "))
}

func (t terminal) Success(s ...string) {
	fmt.Printf("  \u2713 %s\n", strings.Join(s, " "))
}

func (t terminal) Pending(s ...string) {
	fmt.Printf("  … %s ", strings.Join(s, " "))
}

func (t terminal) Done() {
	fmt.Print("\u2713\n")
}
