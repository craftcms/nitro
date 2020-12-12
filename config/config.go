package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	// ErrNoConfigFile is returned when a configuration file cannot be found
	ErrNoConfigFile = fmt.Errorf("there is no config file for the environment")
)

// Envs is used to map a config to a known environment variable that is used
// on the container instances.
// TODO(jasonmccallister) not used
var Envs = map[string]string{
	// PHP specific settings
	"PHP_DISPLAY_ERRORS":      "display_errors",
	"PHP_MEMORY_LIMIT":        "memory_limit",
	"PHP_MAX_EXECUTION_TIME":  "max_execution_time",
	"PHP_UPLOAD_MAX_FILESIZE": "upload_max_filesize",
	"PHP_MAX_INPUT_VARS":      "max_input_vars",
	"PHP_POST_MAX_SIZE":       "post_max_size",
}

// Config represents the nitro-dev.yaml users add for local development.
type Config struct {
	Blackfire Blackfire  `yaml:"blackfire,omitempty"`
	PHP       PHP        `yaml:"php,omitempty"`
	Databases []Database `yaml:"databases,omitempty"`
	Sites     []Site     `yaml:"sites,omitempty"`

	file string
}

// AsEnvs takes a configuration and turns specific options
// such as PHP settings into env vars that can be set on the
// containers environment
func (c *Config) AsEnvs() []string {
	var envs []string

	if c.PHP.DisplayErrors == "" {
		envs = append(envs, "PHP_DISPLAY_ERRORS=on")
	} else {
		envs = append(envs, "PHP_DISPLAY_ERRORS="+c.PHP.DisplayErrors)
	}

	if c.PHP.MemoryLimit == "" {
		envs = append(envs, "PHP_MEMORY_LIMIT=512M")
	} else {
		envs = append(envs, "PHP_MEMORY_LIMIT="+c.PHP.MemoryLimit)
	}

	if c.PHP.MaxExecutionTime == 0 {
		envs = append(envs, "PHP_MAX_EXECUTION_TIME=5000")
	} else {
		envs = append(envs, fmt.Sprintf("%s=%d", "PHP_MAX_EXECUTION_TIME", c.PHP.MaxExecutionTime))
	}

	if c.PHP.UploadMaxFileSize == "" {
		envs = append(envs, "PHP_UPLOAD_MAX_FILESIZE=512M")
	} else {
		envs = append(envs, "PHP_UPLOAD_MAX_FILESIZE="+c.PHP.UploadMaxFileSize)
	}

	if c.PHP.MaxInputVars == 0 {
		envs = append(envs, "PHP_MAX_INPUT_VARS=512M")
	} else {
		envs = append(envs, fmt.Sprintf("%s=%d", "PHP_MAX_INPUT_VARS", c.PHP.MaxInputVars))
	}

	if c.PHP.PostMaxSize == "" {
		envs = append(envs, "PHP_POST_MAX_SIZE=512M")
	} else {
		envs = append(envs, fmt.Sprintf("%s=%s", "PHP_POST_MAX_SIZE", c.PHP.PostMaxSize))
	}

	// handle opcache settings
	if c.PHP.OpcacheEnable == true {
		envs = append(envs, "PHP_OPCACHE_ENABLE=1")
	} else {
		envs = append(envs, "PHP_OPCACHE_ENABLE=0")
	}

	if c.PHP.OpcacheRevalidateFreq > 0 {
		envs = append(envs, fmt.Sprintf("PHP_OPCACHE_REVALIDATE_FREQ=%d", c.PHP.OpcacheRevalidateFreq))
	} else {
		envs = append(envs, "PHP_OPCACHE_REVALIDATE_FREQ=0")
	}

	// TODO(jasonmccallister) add opcache settings
	// "PHP_OPCACHE_VALIDATE_TIMESTAMPS=0",
	// "PHP_OPCACHE_MAX_ACCELERATED_FILES=10000",
	// "PHP_OPCACHE_MEMORY_CONSUMPTION=128",
	// "PHP_OPCACHE_MAX_WASTED_PERCENTAGE=10",
	// "PHP_OPCACHE_INTERNED_STRINGS_BUFFER=16",
	// "PHP_OPCACHE_FAST_SHUTDOWN=1"

	// set the blackfire envs if available
	if c.Blackfire.ServerID != "" {
		envs = append(envs, "BLACKFIRE_SERVER_ID="+c.Blackfire.ServerID)
	}
	if c.Blackfire.ServerToken != "" {
		envs = append(envs, "BLACKFIRE_SERVER_TOKEN="+c.Blackfire.ServerToken)
	}

	return envs
}

// Blackfire allows users to setup their containers to use blackfire locally.
type Blackfire struct {
	ServerID    string `mapstructure:"server_id,omitempty" yaml:"server_id,omitempty"`
	ServerToken string `mapstructure:"server_token,omitempty" yaml:"server_token,omitempty"`
}

// PHP is nested in a configuration and allows setting environment variables
// for sites to override in the local development environment.
type PHP struct {
	DisplayErrors         string `mapstructure:"display_errors,omitempty" yaml:"display_errors,omitempty"`
	MaxExecutionTime      int    `mapstructure:"max_execution_time,omitempty" yaml:"max_execution_time,omitempty"`
	MaxInputVars          int    `mapstructure:"max_input_vars,omitempty" yaml:"max_input_vars,omitempty"`
	MaxInputTime          int    `mapstructure:"max_input_time,omitempty" yaml:"max_input_time,omitempty"`
	MaxFileUpload         string `mapstructure:"max_file_upload,omitempty" yaml:"max_file_upload,omitempty"`
	MemoryLimit           string `mapstructure:"memory_limit,omitempty" yaml:"memory_limit,omitempty"`
	OpcacheEnable         bool   `mapstructure:"opcache_enable,omitempty" yaml:"opcache_enable,omitempty"`
	OpcacheRevalidateFreq int    `mapstructure:"opcache_revalidate_freq,omitempty" yaml:"opcache_revalidate_freq,omitempty"`
	PostMaxSize           string `mapstructure:"post_max_size,omitempty" yaml:"post_max_size,omitempty"`
	UploadMaxFileSize     string `mapstructure:"upload_max_file_size,omitempty" yaml:"upload_max_file_size,omitempty"`
}

// Load is used to return the environment name, unmarshalled config, and
// returns an error when trying to get the users home directory or
// while marshalling the config.
func Load(home, env string) (*Config, error) {
	viper.AddConfigPath(fmt.Sprintf("%s%c%s", home, os.PathSeparator, ".nitro"))
	viper.SetConfigType("yaml")

	// set the config file
	if env == "" {
		env = "nitro-dev"
	}

	viper.SetConfigName(env)

	if err := viper.ReadInConfig(); err != nil {
		return nil, ErrNoConfigFile
	}

	if viper.ConfigFileUsed() == "" {
		return nil, ErrNoConfigFile
	}

	cfg := &Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// set the file being used on the config
	cfg.file = viper.ConfigFileUsed()

	// return the config
	return cfg, nil
}

// AddSite takes a site and adds it to the config
func (c *Config) AddSite(s Site) error {
	// if there are no aliases
	if len(s.Aliases) == 0 {
		s.Aliases = nil
	}

	// check existing sites
	for _, e := range c.Sites {
		// does the hostname match
		if e.Hostname == s.Hostname {
			return fmt.Errorf("hostname already exists")
		}

		// does the path match
		if e.Path == s.PHP {
			return fmt.Errorf("site path already exists")
		}
	}

	// add the site to the list
	c.Sites = append(c.Sites, s)

	return nil
}

// Save takes a file path and marshals the config into a file.
func (c *Config) Save() error {
	if c.file == "" {
		c.file = viper.ConfigFileUsed()
	}

	// open the file
	f, err := os.OpenFile(c.file, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	// marshal into yaml
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	// save the file
	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}
