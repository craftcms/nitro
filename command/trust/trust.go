package trust

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/certinstall"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	// ErrNoContainers is returned when no containers are running for an environment
	ErrNoContainers = fmt.Errorf("there are no running containers")
)

const (
	certificatePath = "/data/caddy/pki/authorities/local/root.crt"
	exampleText     = `  # get the root certificate for the proxy
  nitro trust`
)

// NewCommand returns `trust` to retrieve the certificates from the nitro proxy and install on the
// host machine. The CA is used to sign certificates for websites and adding the certificate
// to the system allows TLS connections to be considered valid and trusted from the container.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "trust",
		Short:   "Trust certificates for environment",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				// when we call commands from other commands (e.g. create)
				// the context could be nil, so we set it to the parent
				// context just in case.
				ctx = cmd.Parent().Context()
			}

			// find the nitro proxy for the environment
			filter := filters.NewArgs()
			filter.Add("label", labels.Nitro)
			filter.Add("label", labels.Proxy+"=true")

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
			output.Pending("getting Nitro’s root site certificate…")

			// verify the file exists in the container
			for {
				stat, err := docker.ContainerStatPath(ctx, containerID, certificatePath)
				if err != nil {
					continue
				}

				if stat.Name != "" {
					break
				}
			}

			// copy the file from the container
			rdr, stat, err := docker.CopyFromContainer(ctx, containerID, certificatePath)
			if err != nil || !stat.Mode.IsRegular() {
				output.Warning()
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

				if _, err := buf.ReadFrom(tr); err != nil {
					return err
				}
			}

			// if we are only outputting the certificate to stdout
			if cmd.Flag("output-only").Value.String() == "true" {
				output.Done()

				output.Info(buf.String())

				return nil
			}

			// create a temp file
			temp, err := ioutil.TempFile(os.TempDir(), "nitro-local-root-ca")
			if err != nil {
				return fmt.Errorf("unable to create a temporary file, %w", err)
			}
			defer temp.Close()

			// write the certificate to the temporary file
			if _, err := temp.Write(buf.Bytes()); err != nil {
				return fmt.Errorf("unable to write the certificate to the temporary file, %w", err)
			}

			output.Done()

			// copy the certificate into the nitro dir
			cert, err := os.Create(filepath.Join(home, ".nitro", "nitro.crt"))
			if err != nil {
				return err
			}
			defer cert.Close()

			// copy the contents
			if _, err := io.Copy(cert, buf); err != nil {
				return err
			}

			output.Info("Installing certificate (you might be prompted for your password)")

			// install the certificate
			if err := certinstall.Install(temp.Name(), runtime.GOOS); err != nil {
				return err
			}

			output.Info("Nitro certificates are now trusted 🔒")

			return nil
		},
	}

	cmd.Flags().Bool("output-only", false, "show the certificate without importing")

	return cmd
}
