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
				"PHP_DISPLAY_ERRORS=off",
				"PHP_MEMORY_LIMIT=256M",
				"PHP_MAX_EXECUTION_TIME=3000",
				"PHP_UPLOAD_MAX_FILESIZE=128M",
				"PHP_MAX_INPUT_VARS=2000",
				"PHP_POST_MAX_SIZE=128M",
				"PHP_OPCACHE_ENABLE=1",
				"PHP_OPCACHE_REVALIDATE_FREQ=60",
				"PHP_OPCACHE_VALIDATE_TIMESTAMPS=0",
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
				"PHP_DISPLAY_ERRORS=off",
				"PHP_MEMORY_LIMIT=256M",
				"PHP_MAX_EXECUTION_TIME=3000",
				"PHP_UPLOAD_MAX_FILESIZE=128M",
				"PHP_MAX_INPUT_VARS=2000",
				"PHP_POST_MAX_SIZE=128M",
				"PHP_OPCACHE_ENABLE=1",
				"PHP_OPCACHE_REVALIDATE_FREQ=60",
				"PHP_OPCACHE_VALIDATE_TIMESTAMPS=0",
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
				"PHP_DISPLAY_ERRORS=off",
				"PHP_MEMORY_LIMIT=256M",
				"PHP_MAX_EXECUTION_TIME=3000",
				"PHP_UPLOAD_MAX_FILESIZE=128M",
				"PHP_MAX_INPUT_VARS=2000",
				"PHP_POST_MAX_SIZE=128M",
				"PHP_OPCACHE_ENABLE=1",
				"PHP_OPCACHE_REVALIDATE_FREQ=60",
				"PHP_OPCACHE_VALIDATE_TIMESTAMPS=0",
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
				"PHP_DISPLAY_ERRORS=on",
				"PHP_MEMORY_LIMIT=512M",
				"PHP_MAX_EXECUTION_TIME=5000",
				"PHP_UPLOAD_MAX_FILESIZE=512M",
				"PHP_MAX_INPUT_VARS=5000",
				"PHP_POST_MAX_SIZE=512M",
				"PHP_OPCACHE_ENABLE=0",
				"PHP_OPCACHE_REVALIDATE_FREQ=0",
				"PHP_OPCACHE_VALIDATE_TIMESTAMPS=0",
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
				File: filepath.Join(testdir, DirectoryName, FileName),
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

func TestConfig_SetPHPStrSetting(t *testing.T) {
	type fields struct {
		Sites []Site
	}
	type args struct {
		hostname string
		setting  string
		value    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Config
		wantErr bool
	}{
		// upload_max_file_size
		{
			name: "can change a sites php upload_max_file_size setting",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "siteone.nitro",
					},
				},
			},
			args: args{
				hostname: "siteone.nitro",
				setting:  "upload_max_file_size",
				value:    "24M",
			},
			wantErr: false,
		},
		{
			name: "can change a sites php post max size setting",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "siteone.nitro",
					},
				},
			},
			args: args{
				hostname: "siteone.nitro",
				setting:  "post_max_size",
				value:    "1024M",
			},
			wantErr: false,
		},
		{
			name: "can change a sites php max file upload setting",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "siteone.nitro",
					},
				},
			},
			args: args{
				hostname: "siteone.nitro",
				setting:  "max_file_upload",
				value:    "1024M",
			},
			wantErr: false,
		},
		{
			name: "can change a sites php memory limit setting",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "siteone.nitro",
					},
				},
			},
			args: args{
				hostname: "siteone.nitro",
				setting:  "memory_limit",
				value:    "1024M",
			},
			wantErr: false,
		},
		{
			name: "unknown settings return an error",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "siteone.nitro",
					},
				},
			},
			args: args{
				hostname: "siteone.nitro",
				setting:  "new_setting_who_dis",
				value:    "1024M",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Sites: tt.fields.Sites,
			}

			if err := c.SetPHPStrSetting(tt.args.hostname, tt.args.setting, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Config.SetPHPStrSetting() error = %v, wantErr %v", err, tt.wantErr)
			}

			// find the site
			var site Site
			for _, s := range c.Sites {
				if s.Hostname == tt.args.hostname {
					site = s
				}
			}

			switch tt.args.setting {
			case "memory_limit":
				if site.PHP.MemoryLimit != tt.args.value {
					t.Errorf("expected the setting to be %s, got %s", tt.args.value, site.PHP.MemoryLimit)
				}
			case "post_max_size":
				if site.PHP.PostMaxSize != tt.args.value {
					t.Errorf("expected the setting to be %s, got %s", tt.args.value, site.PHP.PostMaxSize)
				}
			case "max_file_upload":
				if site.PHP.MaxFileUpload != tt.args.value {
					t.Errorf("expected the setting to be %s, got %s", tt.args.value, site.PHP.MaxFileUpload)
				}
			}
		})
	}
}

func TestConfig_SetPHPBoolSetting(t *testing.T) {
	type fields struct {
		Sites []Site
	}
	type args struct {
		hostname string
		setting  string
		value    bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "can change a sites php opcache enable setting",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "siteone.nitro",
					},
				},
			},
			args: args{
				hostname: "siteone.nitro",
				setting:  "opcache_enable",
				value:    true,
			},
			wantErr: false,
		},
		{
			name: "can change a sites php post max size setting",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "siteone.nitro",
					},
				},
			},
			args: args{
				hostname: "siteone.nitro",
				setting:  "display_errors",
				value:    false,
			},
			wantErr: false,
		},
		{
			name: "unknown settings return an error",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "siteone.nitro",
					},
				},
			},
			args: args{
				hostname: "siteone.nitro",
				setting:  "new_setting_who_dis",
				value:    false,
			},
			wantErr: true,
		},
		{
			name: "missing site returns an error",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "siteone.nitro",
					},
				},
			},
			args: args{
				hostname: "newsite.nitro",
				setting:  "display_errors",
				value:    false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Sites: tt.fields.Sites,
			}

			if err := c.SetPHPBoolSetting(tt.args.hostname, tt.args.setting, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Config.SetPHPBoolSetting() error = %v, wantErr %v", err, tt.wantErr)
			}

			// find the site
			var site Site
			for _, s := range c.Sites {
				if s.Hostname == tt.args.hostname {
					site = s
				}
			}

			switch tt.args.setting {
			case "display_errors":
				if site.PHP.DisplayErrors != tt.args.value {
					t.Errorf("expected the setting to be %v, got %v", tt.args.value, site.PHP.DisplayErrors)
				}
			case "opcache_enable":
				if site.PHP.OpcacheEnable != tt.args.value {
					t.Errorf("expected the setting to be %v, got %v", tt.args.value, site.PHP.OpcacheEnable)
				}
			}
		})
	}
}

func TestConfig_SetPHPIntSetting(t *testing.T) {
	type fields struct {
		Sites []Site
	}
	type args struct {
		hostname string
		setting  string
		value    int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "can change a sites php max_execution_time setting",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "siteone.nitro",
					},
				},
			},
			args: args{
				hostname: "siteone.nitro",
				setting:  "max_execution_time",
				value:    7000,
			},
			wantErr: false,
		},
		{
			name: "can change a sites php max_input_vars setting",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "siteone.nitro",
					},
				},
			},
			args: args{
				hostname: "siteone.nitro",
				setting:  "max_input_vars",
				value:    13000,
			},
			wantErr: false,
		},
		{
			name: "can change a sites php max_input_time setting",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "siteone.nitro",
					},
				},
			},
			args: args{
				hostname: "siteone.nitro",
				setting:  "max_input_time",
				value:    13000,
			},
			wantErr: false,
		},
		{
			name: "can change a sites php opcache_revalidate_freq setting",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "siteone.nitro",
					},
				},
			},
			args: args{
				hostname: "siteone.nitro",
				setting:  "opcache_revalidate_freq",
				value:    30,
			},
			wantErr: false,
		},
		{
			name: "unknown settings return an error",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "siteone.nitro",
					},
				},
			},
			args: args{
				hostname: "siteone.nitro",
				setting:  "new_setting_who_dis",
				value:    0,
			},
			wantErr: true,
		},
		{
			name: "missing site returns an error",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "siteone.nitro",
					},
				},
			},
			args: args{
				hostname: "newsite.nitro",
				setting:  "display_errors",
				value:    0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Sites: tt.fields.Sites,
			}

			if err := c.SetPHPIntSetting(tt.args.hostname, tt.args.setting, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Config.SetPHPIntSetting() error = %v, wantErr %v", err, tt.wantErr)
			}

			// find the site
			var site Site
			for _, s := range c.Sites {
				if s.Hostname == tt.args.hostname {
					site = s
				}
			}

			switch tt.args.setting {
			case "max_execution_time":
				if site.PHP.MaxExecutionTime != tt.args.value {
					t.Errorf("expected the setting to be %v, got %v", tt.args.value, site.PHP.MaxExecutionTime)
				}
			case "max_input_vars":
				if site.PHP.MaxInputVars != tt.args.value {
					t.Errorf("expected the setting to be %v, got %v", tt.args.value, site.PHP.MaxInputVars)
				}
			case "max_input_time":
				if site.PHP.MaxInputTime != tt.args.value {
					t.Errorf("expected the setting to be %v, got %v", tt.args.value, site.PHP.MaxInputTime)
				}
			case "opcache_revalidate_freq":
				if site.PHP.OpcacheRevalidateFreq != tt.args.value {
					t.Errorf("expected the setting to be %v, got %v", tt.args.value, site.PHP.OpcacheRevalidateFreq)
				}
			}
		})
	}
}

func TestSite_GetContainerPath(t *testing.T) {
	type fields struct {
		Webroot string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "handles more than one level of nesting",
			fields: fields{
				Webroot: "another-dir/another-site/web",
			},
			want: "another-dir/another-site",
		},
		{
			name: "returns the correct directory",
			fields: fields{
				Webroot: "another-site/web",
			},
			want: "another-site",
		},
		{
			name: "returns the correct directory for a nest site that has a trailing slash",
			fields: fields{
				Webroot: "another-site/web/",
			},
			want: "another-site",
		},
		{
			name: "default web roots with a trailing slash return the correct value",
			fields: fields{
				Webroot: "web/",
			},
			want: "",
		},
		{
			name: "defaults returns an empty string",
			fields: fields{
				Webroot: "web",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Site{
				Webroot: tt.fields.Webroot,
			}

			if got := s.GetContainerPath(); got != tt.want {
				t.Errorf("Site.GetContainerPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_ListOfSitesByDirectory(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		Sites []Site
	}
	type args struct {
		home string
		wd   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Site
	}{
		{
			name: "all sites suggested in lieu of match",
			args: args{
				home: filepath.Join(wd),
				// we don’t have a site for `home/sites/broccoli`
				wd: filepath.Join(wd, "testdata", "home", "sites", "broccoli"),
			},
			fields: fields{
				Sites: []Site{
					{
						Webroot: "web",
						Path:    filepath.Join(wd, "testdata", "home", "sites", "apple"),
					},
					{
						Webroot: "public",
						Path:    filepath.Join(wd, "testdata", "home", "sites", "banana"),
					},
					{
						Webroot: "web",
						Path:    filepath.Join(wd, "testdata", "home", "sites", "cherry"),
					},
					{
						Webroot: "web",
						Path:    filepath.Join(wd, "testdata", "home", "sites", "cherry", "dragonfruit"),
					},
					{
						Webroot: "",
						Path:    filepath.Join(wd, "testdata", "home", "plugins", "thinginator"),
					},
				},
			},
			want: []Site{
				{
					Webroot: "web",
					Path:    filepath.Join(wd, "testdata", "home", "sites", "apple"),
				},
				{
					Webroot: "public",
					Path:    filepath.Join(wd, "testdata", "home", "sites", "banana"),
				},
				{
					Webroot: "web",
					Path:    filepath.Join(wd, "testdata", "home", "sites", "cherry"),
				},
				{
					Webroot: "web",
					Path:    filepath.Join(wd, "testdata", "home", "sites", "cherry", "dragonfruit"),
				},
				{
					Webroot: "",
					Path:    filepath.Join(wd, "testdata", "home", "plugins", "thinginator"),
				},
			},
		},
		{
			name: "multiple suggestions when working directory is top-level path",
			args: args{
				home: filepath.Join(wd),
				wd:   filepath.Join(wd, "testdata", "home", "sites"),
			},
			fields: fields{
				Sites: []Site{
					{
						Webroot: "web",
						Path:    filepath.Join(wd, "testdata", "home", "sites", "apple"),
					},
					{
						Webroot: "public",
						Path:    filepath.Join(wd, "testdata", "home", "sites", "banana"),
					},
					// this site is an exact path match, but our working directory is too vague
					// for it to be the only suggestion
					{
						Webroot: "cherry/web",
						Path:    filepath.Join(wd, "testdata", "home", "sites"),
					},
					// this site is in a different top-level directory and should not be a suggestion
					{
						Webroot: "",
						Path:    filepath.Join(wd, "testdata", "home", "plugins"),
					},
				},
			},
			want: []Site{
				{
					Webroot: "web",
					Path:    filepath.Join(wd, "testdata", "home", "sites", "apple"),
				},
				{
					Webroot: "public",
					Path:    filepath.Join(wd, "testdata", "home", "sites", "banana"),
				},
				{
					Webroot: "cherry/web",
					Path:    filepath.Join(wd, "testdata", "home", "sites"),
				},
			},
		},
		{
			name: "single suggestion when working directory is site path",
			args: args{
				home: filepath.Join(wd),
				wd:   filepath.Join(wd, "testdata", "home", "sites", "banana"),
			},
			fields: fields{
				Sites: []Site{
					{
						Webroot: "web",
						Path:    filepath.Join(wd, "testdata", "home", "sites", "apple"),
					},
					{
						Webroot: "public",
						Path:    filepath.Join(wd, "testdata", "home", "sites", "banana"),
					},
				},
			},
			want: []Site{
				{
					Path:    filepath.Join(wd, "testdata", "home", "sites", "banana"),
					Webroot: "public",
				},
			},
		},
		{
			name: "single suggestion when working directory is single exact match to site container path",
			args: args{
				home: filepath.Join(wd),
				wd:   filepath.Join(wd, "testdata", "home", "sites", "cherry"),
			},
			fields: fields{
				Sites: []Site{
					{
						Webroot: "web",
						Path:    filepath.Join(wd, "testdata", "home", "sites", "apple"),
					},
					{
						Webroot: "public",
						Path:    filepath.Join(wd, "testdata", "home", "sites", "banana"),
					},
					{
						Webroot: "cherry/web",
						Path:    filepath.Join(wd, "testdata", "home", "sites"),
					},
					{
						Webroot: "cherry/dragonfruit/web",
						Path:    filepath.Join(wd, "testdata", "home", "sites"),
					},
				},
			},
			want: []Site{
				{
					Webroot: "cherry/web",
					Path:    filepath.Join(wd, "testdata", "home", "sites"),
				},
			},
		},
		{
			name: "single suggestion when path values match but container directory is unique",
			args: args{
				home: filepath.Join(wd),
				wd:   filepath.Join(wd, "testdata", "home", "sites", "cherry"),
			},
			fields: fields{
				Sites: []Site{
					{
						Webroot: "apple/web",
						Path:    filepath.Join(wd, "testdata", "home", "sites"),
					},
					{
						Webroot: "banana/public",
						Path:    filepath.Join(wd, "testdata", "home", "sites"),
					},
					{
						Webroot: "cherry/web",
						Path:    filepath.Join(wd, "testdata", "home", "sites"),
					},
					{
						Webroot: "cherry/dragonfruit/web",
						Path:    filepath.Join(wd, "testdata", "home", "sites"),
					},
				},
			},
			want: []Site{
				{
					Path:    filepath.Join(wd, "testdata", "home", "sites"),
					Webroot: "cherry/web",
				},
			},
		},
		{
			name: "multiple suggestions when sites use same directory",
			args: args{
				home: filepath.Join(wd),
				wd:   filepath.Join(wd, "testdata", "home", "sites", "cherry"),
			},
			fields: fields{
				Sites: []Site{
					{
						Webroot:  "apple/web",
						Path:     filepath.Join(wd, "testdata", "home", "sites"),
						Hostname: "apple.nitro",
					},
					{
						Webroot:  "banana/public",
						Path:     filepath.Join(wd, "testdata", "home", "sites"),
						Hostname: "banana.nitro",
					},
					{
						Webroot:  "cherry/web",
						Path:     filepath.Join(wd, "testdata", "home", "sites"),
						Hostname: "cherry.nitro",
					},
					{
						Webroot:  "cherry/dragonfruit/web",
						Path:     filepath.Join(wd, "testdata", "home", "sites"),
						Hostname: "dragonfruit.nitro",
					},
					// this site uses the exact same path as cherry.nitro
					{
						Webroot:  "cherry/web",
						Path:     filepath.Join(wd, "testdata", "home", "sites"),
						Hostname: "doppelganger.nitro",
					},
				},
			},
			want: []Site{
				{
					Webroot:  "cherry/web",
					Path:     filepath.Join(wd, "testdata", "home", "sites"),
					Hostname: "cherry.nitro",
				},
				{
					Webroot:  "cherry/web",
					Path:     filepath.Join(wd, "testdata", "home", "sites"),
					Hostname: "doppelganger.nitro",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Sites: tt.fields.Sites,
			}
			if got := c.ListOfSitesByDirectory(tt.args.home, tt.args.wd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.ListOfSitesByDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_AllSitesWithHostnames(t *testing.T) {
	type fields struct {
		Containers []Container
		Blackfire  Blackfire
		Databases  []Database
		Services   Services
		Sites      []Site
		File       string
	}
	type args struct {
		site Site
		addr string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string][]string
	}{
		{
			name: "can get all of the sites with the address",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "example.com",
						Aliases:  []string{"example.net"},
					},
					{
						Hostname: "craftcms.com",
						Aliases:  []string{"craftcms.net"},
					},
				},
			},
			args: args{
				site: Site{
					Hostname: "example.com",
					Aliases:  []string{"example.net"},
				},
				addr: "127.0.0.1",
			},
			want: map[string][]string{
				"127.0.0.1": {"craftcms.net", "craftcms.com"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Containers: tt.fields.Containers,
				Blackfire:  tt.fields.Blackfire,
				Databases:  tt.fields.Databases,
				Services:   tt.fields.Services,
				Sites:      tt.fields.Sites,
				File:       tt.fields.File,
			}
			if got := c.AllSitesWithHostnames(tt.args.site, tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.AllSitesWithHostnames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_IsEmpty(t *testing.T) {
	type args struct {
		home string
		file string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantFile string
	}{
		{
			name: "empty file returns an error",
			args: args{
				home: filepath.Clean("testdata"),
				file: "empty.yaml",
			},
			wantErr: true,
		},
		{
			name: "non-empty file returns an error",
			args: args{
				home: filepath.Clean("testdata"),
			},
			wantErr:  false,
			wantFile: filepath.Join("testdata", DirectoryName, "nitro.yaml"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.file != "" {
				FileName = tt.args.file
			}

			file, err := IsEmpty(tt.args.home)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsEmpty() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantFile != file {
				t.Errorf("IsEmpty() file = %v, wantFile %v", file, tt.wantFile)
			}

			// set the filename back to the original
			FileName = "nitro.yaml"
		})
	}
}
