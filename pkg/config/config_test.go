package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSite_AsEnvs(t *testing.T) {
	type fields struct {
		Hostname string
		Aliases  []string
		Path     string
		Version  string
		PHP      PHP
		Webroot  string
		Xdebug   bool
	}
	type args struct {
		addr string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "xdebug 2 options are set if enabled",
			fields: fields{
				Hostname: "somewebsite.nitro",
				PHP: PHP{
					DisplayErrors:         true,
					MemoryLimit:           "256M",
					MaxExecutionTime:      3000,
					UploadMaxFileSize:     "128M",
					MaxInputVars:          2000,
					PostMaxSize:           "128M",
					OpcacheEnable:         true,
					OpcacheRevalidateFreq: 60,
				},
				Version: "7.1",
				Xdebug:  true,
			},
			args: args{
				addr: "host.docker.internal",
			},
			want: []string{
				"COMPOSER_HOME=/tmp",
				"PHP_DISPLAY_ERRORS=off",
				"PHP_MEMORY_LIMIT=256M",
				"PHP_MAX_EXECUTION_TIME=3000",
				"PHP_UPLOAD_MAX_FILESIZE=128M",
				"PHP_MAX_INPUT_VARS=2000",
				"PHP_POST_MAX_SIZE=128M",
				"PHP_OPCACHE_ENABLE=1",
				"PHP_OPCACHE_REVALIDATE_FREQ=60",
				"XDEBUG_SESSION=PHPSTORM",
				"PHP_IDE_CONFIG=serverName=somewebsite.nitro",
				"XDEBUG_CONFIG=idekey=PHPSTORM remote_host=host.docker.internal profiler_enable=1 remote_port=9000 remote_autostart=1 remote_enable=1",
				"XDEBUG_MODE=xdebug2",
			},
		},
		{
			name: "xdebug 3 options are set if enabled",
			fields: fields{
				Hostname: "somewebsite.nitro",
				PHP: PHP{
					DisplayErrors:         true,
					MemoryLimit:           "256M",
					MaxExecutionTime:      3000,
					UploadMaxFileSize:     "128M",
					MaxInputVars:          2000,
					PostMaxSize:           "128M",
					OpcacheEnable:         true,
					OpcacheRevalidateFreq: 60,
				},
				Version: "7.4",
				Xdebug:  true,
			},
			args: args{
				addr: "host.docker.internal",
			},
			want: []string{
				"COMPOSER_HOME=/tmp",
				"PHP_DISPLAY_ERRORS=off",
				"PHP_MEMORY_LIMIT=256M",
				"PHP_MAX_EXECUTION_TIME=3000",
				"PHP_UPLOAD_MAX_FILESIZE=128M",
				"PHP_MAX_INPUT_VARS=2000",
				"PHP_POST_MAX_SIZE=128M",
				"PHP_OPCACHE_ENABLE=1",
				"PHP_OPCACHE_REVALIDATE_FREQ=60",
				"XDEBUG_SESSION=PHPSTORM",
				"PHP_IDE_CONFIG=serverName=somewebsite.nitro",
				"XDEBUG_CONFIG=client_host=host.docker.internal client_port=9003",
				"XDEBUG_MODE=develop,debug",
			},
		},
		{
			name: "defaults are overridden when set on the site",
			fields: fields{
				Hostname: "somewebsite.nitro",
				PHP: PHP{
					DisplayErrors:         true,
					MemoryLimit:           "256M",
					MaxExecutionTime:      3000,
					UploadMaxFileSize:     "128M",
					MaxInputVars:          2000,
					PostMaxSize:           "128M",
					OpcacheEnable:         true,
					OpcacheRevalidateFreq: 60,
				},
			},
			want: []string{
				"COMPOSER_HOME=/tmp",
				"PHP_DISPLAY_ERRORS=off",
				"PHP_MEMORY_LIMIT=256M",
				"PHP_MAX_EXECUTION_TIME=3000",
				"PHP_UPLOAD_MAX_FILESIZE=128M",
				"PHP_MAX_INPUT_VARS=2000",
				"PHP_POST_MAX_SIZE=128M",
				"PHP_OPCACHE_ENABLE=1",
				"PHP_OPCACHE_REVALIDATE_FREQ=60",
				"XDEBUG_SESSION=PHPSTORM",
				"PHP_IDE_CONFIG=serverName=somewebsite.nitro",
				"XDEBUG_MODE=off",
			},
		},
		{
			name: "can get the defaults that are expected",
			fields: fields{
				Hostname: "somewebsite.nitro",
			},
			want: []string{
				"COMPOSER_HOME=/tmp",
				"PHP_DISPLAY_ERRORS=on",
				"PHP_MEMORY_LIMIT=512M",
				"PHP_MAX_EXECUTION_TIME=5000",
				"PHP_UPLOAD_MAX_FILESIZE=512M",
				"PHP_MAX_INPUT_VARS=5000",
				"PHP_POST_MAX_SIZE=512M",
				"PHP_OPCACHE_ENABLE=0",
				"PHP_OPCACHE_REVALIDATE_FREQ=0",
				"XDEBUG_SESSION=PHPSTORM",
				"PHP_IDE_CONFIG=serverName=somewebsite.nitro",
				"XDEBUG_MODE=off",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Site{
				Hostname: tt.fields.Hostname,
				Aliases:  tt.fields.Aliases,
				Path:     tt.fields.Path,
				Version:  tt.fields.Version,
				PHP:      tt.fields.PHP,
				Webroot:  tt.fields.Webroot,
				Xdebug:   tt.fields.Xdebug,
			}
			if got := s.AsEnvs(tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Site.AsEnvs() = \ngot:\n%v, \nwant:\n%v", got, tt.want)
			}
		})
	}
}

func TestSite_cleanPath(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		home string
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "can get the full path using a tilde",
			args: args{
				home: wd,
				path: "~/testdata",
			},
			want:    filepath.Join(wd, "testdata"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cleanPath(tt.args.home, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Site.cleanPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Site.cleanPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_GetHostname(t *testing.T) {
	type fields struct {
		Engine    string
		Version   string
		Port      string
		Ephemeral bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name:    "can get the hostname for a database container",
			fields:  fields{Engine: "mysql", Version: "5.7", Port: "3306"},
			want:    "mysql-5.7-3306.database.nitro",
			wantErr: false,
		},
		{
			name:    "empty values return an error",
			fields:  fields{Engine: "mysql", Port: "3306"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				Engine:  tt.fields.Engine,
				Version: tt.fields.Version,
				Port:    tt.fields.Port,
			}
			got, err := d.GetHostname()
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.GetHostname() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Database.GetHostname() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	// get the working dir for the test path
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testdir := filepath.Join(wd, "testdata")

	type args struct {
		home string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "can load a config file",
			args: args{
				home: testdir,
			},
			want: &Config{
				File: filepath.Join(testdir, ".nitro", FileName),
				Blackfire: Blackfire{
					ServerID:    "my-id",
					ServerToken: "my-token",
				},
				Databases: []Database{
					{
						Engine:  "mysql",
						Version: "8.0",
						Port:    "3306",
					},
					{
						Engine:  "postgres",
						Version: "13",
						Port:    "5432",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing file returns an error",
			args: args{
				home: filepath.Join(testdir, "something"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.args.home)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				// check blackfire
				if !reflect.DeepEqual(got.Blackfire, tt.want.Blackfire) {
					t.Errorf("Load() = \ngot:\n%v,\nwant\n%v", got.Blackfire, tt.want.Blackfire)
				}

				// check databases
				if !reflect.DeepEqual(got.Databases, tt.want.Databases) {
					t.Errorf("Load() = \ngot:\n%v,\nwant\n%v", got.Databases, tt.want.Databases)
				}

				// check services
				if !reflect.DeepEqual(got.Services, tt.want.Services) {
					t.Errorf("Load() = \ngot:\n%v,\nwant\n%v", got.Services, tt.want.Services)
				}

				// check sites
				if !reflect.DeepEqual(got.Sites, tt.want.Sites) {
					t.Errorf("Load() = \ngot:\n%v,\nwant\n%v", got.Sites, tt.want.Sites)
				}

				t.Errorf("Load() = \ngot\n%v,\nwant\n%v", got, tt.want)
			}
		})
	}
}

func TestConfig_EnableXdebug(t *testing.T) {
	type fields struct {
		Blackfire Blackfire
		Databases []Database
		Services  Services
		Sites     []Site
		File      string
	}
	type args struct {
		site string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "struct is not modified when xdebug is already true",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "somesite",
						Xdebug:   true,
					},
					{
						Hostname: "anothersite",
						Xdebug:   true,
					},
				},
			},
			args:    args{site: "somesite"},
			wantErr: false,
		},
		{
			name: "can enable xdebug for a site",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "somesite",
						Xdebug:   false,
					},
					{
						Hostname: "anothersite",
						Xdebug:   true,
					},
				},
			},
			args:    args{site: "somesite"},
			wantErr: false,
		},
		{
			name:    "sites that don't exist return an error",
			args:    args{site: "idontexist"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Blackfire: tt.fields.Blackfire,
				Databases: tt.fields.Databases,
				Services:  tt.fields.Services,
				Sites:     tt.fields.Sites,
				File:      tt.fields.File,
			}
			if err := c.EnableXdebug(tt.args.site); (err != nil) != tt.wantErr {
				t.Errorf("Config.EnableXdebug() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_DisableXdebug(t *testing.T) {
	type fields struct {
		Blackfire Blackfire
		Databases []Database
		Services  Services
		Sites     []Site
		File      string
	}
	type args struct {
		site string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "can enable xdebug for a site",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "somesite",
						Xdebug:   false,
					},
					{
						Hostname: "anothersite",
						Xdebug:   true,
					},
				},
			},
			args:    args{site: "somesite"},
			wantErr: false,
		},
		{
			name:    "sites that don't exist return an error",
			args:    args{site: "idontexist"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Blackfire: tt.fields.Blackfire,
				Databases: tt.fields.Databases,
				Services:  tt.fields.Services,
				Sites:     tt.fields.Sites,
				File:      tt.fields.File,
			}
			if err := c.DisableXdebug(tt.args.site); (err != nil) != tt.wantErr {
				t.Errorf("Config.DisableXdebug() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_AddSite(t *testing.T) {
	type fields struct {
		Blackfire Blackfire
		Databases []Database
		Services  Services
		Sites     []Site
		File      string
	}
	type args struct {
		s Site
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "existing hostnames returns an error",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "example",
					},
				},
			},
			args: args{
				s: Site{Hostname: "example"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Blackfire: tt.fields.Blackfire,
				Databases: tt.fields.Databases,
				Services:  tt.fields.Services,
				Sites:     tt.fields.Sites,
				File:      tt.fields.File,
			}
			if err := c.AddSite(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("Config.AddSite() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSite_GetAbsPath(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		Hostname string
		Aliases  []string
		Path     string
		Version  string
		PHP      PHP
		Webroot  string
		Xdebug   bool
	}
	type args struct {
		home string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "existing paths return the complete path",
			fields: fields{
				Path: filepath.Join(wd, "testdata"),
			},
			args: args{
				home: wd,
			},
			want:    filepath.Join(wd, "testdata"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Site{
				Hostname: tt.fields.Hostname,
				Aliases:  tt.fields.Aliases,
				Path:     tt.fields.Path,
				Version:  tt.fields.Version,
				PHP:      tt.fields.PHP,
				Webroot:  tt.fields.Webroot,
				Xdebug:   tt.fields.Xdebug,
			}
			got, err := s.GetAbsPath(tt.args.home)
			if (err != nil) != tt.wantErr {
				t.Errorf("Site.GetAbsPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Site.GetAbsPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
