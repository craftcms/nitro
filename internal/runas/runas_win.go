// +build !linux, !darwin, windows

package runas

// We require Windows users to run the CLI as an administrator
func Elevated(machine string, args []string) error {
	return nil
}
