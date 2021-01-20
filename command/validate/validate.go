package validate

import (
	"fmt"
	"os"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/craftcms/nitro/pkg/validate"
)

var (
	// ErrExample is used when we want to share an error
	ErrExample = fmt.Errorf("some example error")
)

const exampleText = `  # validate a config file
  nitro validate`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "validate",
		Short:   "Validate the config",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			output.Info("Validating…")

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

					// validate the php version
					phpvalidator := validate.PHPVersionValidator{}
					if err := phpvalidator.Validate(s.Version); err != nil {
						siteErrs = append(siteErrs, fmt.Errorf("invalid php version %s", s.Version))
					}
				}

				if len(siteErrs) > 0 {
					output.Warning()
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
