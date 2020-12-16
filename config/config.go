package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

var (
	// ErrNoConfigFile is returned when a configuration file cannot be found
	ErrNoConfigFile = fmt.Errorf("there is no config file for the environment")
)

// Config represents the nitro-dev.yaml users add for local development.
type Config struct {
	Blackfire Blackfire  `mapstructure:"blackfire,omitempty" yaml:"blackfire,omitempty"`
	PHP       PHP        `mapstructure:"php,omitempty" yaml:"php,omitempty"`
	Databases []Database `yaml:"databases,omitempty"`
	Services  Services   `mapstructure:"services,omitempty" yaml:"services,omitempty"`
	Sites     []Site     `yaml:"sites,omitempty"`

	file string
}

// Blackfire allows users to setup their containers to use blackfire locally.
type Blackfire struct {
	ServerID    string `mapstructure:"server_id,omitempty" yaml:"server_id,omitempty"`
	ServerToken string `mapstructure:"server_token,omitempty" yaml:"server_token,omitempty"`
}

// PHP is nested in a configuration and allows setting environment variables
// for sites to override in the local development environment.
type PHP struct {
	DisplayErrors         bool   `mapstructure:"display_errors,omitempty" yaml:"display_errors,omitempty"`
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

// Services define common tools for development that should run as containers. We don't expose the volumes, ports, and
// networking options for these types of services. We plan to support "custom" container options to make local users
// development even better.
type Services struct {
	Blackfire bool `yaml:"blackfire"`
	DynamoDB  bool `yaml:"dynamodb"`
	Mailhog   bool `yaml:"mailhog"`
	Minio     bool `yaml:"minio"`
	Redis     bool `yaml:"redis"`
}

// Load is used to return the environment name, unmarshalled config, and
// returns an error when trying to get the users home directory or
// while marshalling the config.
func Load(home, env string) (*Config, error) {
	v := viper.New()
	v.AddConfigPath(home)
	v.SetConfigName(env + ".yaml")

	// set the config file
	if env == "" {
		env = "nitro-dev"
	}

	// set the config file
	v.SetConfigFile(filepath.Join(home, ".nitro", env+".yaml"))

	// read the config
	if err := v.ReadInConfig(); err != nil {
		return nil, ErrNoConfigFile
	}

	if v.ConfigFileUsed() == "" {
		return nil, ErrNoConfigFile
	}

	cfg := &Config{}
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// set the file being used on the config
	cfg.file = v.ConfigFileUsed()

	// return the config
	return cfg, nil
}

// SetFile is used the set a file for the config and
// is only really used when first time setup is used.
func (c *Config) SetFile(file string) {
	c.file = file
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
	viper.SetConfigFile(c.file)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	cfg := Config{}
	viper.Unmarshal(&cfg)

	fmt.Println(cfg)

	if err := viper.WriteConfigAs("testing.yaml"); err != nil {
		return err
	}

	return nil
}
