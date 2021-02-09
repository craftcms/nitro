package match

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/labels"
)

var (
	ErrMisMatchedImage  = fmt.Errorf("container image does not match")
	ErrMisMatchedLabel  = fmt.Errorf("container label does not match")
	ErrEnvFileNotFound  = fmt.Errorf("unable to find the containers env file")
	ErrMisMatchedEnvVar = fmt.Errorf("container environment variables do not match")
)

// Container checks if a custom container is up to date with the configuration
func Container(home string, container config.Container, details types.ContainerJSON) error {
	// check if the image does not match - this uses the image name, not ref
	if fmt.Sprintf("%s:%s", container.Image, container.Tag) != details.Config.Image {
		return ErrMisMatchedImage
	}

	// check the name has been changed
	if details.Config.Labels[labels.NitroContainer] != container.Name {
		return ErrMisMatchedLabel
	}

	if container.EnvFile != "" {
		customEnvs := make(map[string]string)

		content, err := ioutil.ReadFile(filepath.Join(home, ".nitro", "."+container.Name))
		if err != nil {
			return ErrEnvFileNotFound
		}

		for _, line := range strings.Split(string(content), "\n") {
			parts := strings.Split(line, "=")
			customEnvs[parts[0]] = parts[1]
		}

		// check the containers env against the file and merge
		for _, e := range details.Config.Env {
			parts := strings.Split(e, "=")
			env := parts[0]
			val := parts[1]

			// is there a custom env val for the variable?
			if custom, ok := customEnvs[env]; ok {
				if val != custom {
					return ErrMisMatchedEnvVar
				}
			}
		}
	}

	// TODO(jasonmccallister) check the port mappings
	// TODO(jasonmccallister) check the volumes

	return nil
}

// Site takes the home directory, site, and a container to determine if they
// match whats expected.
func Site(home string, site config.Site, container types.ContainerJSON, blackfire config.Blackfire) bool {
	// check if the image does not match - this uses the image name, not ref
	if fmt.Sprintf("docker.io/craftcms/nginx:%s-dev", site.Version) != container.Config.Image {
		return false
	}

	// check the sites hostname using the label
	if container.Config.Labels[labels.Host] != site.Hostname {
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

	// check the path
	if len(container.Mounts) > 0 {
		if path != container.Mounts[0].Source {
			return false
		}
	}

	// TODO(jasonmccallister) check the labels for php extensions and write tests
	switch len(site.Extensions) > 0 {
	case false:
		if container.Config.Labels[labels.Extensions] != "" {
			return false
		}
	default:
		if container.Config.Labels[labels.Extensions] != strings.Join(site.Extensions, ",") {
			return false
		}
	}

	// run the final check on the environment variables
	return checkEnvs(site.PHP, site.Xdebug, container.Config.Env, blackfire)
}

func checkEnvs(php config.PHP, xdebug bool, envs []string, blackfire config.Blackfire) bool {
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
			case "BLACKFIRE_SERVER_ID":
				if blackfire.ServerID != val {
					return false
				}
			case "BLACKFIRE_SERVER_TOKEN":
				if blackfire.ServerToken != val {
					return false
				}
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
