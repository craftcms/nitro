package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const FileName = "nitro.yml"

var (
	// ErrNoConfigFile is returned when a configuration file cannot be found
	ErrNoConfigFile = fmt.Errorf("there is no config file for the environment")
)

// Config represents the nitro-dev.yaml users add for local development.
type Config struct {
	Blackfire Blackfire  `yaml:"blackfire,omitempty"`
	Databases []Database `yaml:"databases,omitempty"`
	Mounts    []Mount    `yaml:"mounts,omitempty"`
	Services  Services   `yaml:"services"`
	Sites     []Site     `yaml:"sites,omitempty"`
	File      string     `yaml:"-"`
}

// Blackfire allows users to setup their containers to use blackfire locally.
type Blackfire struct {
	ServerID    string `yaml:"server_id,omitempty"`
	ServerToken string `yaml:"server_token,omitempty"`
}

// Mount represents a docker container that is not mounted in an nginx container
// and does not accept routing through the proxy. It is however added to the nitro
// network so it can access other resources.
type Mount struct {
	Path    string `yaml:"path"`
	Version string `yaml:"version"`
	PHP     PHP    `yaml:"php,omitempty"`
	Xdebug  bool   `yaml:"xdebug,omitempty"`
}

// GetAbsPath gets the directory for a mount.Path,
// It is used to create the mount for a container.
func (m *Mount) GetAbsPath(home string) (string, error) {
	return cleanPath(home, m.Path)
}

// AsEnvs takes a gateway addr and turns specific options
// such as PHP settings into env vars that can be set on the
// containers environment
func (m *Mount) AsEnvs(addr string) []string {
	var envs []string

	if addr == "" {
		addr = "host.docker.internal"
	}

	// set the php vars
	envs = append(envs, phpVars(envs, m.PHP, m.Version)...)

	// get the xdebug vars
	envs = append(envs, xdebugVars(envs, m.PHP, m.Xdebug, m.Version, addr)...)

	// set the blackfire envs if available
	// if s.Blackfire.ServerID != "" {
	// 	envs = append(envs, "BLACKFIRE_SERVER_ID="+s.Blackfire.ServerID)
	// }
	// if s.Blackfire.ServerToken != "" {
	// 	envs = append(envs, "BLACKFIRE_SERVER_TOKEN="+s.Blackfire.ServerToken)
	// }

	return envs
}

// PHP is nested in a configuration and allows setting environment variables
// for sites to override in the local development environment.
type PHP struct {
	DisplayErrors         bool   `yaml:"display_errors,omitempty"`
	MaxExecutionTime      int    `yaml:"max_execution_time,omitempty"`
	MaxInputVars          int    `yaml:"max_input_vars,omitempty"`
	MaxInputTime          int    `yaml:"max_input_time,omitempty"`
	MaxFileUpload         string `yaml:"max_file_upload,omitempty"`
	MemoryLimit           string `yaml:"memory_limit,omitempty"`
	OpcacheEnable         bool   `yaml:"opcache_enable,omitempty"`
	OpcacheRevalidateFreq int    `yaml:"opcache_revalidate_freq,omitempty"`
	PostMaxSize           string `yaml:"post_max_size,omitempty"`
	UploadMaxFileSize     string `yaml:"upload_max_file_size,omitempty"`
}

// Services define common tools for development that should run as containers. We don't expose the volumes, ports, and
// networking options for these types of services. We plan to support "custom" container options to make local users
// development even better.
type Services struct {
	DynamoDB bool `yaml:"dynamodb"`
	Mailhog  bool `yaml:"mailhog"`
	Minio    bool `yaml:"minio"`
	Redis    bool `yaml:"redis"`
}

// Load is used to return the unmarshalled config, and
// returns an error when trying to get the users home directory or
// while marshalling the config.
func Load(home string) (*Config, error) {
	// set the config file
	file := filepath.Join(home, ".nitro", FileName)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, ErrNoConfigFile
	}

	// create the config
	c := &Config{
		File: file,
	}

	// read the file
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// unmarshal
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	// return the config
	return c, nil
}

// AddSite takes a site and adds it to the config
func (c *Config) AddSite(s Site) error {
	// check existing sites
	for _, e := range c.Sites {
		// does the hostname match
		if e.Hostname == s.Hostname {
			return fmt.Errorf("hostname already exists")
		}
	}

	// add the site to the list
	c.Sites = append(c.Sites, s)

	return nil
}

// DisableXdebug takes a sites hostname and sets the xdebug option
// to false. If the site cannot be found, it returns an error.
func (c *Config) DisableXdebug(site string) error {
	// find the site by the hostname
	for i, s := range c.Sites {
		if s.Hostname == site {
			// replace the site
			s.Xdebug = false
			c.Sites = append(c.Sites[:i], s)
			return nil
		}
	}

	return fmt.Errorf("unknown site, %s", site)
}

// EnableXdebug takes a sites hostname and sets the xdebug option
// to true. If the site cannot be found, it returns an error.
func (c *Config) EnableXdebug(site string) error {
	// find the site by the hostname
	for i, s := range c.Sites {
		if s.Hostname == site {
			// replace the site
			s.Xdebug = true
			c.Sites = append(c.Sites[:i], s)
			return nil
		}
	}

	return fmt.Errorf("unknown site, %s", site)
}

// Save takes a file path and marshals the config into a file.
func (c *Config) Save() error {
	// make sure the file exists
	if _, err := os.Stat(c.File); os.IsNotExist(err) {
		// otherwise create it
		f, err := os.Create(c.File)
		if err != nil {
			return err
		}
		defer f.Close()

		f.Chown(os.Geteuid(), os.Getuid())
	}

	// unmarshal
	data, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	// open the file
	f, err := os.OpenFile(c.File, os.O_SYNC|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	// write the content
	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}

// GetFile returns the file location for the config
func (c *Config) GetFile() string {
	return c.File
}
