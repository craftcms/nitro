package certinstall

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/craftcms/nitro/pkg/sudo"
)

var (
	certificatePaths = map[string]string{
		"arch":   "/etc/ca-certificates/trust-source/anchors",
		"debian": "/usr/local/share/ca-certificates/",
	}
	certificateTools = map[string]string{
		"arch":   "update-ca-trust",
		"debian": "update-ca-certificates",
	}
)

// Install is responsible for taking a path to a root certificate and the runtime.GOOS as the system
// and finding the distribution and tools to install a root certificate.
func Install(file, system string) error {
	switch system {
	case "linux":
		// find the release tool
		lsb, err := exec.LookPath("lsb_release")
		if err != nil || lsb == "" {
			return fmt.Errorf("unable to find the lsb_release path: %w", err)
		}

		// setup the command
		cmd := exec.Command(lsb, "--description")

		// capture the output into a temp file
		buf := bytes.NewBufferString("")
		cmd.Stdout = buf

		if err := cmd.Start(); err != nil {
			return err
		}

		if err := cmd.Wait(); err != nil {
			return err
		}

		// find the linux distro
		dist, err := findDistro(buf.String())
		if err != nil {
			return err
		}

		// get the certpath
		certPath, ok := certificatePaths[dist]
		if !ok {
			return fmt.Errorf("unable to find the certificate path for %s", dist)
		}

		// get the cert tool
		certTool, ok := certificateTools[dist]
		if !ok {
			return fmt.Errorf("unable to find the certificate tool for %s", dist)
		}

		if err := sudo.Run("mv", "mv", file, fmt.Sprintf("%s%s.crt", certPath, "nitro")); err != nil {
			return err
		}

		// update the ca certs
		if err := sudo.Run(certTool, certTool); err != nil {
			return err
		}

		// is this a wsl machine?
		if dist, exists := os.LookupEnv("WSL_DISTRO_NAME"); exists {
			user := os.Getenv("USER")
			fmt.Println("Users on WSL will need to open an elevated (run as administrator) Command Prompt or terminal on Windows and run the following command:")
			fmt.Printf(`certutil -addstore -f "Root" \\wsl$\%s\home\%s\.nitro\nitro.crt\n`, dist, user)
		}
	default:
		// add the certificate to the macOS keychain
		if err := sudo.Run("security", "security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", "/Library/Keychains/System.keychain", file); err != nil {
			return nil
		}
	}

	return nil
}

func findDistro(description string) (string, error) {
	// detect arch systems
	if strings.Contains(description, "Manjaro") || strings.Contains(description, "Arch Linux") {
		return "arch", nil
	}

	// detect debian systems
	if strings.Contains(description, "Ubuntu") || strings.Contains(description, "Pop!_OS") {
		return "debian", nil
	}

	return "", fmt.Errorf("unable to find the distribution from the description: %s", description)
}
