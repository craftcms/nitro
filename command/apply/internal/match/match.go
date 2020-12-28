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
	// check if the image does not match - this uses the image name, not ref
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

	// run the final check on the environment variables
	return checkEnvs(site, container.Config.Env)
}

func checkEnvs(site config.Site, envs []string) bool {
	// check the environment variables
	for _, e := range envs {
		sp := strings.Split(e, "=")

		// show only the environment variables we know about/support
		if _, ok := config.DefaultEnvs[sp[0]]; ok {
			env := sp[0]
			val := sp[1]

			// check the value of each environment variable
			// we want to ensure the site.php.config is not the "default" value and that the
			// default value
			switch env {
			case "PHP_DISPLAY_ERRORS":
				// if there is a custom value
				if site.PHP.DisplayErrors == false && val != config.DefaultEnvs[env] {
					return false
				}
			case "PHP_MEMORY_LIMIT":
				if site.PHP.MemoryLimit != "" && val != config.DefaultEnvs[env] {
					return false
				}
			case "PHP_MAX_EXECUTION_TIME":
				if site.PHP.MaxExecutionTime != 0 && val != config.DefaultEnvs[env] {
					return false
				}
			case "PHP_UPLOAD_MAX_FILESIZE":
				if site.PHP.MaxFileUpload != "" && val != config.DefaultEnvs[env] {
					return false
				}
			case "PHP_MAX_INPUT_VARS":
				if site.PHP.MaxInputVars != 0 && val != config.DefaultEnvs[env] {
					return false
				}
			case "PHP_POST_MAX_SIZE":
				if site.PHP.PostMaxSize != "" && val != config.DefaultEnvs[env] {
					return false
				}
			case "PHP_OPCACHE_ENABLE":
				// TODO(jasonmccallister) verify this logic
				if site.PHP.OpcacheEnable && val == config.DefaultEnvs[env] {
					return false
				}
			case "PHP_OPCACHE_REVALIDATE_FREQ":
				if site.PHP.OpcacheRevalidateFreq != 0 && val != config.DefaultEnvs[env] {
					return false
				}
			}
		}
	}

	return true
}
