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
		Dir      string
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
				"XDEBUG_CONFIG=client_host=host.docker.internal start_with_request=yes discover_client_host=1",
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
				Dir:      tt.fields.Dir,
				Xdebug:   tt.fields.Xdebug,
			}
			if got := s.AsEnvs(tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Site.AsEnvs() = %v, want %v", got, tt.want)
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
			s := &Site{}
			got, err := s.cleanPath(tt.args.home, tt.args.path)
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
