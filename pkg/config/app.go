package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type App struct {
	Config     string   `yaml:"config,omitempty"`
	Dockerfile bool     `yaml:"dockerfile,omitempty"`
	Hostname   string   `yaml:"hostname,omitempty"`
	Aliases    []string `yaml:"aliases,omitempty"`
	Path       string   `yaml:"path,omitempty"`
	Webroot    string   `yaml:"webroot,omitempty"`
	PHPVersion string   `yaml:"php_version,omitempty"`
	PHP        PHP      `yaml:"php,omitempty"`
	Extensions []string `yaml:"extensions,omitempty"`
	Xdebug     bool     `yaml:"xdebug,omitempty"`
	Blackfire  bool     `yaml:"blackfire,omitempty"`
	Suspended  bool     `yaml:"suspended,omitempty"`
	Database   struct {
		Engine  string `yaml:"engine,omitempty"`
		Version string `yaml:"version,omitempty"`
	} `yaml:"database,omitempty"`

	// we store a location to the users HomeDir directory to keep the code a little cleaner
	HomeDir string
}

// AsEnvs takes a gateway addr and turns specific options
// such as PHP settings into env vars that can be set on the
// containers environment
func (s *App) AsEnvs(addr string) []string {
	var envs []string

	// set the php vars
	envs = append(envs, phpVars(s.PHP, s.PHPVersion)...)

	return append(envs, xdebugVars(s.PHP, s.Xdebug, s.PHPVersion, s.GetHostname(), addr)...)
}

// GetHostname will get an apps hostname and ensure the top-level
// configuration takes precedence.
func (a App) GetHostname() string {
	// if the config path is defined read the file and get the hostname
	if a.Config != "" {
		path, err := cleanPath(a.HomeDir, a.Config)
		if err != nil {
			return a.Hostname
		}

		app := marshalAppConfigFrom(path)
		if app.Hostname != "" {
			return app.Hostname
		}
	}

	return a.Hostname
}

func (a App) GetAbsPath() (string, error) {
	return cleanPath(a.HomeDir, a.Path)
}

func marshalAppConfigFrom(path string) App {
	var app App
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return app
	}

	err = yaml.Unmarshal(f, &app)
	if err != nil {
		return app
	}

	return app
}
