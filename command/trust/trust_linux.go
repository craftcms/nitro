// +build linux, !darwin, !windows

package trust

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

	"github.com/craftcms/nitro/labels"
	"github.com/craftcms/nitro/pkg/sudo"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var (
	// ErrNoContainers is returned when no containers are running for an environment
	ErrNoContainers = fmt.Errorf("there are no running containers")
)

const exampleText = `  # get the certificates for an environment
  nitro trust`

// New returns `trust` to retreive the certificates from the nitro proxy and install on the
// host machine. The CA is used to sign certificates for websites and adding the certificate
// to the system allows TLS connections to be considered valid and trusted from the container.
func New(docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "trust",
		Short:   "Trust certificates for environment",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			env := cmd.Flag("environment").Value.String()
			ctx := cmd.Context()

			// find the nitro proxy for the environment
			filter := filters.NewArgs()
			filter.Add("label", labels.Environment+"="+env)
			filter.Add("label", labels.Proxy+"="+env)

			// find the container, should only be one
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
			if err != nil {
				return fmt.Errorf("unable to get the list of containers, %w", err)
			}

			// make sure there is at least one container
			if len(containers) == 0 {
				return ErrNoContainers
			}

			containerID := containers[0].ID

			// get the contents of the certificate from the container
			output.Pending(fmt.Sprintf("getting certificate for %s...", env))

			// copy the file from the container
			rdr, stat, err := docker.CopyFromContainer(ctx, containerID, "/data/caddy/pki/authorities/local/root.crt")
			if err != nil || stat.Mode.IsRegular() == false { // make sure it is a file not a directory
				return fmt.Errorf("unable to get the certificate from the container, %w", err)
			}

			// the file is in a tar format
			buf := new(bytes.Buffer)
			tr := tar.NewReader(rdr)
			for {
				_, err := tr.Next()
				// if end of tar archive
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}

				buf.ReadFrom(tr)
			}

			// create a temp file
			f, err := ioutil.TempFile(os.TempDir(), "nitro-cert")
			if err != nil {
				return fmt.Errorf("unable to create a temporary file, %w", err)
			}

			// write the certificate to the temporary file
			if _, err := f.Write(buf.Bytes()); err != nil {
				return fmt.Errorf("unable to write the certificate to the temporary file, %w", err)
			}
			defer f.Close()

			output.Done()

			output.Info("Installing certificate (you might be prompted for your password)")
			switch runtime.GOOS {
			case "linux":
				// using the reference from: https://askubuntu.com/questions/645818/how-to-install-certificates-for-command-line
				if err := sudo.Run("mv", "mv", f.Name(), fmt.Sprintf("/usr/local/share/ca-certificates/%s.crt", env)); err != nil {
					output.Info("Unable to automatically add certificate\n")
					output.Info("To install the certificate, run the following command:")
					output.Info(fmt.Sprintf("  sudo mv %s /usr/local/share/ca-certificates/%s.crt", f.Name(), env))
					output.Info("  sudo update-ca-certificates")

					return nil
				}

				// update the ca certs
				if err := sudo.Run("update-ca-certificates", "update-ca-certificates"); err != nil {
					output.Info("Unable to automatically add certificate\n")
					output.Info("To install the certificate, run the following command:")
					output.Info("  sudo update-ca-certificates")

					return nil
				}
				// TODO(jasonmccallister) trust the certificate for firefox and edge on linux (https://javorszky.co.uk/2019/11/06/get-firefox-to-trust-your-self-signed-certificates/)
			// windows
			case "windows":
				// automate the certificate installation from this article: https://superuser.com/a/1506481/215387
				// we cannot assume PowerShell is enabled, so we use certutil.exe
				if err := exec.Command("certutil.exe", "-addstore", "root", f.Name()).Run(); err != nil {
					output.Info("Unable to automatically add certificate\n")
					output.Info("To install the certificate, run the following command:")
					output.Info(fmt.Sprintf("  certutil.exe -addstore root %s", f.Name()))

					return nil
				}
			// macos
			default:
				// add the certificate to the keychain
				if err := sudo.Run("security", "security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", "/Library/Keychains/System.keychain", f.Name()); err != nil {
					output.Info("Unable to automatically add certificate\n")
					output.Info("To install the certificate, run the following command:")
					output.Info(fmt.Sprintf("  sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain %s", f.Name()))

					return nil
				}
			}

			// clean up
			output.Pending("cleaning up")

			// remove the temp file
			if err := os.Remove(f.Name()); err != nil {
				output.Warning()

				output.Info("unable to remove temporary file, it will be automatically removed on reboot")
			} else {
				output.Done()
			}

			output.Info(env, "certificates are now trusted 🔒")

			return nil
		},
	}

	return cmd
}
