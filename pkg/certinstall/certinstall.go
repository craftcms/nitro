package certinstall

import (
	"bytes"
	"fmt"
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
		dist, err := findLinuxDistribution(buf.String())
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
	default:
		// add the certificate to the macOS keychain
		if err := sudo.Run("security", "security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", "/Library/Keychains/System.keychain", file); err != nil {
			return nil
		}
	}

	return nil
}

func findLinuxDistribution(description string) (string, error) {
	if strings.Contains(description, "Manjaro") {
		return "arch", nil
	}

	if strings.Contains(description, "Ubuntu") {
		return "debian", nil
	}

	return "", fmt.Errorf("unable to find the distribution from the description: %s", description)
}
