package terminal

import (
	"fmt"
	"strings"
)

type Terminal struct{}

func New() Terminal {
	return Terminal{}
}

func (t *Terminal) Info(s ...string) {
	fmt.Printf("%s\n", strings.Join(s, " "))
}

func (t *Terminal) Success(s ...string) {
	fmt.Printf("  \u2713 %s\n", strings.Join(s, " "))
}

func (t *Terminal) Pending(s ...string) {
	fmt.Printf("  … %s ", strings.Join(s, " "))
}

func (t *Terminal) Done() {
	fmt.Print("\u2713\n")
}
