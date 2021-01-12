package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/craftcms/nitro/pkg/helpers"
	"gopkg.in/yaml.v3"
)

const FileName = "nitro.yaml"

var (
	// ErrNoConfigFile is returned when a configuration file cannot be found
	ErrNoConfigFile = fmt.Errorf("there is no config file for the environment")

	// DefaultEnvs is used to map a config to a known environment variable that is used
	// on the container instances to their default values
	DefaultEnvs = map[string]string{
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
)

// Config represents the nitro-dev.yaml users add for local development.
type Config struct {
	Blackfire Blackfire  `yaml:"blackfire,omitempty"`
	Databases []Database `yaml:"databases,omitempty"`
	Services  Services   `yaml:"services"`
	Sites     []Site     `yaml:"sites,omitempty"`
	File      string     `yaml:"-"`
}

// Blackfire allows users to setup their containers to use blackfire locally.
type Blackfire struct {
	ServerID    string `yaml:"server_id,omitempty"`
	ServerToken string `yaml:"server_token,omitempty"`
}

// Database is the struct used to represent a database engine
// that is a combination of a engine (e.g. mariadb, mysql, or
// postgres), the version number, and the port. The engine
// and version are directly related to the official docker
// images on the docker hub.
type Database struct {
	Engine  string `yaml:"engine,omitempty"`
	Version string `yaml:"version,omitempty"`
	Port    string `yaml:"port,omitempty"`
}

// GetHostname returns a friendly and predictable name for a database
// container. It is used for accessing a database by hostname. For
// example, mysql-8.0-3306 would be the hostname to use in the .env
// for DB_HOST.
func (d *Database) GetHostname() (string, error) {
	if d.Engine == "" || d.Version == "" || d.Port == "" {
		return "", fmt.Errorf("the engine, version, and port must be defined for the database")
	}

	return fmt.Sprintf("%s-%s-%s.nitro", d.Engine, d.Version, d.Port), nil
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

// Site represents a web application. It has a hostname, aliases (which
// are alternate domains), the local path to the site, additional mounts
// to add to the container, and the directory the index.php is located.
type Site struct {
	Hostname string   `yaml:"hostname"`
	Aliases  []string `yaml:"aliases,omitempty"`
	Path     string   `yaml:"path"`
	Version  string   `yaml:"version"`
	PHP      PHP      `yaml:"php,omitempty"`
	Webroot  string   `yaml:"webroot"`
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
	envs = append(envs, xdebugVars(s.PHP, s.Xdebug, s.Version, s.Hostname, addr)...)

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

// RemoveSite takes a hostname and will remove the site by its
// hostname from the config file.
func (c *Config) RemoveSite(hostname string) error {
	for i := len(c.Sites) - 1; i >= 0; i-- {
		site := c.Sites[i]
		if site.Hostname == hostname {
			c.Sites = append(c.Sites[:i], c.Sites[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("unable to find the hostname %q to remove", hostname)
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
			// replace the site if
			if s.Xdebug {
				return nil
			}

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
		// create the .nitro directory if it does not exist
		dir, _ := filepath.Split(c.File)
		if err := helpers.MkdirIfNotExists(dir); err != nil {
			return err
		}

		// otherwise create it
		f, err := os.Create(c.File)
		if err != nil {
			return err
		}
		defer f.Close()

		// try to chown otherwise be quiet
		_ = f.Chown(os.Geteuid(), os.Getuid())
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

func xdebugVars(php PHP, xdebug bool, version, hostname, addr string) []string {
	envs := []string{}

	// always set the session
	envs = append(envs, "XDEBUG_SESSION=PHPSTORM")

	// set the site name for xdebug clients
	envs = append(envs, fmt.Sprintf("PHP_IDE_CONFIG=serverName=%s", hostname))

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
