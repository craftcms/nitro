package match

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/craftcms/nitro/pkg/bindmounts"
	"github.com/docker/docker/api/types"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
)

// DEFAULT_MOUNTS is the number of mounts for a standard container
const DEFAULT_MOUNTS = 3

var (
	ErrMisMatchedImage  = fmt.Errorf("container image does not match")
	ErrMisMatchedLabel  = fmt.Errorf("container label does not match")
	ErrEnvFileNotFound  = fmt.Errorf("unable to find the containers env file")
	ErrMisMatchedEnvVar = fmt.Errorf("container environment variables do not match")
	// SiteImage is the image used for sites, with the PHP version
	SiteImage = "docker.io/craftcms/nitro:%s"
)

// Container checks if a custom container is up-to-date with the configuration
func Container(home string, container config.Container, details types.ContainerJSON) error {
	// check if the image does not match - this uses the image name, not ref
	if fmt.Sprintf("%s:%s", container.Image, container.Tag) != details.Config.Image {
		return ErrMisMatchedImage
	}

	// check the name has been changed
	if details.Config.Labels[containerlabels.NitroContainer] != container.Name {
		return ErrMisMatchedLabel
	}

	if container.EnvFile != "" {
		customEnvs := make(map[string]string)

		content, err := ioutil.ReadFile(filepath.Join(home, config.DirectoryName, "."+container.Name))
		if err != nil {
			return ErrEnvFileNotFound
		}

		for _, line := range strings.Split(string(content), "\n") {
			parts := strings.Split(line, "=")
			if len(parts) > 2 {
				customEnvs[parts[0]] = parts[1]
			}
		}

		// check the containers' env against the file and merge
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
// match what's expected.
func Site(home string, site config.Site, container types.ContainerJSON, blackfire config.Blackfire) bool {
	// check if nitro development is defined and override the image
	if _, ok := os.LookupEnv("NITRO_DEVELOPMENT"); ok {
		SiteImage = "craftcms/nitro:%s"
	}

	// check if the image does not match - this uses the image name, not ref
	if fmt.Sprintf(SiteImage, site.Version) != container.Config.Image {
		return false
	}

	// check the web root is defined and they match
	if container.Config.Labels[containerlabels.Webroot] != site.Webroot {
		return false
	}

	// check the sites' hostname using the label
	if container.Config.Labels[containerlabels.Host] != site.Hostname {
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

	// check the bind mounts for the site
	mounts, err := bindmounts.ForSite(site, home)
	if err != nil {
		return false
	}

	// if there are more than 1 mount - the site is using excludes
	switch len(mounts) == 1 {
	case true:
		for _, mount := range container.Mounts {
			// make sure the only bind mount matches the path
			if (mount.Type == "bind") && path != mount.Source {
				return false
			}
		}
	default:
		// check the number of binds matches the number of container binds (we exclude the user home, certs, and nginx configs since they are volumes)
		if len(mounts) != len(container.Mounts)-DEFAULT_MOUNTS {
			return false
		}
	}

	// TODO(jasonmccallister) check the labels for php extensions and write tests
	switch len(site.Extensions) > 0 {
	case false:
		if container.Config.Labels[containerlabels.Extensions] != "" {
			return false
		}
	default:
		if container.Config.Labels[containerlabels.Extensions] != strings.Join(site.Extensions, ",") {
			return false
		}
	}

	// run the final check on the environment variables
	return checkEnvs(site, blackfire, container.Config.Env)
}

func checkEnvs(site config.Site, blackfire config.Blackfire, envs []string) bool {
	// check the environment variables
	for _, e := range envs {
		sp := strings.Split(e, "=")
		env := sp[0]
		val := sp[1]

		// TODO(jasonmccallister) consider adding checks for if blackfire is
		// enabled for this site
		if env == "BLACKFIRE_SERVER_ID" && blackfire.ServerID != val {
			return false
		}
		if env == "BLACKFIRE_SERVER_TOKEN" && blackfire.ServerToken != val {
			return false
		}

		// show only the environment variables we know about/support
		if _, ok := config.DefaultEnvs[env]; ok {
			// check the value of each environment variable we want to ensure the php config is not the "default" value and that the
			// current value from the container match
			switch env {
			case "PHP_DISPLAY_ERRORS":
				// if there is a custom value
				if !site.PHP.DisplayErrors && val != config.DefaultEnvs[env] {
					return false
				}
			case "PHP_MEMORY_LIMIT":
				if (site.PHP.MemoryLimit == "" && val != config.DefaultEnvs[env]) || (site.PHP.MemoryLimit != "" && val != site.PHP.MemoryLimit) {
					return false
				}
			case "PHP_MAX_EXECUTION_TIME":
				if (site.PHP.MaxExecutionTime == 0 && val != config.DefaultEnvs[env]) || (site.PHP.MaxExecutionTime != 0 && val != strconv.Itoa(site.PHP.MaxExecutionTime)) {
					return false
				}
			case "PHP_UPLOAD_MAX_FILESIZE":
				if (site.PHP.MaxFileUpload == "" && val != config.DefaultEnvs[env]) || (site.PHP.MaxFileUpload != "" && val != site.PHP.MaxFileUpload) {
					return false
				}
			case "PHP_MAX_INPUT_VARS":
				if (site.PHP.MaxInputVars == 0 && val != config.DefaultEnvs[env]) || (site.PHP.MaxInputVars != 0 && val != strconv.Itoa(site.PHP.MaxInputVars)) {
					return false
				}
			case "PHP_POST_MAX_SIZE":
				if (site.PHP.PostMaxSize == "" && val != config.DefaultEnvs[env]) || (site.PHP.PostMaxSize != "" && val != site.PHP.PostMaxSize) {
					return false
				}
			case "PHP_OPCACHE_ENABLE":
				if (site.PHP.OpcacheEnable && val == config.DefaultEnvs[env]) || (!site.PHP.OpcacheEnable && val != config.DefaultEnvs[env]) {
					return false
				}
			case "PHP_OPCACHE_REVALIDATE_FREQ":
				if (site.PHP.OpcacheRevalidateFreq == 0 && val != config.DefaultEnvs[env]) || (site.PHP.OpcacheRevalidateFreq != 0 && val != strconv.Itoa(site.PHP.OpcacheRevalidateFreq)) {
					return false
				}
			case "PHP_OPCACHE_VALIDATE_TIMESTAMPS":
				// if there is a custom value
				if !site.PHP.OpcacheValidateTimestamps && val != config.DefaultEnvs[env] {
					return false
				}
			case "XDEBUG_MODE":
				if site.Xdebug && val == config.DefaultEnvs[env] {
					return false
				}

				if !site.Xdebug && val != config.DefaultEnvs[env] {
					return false
				}
			}
		}
	}

	return true
}
