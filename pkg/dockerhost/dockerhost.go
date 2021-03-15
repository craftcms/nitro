package dockerhost

import "github.com/craftcms/nitro/pkg/wsl"

const host = "host.docker.internal"

// Determine takes the runtime.GOOS and determines the
// hostname to use to allow Docker to access the host machine.
func Determine(s string) string {
	switch s {
	case "linux":
		if wsl.IsWSL() {
			return host
		}

		return "127.0.0.1"
	}

	return host
}
