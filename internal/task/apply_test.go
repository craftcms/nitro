package task

import (
	"reflect"
	"testing"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

func TestApply(t *testing.T) {
	type args struct {
		machine             string
		configFile          config.Config
		fromMultipassMounts []config.Mount
		sites               []config.Site
		dbs                 []config.Database
		php                 string
	}
	tests := []struct {
		name    string
		args    args
		want    []nitro.Action
		wantErr bool
	}{
		{
			name: "changing a sites webroot will update the nginx configuration",
			args: args{
				machine: "mytestmachine",
				configFile: config.Config{
					PHP: "7.4",
					Mounts: []config.Mount{
						{
							Source: "./testdata/existing-mount",
							Dest:   "/nitro/sites/existing-site",
						},
					},
					Sites: []config.Site{
						{
							Hostname: "existing-site",
							Webroot:  "/nitro/sites/existing-site/public", // renamed in the config to public
						},
					},
				},
				fromMultipassMounts: []config.Mount{
					{
						Source: "./testdata/existing-mount",
						Dest:   "/nitro/sites/existing-site",
					},
				},
				sites: []config.Site{
					{
						Hostname: "existing-site",
						Webroot:  "/nitro/sites/existing-site/web",
					},
				},
				php: "7.4",
			},
			want: []nitro.Action{
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "rm", "/etc/nginx/sites-available/existing-site"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "service", "nginx", "restart"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "cp", "/opt/nitro/nginx/template.conf", "/etc/nginx/sites-available/existing-site"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "sed", "-i", "s|CHANGEWEBROOTDIR|/nitro/sites/existing-site/public|g", "/etc/nginx/sites-available/existing-site"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "sed", "-i", "s|CHANGESERVERNAME|existing-site|g", "/etc/nginx/sites-available/existing-site"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "sed", "-i", "s|CHANGEPHPVERSION|7.4|g", "/etc/nginx/sites-available/existing-site"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "ln", "-s", "/etc/nginx/sites-available/existing-site", "/etc/nginx/sites-enabled/"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "service", "nginx", "restart"},
				},
			},
			wantErr: false,
		},
		{
			name: "mismatched versions of PHP installs the request version from the config file",
			args: args{
				machine:    "mytestmachine",
				configFile: config.Config{PHP: "7.4"},
				php:        "7.2",
			},
			want: []nitro.Action{
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "apt-get", "install", "-y", "php7.4", "php7.4-mbstring", "php7.4-cli", "php7.4-curl", "php7.4-fpm", "php7.4-gd", "php7.4-intl", "php7.4-json", "php7.4-mysql", "php7.4-pgsql", "php7.4-zip", "php7.4-xml", "php7.4-soap", "php7.4-bcmath", "php7.4-gmp", "php-xdebug", "php-imagick", "blackfire-agent", "blackfire-php"},
				},
			},
		},
		{
			name: "new databases that are in the config are created",
			args: args{
				machine: "mytestmachine",
				configFile: config.Config{
					Databases: []config.Database{
						{
							Engine:  "mysql",
							Version: "5.7",
							Port:    "3306",
						},
					},
				},
				dbs: nil,
			},
			want: []nitro.Action{
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "docker", "volume", "create", "mysql_5.7_3306"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "docker", "run", "-v", "/home/ubuntu/.nitro/databases/mysql/setup.sql:/docker-entrypoint-initdb.d/setup.sql", "-v", "/home/ubuntu/.nitro/databases/mysql/conf.d/:/etc/mysql/conf.d", "-v", "mysql_5.7_3306:/var/lib/mysql", "--name", "mysql_5.7_3306", "-d", "--restart=always", "-p", "3306:3306", "-e", "MYSQL_ROOT_PASSWORD=nitro", "-e", "MYSQL_USER=nitro", "-e", "MYSQL_PASSWORD=nitro", "mysql:5.7"},
				},
			},
		},
		{
			name: "new databases are created but the ones in the config are kept",
			args: args{
				machine: "mytestmachine",
				configFile: config.Config{
					Databases: []config.Database{
						{
							Engine:  "mysql",
							Version: "5.7",
							Port:    "3306",
						},
						{
							Engine:  "postgres",
							Version: "11",
							Port:    "5432",
						},
					},
				},
				dbs: []config.Database{
					{
						Engine:  "mysql",
						Version: "5.7",
						Port:    "3306",
					},
					{
						Engine:  "postgres",
						Version: "11",
						Port:    "5432",
					},
					{
						Engine:  "postgres",
						Version: "12",
						Port:    "54321",
					},
				},
			},
			want: []nitro.Action{
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "docker", "rm", "-v", "postgres_12_54321", "-f"},
				},
			},
		},
		{
			name: "databases that are not in the config are removed",
			args: args{
				machine: "mytestmachine",
				configFile: config.Config{
					Databases: []config.Database{
						{
							Engine:  "mysql",
							Version: "5.7",
							Port:    "3306",
						},
					},
				},
				dbs: []config.Database{
					{
						Engine:  "mysql",
						Version: "5.7",
						Port:    "3306",
					},
					{
						Engine:  "postgres",
						Version: "11",
						Port:    "5432",
					},
				},
			},
			want: []nitro.Action{
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "docker", "rm", "-v", "postgres_11_5432", "-f"},
				},
			},
		},
		{
			name: "sites that exist but are not in the config file are removed",
			args: args{
				machine:    "mytestmachine",
				configFile: config.Config{},
				fromMultipassMounts: []config.Mount{
					{
						Source: "./testdata/existing/mount",
						Dest:   "/nitro/sites/leftoversite.test",
					},
				},
				sites: []config.Site{
					{
						Hostname: "leftoversite.test",
						Webroot:  "/nitro/sites/leftoversite.test/web",
					},
				},
			},
			want: []nitro.Action{
				{
					Type:       "umount",
					UseSyscall: false,
					Args:       []string{"umount", "mytestmachine:/nitro/sites/leftoversite.test"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "rm", "/etc/nginx/sites-available/leftoversite.test"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "service", "nginx", "restart"},
				},
			},
		},
		{
			name: "new sites without a parent mount in the config are added to the machine and mounted",
			args: args{
				machine: "mytestmachine",
				configFile: config.Config{
					PHP: "7.4",
					Mounts: []config.Mount{
						{
							Source: "./testdata/existing-mount",
							Dest:   "/nitro/sites/existing-site",
						},
					},
					Sites: []config.Site{
						{
							Hostname: "existing-site",
							Webroot:  "/nitro/sites/existing-site",
						},
						{
							Hostname: "new-site",
							Webroot:  "/nitro/sites/new-site",
						},
					},
				},
				fromMultipassMounts: []config.Mount{
					{
						Source: "./testdata/existing-mount",
						Dest:   "/nitro/sites/existing-site",
					},
				},
				sites: []config.Site{
					{
						Hostname: "existing-site",
						Webroot:  "/nitro/sites/existing-site",
					},
				},
				php: "7.4",
			},
			want: []nitro.Action{
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "cp", "/opt/nitro/nginx/template.conf", "/etc/nginx/sites-available/new-site"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "sed", "-i", "s|CHANGEWEBROOTDIR|/nitro/sites/new-site|g", "/etc/nginx/sites-available/new-site"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "sed", "-i", "s|CHANGESERVERNAME|new-site|g", "/etc/nginx/sites-available/new-site"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "sed", "-i", "s|CHANGEPHPVERSION|7.4|g", "/etc/nginx/sites-available/new-site"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "ln", "-s", "/etc/nginx/sites-available/new-site", "/etc/nginx/sites-enabled/"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "service", "nginx", "restart"},
				},
			},
			wantErr: false,
		},
		{
			name: "new sites using parent mounts in the config are added to the machine",
			args: args{
				machine: "mytestmachine",
				configFile: config.Config{
					PHP: "7.4",
					Mounts: []config.Mount{
						{
							Source: "./testdata/existing-mount",
							Dest:   "/nitro/sites",
						},
					},
					Sites: []config.Site{
						{
							Hostname: "existing-site",
							Webroot:  "/nitro/sites/existing-site",
						},
						{
							Hostname: "new-site",
							Webroot:  "/nitro/sites/new-site",
						},
					},
				},
				fromMultipassMounts: []config.Mount{
					{
						Source: "./testdata/existing-mount",
						Dest:   "/nitro/sites",
					},
				},
				sites: []config.Site{
					{
						Hostname: "existing-site",
						Webroot:  "/nitro/sites/existing-site",
					},
				},
				php: "7.4",
			},
			want: []nitro.Action{
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "cp", "/opt/nitro/nginx/template.conf", "/etc/nginx/sites-available/new-site"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "sed", "-i", "s|CHANGEWEBROOTDIR|/nitro/sites/new-site|g", "/etc/nginx/sites-available/new-site"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "sed", "-i", "s|CHANGESERVERNAME|new-site|g", "/etc/nginx/sites-available/new-site"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "sed", "-i", "s|CHANGEPHPVERSION|7.4|g", "/etc/nginx/sites-available/new-site"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "ln", "-s", "/etc/nginx/sites-available/new-site", "/etc/nginx/sites-enabled/"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mytestmachine", "--", "sudo", "service", "nginx", "restart"},
				},
			},
			wantErr: false,
		},
		{
			name: "new mounts return actions to create mounts",
			args: args{
				machine: "mytestmachine",
				configFile: config.Config{
					Mounts: []config.Mount{
						{
							Source: "./testdata/existing-mount",
							Dest:   "/nitro/sites/example-site",
						},
						{
							Source: "./testdata/new-mount",
							Dest:   "/nitro/sites/new-site",
						},
					},
				},
				fromMultipassMounts: []config.Mount{
					{
						Source: "./testdata/existing-mount",
						Dest:   "/nitro/sites/example-site",
					},
				},
			},
			want: []nitro.Action{
				{
					Type:       "mount",
					UseSyscall: false,
					Args:       []string{"mount", "./testdata/new-mount", "mytestmachine:/nitro/sites/new-site"},
				},
			},
			wantErr: false,
		},
		{
			name: "removed mounts return actions to remove mounts",
			args: args{
				machine: "mytestmachine",
				configFile: config.Config{
					Mounts: []config.Mount{
						{
							Source: "./testdata/new-mount",
							Dest:   "/nitro/sites/new-site",
						},
					},
				},
				fromMultipassMounts: []config.Mount{
					{
						Source: "./testdata/new-mount",
						Dest:   "/nitro/sites/new-site",
					},
					{
						Source: "./testdata/existing-mount",
						Dest:   "/nitro/sites/example-site",
					},
				},
				php: "",
			},
			want: []nitro.Action{
				{
					Type:       "umount",
					UseSyscall: false,
					Args:       []string{"umount", "mytestmachine:/nitro/sites/example-site"},
				},
			},
			wantErr: false,
		},
		{
			name: "renamed mounts get removed and added",
			args: args{
				machine: "mytestmachine",
				configFile: config.Config{
					Mounts: []config.Mount{
						{
							Source: "./testdata/new-mount",
							Dest:   "/nitro/sites/new-site",
						},
					},
				},
				fromMultipassMounts: []config.Mount{
					{
						Source: "./testdata/existing-mount",
						Dest:   "/nitro/sites/existing-site",
					},
				},
			},
			want: []nitro.Action{
				{
					Type:       "umount",
					UseSyscall: false,
					Args:       []string{"umount", "mytestmachine:/nitro/sites/existing-site"},
				},
				{
					Type:       "mount",
					UseSyscall: false,
					Args:       []string{"mount", "./testdata/new-mount", "mytestmachine:/nitro/sites/new-site"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Apply(tt.args.machine, tt.args.configFile, tt.args.fromMultipassMounts, tt.args.sites, tt.args.dbs, tt.args.php)
			if (err != nil) != tt.wantErr {
				t.Errorf("Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(tt.want) != len(got) {
				t.Errorf("expected the number of actions to be equal for Apply(); got %d, want %d", len(got), len(tt.want))
				return
			}

			if tt.want != nil {
				for i, action := range tt.want {
					if !reflect.DeepEqual(action, got[i]) {
						t.Errorf("Apply() got = \n%v, \nwant \n%v", got[i], tt.want[i])
					}
				}
			}
		})
	}
}
