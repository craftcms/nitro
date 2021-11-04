package actions

import (
	"reflect"
	"testing"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/docker/docker/api/types/container"
)

func TestAppToContainerConfig(t *testing.T) {
	type args struct {
		app config.App
	}
	tests := []struct {
		name string
		args args
		want *container.Config
	}{
		{
			name: "uses the custom container image if defined",
			args: args{
				app: config.App{
					Hostname:   "mysite.nitro",
					Dockerfile: true,
					PHPVersion: "8.0",
					PHP: config.PHP{
						DisplayErrors:         true,
						MemoryLimit:           "256M",
						MaxExecutionTime:      3000,
						UploadMaxFileSize:     "128M",
						MaxInputVars:          2000,
						PostMaxSize:           "128M",
						OpcacheEnable:         true,
						OpcacheRevalidateFreq: 60,
					},
					Webroot: "web",
				},
			},
			want: &container.Config{
				Image: "mysite.nitro:local",
				Labels: map[string]string{
					"com.craftcms.nitro":            "true",
					"com.craftcms.nitro.host":       "mysite.nitro",
					"com.craftcms.nitro.webroot":    "web",
					"com.craftcms.nitro.type":       "app",
					"com.craftcms.nitro.dockerfile": "true",
				},
				Env: []string{
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
					"PHP_IDE_CONFIG=serverName=mysite.nitro",
					"XDEBUG_MODE=off",
				},
				Hostname: "mysite.nitro",
			},
		},
		{
			name: "uses the default container images",
			args: args{
				app: config.App{
					Hostname:   "mysite.nitro",
					PHPVersion: "8.0",
					PHP: config.PHP{
						DisplayErrors:         true,
						MemoryLimit:           "256M",
						MaxExecutionTime:      3000,
						UploadMaxFileSize:     "128M",
						MaxInputVars:          2000,
						PostMaxSize:           "128M",
						OpcacheEnable:         true,
						OpcacheRevalidateFreq: 60,
					},
					Webroot: "web",
				},
			},
			want: &container.Config{
				Image: "craftcms/nitro:8.0",
				Labels: map[string]string{
					"com.craftcms.nitro":         "true",
					"com.craftcms.nitro.host":    "mysite.nitro",
					"com.craftcms.nitro.webroot": "web",
					"com.craftcms.nitro.type":    "app",
				},
				Env: []string{
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
					"PHP_IDE_CONFIG=serverName=mysite.nitro",
					"XDEBUG_MODE=off",
				},
				Hostname: "mysite.nitro",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppToContainerConfig(tt.args.app); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppToContainerConfig() = \ngot:\n%v\n\nwant:\n%v", got, tt.want)
			}
		})
	}
}
