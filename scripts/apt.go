package scripts

// AptInstall takes a machine name and a list of packages to install onto a machine.
// It will return the complete command to send to multipass
func AptInstall(name string, pkgs []string) []string {
	c := []string{"exec", name, "--", "sudo", "apt", "install", "-y"}
	c = append(c, pkgs...)
	return c
}
