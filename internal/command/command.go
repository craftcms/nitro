package command

import (
	"errors"
	"fmt"

	"github.com/pixelandtonic/prompt"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/helpers"
	"github.com/craftcms/nitro/validate"
)

// Add is not implemented abd will be moved to task.Add
func Add(args []string, cfg *config.Config, pmt *prompt.Prompt, foundWebroot, flagHostname, flagWebroot string, flagDebug bool) error {
	dirName, _, err := helpers.GetDirectoryArg(args)
	if err != nil {
		return err
	}

	// set the vars we need for the add
	var hostname string
	var webroot string

	// ask for the hostname
	if flagHostname == "" {
		hostname, err = pmt.Ask("What should the hostname be", &prompt.InputOptions{
			Default:   dirName,
			Validator: validate.Hostname,
		})
	} else {
		hostname = flagHostname
	}

	if flagWebroot == "" {
		// ask for the webroot
		webroot, err = pmt.Ask("What is the webroot", &prompt.InputOptions{
			Default:   foundWebroot,
			Validator: validate.Path,
		})
	} else {

	}

	fmt.Println(hostname)
	fmt.Println(webroot)

	return errors.New("not implemented")

	// create a mount
	// create a site

	if flagDebug {
		return nil
	}

	return cfg.Save(viper.ConfigFileUsed())
}
