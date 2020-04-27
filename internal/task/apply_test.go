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
