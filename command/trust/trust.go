package trust

import (
	"fmt"
	"io/ioutil"
	"os"
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

			commands := []string{"less", "/data/caddy/pki/authorities/local/root.crt"}

			// create an exec for the container
			exec, err := docker.ContainerExecCreate(ctx, containerID, types.ExecConfig{
				AttachStderr: true,
				AttachStdin:  true,
				AttachStdout: true,
				Cmd:          commands,
			})
			if err != nil {
				return fmt.Errorf("unable to create an execution for container, %w", err)
			}

			// attach to the container
			stream, err := docker.ContainerExecAttach(ctx, exec.ID, types.ExecConfig{
				AttachStdout: true,
				AttachStderr: true,
				AttachStdin:  true,
				Cmd:          commands,
			})
			if err != nil {
				return fmt.Errorf("unable to attach to container, %w", err)
			}
			defer stream.Close()

			// read the stream content
			bytes, err := ioutil.ReadAll(stream.Reader)
			if err != nil || len(bytes) == 0 {
				return fmt.Errorf("unable to read the content from container, %w", err)
			}

			// remove special characters from the output
			var stop int
			for i, s := range bytes {
				if s != 0 && s != 1 {
					stop = i + 1
					break
				}
			}

			// create a temp file
			f, err := ioutil.TempFile(os.TempDir(), "nitro-cert")
			if err != nil {
				return fmt.Errorf("unable to create a temporary file, %w", err)
			}

			// write the certificate to the temporary file
			if _, err := f.Write(bytes[stop+1:]); err != nil {
				return fmt.Errorf("unable to write the certificate to the temporary file, %w", err)
			}
			defer f.Close()

			output.Done()

			output.Info("Installing certificate (you might be prompted for your password)")
			switch runtime.GOOS {
			case "linux":
				// TODO(jasonmccallister) run multiple commands to set permissions
				if err := sudo.Run("cp", "cp", "/usr/local/share/ca-certificates/", f.Name()); err != nil {
					output.Info("Unable to automatically add certificate\n")
					output.Info("To install the certificate, run the following command:")
					output.Info(fmt.Sprintf("  sudo cp %s /usr/local/share/ca-certificates/", f.Name()))
					output.Info("  sudo sudo update-ca-certificates")

					return nil
				}
			// linux
			default:
				if err := sudo.Run("security", "security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", "/Library/Keychains/System.keychain", f.Name()); err != nil {
					output.Info("Unable to automatically add certificate\n")
					output.Info("To install the certificate, run the following command:")
					output.Info(fmt.Sprintf("  sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain %s", f.Name()))
					return nil
				}
			}

			// clean up
			output.Pending("cleaning up")

			if err := os.Remove(f.Name()); err != nil {
				output.Info("unable to remove temporary file, it will be automatically removed on reboot")
			}

			output.Done()

			output.Info(env, "certificates are now trusted 🔒")

			return nil
		},
	}

	return cmd
}
