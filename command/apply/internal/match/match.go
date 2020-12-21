package match

import (
	"fmt"
	"os"
	"strings"

	"github.com/craftcms/nitro/config"
	"github.com/docker/docker/api/types"
)

// Site takes the home directory, site, and a container to determine if they
// match whats expected.
func Site(home string, site config.Site, container types.ContainerJSON) bool {
	// check if the image does not match
	if fmt.Sprintf("docker.io/craftcms/nginx:%s-dev", site.Version) != container.Config.Image {
		return false
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

	// check the environment variables
	for _, e := range container.Config.Env {
		sp := strings.Split(e, "=")

		// show only the environment variables we know about/support
		if _, ok := config.Envs[sp[0]]; ok {
			env := sp[0]
			val := sp[1]

			// check the value of each environment variable
			switch env {
			case "PHP_DISPLAY_ERRORS":
				// if there is a custom value
				if site.PHP.DisplayErrors {
					if val != "on" {
						fmt.Println("display errors is not on")
						return false
					}
				}
			case "PHP_MAX_EXECUTION_TIME":
				if site.PHP.MaxExecutionTime != 0 && val != "5000" {
					return false
				}
			}
		}
	}

	return true
}
