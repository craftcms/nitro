package envvars

import (
	"fmt"

	"github.com/craftcms/nitro/pkg/config"
)

var (
	// Defaults is used to map a config to a known environment variable that is used
	// on the container instances to their default values
	Defaults = map[string]string{
		"PHP_DISPLAY_ERRORS":              "on",
		"PHP_MEMORY_LIMIT":                "512M",
		"PHP_MAX_EXECUTION_TIME":          "5000",
		"PHP_UPLOAD_MAX_FILESIZE":         "512M",
		"PHP_MAX_INPUT_VARS":              "5000",
		"PHP_POST_MAX_SIZE":               "512M",
		"PHP_OPCACHE_ENABLE":              "0",
		"PHP_OPCACHE_REVALIDATE_FREQ":     "0",
		"PHP_OPCACHE_VALIDATE_TIMESTAMPS": "0",
		"XDEBUG_MODE":                     "off",
		"XDEBUG_SESSION":                  "PHPSTORM",
		"XDEBUG_CONFIG":                   "",
		"BLACKFIRE_SERVER_ID":             "",
		"BLACKFIRE_SERVER_TOKEN":          "",
	}
)

func BlackfireCredentials(c config.Config) ([]string, error) {
	if c.Blackfire.ClientID == "" || c.Blackfire.ClientToken == "" {
		return nil, fmt.Errorf("no blackfire client credentials provided")
	}

	var envs []string
	envs = append(envs, "BLACKFIRE_CLIENT_ID="+c.Blackfire.ClientID)
	envs = append(envs, "BLACKFIRE_CLIENT_TOKEN="+c.Blackfire.ClientToken)

	return envs, nil
}

// ForApp takes an app and creates the default environment variables for a container
func ForApp(app config.App, addr string) ([]string, error) {
	var envs []string

	// set the php vars
	php := app.PHP

	// if they do not specify the error... false means on
	if !php.DisplayErrors {
		envs = append(envs, "PHP_DISPLAY_ERRORS="+Defaults["PHP_DISPLAY_ERRORS"])
	} else {
		envs = append(envs, "PHP_DISPLAY_ERRORS=off")
	}

	if php.MemoryLimit == "" {
		envs = append(envs, "PHP_MEMORY_LIMIT="+Defaults["PHP_MEMORY_LIMIT"])
	} else {
		envs = append(envs, "PHP_MEMORY_LIMIT="+php.MemoryLimit)
	}

	if php.MaxExecutionTime == 0 {
		envs = append(envs, "PHP_MAX_EXECUTION_TIME="+Defaults["PHP_MAX_EXECUTION_TIME"])
	} else {
		envs = append(envs, fmt.Sprintf("%s=%d", "PHP_MAX_EXECUTION_TIME", php.MaxExecutionTime))
	}

	if php.UploadMaxFileSize == "" {
		envs = append(envs, "PHP_UPLOAD_MAX_FILESIZE="+Defaults["PHP_UPLOAD_MAX_FILESIZE"])
	} else {
		envs = append(envs, "PHP_UPLOAD_MAX_FILESIZE="+php.UploadMaxFileSize)
	}

	if php.MaxInputVars == 0 {
		envs = append(envs, "PHP_MAX_INPUT_VARS="+Defaults["PHP_MAX_INPUT_VARS"])
	} else {
		envs = append(envs, fmt.Sprintf("%s=%d", "PHP_MAX_INPUT_VARS", php.MaxInputVars))
	}

	if php.PostMaxSize == "" {
		envs = append(envs, "PHP_POST_MAX_SIZE="+Defaults["PHP_POST_MAX_SIZE"])
	} else {
		envs = append(envs, fmt.Sprintf("%s=%s", "PHP_POST_MAX_SIZE", php.PostMaxSize))
	}

	// handle opcache settings
	if php.OpcacheEnable {
		envs = append(envs, "PHP_OPCACHE_ENABLE=1")
	} else {
		envs = append(envs, "PHP_OPCACHE_ENABLE="+Defaults["PHP_OPCACHE_ENABLE"])
	}

	if php.OpcacheRevalidateFreq == 0 {
		envs = append(envs, "PHP_OPCACHE_REVALIDATE_FREQ="+Defaults["PHP_OPCACHE_REVALIDATE_FREQ"])
	} else {
		envs = append(envs, fmt.Sprintf("PHP_OPCACHE_REVALIDATE_FREQ=%d", php.OpcacheRevalidateFreq))
	}

	if php.OpcacheValidateTimestamps {
		envs = append(envs, "PHP_OPCACHE_VALIDATE_TIMESTAMPS="+Defaults["PHP_OPCACHE_VALIDATE_TIMESTAMPS"])
	} else {
		envs = append(envs, "PHP_OPCACHE_VALIDATE_TIMESTAMPS="+Defaults["PHP_OPCACHE_VALIDATE_TIMESTAMPS"])
	}

	// set the xdebug vars

	// always set the session
	envs = append(envs, "XDEBUG_SESSION=PHPSTORM")

	// set the site name for xdebug clients
	envs = append(envs, fmt.Sprintf("PHP_IDE_CONFIG=serverName=%s", app.Hostname))

	// if xdebug is not enabled
	if !app.Xdebug {
		envs = append(envs, "XDEBUG_MODE=off")
	} else {
		switch app.PHPVersion {
		case "8.0", "7.4", "7.3", "7.2":
			envs = append(envs, fmt.Sprintf(`XDEBUG_CONFIG=client_host=%s client_port=9003`, addr))
			envs = append(envs, "XDEBUG_MODE=develop,debug")
		default:
			// use legacy xdebug settings to support older versions of php
			envs = append(envs, fmt.Sprintf(`XDEBUG_CONFIG=idekey=PHPSTORM remote_host=%s profiler_enable=1 remote_port=9000 remote_autostart=1 remote_enable=1`, addr))
			envs = append(envs, "XDEBUG_MODE=xdebug2")
		}
	}

	return envs, nil
}
