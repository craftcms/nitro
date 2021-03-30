package iniset

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/craftcms/nitro/pkg/validate"
)

var (
	// ErrUnknownSetting is used when an unknown service is requested
	ErrUnknownSetting = fmt.Errorf("unknown setting requested")
)

const exampleText = `  # change PHP settings for a site
  nitro iniset`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "iniset",
		Short:   "Change PHP setting",
		Example: exampleText,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			cfg, err := config.Load(home)
			if err != nil {
				return nil, cobra.ShellCompDirectiveDefault
			}

			var options []string
			for _, s := range cfg.Sites {
				options = append(options, s.Hostname)
			}

			return options, cobra.ShellCompDirectiveDefault
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return prompt.RunApply(cmd, args, false, output)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			// load the configuration
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

			// get all of the sites
			var sites, found []string
			for _, s := range cfg.Sites {
				p, _ := s.GetAbsPath(home)

				// check if the path matches a sites path, then we are in a known site
				if strings.Contains(wd, p) {
					found = append(found, s.Hostname)
				}

				// add the site to the list in case we cannot find the directory
				sites = append(sites, s.Hostname)
			}

			// if there are found sites we want to show or connect to the first one, otherwise prompt for
			// which site to connect to.
			switch len(found) {
			case 0:
				// prompt for the site to ssh into
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", sites)
				if err != nil {
					return err
				}

				// add the label to get the site
				filter.Add("label", containerlabels.Host+"="+sites[selected])
			case 1:
				// add the label to get the site
				filter.Add("label", containerlabels.Host+"="+found[0])
			default:
				// prompt for the site to ssh into
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", found)
				if err != nil {
					return err
				}

				// add the label to get the site
				filter.Add("label", containerlabels.Host+"="+found[selected])
			}

			// find the containers but limited to the site label
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter, All: true})
			if err != nil {
				return err
			}

			// are there any containers??
			if len(containers) == 0 {
				return fmt.Errorf("unable to find an matching site")
			}

			// start the container if its not running
			if containers[0].State != "running" {
				for _, command := range cmd.Root().Commands() {
					if command.Use == "start" {
						if err := command.RunE(cmd, []string{}); err != nil {
							return err
						}
					}
				}
			}

			// set the hostname of the site based on the container name
			hostname := strings.TrimLeft(containers[0].Names[0], "/")

			settings := []string{
				"display_errors",
				"max_execution_time",
				"max_input_vars",
				"max_input_time",
				"max_file_upload",
				"memory_limit",
				"opcache_enable",
				"opcache_revalidate_freq",
				"opcache_validate_timestamps",
				"post_max_size",
				"upload_max_file_size",
			}

			// which setting to change
			selected, err := output.Select(cmd.InOrStdin(), "Which PHP setting would you like to change for "+hostname+"?", settings)
			if err != nil {
				return err
			}

			// get the specific setting to change
			setting := settings[selected]

			// prompt the user for the setting to change
			switch setting {
			case "display_errors":
				value, err := output.Ask("Should we display PHP errors", "true", "?", &validate.IsBoolean{})
				if err != nil {
					return err
				}

				// convert to bool
				display, err := strconv.ParseBool(value)
				if err != nil {
					return err
				}

				// change the value because its validated
				if err := cfg.SetPHPBoolSetting(hostname, setting, display); err != nil {
					return err
				}
			case "max_execution_time":
				value, err := output.Ask("What should the max execution time be", config.DefaultEnvs["PHP_MAX_EXECUTION_TIME"], "?", &validate.MaxExecutionTime{})
				if err != nil {
					return err
				}

				v, err := strconv.Atoi(value)
				if err != nil {
					return err
				}

				// change the value because its validated
				if err := cfg.SetPHPIntSetting(hostname, setting, v); err != nil {
					return err
				}
			case "max_input_vars":
				value, err := output.Ask("What should the max input vars be", config.DefaultEnvs["PHP_MAX_INPUT_VARS"], "?", &validate.MaxExecutionTime{})
				if err != nil {
					return err
				}

				v, err := strconv.Atoi(value)
				if err != nil {
					return err
				}

				// change the value because its validated
				if err := cfg.SetPHPIntSetting(hostname, setting, v); err != nil {
					return err
				}
			case "max_input_time":
				value, err := output.Ask("What should the max input time be", config.DefaultEnvs["PHP_MAX_INPUT_TIME"], "?", &validate.MaxExecutionTime{})
				if err != nil {
					return err
				}

				v, err := strconv.Atoi(value)
				if err != nil {
					return err
				}

				// change the value because its validated
				if err := cfg.SetPHPIntSetting(hostname, setting, v); err != nil {
					return err
				}
			case "max_file_upload":
				value, err := output.Ask("What should the new max file upload be", config.DefaultEnvs["PHP_UPLOAD_MAX_FILESIZE"], "?", &validate.IsMegabyte{})
				if err != nil {
					return err
				}

				// change the value because its validated
				if err := cfg.SetPHPStrSetting(hostname, setting, value); err != nil {
					return err
				}
			case "memory_limit":
				value, err := output.Ask("What should the new memory limit be", config.DefaultEnvs["PHP_MEMORY_LIMIT"], "?", &validate.IsMegabyte{})
				if err != nil {
					return err
				}

				// change the value because its validated
				if err := cfg.SetPHPStrSetting(hostname, setting, value); err != nil {
					return err
				}
			case "opcache_enable":
				value, err := output.Ask("Should we enable OPCache", "false", "?", &validate.IsBoolean{})
				if err != nil {
					return err
				}

				// convert to bool
				display, err := strconv.ParseBool(value)
				if err != nil {
					return err
				}

				// change the value because its validated
				if err := cfg.SetPHPBoolSetting(hostname, setting, display); err != nil {
					return err
				}
			case "opcache_validate_timestamps":
				value, err := output.Ask("Should we validate timestamps with OPCache", "false", "?", &validate.IsBoolean{})
				if err != nil {
					return err
				}

				// convert to bool
				display, err := strconv.ParseBool(value)
				if err != nil {
					return err
				}

				// change the value because its validated
				if err := cfg.SetPHPBoolSetting(hostname, setting, display); err != nil {
					return err
				}
			case "opcache_revalidate_freq":
				value, err := output.Ask("What should the OPCache revalidate frequency be", config.DefaultEnvs["PHP_OPCACHE_REVALIDATE_FREQ"], "?", &validate.MaxExecutionTime{})
				if err != nil {
					return err
				}

				v, err := strconv.Atoi(value)
				if err != nil {
					return err
				}

				// change the value because its validated
				if err := cfg.SetPHPIntSetting(hostname, setting, v); err != nil {
					return err
				}
			case "post_max_size":
				value, err := output.Ask("What should post max size be", config.DefaultEnvs["PHP_POST_MAX_SIZE"], "?", &validate.IsMegabyte{})
				if err != nil {
					return err
				}

				// change the value because its validated
				if err := cfg.SetPHPStrSetting(hostname, setting, value); err != nil {
					return err
				}
			case "upload_max_file_size":
				value, err := output.Ask("What should upload maximum file size be", config.DefaultEnvs["PHP_UPLOAD_MAX_FILESIZE"], "?", &validate.IsMegabyte{})
				if err != nil {
					return err
				}

				// change the value because its validated
				if err := cfg.SetPHPStrSetting(hostname, setting, value); err != nil {
					return err
				}
			default:
				return ErrUnknownSetting
			}

			// save the config file
			if err := cfg.Save(); err != nil {
				return fmt.Errorf("unable to save config, %w", err)
			}

			return nil
		},
	}

	return cmd
}
