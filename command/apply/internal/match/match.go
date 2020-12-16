package match

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/craftcms/nitro/config"
	"github.com/docker/docker/api/types"
)

// Site takes the home directory, site, php config, and a container to determine if they
// match whats expected.
func Site(home string, site config.Site, php config.PHP, container types.ContainerJSON) bool {
	// check if the image does not match
	if fmt.Sprintf("docker.io/craftcms/nginx:%s-dev", site.PHP) != container.Config.Image {
		return false
	}

	// check the environment variables
	for _, e := range container.Config.Env {
		sp := strings.Split(e, "=")

		// show only the environment variables we know about/support
		if _, ok := config.Envs[sp[0]]; ok {
			// TODO(jasonmccallister) check if the value matches the config option
			e := sp[0]
			v := sp[1]

			// check the value of each environment variable
			switch e {
			case "PHP_DISPLAY_ERRORS":
				// if there is a custom value
				if php.DisplayErrors {
					b, err := strconv.ParseBool(v)
					if err != nil {
						return false
					}

					if php.DisplayErrors != b {
						return false
					}
				}
			case "PHP_MAX_EXECUTION_TIME":
				if php.MaxExecutionTime != 0 && v != "5000" {
					return false
				}
			}
		}
	}

	// get the main site path (e.g. ~/dev/craft-dev)
	path, err := site.GetAbsPath(home)
	if err != nil {
		return false
	}

	// check if the path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
