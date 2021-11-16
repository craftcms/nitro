package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/craftcms/nitro/pkg/helpers"
	"github.com/craftcms/nitro/pkg/paths"

	"gopkg.in/yaml.v3"
)

var (
	// DirectoryName is the name of the directory to store nitro configs
	DirectoryName = ".nitro"

	// ErrNoConfigFile is returned when a configuration file cannot be found
	ErrNoConfigFile = fmt.Errorf("there is no config file for the environment")

	// ErrEmptyFile is returned when a config file is empty
	ErrEmptyFile = fmt.Errorf("the config file appears to be empty")

	// ErrDeprecatedConfig is used when a config contains sites
	ErrDeprecatedConfig = fmt.Errorf("sites are deprecated in v3.0. See the upgrade guide https://craftcms.com/docs/nitro/3.x/upgrade.html")

	// FileName is the default name for the yaml file
	FileName = "nitro.yaml"

	// DefaultEnvs is used to map a config to a known environment variable that is used
	// on the container instances to their default values
	DefaultEnvs = map[string]string{
		"PHP_DISPLAY_ERRORS":              "on",
		"PHP_MEMORY_LIMIT":                "512M",
		"PHP_MAX_EXECUTION_TIME":          "5000",
		"PHP_UPLOAD_MAX_FILESIZE":         "512M",
		"PHP_MAX_INPUT_VARS":              "5000",
		"PHP_POST_MAX_SIZE":               "512M",
		"PHP_OPCACHE_ENABLE":              "0",
		"PHP_OPCACHE_REVALIDATE_FREQ":     "0",
		"PHP_OPCACHE_VALIDATE_TIMESTAMPS": "0",
		"XDEBUG_MODE":                     "off",
		"XDEBUG_SESSION":                  "PHPSTORM",
		"XDEBUG_CONFIG":                   "",
		"BLACKFIRE_SERVER_ID":             "",
		"BLACKFIRE_SERVER_TOKEN":          "",
	}
)

// Config represents the nitro-dev.yaml users add for local development.
type Config struct {
	Containers []Container `json:"containers,omitempty" yaml:"containers,omitempty"`
	Blackfire  Blackfire   `json:"blackfire,omitempty" yaml:"blackfire,omitempty"`
	Databases  []Database  `json:"databases,omitempty" yaml:"databases,omitempty"`
	Services   Services    `json:"services" yaml:"services"`
	Sites      []Site      `json:"sites,omitempty" yaml:"sites,omitempty"`
	Apps       []App       `yaml:"apps,omitempty"`

	// ParsedApps is the representation of the apps loading from config files and the users global nitro
	ParsedApps    []App  `yaml:"-"`
	File          string `json:"-" yaml:"-"`
	HomeDirectory string `yaml:"-"`
}

// AllSitesWithHostnames takes the address, which is the nitro-proxy
// ip address, and the current site and returns a list of all the
func (c *Config) AllSitesWithHostnames(site Site, addr string) map[string][]string {
	hostnames := make(map[string][]string)
	for _, s := range c.Sites {
		// don't add the current site, since we can use the 127.0.0.1 address
		if site.Hostname == s.Hostname {
			continue
		}

		// add the sites hostname and aliases to the list
		hostnames[addr] = append(s.Aliases, s.Hostname)
	}

	return hostnames
}

// FindContainerByName takes a name and returns the container if name matches.
func (c *Config) FindContainerByName(name string) (*Container, error) {
	// find the site by the hostname
	for _, c := range c.Containers {
		if c.Name == name {
			return &c, nil
		}
	}

	return nil, fmt.Errorf("unable to find container with name %s", name)
}

// FindAppByHostname takes a hostname and returns the site if the hostnames match.
func (c *Config) FindAppByHostname(hostname string) (*App, error) {
	// find the app by the hostname
	for _, a := range c.ParsedApps {
		if a.Hostname == hostname {
			return &a, nil
		}
	}

	return nil, fmt.Errorf("unable to find app with hostname %s", hostname)
}

// FindSiteByHostName takes a hostname and returns the site if the hostnames match.
func (c *Config) FindSiteByHostName(hostname string) (*Site, error) {
	// find the site by the hostname
	for _, s := range c.Sites {
		if s.Hostname == hostname {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("unable to find site with hostname %s", hostname)
}

// ListOfSitesByDirectory takes the user’s home directory and the current
// working directory and returns a list of sites within that context.
func (c *Config) ListOfSitesByDirectory(home, wd string) []Site {
	var found []Site

	// Collect the broadest-possible options: sites having container paths
	// within the working directory.
	for _, s := range c.Sites {
		p, _ := s.GetAbsContainerPath(home)

		if strings.Contains(p, wd) {
			found = append(found, s)
		}
	}

	// Narrow matches further if possible and return the subset.
	if len(found) > 0 {
		var exactMatches []Site
		var bestMatches = 0
		var b Site
		var m = 0

		for _, s := range found {
			containerPath := s.GetContainerPath()
			p, _ := s.GetAbsContainerPath(home)
			segments := strings.Split(containerPath, "/")
			var matchingSegments = len(segments)

			if wd == p {
				// Append because more than one site *could* use the exact same path.
				exactMatches = append(exactMatches, s)
			} else if matchingSegments > 1 && matchingSegments > m {
				bestMatches += 1
				m = matchingSegments
				b = s
			}
		}

		// Return sites with container path(s) matching the working directory.
		if len(exactMatches) > 0 {
			return exactMatches
		}

		// Return the single, most specific found item if we have one.
		if bestMatches == 1 {
			return []Site{b}
		}

		return found
	}

	// Otherwise just return all the sites.
	return c.Sites
}

type App struct {
	Config     string   `yaml:"config,omitempty"`
	Hostname   string   `yaml:"hostname,omitempty"`
	Aliases    []string `yaml:"aliases,omitempty"`
	Path       string   `yaml:"path,omitempty"`
	Webroot    string   `yaml:"webroot,omitempty"`
	PHPVersion string   `yaml:"php_version,omitempty"`
	Dockerfile bool     `yaml:"dockerfile,omitempty"`
	Excludes   []string `yaml:"excludes,omitempty"`
	PHP        PHP      `yaml:"php,omitempty"`
	Extensions []string `yaml:"extensions,omitempty"`
	Xdebug     bool     `yaml:"xdebug,omitempty"`
	Blackfire  bool     `yaml:"blackfire,omitempty"`
	Suspended  bool     `yaml:"suspended,omitempty"`
	Database   struct {
		Engine  string `yaml:"engine,omitempty"`
		Version string `yaml:"version,omitempty"`
	} `yaml:"database,omitempty"`
}

func (c *Config) AddApp(app App) error {
	c.Apps = append(c.Apps, app)
	return nil
}

// Blackfire allows users to setup their containers to use blackfire locally.
type Blackfire struct {
	ClientID    string `json:"client_id,omitempty" yaml:"client_id,omitempty"`
	ClientToken string `json:"client_token,omitempty" yaml:"client_token,omitempty"`
	ServerID    string `json:"server_id,omitempty" yaml:"server_id,omitempty"`
	ServerToken string `json:"server_token,omitempty" yaml:"server_token,omitempty"`
}

// Container represents a custom container to add to nitro. Containers can be
// publicly hosted on Docker Hub.
type Container struct {
	// Name is a unique name, with no spaces or special characters and is used to generate the hostname
	Name string `json:"name" yaml:"name"`

	// Image the is canonical name of the image to use for the container `docker.elastic.co/elasticsearch/elasticsearch`
	Image string `json:"image" yaml:"image"`

	// Tag tells Nitro which docker image tag to use, it defaults to latest.
	Tag string `json:"tag" yaml:"tag"`

	// Ports is used to expose ports on the host machine to the
	// containers port in the <host>:<container> syntax
	Ports []string `json:"ports,omitempty" yaml:"ports,omitempty"`

	// Volume stores the volumes we should create and maintain for the container (e.g. <name>_container_<vol>_nitro_volume)
	Volumes []string `json:"volumes,omitempty" yaml:"volumes,omitempty"`

	WebGui  int    `json:"web_gui,omitempty" yaml:"web_gui,omitempty"`
	EnvFile string `json:"env_file,omitempty" yaml:"env_file,omitempty"`
}

// AddContainer adds a new container config to an config. It will validate there are no other
// container names to avoid colision or duplicate ports.
func (c *Config) AddContainer(container Container) error {
	for _, e := range c.Containers {
		// check the containers name
		if e.Name == container.Name {
			return fmt.Errorf("an existing container %q already exists", e.Name)
		}

		// TODO(jasonmccallister) check is a port already is in use
	}

	c.Containers = append(c.Containers, container)

	sort.SliceStable(c.Containers, func(i, j int) bool {
		return c.Containers[i].Name < c.Containers[j].Name
	})

	return nil
}

// Database is the struct used to represent a database engine
// that is a combination of a engine (e.g. mariadb, mysql, or
// postgres), the version number, and the port. The engine
// and version are directly related to the official docker
// images on the docker hub.
type Database struct {
	Engine  string `json:"engine" yaml:"engine"`
	Version string `json:"version" yaml:"version"`
	Port    string `json:"port" yaml:"port"`
}

// GetHostname returns a friendly and predictable name for a database
// container. It is used for accessing a database by hostname. For
// example, mysql-8.0-3306 would be the hostname to use in the .env
// for DB_HOST.
func (d *Database) GetHostname() (string, error) {
	if d.Engine == "" || d.Version == "" || d.Port == "" {
		return "", fmt.Errorf("the engine, version, and port must be defined for the database")
	}

	return fmt.Sprintf("%s-%s-%s.database.nitro", d.Engine, d.Version, d.Port), nil
}

// Services define common tools for development that should run as containers. We don't expose the volumes, ports, and
// networking options for these types of services. We plan to support "custom" container options to make local users
// development even better.
type Services struct {
	Blackfire bool `json:"blackfire"`
	DynamoDB  bool `json:"dynamodb"`
	Mailhog   bool `json:"mailhog"`
	Minio     bool `json:"minio"`
	Redis     bool `json:"redis"`
}

// Site represents a web application. It has a hostname, aliases (which
// are alternate domains), the local path to the site, additional mounts
// to add to the container, and the directory the index.php is located.
type Site struct {
	Hostname   string   `json:"hostname" yaml:"hostname"`
	Aliases    []string `json:"aliases,omitempty" yaml:"aliases,omitempty"`
	Path       string   `json:"path" yaml:"path"`
	Version    string   `json:"version" yaml:"version"`
	PHP        PHP      `json:"php,omitempty" yaml:"php,omitempty"`
	Extensions []string `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	Webroot    string   `json:"webroot" yaml:"webroot"`
	Xdebug     bool     `json:"xdebug" yaml:"xdebug"`
	Blackfire  bool     `json:"blackfire" yaml:"blackfire"`
	Excludes   []string `json:"excludes,omitempty" yaml:"excludes,omitempty"`
}

// GetAbsPath gets the directory for a site.Path,
// It is used to create the mount for a sites
// container.
func (s *Site) GetAbsPath(home string) (string, error) {
	return cleanPath(home, s.Path)
}

// GetAbsContainerPath gets the directory for a site’s
// container path.
func (s *Site) GetAbsContainerPath(home string) (string, error) {
	return cleanPath(home, s.Path+string(os.PathSeparator)+s.GetContainerPath())
}

// GetContainerPath is responsible for looking at the
// site’s web root and determining the correct path in the
// container. This is used for the craft and queue
// commands to identify the location of the "craft"
// executable.
func (s *Site) GetContainerPath() string {
	// trim trailing slashes
	webroot := strings.TrimRight(s.Webroot, "/")

	// is there a path separator?
	if strings.Contains(webroot, "/") {
		parts := strings.Split(webroot, "/")

		if len(parts) >= 2 {
			return strings.Join(parts[:len(parts)-1], "/")
		}
	}

	return ""
}

// AsEnvs takes a gateway addr and turns specific options
// such as PHP settings into env vars that can be set on the
// containers environment
func (s *Site) AsEnvs(addr string) []string {
	var envs []string

	// set the php vars
	envs = append(envs, phpVars(s.PHP, s.Version)...)

	return append(envs, xdebugVars(s.PHP, s.Xdebug, s.Version, s.Hostname, addr)...)
}

// SetPHPBoolSetting is used to set php settings that are bool. It will look
// for the site by its hostname and change the setting. If it cannot find the
// site or setting it will return an error.
func (c *Config) SetPHPBoolSetting(hostname, setting string, value bool) error {
	for i, s := range c.Sites {
		if s.Hostname == hostname {
			switch setting {
			case "display_errors":
				c.Sites[i].PHP.DisplayErrors = value

				return nil
			case "opcache_enable":
				c.Sites[i].PHP.OpcacheEnable = value

				return nil
			case "opcache_validate_timestamps":
				c.Sites[i].PHP.OpcacheValidateTimestamps = value

				return nil
			default:
				return fmt.Errorf("unknown php setting %s", setting)
			}
		}
	}

	return fmt.Errorf("unable to find the site: %s", hostname)
}

// SetAppAliases is used to add an alias domain to an app. If
// the app cannot be found or the alias is already set it
// will return an error.
func (c *Config) SetAppAliases(hostname, alias string) error {

	// check the parsed apps for the hostname, but set the alias on the app index
	for i, a := range c.ParsedApps {
		// if it's not the right hostname
		if a.Hostname != hostname {
			continue
		}

		// check how many aliases are set
		switch len(c.ParsedApps[i].Aliases) == 0 {
		case false:
			for _, a := range c.ParsedApps[i].Aliases {
				// make sure it's not already set
				if a == alias {
					return fmt.Errorf("alias %s is already set for %s", alias, hostname)
				}

				// add the alias
				c.Apps[i].Aliases = append(c.Apps[i].Aliases, alias)

				// sort aliases
				sort.Strings(c.Apps[i].Aliases)

				return nil
			}
		default:
			c.Apps[i].Aliases = append(c.Apps[i].Aliases, alias)

			return nil
		}
	}

	return fmt.Errorf("unable to find the app %s", hostname)
}

// SetSiteAlias is used to add an alias domain to a site. If
// the site cannot be found or the alias is already set it
// will return an error.
func (c *Config) SetSiteAlias(hostname, alias string) error {
	for i, s := range c.Sites {
		// if its not the right hostname
		if s.Hostname != hostname {
			continue
		}

		// check how many aliases are set
		switch len(c.Sites[i].Aliases) == 0 {
		case false:
			for _, a := range c.Sites[i].Aliases {
				// make sure its not already set
				if a == alias {
					return fmt.Errorf("alias %s is already set for %s", alias, hostname)
				}

				// add the alias
				c.Sites[i].Aliases = append(c.Sites[i].Aliases, alias)

				// sort aliases
				sort.Strings(c.Sites[i].Aliases)

				return nil
			}
		default:
			c.Sites[i].Aliases = append(c.Sites[i].Aliases, alias)

			return nil
		}
	}

	return fmt.Errorf("unable to find the site: %s", hostname)
}

// SetPHPExtension is used to set php settings that are bool. It will look
// for the site by its hostname and change the setting. If it cannot find the
// site or setting it will return an error.
func (c *Config) SetPHPExtension(hostname, extension string) error {
	for i, s := range c.Sites {
		if s.Hostname == hostname {
			// if the extension is already set, we return an error
			for _, e := range c.Sites[i].Extensions {
				if e == extension {
					return fmt.Errorf("extension %s is already set for %s", extension, hostname)
				}
			}

			// add the extension to the list
			c.Sites[i].Extensions = append(c.Sites[i].Extensions, extension)

			// sort the extensions by alpha
			sort.Strings(c.Sites[i].Extensions)

			return nil
		}
	}

	return fmt.Errorf("unable to find the site: %s", hostname)
}

// SetPHPIntSetting is used to set php settings that are ints. It will look
// for the site by its hostname and change the setting. If it cannot find the
// site or setting it will return an error.
func (c *Config) SetPHPIntSetting(hostname, setting string, value int) error {
	for i, s := range c.Sites {
		if s.Hostname == hostname {
			switch setting {
			case "max_execution_time":
				c.Sites[i].PHP.MaxExecutionTime = value

				return nil
			case "max_input_vars":
				c.Sites[i].PHP.MaxInputVars = value

				return nil
			case "max_input_time":
				c.Sites[i].PHP.MaxInputTime = value

				return nil
			case "opcache_revalidate_freq":
				c.Sites[i].PHP.OpcacheRevalidateFreq = value

				return nil
			default:
				return fmt.Errorf("unknown php setting %s", setting)
			}
		}
	}

	return fmt.Errorf("unable to find the site: %s", hostname)
}

// SetPHPStrSetting is used to set php settings that are strings. It will look
// for the site by its hostname and change the setting. If it cannot find the
// site or setting it will return an error.
func (c *Config) SetPHPStrSetting(hostname, setting, value string) error {
	for i, s := range c.Sites {
		if s.Hostname == hostname {
			switch setting {
			case "post_max_size":
				c.Sites[i].PHP.PostMaxSize = value

				return nil
			case "max_file_upload":
				c.Sites[i].PHP.MaxFileUpload = value

				return nil
			case "memory_limit":
				c.Sites[i].PHP.MemoryLimit = value

				return nil
			case "upload_max_file_size":
				c.Sites[i].PHP.UploadMaxFileSize = value

				return nil
			default:
				return fmt.Errorf("unknown php setting %s", setting)
			}
		}
	}

	return fmt.Errorf("unable to find the site: %s", hostname)
}

// PHP is nested in a configuration and allows setting environment variables
// for sites to override in the local development environment.
type PHP struct {
	DisplayErrors             bool   `json:"display_errors,omitempty" yaml:"display_errors,omitempty"`
	MaxExecutionTime          int    `json:"max_execution_time,omitempty" yaml:"max_execution_time,omitempty"`
	MaxInputVars              int    `json:"max_input_vars,omitempty" yaml:"max_input_vars,omitempty"`
	MaxInputTime              int    `json:"max_input_time,omitempty" yaml:"max_input_time,omitempty"`
	MaxFileUpload             string `json:"max_file_upload,omitempty" yaml:"max_file_upload,omitempty"`
	MemoryLimit               string `json:"memory_limit,omitempty" yaml:"memory_limit,omitempty"`
	OpcacheEnable             bool   `json:"opcache_enable,omitempty" yaml:"opcache_enable,omitempty"`
	OpcacheRevalidateFreq     int    `json:"opcache_revalidate_freq,omitempty" yaml:"opcache_revalidate_freq,omitempty"`
	OpcacheValidateTimestamps bool   `json:"opcache_validate_timestamps,omitempty" yaml:"opcache_validate_timestamps,omitempty"`
	PostMaxSize               string `json:"post_max_size,omitempty" yaml:"post_max_size,omitempty"`
	UploadMaxFileSize         string `json:"upload_max_file_size,omitempty" yaml:"upload_max_file_size,omitempty"`
}

// Load is used to return the unmarshalled config, and
// returns an error when trying to get the users home directory or
// while marshalling the config.
func Load(home string) (*Config, error) {
	file, err := IsEmpty(home)
	if err != nil {
		return nil, err
	}

	// create the config
	c := &Config{
		File:          file,
		HomeDirectory: home,
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

	// does the config contain sites?
	if len(c.Sites) > 0 {
		return nil, ErrDeprecatedConfig
	}

	// copy the user config into the parsed apps as a starting point
	c.ParsedApps = make([]App, len(c.Apps))

	// load each of the apps
	for i, global := range c.Apps {
		// if there is a config file, load it
		if global.Config != "" {
			// set the config
			c.ParsedApps[i].Config = global.Config

			// load the file
			p, err := paths.Clean(c.HomeDirectory, global.Config)
			if err != nil {
				return nil, err
			}

			// lots of things rely on the app path, so get the path from the config file and set it
			path, _ := filepath.Split(p)
			c.ParsedApps[i].Path = paths.MakeRelative(home, path, true)

			// read the local app config file
			local, err := unmarshalAppConfigFrom(p)
			if err != nil {
				return nil, err
			}

			// parse the values but global values override the global config if present

			// check the hostname
			if hostname, err := parseStringValue(global, local, "hostname"); err == nil {
				c.ParsedApps[i].Hostname = hostname
			} else if err != nil {
				return c, err
			}

			// check the aliases
			if global.Aliases != nil {
				c.ParsedApps[i].Aliases = global.Aliases
			} else if local.Aliases != nil {
				c.ParsedApps[i].Aliases = local.Aliases
			}

			// check the webroot
			if webroot, err := parseStringValue(global, local, "webroot"); err == nil {
				c.ParsedApps[i].Webroot = webroot
			} else if err != nil {
				return c, err
			}

			// check the php_version
			if phpVersion, err := parseStringValue(global, local, "php_version"); err == nil {
				c.ParsedApps[i].PHPVersion = phpVersion
			} else if err != nil {
				return c, err
			}

			// check the dockerfile
			if global.Dockerfile != local.Dockerfile {
				c.ParsedApps[i].Dockerfile = global.Dockerfile
			} else {
				c.ParsedApps[i].Dockerfile = local.Dockerfile
			}

			// TODO(jasonmccallister) check the php settings

			// check the php extensions
			if global.Extensions != nil {
				c.ParsedApps[i].Extensions = global.Extensions
			} else if local.Extensions != nil {
				c.ParsedApps[i].Extensions = local.Extensions
			}

			// check xdebug
			if global.Xdebug != local.Xdebug {
				c.ParsedApps[i].Xdebug = global.Xdebug
			} else {
				c.ParsedApps[i].Xdebug = local.Xdebug
			}

			// check blackfire
			if global.Blackfire != local.Blackfire {
				c.ParsedApps[i].Blackfire = global.Blackfire
			} else {
				c.ParsedApps[i].Blackfire = local.Blackfire
			}

			// check suspend
			if global.Suspended != local.Suspended {
				c.ParsedApps[i].Suspended = global.Suspended
			} else {
				c.ParsedApps[i].Suspended = local.Suspended
			}

			// check the database engine
			if global.Database.Engine != "" {
				c.ParsedApps[i].Database.Engine = global.Database.Engine
			} else if local.Database.Engine != "" {
				c.ParsedApps[i].Database.Engine = local.Database.Engine
			}

			// check the database version
			if global.Database.Version != "" {
				c.ParsedApps[i].Database.Version = global.Database.Version
			} else if local.Database.Version != "" {
				c.ParsedApps[i].Database.Version = local.Database.Version
			}

			break
		}

		// assign what we found
		c.ParsedApps[i] = c.Apps[i]
	}

	// return the config
	return c, nil
}

// IsEmpty is used to check if the config file is empty
func IsEmpty(home string) (string, error) {
	// verify the file exists
	file := filepath.Join(home, DirectoryName, FileName)
	stat, err := os.Stat(file)
	if os.IsNotExist(err) {
		return "", ErrNoConfigFile
	}

	// check if the file is empty
	if stat.Size() == 0 {
		return "", ErrEmptyFile
	}

	return file, nil
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

	sort.SliceStable(c.Sites, func(i, j int) bool {
		return c.Sites[i].Hostname < c.Sites[j].Hostname
	})

	return nil
}

// RemoveContainer takes a name and will remove the container by its
// name from the config file.
func (c *Config) RemoveContainer(container *Container) error {
	for k, v := range c.Containers {
		if container.Name == v.Name {
			c.Containers = append(c.Containers[:k], c.Containers[k+1:]...)
			return nil
		}
	}

	return fmt.Errorf("unknown container %q", container.Name)
}

// RemoveDatabase is used to destroy or remove a database
// engine from the config.
func (c *Config) RemoveDatabase(database Database) error {
	hostname, _ := database.GetHostname()

	for k, v := range c.Databases {
		h, _ := v.GetHostname()
		if hostname == h {
			c.Databases = append(c.Databases[:k], c.Databases[k+1:]...)
			return nil
		}
	}

	return fmt.Errorf("unknown database %q", hostname)
}

// RemoveSite takes a hostname and will remove the site by its
// hostname from the config file.
func (c *Config) RemoveSite(site *Site) error {
	for i, s := range c.Sites {
		if site.Hostname == s.Hostname {
			c.Sites = append(c.Sites[:i], c.Sites[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("unknown site %q", site.Hostname)
}

// DisableBlackfire takes a sites hostname and sets the blackfire option
// to false. If the site cannot be found, it returns an error.
func (c *Config) DisableBlackfire(site string) error {
	// find the site by the hostname
	for i, s := range c.Sites {
		if s.Hostname == site {
			// only toggle if the setting is true
			if c.Sites[i].Blackfire {
				c.Sites[i].Blackfire = false
			}

			return nil
		}
	}

	return fmt.Errorf("unknown site, %s", site)
}

// DisableXdebug takes an apps hostname and sets the xdebug option
// to false. If the app cannot be found, it returns an error.
func (c *Config) DisableXdebug(hostname string) error {
	// find the site by the hostname
	for i, a := range c.ParsedApps {
		if a.Hostname == hostname {
			// only toggle if the setting is true
			if c.Apps[i].Xdebug {
				c.Apps[i].Xdebug = false
			}

			return nil
		}
	}

	return fmt.Errorf("unknown app, %s", hostname)
}

// EnableBlackfire takes an app hostname and sets the xdebug option
// to true. If the app cannot be found, it returns an error.
func (c *Config) EnableBlackfire(hostname string) error {
	// find the site by the hostname
	for i, a := range c.ParsedApps {
		if a.Hostname == hostname {
			if !c.Apps[i].Blackfire {
				c.Apps[i].Blackfire = true
			}

			return nil
		}
	}

	return fmt.Errorf("unknown app, %s", hostname)
}

// EnableXdebug takes an app hostname and sets the xdebug option
// to true. If the app cannot be found, it returns an error.
func (c *Config) EnableXdebug(hostname string) error {
	// find the site by the hostname
	for i, a := range c.ParsedApps {
		if a.Hostname == hostname {
			if !c.Apps[i].Xdebug {
				c.Apps[i].Xdebug = true
			}

			return nil
		}
	}

	return fmt.Errorf("unknown app, %s", hostname)
}

// Save takes a file path and marshals the config into a file.
func (c *Config) Save() error {
	// make sure the file exists
	if _, err := os.Stat(c.File); os.IsNotExist(err) {
		dir, _ := filepath.Split(c.File)

		if err := c.createFile(dir); err != nil {
			return err
		}
	}

	// open the file
	f, err := os.OpenFile(c.File, os.O_TRUNC|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	// unmarshal
	data, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	// write the content
	if _, err := f.Write(data); err != nil {
		return err
	}

	return f.Close()
}

func (c *Config) createFile(dir string) error {
	// create the .nitro directory if it does not exist
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
	if err := f.Chown(os.Geteuid(), os.Getuid()); err != nil {
		return nil
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
	var envs []string

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

	if php.OpcacheValidateTimestamps {
		envs = append(envs, "PHP_OPCACHE_VALIDATE_TIMESTAMPS="+DefaultEnvs["PHP_OPCACHE_VALIDATE_TIMESTAMPS"])
	} else {
		envs = append(envs, "PHP_OPCACHE_VALIDATE_TIMESTAMPS="+DefaultEnvs["PHP_OPCACHE_VALIDATE_TIMESTAMPS"])
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

func parseStringValue(global, local App, key string) (string, error) {
	switch key {
	case "hostname":
		if global.Hostname != "" {
			return global.Hostname, nil
		}

		if local.Hostname != "" {
			return local.Hostname, nil
		}

		return "", fmt.Errorf("hostname must be defined for an app")
	case "php_version":
		if global.PHPVersion != "" {
			return global.PHPVersion, nil
		}

		if local.PHPVersion != "" {
			return local.PHPVersion, nil
		}

		return "", fmt.Errorf("php_version must be defined for an app")
	case "webroot":
		if global.Webroot != "" {
			return global.Webroot, nil
		}

		if local.Webroot != "" {
			return local.Webroot, nil
		}

		return "", fmt.Errorf("webroot must be defined for an app")
	}

	return "", fmt.Errorf("unkown key %q provided", key)
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
