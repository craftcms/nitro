// +build darwin, !linux

package certinstall

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/craftcms/nitro/pkg/sudo"
)

// Install is responsible for taking a path to a root certificate and the runtime.GOOS as the system
// and finding the distribution and tools to install a root certificate.
func Install(file, system string) error {
	if err := sudo.Run("security", "security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", "/Library/Keychains/System.keychain", file); err != nil {
		return fmt.Errorf("unable to install the certificate, %w", err)
	}

	return nil
}
