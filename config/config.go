package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/craftcms/nitro/internal/helpers"
)

type Config struct {
	Name      string     `yaml:"name"`
	PHP       string     `yaml:"php"`
	CPUs      string     `yaml:"cpus"`
	Disk      string     `yaml:"disk"`
	Memory    string     `yaml:"memory"`
	Mounts    []Mount    `yaml:"mounts"`
	Databases []Database `yaml:"databases"`
	Sites     []Site     `yaml:"sites"`
}

type Mount struct {
	Source string `yaml:"source"`
	Dest   string `yaml:"dest"`
}

type Database struct {
	Engine  string `yaml:"engine"`
	Version string `yaml:"version"`
	Port    string `yaml:"port"`
}

type Site struct {
	Hostname string `yaml:"domain"`
	Path     string `yaml:"path"`
	Webroot  string `yaml:"docroot"`
}

func (c *Config) AddSite(site Site) error {
	c.Sites = append(c.Sites, site)
	return nil
}

func (c *Config) AddMount(m Mount) error {
	// replace the homedir with the tilde
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	prefix := m.Source[0:2]
	switch prefix {
	case "./":
		fp, err := filepath.Abs(m.Source)
		if err != nil {
			return err
		}

		m.Source = strings.Replace(fp, home, "~", 1)
	case "~/":
		m.Source = strings.Replace(m.Source, home, "~", 1)
	default:
		fp, err := filepath.Abs(m.Source)
		if err != nil {
			return err
		}

		m.Source = strings.Replace(fp, home, "~", 1)
	}

	if m.Dest == "" {
		dirname, err := helpers.PathName(m.Source)
		if err != nil {
			return err
		}
		m.Dest = "/nitro/sites/" + dirname
	}

	c.Mounts = append(c.Mounts, m)

	return nil
}

func (c *Config) RemoveSite(site string) error {
	for i, s := range c.Sites {
		if s.Hostname == site {
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
