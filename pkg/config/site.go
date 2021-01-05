package config

import (
	"fmt"
	"path/filepath"
	"strings"
)

// DefaultEnvs is used to map a config to a known environment variable that is used
// on the container instances to their default values
var DefaultEnvs = map[string]string{
	"PHP_DISPLAY_ERRORS":          "on",
	"PHP_MEMORY_LIMIT":            "512M",
	"PHP_MAX_EXECUTION_TIME":      "5000",
	"PHP_UPLOAD_MAX_FILESIZE":     "512M",
	"PHP_MAX_INPUT_VARS":          "5000",
	"PHP_POST_MAX_SIZE":           "512M",
	"PHP_OPCACHE_ENABLE":          "0",
	"PHP_OPCACHE_REVALIDATE_FREQ": "0",
	"XDEBUG_MODE":                 "off",
	"XDEBUG_SESSION":              "PHPSTORM",
	"XDEBUG_CONFIG":               "",
}

// Site represents a web application. It has a hostname, aliases (which
// are alternate domains), the local path to the site, additional mounts
// to add to the container, and the directory the index.php is located.
type Site struct {
	Hostname string   `yaml:"hostname"`
	Aliases  []string `yaml:"aliases,omitempty"`
	Path     string   `yaml:"path"`
	Version  string   `yaml:"version"`
	PHP      PHP      `yaml:"php,omitempty"`
	Dir      string   `yaml:"dir"`
	Xdebug   bool     `yaml:"xdebug"`
}

// GetAbsPath gets the directory for a site.Path,
// It is used to create the mount for a sites
// container.
func (s *Site) GetAbsPath(home string) (string, error) {
	return cleanPath(home, s.Path)
}

// AsEnvs takes a gateway addr and turns specific options
// such as PHP settings into env vars that can be set on the
// containers environment
func (s *Site) AsEnvs(addr string) []string {
	var envs []string

	if addr == "" {
		addr = "host.docker.internal"
	}

	// set the php vars
	envs = append(envs, phpVars(s.PHP, s.Version)...)

	// get the xdebug vars
	envs = append(envs, xdebugVars(s.PHP, s.Xdebug, s.Version, addr)...)

	// set the blackfire envs if available
	// if s.Blackfire.ServerID != "" {
	// 	envs = append(envs, "BLACKFIRE_SERVER_ID="+s.Blackfire.ServerID)
	// }
	// if s.Blackfire.ServerToken != "" {
	// 	envs = append(envs, "BLACKFIRE_SERVER_TOKEN="+s.Blackfire.ServerToken)
	// }

	return envs
}

func phpVars(php PHP, version string) []string {
	// set the composer home so we can install plugins and
	// updates from the control panel
	envs := []string{"COMPOSER_HOME=/tmp"}

	// if they do not specify the error... false means on
	if !php.DisplayErrors {
		envs = append(envs, "PHP_DISPLAY_ERRORS="+DefaultEnvs["PHP_DISPLAY_ERRORS"])
	} else {
		envs = append(envs, "PHP_DISPLAY_ERRORS=off")
	}

	if php.MemoryLimit == "" {
		envs = append(envs, "PHP_MEMORY_LIMIT="+DefaultEnvs["PHP_MEMORY_LIMIT"])
	} else {
		envs = append(envs, "PHP_MEMORY_LIMIT="+php.MemoryLimit)
	}

	if php.MaxExecutionTime == 0 {
		envs = append(envs, "PHP_MAX_EXECUTION_TIME="+DefaultEnvs["PHP_MAX_EXECUTION_TIME"])
	} else {
		envs = append(envs, fmt.Sprintf("%s=%d", "PHP_MAX_EXECUTION_TIME", php.MaxExecutionTime))
	}

	if php.UploadMaxFileSize == "" {
		envs = append(envs, "PHP_UPLOAD_MAX_FILESIZE="+DefaultEnvs["PHP_UPLOAD_MAX_FILESIZE"])
	} else {
		envs = append(envs, "PHP_UPLOAD_MAX_FILESIZE="+php.UploadMaxFileSize)
	}

	if php.MaxInputVars == 0 {
		envs = append(envs, "PHP_MAX_INPUT_VARS="+DefaultEnvs["PHP_MAX_INPUT_VARS"])
	} else {
		envs = append(envs, fmt.Sprintf("%s=%d", "PHP_MAX_INPUT_VARS", php.MaxInputVars))
	}

	if php.PostMaxSize == "" {
		envs = append(envs, "PHP_POST_MAX_SIZE="+DefaultEnvs["PHP_POST_MAX_SIZE"])
	} else {
		envs = append(envs, fmt.Sprintf("%s=%s", "PHP_POST_MAX_SIZE", php.PostMaxSize))
	}

	// handle opcache settings
	if php.OpcacheEnable {
		envs = append(envs, "PHP_OPCACHE_ENABLE=1")
	} else {
		envs = append(envs, "PHP_OPCACHE_ENABLE="+DefaultEnvs["PHP_OPCACHE_ENABLE"])
	}

	if php.OpcacheRevalidateFreq == 0 {
		envs = append(envs, "PHP_OPCACHE_REVALIDATE_FREQ="+DefaultEnvs["PHP_OPCACHE_REVALIDATE_FREQ"])
	} else {
		envs = append(envs, fmt.Sprintf("PHP_OPCACHE_REVALIDATE_FREQ=%d", php.OpcacheRevalidateFreq))

	}

	return envs
}

func xdebugVars(php PHP, xdebug bool, version, addr string) []string {
	envs := []string{}

	// always set the session
	envs = append(envs, "XDEBUG_SESSION=PHPSTORM")

	// if xdebug is not enabled
	if !xdebug {
		return append(envs, "XDEBUG_MODE=off")
	}

	switch version {
	case "8.0", "7.4", "7.3", "7.2":
		envs = append(envs, fmt.Sprintf(`XDEBUG_CONFIG=client_host=%s client_port=9003`, addr))
		envs = append(envs, "XDEBUG_MODE=develop,debug")
	default:
		// use legacy xdebug settings to support older versions of php
		envs = append(envs, fmt.Sprintf(`XDEBUG_CONFIG=idekey=PHPSTORM remote_host=%s profiler_enable=1 remote_port=9000 remote_autostart=1 remote_enable=1`, addr))
		envs = append(envs, "XDEBUG_MODE=xdebug2")
	}

	return envs
}

func cleanPath(home, path string) (string, error) {
	p := path
	if strings.Contains(p, "~") {
		p = strings.Replace(p, "~", home, -1)
	}

	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}

	return filepath.Clean(abs), nil
}
