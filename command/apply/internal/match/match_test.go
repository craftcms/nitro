package match

import (
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/labels"
)

func Test_checkEnvs(t *testing.T) {
	type args struct {
		site config.Site
		envs []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "xdebug returns false if disable on the site but not for the container",
			args: args{
				site: config.Site{
					Version: "7.4",
					Xdebug:  false,
				},
				envs: []string{
					"XDEBUG_MODE=debug,develop",
				},
			},
			want: false,
		},
		{
			name: "xdebug returns false if enabled on the site",
			args: args{
				site: config.Site{
					Version: "7.4",
					Xdebug:  true,
				},
				envs: []string{
					"XDEBUG_MODE=off",
				},
			},
			want: false,
		},
		{
			name: "opcache_revalidate_freq returns false",
			args: args{
				site: config.Site{
					PHP: config.PHP{
						OpcacheRevalidateFreq: 2,
					},
				},
				envs: []string{
					"PHP_OPCACHE_REVALIDATE_FREQ=1",
				},
			},
			want: false,
		},
		{
			name: "opcache_enable returns false",
			args: args{
				site: config.Site{
					PHP: config.PHP{
						OpcacheEnable: true,
					},
				},
				envs: []string{
					"PHP_OPCACHE_ENABLE=0",
				},
			},
			want: false,
		},
		{
			name: "post_max_size returns false",
			args: args{
				site: config.Site{
					PHP: config.PHP{
						PostMaxSize: "256M",
					},
				},
				envs: []string{
					"PHP_POST_MAX_SIZE=128M",
				},
			},
			want: false,
		},
		{
			name: "max_input_vars returns false",
			args: args{
				site: config.Site{
					PHP: config.PHP{
						MaxInputVars: 2000,
					},
				},
				envs: []string{
					"PHP_MAX_INPUT_VARS=10000",
				},
			},
			want: false,
		},
		{
			name: "upload_max_filesize returns false",
			args: args{
				site: config.Site{
					PHP: config.PHP{
						MaxFileUpload: "1024M",
					},
				},
				envs: []string{
					"PHP_UPLOAD_MAX_FILESIZE=2048M",
				},
			},
			want: false,
		},
		{
			name: "memory_limit returns false",
			args: args{
				site: config.Site{
					PHP: config.PHP{
						MemoryLimit: "1024M",
					},
				},
				envs: []string{
					"PHP_MEMORY_LIMIT=2048M",
				},
			},
			want: false,
		},
		{
			name: "max_execution_time returns false",
			args: args{
				site: config.Site{
					PHP: config.PHP{
						MaxExecutionTime: 10000,
					},
				},
				envs: []string{
					"PHP_MAX_EXECUTION_TIME=3000",
				},
			},
			want: false,
		},
		{
			name: "display_errors returns false",
			args: args{
				site: config.Site{
					PHP: config.PHP{
						DisplayErrors: false,
					},
				},
				envs: []string{
					"PHP_DISPLAY_ERRORS=off",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkEnvs(tt.args.site.PHP, tt.args.site.Xdebug, tt.args.envs); got != tt.want {
				t.Errorf("checkEnvs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSite(t *testing.T) {
	type args struct {
		home      string
		site      config.Site
		container types.ContainerJSON
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "hostname updates return false using labels",
			args: args{
				home: "testdata/example-site",
				site: config.Site{
					Hostname: "newname",
					Path:     "testdata/example-site",
					Version:  "7.4",
				},
				container: types.ContainerJSON{
					Config: &container.Config{
						Image: "docker.io/craftcms/nginx:7.4-dev",
						Labels: map[string]string{
							labels.Host: "oldname",
						},
					},
				},
			},
			want: false,
		},
		{
			name: "path mismatches return false",
			args: args{
				home: "testdata/example-site",
				site: config.Site{
					Path:    "testdata/new-site",
					Version: "7.4",
				},
				container: types.ContainerJSON{
					Config: &container.Config{
						Image: "docker.io/craftcms/nginx:7.4-dev",
					},
				},
			},
			want: false,
		},
		{
			name: "mismatched images return false",
			args: args{
				site: config.Site{
					Version: "7.4",
				},
				container: types.ContainerJSON{
					Config: &container.Config{
						Image: "docker.io/craftcms/nginx:7.3-dev",
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Site(tt.args.home, tt.args.site, tt.args.container); got != tt.want {
				t.Errorf("Site() = %v, want %v", got, tt.want)
			}
		})
	}
}
