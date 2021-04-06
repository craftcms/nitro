package match

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func Test_checkEnvs(t *testing.T) {
	type args struct {
		site      config.Site
		blackfire config.Blackfire
		envs      []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "blackfire server token returns false if there are no credentials but the environment variables are set",
			args: args{
				site: config.Site{
					Version: "7.4",
				},
				envs: []string{
					"BLACKFIRE_SERVER_TOKEN=someid",
				},
			},
			want: false,
		},
		{
			name: "blackfire server token returns false if there are credentials but the environment variables are not set",
			args: args{
				site: config.Site{
					Version: "7.4",
				},
				envs: []string{
					"BLACKFIRE_SERVER_TOKEN=",
				},
				blackfire: config.Blackfire{
					ServerToken: "someid",
				},
			},
			want: false,
		},
		{
			name: "blackfire server id returns false if there are no credentials but the environment variables are set",
			args: args{
				site: config.Site{
					Version: "7.4",
				},
				envs: []string{
					"BLACKFIRE_SERVER_ID=someid",
				},
			},
			want: false,
		},
		{
			name: "blackfire server id returns false if there are credentials but the environment variables are not set",
			args: args{
				site: config.Site{
					Version: "7.4",
				},
				envs: []string{
					"BLACKFIRE_SERVER_ID=",
				},
				blackfire: config.Blackfire{
					ServerID: "someid",
				},
			},
			want: false,
		},
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
			name: "opcache_disable returns false",
			args: args{
				site: config.Site{
					PHP: config.PHP{
						OpcacheEnable: false,
					},
				},
				envs: []string{
					"PHP_OPCACHE_ENABLE=1",
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
			if got := checkEnvs(tt.args.site, tt.args.blackfire, tt.args.envs); got != tt.want {
				t.Errorf("checkEnvs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSite(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		home      string
		site      config.Site
		container types.ContainerJSON
		blackfire config.Blackfire
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "mismatched webroot returns false",
			args: args{
				home: "testdata/example-site",
				site: config.Site{
					Hostname: "newname",
					Path:     "testdata/example-site",
					Version:  "7.4",
					Webroot:  "web",
				},
				container: types.ContainerJSON{
					Config: &container.Config{
						Image: "docker.io/craftcms/nginx:7.4-dev",
						Labels: map[string]string{
							containerlabels.Host:    "newname",
							containerlabels.Webroot: "public",
						},
					},
					Mounts: []types.MountPoint{
						{
							Source: filepath.Join(wd, "testdata", "new-path"),
						},
					},
				},
			},
			want: false,
		},
		{
			name: "mismatched paths return false",
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
							containerlabels.Host: "newname",
						},
					},
					Mounts: []types.MountPoint{
						{
							Source: filepath.Join(wd, "testdata", "new-path"),
						},
					},
				},
			},
			want: false,
		},
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
							containerlabels.Host: "oldname",
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
			if got := Site(tt.args.home, tt.args.site, tt.args.container, tt.args.blackfire); got != tt.want {
				t.Errorf("Site() = %v, want %v", got, tt.want)
			}
		})
	}
}
