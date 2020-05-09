// +build !linux, !darwin, windows

package runas

import (
	"golang.org/x/sys/windows"
	"os"
	"strings"
	"syscall"
)

func Elevated(args []string) error {
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
