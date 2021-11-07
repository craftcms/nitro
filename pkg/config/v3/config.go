package v3

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/craftcms/nitro/pkg/paths"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Containers []Container `yaml:"containers,omitempty"`
	Blackfire  Blackfire   `yaml:"blackfire,omitempty"`
	Databases  []Database  `yaml:"databases,omitempty"`
	Services   Services    `yaml:"services,omitempty"`
	Apps       []App       `yaml:"apps,omitempty"`

	// store the users home directory
	HomeDir    string `yaml:"-"`
	ConfigFile string `yaml:"-"`
}

type App struct {
	Config     string   `yaml:"config,omitempty"`
	Dockerfile bool     `yaml:"dockerfile,omitempty"`
	Hostname   string   `yaml:"hostname,omitempty"`
	Aliases    []string `yaml:"aliases,omitempty"`
	Path       string   `yaml:"path,omitempty"`
	Webroot    string   `yaml:"webroot,omitempty"`
	PHPVersion string   `yaml:"php_version,omitempty"`
	PHP        struct {
		DisplayErrors             bool   `yaml:"display_errors,omitempty"`
		MaxExecutionTime          int    `yaml:"max_execution_time,omitempty"`
		MaxInputVars              int    `yaml:"max_input_vars,omitempty"`
		MaxInputTime              int    `yaml:"max_input_time,omitempty"`
		MaxFileUpload             string `yaml:"max_file_upload,omitempty"`
		MemoryLimit               string `yaml:"memory_limit,omitempty"`
		OpcacheEnable             bool   `yaml:"opcache_enable,omitempty"`
		OpcacheRevalidateFreq     int    `yaml:"opcache_revalidate_freq,omitempty"`
		OpcacheValidateTimestamps bool   `yaml:"opcache_validate_timestamps,omitempty"`
		PostMaxSize               string `yaml:"post_max_size,omitempty"`
		UploadMaxFileSize         string `yaml:"upload_max_file_size,omitempty"`
	} `yaml:"php,omitempty"`
	Extensions []string `yaml:"extensions,omitempty"`
	Xdebug     bool     `yaml:"xdebug,omitempty"`
	Blackfire  bool     `yaml:"blackfire,omitempty"`
	Suspended  bool     `yaml:"suspended,omitempty"`
	Database   struct {
		Engine  string `yaml:"engine,omitempty"`
		Version string `yaml:"version,omitempty"`
	} `yaml:"database,omitempty"`
}

type Blackfire struct {
	ClientID    string `yaml:"client_id,omitempty"`
	ClientToken string `yaml:"client_token,omitempty"`
	ServerID    string `yaml:"server_id,omitempty"`
	ServerToken string `yaml:"server_token,omitempty"`
}

type Database struct {
	Engine  string `yaml:"engine"`
	Version string `yaml:"version"`
	Port    string `yaml:"port"`
}

type Services struct {
	Blackfire bool `yaml:"blackfire"`
	DynamoDB  bool `yaml:"dynamodb"`
	Mailhog   bool `yaml:"mailhog"`
	Minio     bool `yaml:"minio"`
	Redis     bool `yaml:"redis"`
}

func Load() (Config, error) {
	c := Config{}
	home, err := os.UserHomeDir()
	if err != nil {
		return c, err
	}

	c.HomeDir = home

	return c, nil
}

type Container struct {
	Name    string   `yaml:"name"`
	Image   string   `yaml:"image"`
	Tag     string   `yaml:"tag"`
	Ports   []string `yaml:"ports,omitempty"`
	WebUI   int      `yaml:"web_ui,omitempty"`
	EnvFile string   `yaml:"env_file,omitempty"`
}

func (c Config) GetAppHostName(hostname string) (string, error) {
	for _, user := range c.Apps {
		// is there a user hostname defined and does it match?
		if user.Hostname == hostname {
			return user.Hostname, nil
		}

		// is user hostname not defined and a config file is present?
		if user.Hostname == "" && user.Config != "" {
			// is there a config?
			p, _ := cleanPath(c.HomeDir, user.Config)

			local, err := unmarshalAppConfigFrom(p)
			if err != nil {
				return "", err
			}

			return local.Hostname, nil
		}
	}

	return "", fmt.Errorf("unable to find app with hostname %q", hostname)
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

func unmarshalAppConfigFrom(path string) (App, error) {
	var app App
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return app, fmt.Errorf("unable to read file at %v", path)
	}

	err = yaml.Unmarshal(f, &app)
	if err != nil {
		return app, fmt.Errorf("unable to marshal file, %v", err)
	}

	return app, nil
}
