// +build !linux, !darwin, windows

package runas

import (
	"os"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

// Elevated allows the command to be run as a administrator
// user. We explicit pass the path to the nitro binary, the name
// of the machine, and args that we are going to use
// (e.g runas nitro -m machine-name hosts remove)
func Elevated(nitro, machine, string, args []string) error {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(strings.Join(args, " "))

	var showCmd int32 = 1 //SW_NORMAL

	return windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
}
