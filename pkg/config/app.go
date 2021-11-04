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
	Extensions []string `yaml:"extensions,omitempty"`
	Xdebug     bool     `yaml:"xdebug,omitempty"`
	Blackfire  bool     `yaml:"blackfire,omitempty"`
	Database   struct {
		Engine  string `yaml:"engine,omitempty"`
		Version string `yaml:"version,omitempty"`
	} `yaml:"database,omitempty"`
}

// GetHostname will get an apps hostname and ensure the top-level
// configuration takes precedence.
func (a App) GetHostname(home string) string {
	// if the config path is defined read the file and get the hostname
	if a.Config != "" {
		path, err := cleanPath(home, a.Config)
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

func (a App) GetAbsPath(home string) (string, error) {
	return cleanPath(home, a.Path)
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
