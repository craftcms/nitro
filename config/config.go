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
	"github.com/craftcms/nitro/internal/resolve"
)

type Config struct {
	PHP       string     `yaml:"php"`
	Mounts    []Mount    `yaml:"mounts,omitempty"`
	Databases []Database `yaml:"databases"`
	Sites     []Site     `yaml:"sites,omitempty"`
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

// AlreadyMounted takes a new mount and will check if the
// mount source is already mounted to the virtual machine
// and will also check if the new mount is a parent mount
func (c *Config) AlreadyMounted(mount Mount) bool {
	// get the home directory
	home, err := homedir.Dir()
	if err != nil {
		return false
	}

	// get the local path of the mount
	newLocal, err := resolve.AbsPath(mount.Source, home)
	if err != nil {
		return false
	}

	// check each of the mounts in the config
	for _, m := range c.Mounts {
		// get the abs path of the existing mount
		existingLocal, err := resolve.AbsPath(m.Source, home)
		if err != nil {
			continue
		}

		// if it is an exact match
		if existingLocal == newLocal {
			return true
		}

		// if it is a sub folder of the mount
		if strings.Contains(existingLocal, newLocal) {
			return true
		}
	}

	return false
}

// GetExpandedMounts will take all of the mounts in a config file
// and "expand" or get the full path mount source and return
// a slice of mounts
func (c *Config) GetExpandedMounts() []Mount {
	var mounts []Mount
	for _, m := range c.Mounts {
		mounts = append(mounts, Mount{Source: m.AbsSourcePath(), Dest: m.Dest})
	}
	return mounts
}

// MountExists will check if a mount exists by checking if it is an exact
// dest or a parent of an existing dest
func (c *Config) MountExists(dest string) bool {
	for _, mount := range c.Mounts {
		// TODO expand the paths for the mounts?
		if mount.IsExact(dest) || mount.IsParent(dest) {
			return true
		}
	}

	return false
}

// SiteExists will check if a site exists withing the current config,
// it uses .IsExact to verify the site hostname and webroot.
func (c *Config) SiteExists(site Site) bool {
	for _, s := range c.Sites {
		if s.IsExact(site) {
			return true
		}
	}

	return false
}

// DatabaseExists check a provided database against the config file
// to see if the database already exists.
func (c *Config) DatabaseExists(database Database) bool {
	for _, d := range c.Databases {
		if d.Engine == database.Engine && d.Version == database.Version && d.Port == database.Port {
			return true
		}
	}

	return false
}

// SitesAsList returns the sites a slice of strings
// which is useful for select lists.
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
		m.Dest = "/home/ubuntu/sites/" + dirname
	}

	c.Mounts = append(c.Mounts, m)

	return nil
}

// RenameSite takes a site and a hostname and will replace the
// hostname in the config file for that site. If the site
// cannot be found, it will return an error.
func (c *Config) RenameSite(site Site, hostname string) error {
	for i, s := range c.Sites {
		if s.Hostname == site.Hostname {
			c.Sites[i] = Site{
				Hostname: hostname,
				Webroot:  strings.Replace(s.Webroot, s.Hostname, hostname, 1),
			}

			return nil
		}
	}

	// TODO rename the mount as well if it is a direct mount

	return errors.New("unable to locate the site with the hostname: " + site.Hostname)
}

func (c *Config) RenameMountBySite(site Site) error {
	for i, mount := range c.Mounts {
		sp := strings.Split(site.Webroot, "/")
		siteMount := sp[len(sp)-1]
		if strings.Contains(mount.Dest, siteMount) {
			c.Mounts[i] = Mount{
				Source: mount.Source,
				Dest:   siteMount,
			}

			return nil
		}
	}

	return errors.New("unable to find the mount for the site " + site.Hostname)
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

	return errors.New("unable to find the hostname " + hostname + " to remove")
}

// RemoveMountBySiteWebroot takes a complete webroot, including
// the www,public,public_html,www directory name. It will then
// find the dest by splitting a path and removing the webroot
// directory name. If it cannot find the mount, it errors.
func (c *Config) RemoveMountBySiteWebroot(webroot string) error {
	path := strings.Split(webroot, "/")
	t := path[:len(path)-1]
	dest := strings.Join(t, "/")

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
	path := strings.Split(webroot, "/")
	t := path[:len(path)-1]
	dest := strings.Join(t, "/")

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

// Read is used to read in a config file or
// return an error
func Read() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
