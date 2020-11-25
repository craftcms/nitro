package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var envMap = map[string]string{
	"display_errors": "PHP_DISPLAY_ERRORS",
	"memory_limit":   "PHP_MEMORY_LIMIT",
}

// Config represents the nitro-dev.yaml users add for local development.
type Config struct {
	Blackfire Blackfire  `yaml:"blackfire,omitempty"`
	PHP       PHP        `yaml:"php,omitempty"`
	Sites     []Site     `yaml:"sites,omitempty"`
	Databases []Database `yaml:"databases,omitempty"`
}

// AsEnvs takes a configuration and turns specific options
// such as PHP settings into env vars that can be set on the
// containers environment
func (c *Config) AsEnvs() []string {
	var envs []string

	if c.PHP.DisplayErrors == "" {
		envs = append(envs, envMap["display_errors"]+"=on")
	} else {
		envs = append(envs, envMap["display_errors"]+"="+c.PHP.DisplayErrors)
	}

	if c.PHP.MemoryLimit == "" {
		envs = append(envs, envMap["memory_limit"]+"=512MB")
	} else {
		// TODO(jasonmccallister) add validation
		envs = append(envs, envMap["memory_limit"]+"="+c.PHP.MemoryLimit)
	}

	if c.PHP.MaxExecutionTime == 0 {
		envs = append(envs, "PHP_MAX_EXECUTION_TIME=5000")
	} else {
		// TODO(jasonmccallister) add validation
		envs = append(envs, fmt.Sprintf("PHP_MAX_EXECUTION_TIME=%d", c.PHP.MaxExecutionTime))
	}

	if c.PHP.UploadMaxFileSize == "" {
		envs = append(envs, "PHP_UPLOAD_MAX_FILESIZE=512M")
	} else {
		envs = append(envs, "PHP_UPLOAD_MAX_FILESIZE="+c.PHP.UploadMaxFileSize)
	}

	if c.PHP.MaxInputVars == 0 {
		envs = append(envs, "PHP_MAX_INPUT_VARS=512M")
	} else {
		envs = append(envs, fmt.Sprintf("PHP_MAX_INPUT_VARS=%d", c.PHP.MaxInputVars))
	}

	if c.PHP.PostMaxSize == "" {
		envs = append(envs, "PHP_POST_MAX_SIZE=512M")
	} else {
		envs = append(envs, fmt.Sprintf("PHP_POST_MAX_SIZE=%s", c.PHP.PostMaxSize))
	}

	// set the blackfire envs if available
	if c.Blackfire.ServerID != "" {
		envs = append(envs, "BLACKFIRE_SERVER_ID="+c.Blackfire.ServerID)
	}
	if c.Blackfire.ServerToken != "" {
		envs = append(envs, "BLACKFIRE_SERVER_TOKEN="+c.Blackfire.ServerToken)
	}

	// TODO(jasonmccallister) add opcache settings
	// "PHP_OPCACHE_ENABLE=1",
	// "PHP_OPCACHE_REVALIDATE_FREQ=0",
	// "PHP_OPCACHE_VALIDATE_TIMESTAMPS=0",
	// "PHP_OPCACHE_MAX_ACCELERATED_FILES=10000",
	// "PHP_OPCACHE_MEMORY_CONSUMPTION=128",
	// "PHP_OPCACHE_MAX_WASTED_PERCENTAGE=10",
	// "PHP_OPCACHE_INTERNED_STRINGS_BUFFER=16",
	// "PHP_OPCACHE_FAST_SHUTDOWN=1"

	return envs
}

// Blackfire allows users to setup their containers to use blackfire locally.
type Blackfire struct {
	ServerID    string `mapstructure:"server_id,omitempty"`
	ServerToken string `mapstructure:"server_token,omitempty"`
}

// PHP is nested in a configuration and allows setting environment variables
// for sites to override in the local development environment.
type PHP struct {
	DisplayErrors     string `mapstructure:"display_errors,omitempty"`
	MaxExecutionTime  int    `mapstructure:"max_execution_time,omitempty"`
	MaxInputVars      int    `mapstructure:"max_input_vars,omitempty"`
	MaxInputTime      int    `mapstructure:"max_input_time,omitempty"`
	MaxFileUpload     string `mapstructure:"max_file_upload,omitempty"`
	MemoryLimit       string `mapstructure:"memory_limit,omitempty"`
	PostMaxSize       string `mapstructure:"post_max_size,omitempty"`
	UploadMaxFileSize string `mapstructure:"upload_max_file_size,omitempty"`
}

// Load is used to return the environment name, unmarshalled config, and
// returns an error when trying to get the users home directory or
// while marshalling the config.
func Load() (string, *Config, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", nil, fmt.Errorf("unable to get the home directory, %w", err)
	}

	viper.AddConfigPath(fmt.Sprintf("%s%c%s", home, os.PathSeparator, ".nitro"))
	viper.SetConfigType("yaml")

	// set the default environment name
	def := "nitro-dev"
	if os.Getenv("NITRO_DEFAULT_ENVIRONMENT") != "" {
		def = os.Getenv("NITRO_DEFAULT_ENVIRONMENT")
	}

	// set the config file
	viper.SetConfigName(def)

	if err := viper.ReadInConfig(); err != nil {
		return "", nil, err
	}

	cfg := &Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		return "", nil, err
	}

	// read the config
	return def, cfg, nil
}
