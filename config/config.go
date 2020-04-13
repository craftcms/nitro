package config

import (
	"errors"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Name      string     `yaml:"name"`
	PHP       string     `yaml:"php"`
	CPUs      string     `yaml:"cpus"`
	Disk      string     `yaml:"disk"`
	Memory    string     `yaml:"memory"`
	Databases []Database `yaml:"databases"`
	Sites     []Site     `yaml:"sites"`
}

func (c *Config) AddSite(site Site) error {
	// replace the homedir with the tilde
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	site.Path = strings.Replace(site.Path, home, "~", 1)

	c.Sites = append(c.Sites, site)
	return nil
}

func (c *Config) RemoveSite(site string) error {
	for i, s := range c.Sites {
		if s.Domain == site {
			copy(c.Sites[i:], c.Sites[i+1:])
			c.Sites[len(c.Sites)-1] = Site{}
			c.Sites = c.Sites[:len(c.Sites)-1]
			return nil
		}
	}

	return errors.New("unable to find the domain " + site + " to remove")
}

func (c *Config) Save(filename string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}

func GetString(key, flag string) string {
	if viper.IsSet(key) && flag == "" {
		return viper.GetString(key)
	}

	return flag
}

func GetInt(key string, flag int) int {
	if viper.IsSet(key) && flag == 0 {
		return viper.GetInt(key)
	}

	return flag
}
