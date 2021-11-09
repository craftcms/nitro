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
	Hostname   string   `yaml:"hostname,omitempty"`
	Aliases    []string `yaml:"aliases,omitempty"`
	Path       string   `yaml:"path,omitempty"`
	Webroot    string   `yaml:"webroot,omitempty"`
	PHPVersion string   `yaml:"php_version,omitempty"`
	Dockerfile bool     `yaml:"dockerfile,omitempty"`
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

// Load is responsible for loading the nitro.yaml config file.
// It takes an optional home arg (for testing) and if the
// home arg is not provided it will use os.UserHomeDir to
// find the users home directory
func Load(home string) (*Config, error) {
	var h string
	if home == "" {
		var err error
		h, err = os.UserHomeDir()
		if err != nil {
			return nil, err
		}
	} else {
		h = home
	}

	// create the config and read from the file
	c := &Config{}
	f, err := ioutil.ReadFile(filepath.Join(h, ".nitro", "nitro.yaml"))
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(f, c); err != nil {
		return nil, err
	}

	// set the home directory in case we need it later
	c.HomeDir = h

	// load each of the apps
	for i, global := range c.Apps {
		// if there is a config file, load it
		if global.Config != "" {
			// load the file
			p, err := paths.Clean(c.HomeDir, global.Config)
			if err != nil {
				return nil, err
			}

			// read the local app config file
			local, err := unmarshalAppConfigFrom(p)
			if err != nil {
				return nil, err
			}

			// parse the values but global values override the global config if present

			// check the hostname
			if hostname, err := parseStringValue(global, local, "hostname"); err == nil {
				c.Apps[i].Hostname = hostname
			} else if err != nil {
				return c, err
			}

			// check the aliases
			if global.Aliases != nil {
				c.Apps[i].Aliases = global.Aliases
			} else if local.Aliases != nil {
				c.Apps[i].Aliases = local.Aliases
			}

			// check the webroot
			if webroot, err := parseStringValue(global, local, "webroot"); err == nil {
				c.Apps[i].Webroot = webroot
			} else if err != nil {
				return c, err
			}

			// check the php_version
			if phpVersion, err := parseStringValue(global, local, "php_version"); err == nil {
				c.Apps[i].PHPVersion = phpVersion
			} else if err != nil {
				return c, err
			}

			// check the dockerfile
			if global.Dockerfile != local.Dockerfile {
				c.Apps[i].Dockerfile = global.Dockerfile
			} else {
				c.Apps[i].Dockerfile = local.Dockerfile
			}

			// TODO(jasonmccallister) check the php settings

			// check the php extensions
			if global.Extensions != nil {
				c.Apps[i].Extensions = global.Extensions
			} else if local.Extensions != nil {
				c.Apps[i].Extensions = local.Extensions
			}

			// check xdebug
			if global.Xdebug != local.Xdebug {
				c.Apps[i].Xdebug = global.Xdebug
			} else {
				c.Apps[i].Xdebug = local.Xdebug
			}

			// check blackfire
			if global.Blackfire != local.Blackfire {
				c.Apps[i].Blackfire = global.Blackfire
			} else {
				c.Apps[i].Blackfire = local.Blackfire
			}

			// check suspend
			if global.Suspended != local.Suspended {
				c.Apps[i].Suspended = global.Suspended
			} else {
				c.Apps[i].Suspended = local.Suspended
			}

			// check the database engine
			if global.Database.Engine != "" {
				c.Apps[i].Database.Engine = global.Database.Engine
			} else if local.Database.Engine != "" {
				c.Apps[i].Database.Engine = local.Database.Engine
			}

			// check the database version
			if global.Database.Version != "" {
				c.Apps[i].Database.Version = global.Database.Version
			} else if local.Database.Version != "" {
				c.Apps[i].Database.Version = local.Database.Version
			}
		}
	}

	return c, nil
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

type Container struct {
	Name    string   `yaml:"name"`
	Image   string   `yaml:"image"`
	Tag     string   `yaml:"tag"`
	Ports   []string `yaml:"ports,omitempty"`
	WebUI   int      `yaml:"web_ui,omitempty"`
	EnvFile string   `yaml:"env_file,omitempty"`
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
