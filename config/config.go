package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/craftcms/nitro/internal/helpers"
)

type Config struct {
	PHP       string     `yaml:"-"`
	CPUs      string     `yaml:"-"`
	Disk      string     `yaml:"-"`
	Memory    string     `yaml:"-"`
	Mounts    []Mount    `yaml:"mounts,omitempty"`
	Databases []Database `yaml:"databases"`
	Sites     []Site     `yaml:"sites,omitempty"`
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
	Hostname string   `yaml:"hostname"`
	Webroot  string   `yaml:"webroot"`
	Aliases  []string `yaml:"aliases,omitempty"`
}

func (m *Mount) AbsSourcePath() string {
	home, _ := homedir.Dir()
	return strings.Replace(m.Source, "~", home, 1)
}

func (c *Config) AddSite(site Site) error {
	if len(site.Aliases) == 0 {
		site.Aliases = nil
	}

	c.Sites = append(c.Sites, site)
	return nil
}

func (c *Config) GetSites() []Site {
	return c.Sites
}

func (c *Config) SitesAsList() []string {
	var s []string
	for _, site := range c.Sites {
		s = append(s, site.Hostname)
	}
	return s
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

func (c *Config) RenameSite(site Site, hostname string) error {
	for _, s := range c.Sites {
		if site.Hostname == s.Hostname {
			s.Hostname = hostname

			return nil
		}
	}

	return errors.New("unable to locate the site with the hostname: " + site.Hostname)
}

func (c *Config) RemoveSite(hostname string) error {
	for i := len(c.Sites) - 1; i >= 0; i-- {
		site := c.Sites[i]
		if site.Hostname == hostname {
			c.Sites = append(c.Sites[:i], c.Sites[i+1:]...)
			return nil
		}
	}

	return errors.New("unable to find the hostname " + hostname + " to remove")
}

// RemoveMountBySiteWebroot takes a complete webroot, including
// the www,public,public_html,www directory name. It will then
// find the dest by splitting a path and removing the webroot
// directory name. If it cannot find the mount, it errors.
func (c *Config) RemoveMountBySiteWebroot(webroot string) error {
	path := strings.Split(webroot, string(os.PathSeparator))
	t := path[:len(path)-1]
	dest := strings.Join(t, string(os.PathSeparator))

	for i := len(c.Mounts) - 1; i >= 0; i-- {
		mount := c.Mounts[i]
		if mount.Dest == dest {
			c.Mounts = append(c.Mounts[:i], c.Mounts[i+1:]...)
			return nil
		}
	}

	return errors.New("unable to find the mount")
}

func (c *Config) FindMountBySiteWebroot(webroot string) *Mount {
	path := strings.Split(webroot, string(os.PathSeparator))
	t := path[:len(path)-1]
	dest := strings.Join(t, string(os.PathSeparator))

	for _, mount := range c.Mounts {
		if mount.Dest == dest {
			return &mount
		}
	}

	return nil
}

func (c *Config) Save(filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}

func (c *Config) SaveAs(home, machine string) error {
	nitroDir := home + "/.nitro/"

	if err := helpers.MkdirIfNotExists(nitroDir); err != nil {
		return err
	}
	filename := nitroDir + machine + ".yaml"

	_ = helpers.CreateFileIfNotExist(filename)

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
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
