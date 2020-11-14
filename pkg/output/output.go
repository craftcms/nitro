package output

import (
	"github.com/fatih/color"
)

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

	out.printErr = color.New(color.FgRed).PrintlnFunc()
	out.printInfo = color.New(color.FgCyan, color.Bold).PrintlnFunc()

	return out
}
