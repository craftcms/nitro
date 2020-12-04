package validate

import (
	"fmt"
	"os"

	"github.com/craftcms/nitro/config"
	validator "github.com/craftcms/nitro/internal/validate"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var (
	// ErrExample is used when we want to share an error
	ErrExample = fmt.Errorf("some example error")
)

const exampleText = `  # validate a config file
  nitro validate`

// New
func New(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "validate",
		Short:   "Validate the config",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			env := cmd.Flag("environment").Value.String()
			cfg, err := config.Load(home, env)
			if err != nil {
				return err
			}

			output.Info("Validating...")

			// set errors
			var siteErrs, dbErrs []error

			dbs := cfg.Databases
			sites := cfg.Sites

			if len(dbs) > 0 {
				output.Pending("validating databases")

				for _, d := range dbs {
					if d.Port == "" {
						dbErrs = append(dbErrs, fmt.Errorf("port is not assigned"))
					}
				}

				output.Done()
			}

			// check the site paths
			if len(sites) > 0 {
				output.Pending("validating sites")

				for _, s := range sites {
					// check the site path
					p, err := s.GetAbsPath(home)
					if err != nil {
						siteErrs = append(siteErrs, err)
					}

					if _, err := os.Stat(p); os.IsNotExist(err) {
						siteErrs = append(siteErrs, fmt.Errorf("unable to locate site path %s", p))
					}

					// check if the mounts exist
					mnts, err := s.GetAbsMountPaths(home)
					if err != nil {
						siteErrs = append(siteErrs, err)
					}

					for k, _ := range mnts {
						if _, err := os.Stat(k); os.IsNotExist(err) {
							siteErrs = append(siteErrs, fmt.Errorf("Mount source is missing for %s %q", s.Hostname, k))
						}
					}

					// validate the php version
					if err := validator.PHPVersion(s.PHP); err != nil {
						siteErrs = append(siteErrs, fmt.Errorf("the php version for %s is not valid", s.Hostname))
					}
				}

				if len(siteErrs) > 0 {
					// TODO(jasonmccallister) add a output.Warning()
					output.Info("\u2717")
				} else {
					output.Done()
				}
			}

			// show any errors
			if len(siteErrs) > 0 {
				output.Info("Site Errors:")
				for _, e := range siteErrs {
					output.Info(" \u2610", e.Error())
				}
			}

			return nil
		},
	}

	return cmd
}
