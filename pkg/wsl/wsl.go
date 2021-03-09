package wsl

import "os"

// IsWSL will check for environment variable and determine if the
// system is a WSL installation.
func IsWSL() bool {
	for _, e := range []string{"WSL_DISTRO_NAME", "WSLENV"} {
		if _, is := os.LookupEnv(e); is {
			return true
		}
	}

	return false
}
