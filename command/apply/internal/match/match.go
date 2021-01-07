package match

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"

	"github.com/craftcms/nitro/pkg/config"
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
	return checkEnvs(site.PHP, site.Xdebug, container.Config.Env)
}

func checkEnvs(php config.PHP, xdebug bool, envs []string) bool {
	// check the environment variables
	for _, e := range envs {
		sp := strings.Split(e, "=")

		// show only the environment variables we know about/support
		if _, ok := config.DefaultEnvs[sp[0]]; ok {
			env := sp[0]
			val := sp[1]

			// check the value of each environment variable we want to ensure the php config is not the "default" value and that the
			// current value from the container match
			switch env {
			case "PHP_DISPLAY_ERRORS":
				// if there is a custom value
				if !php.DisplayErrors && val != config.DefaultEnvs[env] {
					return false
				}
			case "PHP_MEMORY_LIMIT":
				if (php.MemoryLimit == "" && val != config.DefaultEnvs[env]) || (php.MemoryLimit != "" && val != php.MemoryLimit) {
					return false
				}
			case "PHP_MAX_EXECUTION_TIME":
				if (php.MaxExecutionTime == 0 && val != config.DefaultEnvs[env]) || (php.MaxExecutionTime != 0 && val != strconv.Itoa(php.MaxExecutionTime)) {
					return false
				}
			case "PHP_UPLOAD_MAX_FILESIZE":
				if (php.MaxFileUpload == "" && val != config.DefaultEnvs[env]) || (php.MaxFileUpload != "" && val != php.MaxFileUpload) {
					return false
				}
			case "PHP_MAX_INPUT_VARS":
				if (php.MaxInputVars == 0 && val != config.DefaultEnvs[env]) || (php.MaxInputVars != 0 && val != strconv.Itoa(php.MaxInputVars)) {
					return false
				}
			case "PHP_POST_MAX_SIZE":
				if (php.PostMaxSize == "" && val != config.DefaultEnvs[env]) || (php.PostMaxSize != "" && val != php.PostMaxSize) {
					return false
				}
			case "PHP_OPCACHE_ENABLE":
				if php.OpcacheEnable && val == config.DefaultEnvs[env] {
					return false
				}
			case "PHP_OPCACHE_REVALIDATE_FREQ":
				if (php.OpcacheRevalidateFreq == 0 && val != config.DefaultEnvs[env]) || (php.OpcacheRevalidateFreq != 0 && val != strconv.Itoa(php.OpcacheRevalidateFreq)) {
					return false
				}
			case "XDEBUG_MODE":
				if xdebug && val == config.DefaultEnvs[env] {
					return false
				}

				if !xdebug && val != config.DefaultEnvs[env] {
					return false
				}
			}
		}
	}

	return true
}
