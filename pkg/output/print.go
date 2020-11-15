package output

import (
	"strings"

	"github.com/fatih/color"
)

func Info(s ...string) {
	c := color.New(color.Bold, color.FgGreen)
	c.Printf("%s\n", strings.Join(s, " "))
}

func SubInfo(s ...string) {
	c := color.New(color.FgGreen)
	c.Printf("  ==> %s\n", strings.Join(s, " "))
}

func Error(s ...string) {
	c := color.New(color.Bold, color.FgRed)
	c.Printf("%s\n", strings.Join(s, " "))
}

func SubError(s ...string) {
	c := color.New(color.FgRed)
	c.Printf("  ==> %s\n", strings.Join(s, " "))
}
