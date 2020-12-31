package config

import (
	"fmt"
	"path/filepath"
	"strings"
)

// DefaultEnvs is used to map a config to a known environment variable that is used
// on the container instances to their default values
var DefaultEnvs = map[string]string{
	// PHP specific settings
	"PHP_DISPLAY_ERRORS":          "on",
	"PHP_MEMORY_LIMIT":            "512M",
	"PHP_MAX_EXECUTION_TIME":      "5000",
	"PHP_UPLOAD_MAX_FILESIZE":     "512M",
	"PHP_MAX_INPUT_VARS":          "5000",
	"PHP_POST_MAX_SIZE":           "512M",
	"PHP_OPCACHE_ENABLE":          "0",
	"PHP_OPCACHE_REVALIDATE_FREQ": "0",
	"XDEBUG_MODE":                 "off",
	"XDEBUG_SESSION":              "nitro",
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
	return s.cleanPath(home, s.Path)
}

// AsEnvs takes a gateway addr and turns specific options
// such as PHP settings into env vars that can be set on the
// containers environment
func (s *Site) AsEnvs(addr string) []string {
	var envs []string

	// if they do not specify the error... false means on
	if !s.PHP.DisplayErrors {
		envs = append(envs, "PHP_DISPLAY_ERRORS="+DefaultEnvs["PHP_DISPLAY_ERRORS"])
	} else {
		envs = append(envs, "PHP_DISPLAY_ERRORS=off")
	}

	if s.PHP.MemoryLimit == "" {
		envs = append(envs, "PHP_MEMORY_LIMIT="+DefaultEnvs["PHP_MEMORY_LIMIT"])
	} else {
		envs = append(envs, "PHP_MEMORY_LIMIT="+s.PHP.MemoryLimit)
	}

	if s.PHP.MaxExecutionTime == 0 {
		envs = append(envs, "PHP_MAX_EXECUTION_TIME="+DefaultEnvs["PHP_MAX_EXECUTION_TIME"])
	} else {
		envs = append(envs, fmt.Sprintf("%s=%d", "PHP_MAX_EXECUTION_TIME", s.PHP.MaxExecutionTime))
	}

	if s.PHP.UploadMaxFileSize == "" {
		envs = append(envs, "PHP_UPLOAD_MAX_FILESIZE="+DefaultEnvs["PHP_UPLOAD_MAX_FILESIZE"])
	} else {
		envs = append(envs, "PHP_UPLOAD_MAX_FILESIZE="+s.PHP.UploadMaxFileSize)
	}

	if s.PHP.MaxInputVars == 0 {
		envs = append(envs, "PHP_MAX_INPUT_VARS="+DefaultEnvs["PHP_MAX_INPUT_VARS"])
	} else {
		envs = append(envs, fmt.Sprintf("%s=%d", "PHP_MAX_INPUT_VARS", s.PHP.MaxInputVars))
	}

	if s.PHP.PostMaxSize == "" {
		envs = append(envs, "PHP_POST_MAX_SIZE="+DefaultEnvs["PHP_POST_MAX_SIZE"])
	} else {
		envs = append(envs, fmt.Sprintf("%s=%s", "PHP_POST_MAX_SIZE", s.PHP.PostMaxSize))
	}

	// handle opcache settings
	if s.PHP.OpcacheEnable {
		envs = append(envs, "PHP_OPCACHE_ENABLE=1")
	} else {
		envs = append(envs, "PHP_OPCACHE_ENABLE="+DefaultEnvs["PHP_OPCACHE_ENABLE"])
	}

	if s.PHP.OpcacheRevalidateFreq == 0 {
		envs = append(envs, "PHP_OPCACHE_REVALIDATE_FREQ="+DefaultEnvs["PHP_OPCACHE_REVALIDATE_FREQ"])
	} else {
		envs = append(envs, fmt.Sprintf("PHP_OPCACHE_REVALIDATE_FREQ=%d", s.PHP.OpcacheRevalidateFreq))

	}

	// check if xdebug is enabled
	switch s.Xdebug {
	case false:
		envs = append(envs, "XDEBUG_MODE=off")
	default:
		// opts.Proxy.NetworkSettings.Networks[opts.Environment].IPAddress
		// opts.Network.IPAM.Config[0].Gateway
		envs = append(envs, fmt.Sprintf(`XDEBUG_CONFIG=client_host=%s start_with_request=yes discover_client_host=1`, addr))
		envs = append(envs, "XDEBUG_SESSION=PHPSTORM")
		envs = append(envs, "XDEBUG_MODE=develop,debug")
	}

	// TODO(jasonmccallister) add opcache settings
	// "PHP_OPCACHE_VALIDATE_TIMESTAMPS=0",
	// "PHP_OPCACHE_MAX_ACCELERATED_FILES=10000",
	// "PHP_OPCACHE_MEMORY_CONSUMPTION=128",
	// "PHP_OPCACHE_MAX_WASTED_PERCENTAGE=10",
	// "PHP_OPCACHE_INTERNED_STRINGS_BUFFER=16",
	// "PHP_OPCACHE_FAST_SHUTDOWN=1"

	// set the blackfire envs if available
	// if s.Blackfire.ServerID != "" {
	// 	envs = append(envs, "BLACKFIRE_SERVER_ID="+s.Blackfire.ServerID)
	// }
	// if s.Blackfire.ServerToken != "" {
	// 	envs = append(envs, "BLACKFIRE_SERVER_TOKEN="+s.Blackfire.ServerToken)
	// }

	return envs
}

func (s *Site) cleanPath(home, path string) (string, error) {
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
