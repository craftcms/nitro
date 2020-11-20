package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Config struct {
	Extensions []string   `yaml:"exts,omitempty"`
	Blackfire  Blackfire  `yaml:"blackfire,omitempty"`
	Sites      []Site     `yaml:"sites,omitempty"`
	Databases  []Database `yaml:"databases,omitempty"`
}

type Blackfire struct {
	ClientID     string `yaml:"client_id,omitempty"`
	ClientSecret string `yaml:"client_secret,omitempty"`
}

type Site struct {
	Hostname string   `yaml:"hostname,omitempty"`
	Aliases  []string `yaml:"aliases,omitempty"`
	Path     string   `yaml:"path,omitempty"`
	PHP      string   `yaml:"php,omitempty"`
	Dir      string   `yaml:"dir,omitempty"`
}

type Database struct {
	Engine  string `yaml:"engine,omitempty"`
	Version string `yaml:"version,omitempty"`
	Port    string `yaml:"port,omitempty"`
}

func (d *Database) GetHostname() (string, error) {
	if d.Engine == "" || d.Version == "" || d.Port == "" {
		return "", fmt.Errorf("the engine, version, and port must be defined for the database")
	}

	return fmt.Sprintf("%s-%s-%s", d.Engine, d.Version, d.Port), nil
}

func Load() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("unable to get the home directory, %w", err)
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

	// read the config
	return def, viper.ReadInConfig()
}

func Umarshal() (*Config, error) {
	c := Config{}

	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
