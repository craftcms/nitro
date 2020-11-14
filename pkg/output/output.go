package output

import (
	"github.com/fatih/color"
)

// Outputer is used to style terminal output in the CLI tool. It is also used to
// help keep indentation easier by allowing nesting error and info output using
// SubError and SubInfo.
type Outputer interface {
	Error(a ...interface{})
	Info(a ...interface{})
	SubError(a ...interface{})
	SubInfo(a ...interface{})
}

type output struct {
	printErr  func(a ...interface{})
	printInfo func(a ...interface{})
}

func (o output) Error(a ...interface{}) {
	o.printErr(a...)
}

func (o output) Info(a ...interface{}) {
	o.printInfo(a...)
}

func (o output) SubInfo(a ...interface{}) {
	o.printInfo("  ==>", a)
}

func (o output) SubError(a ...interface{}) {
	o.printErr("  ==>", a)
}

func New() Outputer {
	out := output{}

	out.printErr = color.New(color.FgRed, color.Bold).PrintlnFunc()
	out.printInfo = color.New(color.FgCyan, color.Bold).PrintlnFunc()

	return out
}
