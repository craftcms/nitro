package task

import (
	"reflect"
	"testing"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

func TestRemoveSite(t *testing.T) {
	type args struct {
		machine    string
		php        string
		oldSite    config.Site
		newSite    config.Site
		configFile config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    []nitro.Action
		wantErr bool
	}{
		{
			name: "if there is an exact mount it unmounts it and adds a new mount",
			args: args{
				machine: "mymachine",
				php:     "7.4",
				oldSite: config.Site{
					Hostname: "oldhostname.test",
					Webroot:  "/nitro/sites/oldhostname.test",
				},
				newSite: config.Site{
					Hostname: "newhostname.test",
					Webroot:  "/nitro/sites/newhostname.test",
				},
				configFile: config.Config{
					Mounts: []config.Mount{
						{
							Source: "./testdata/example-source",
							Dest:   "/nitro/sites/oldhostname.test",
						},
					},
				},
			},
			want: []nitro.Action{
				{
					Type:       "umount",
					UseSyscall: false,
					Args:       []string{"umount", "mymachine" + ":/app/sites/oldhostname.test"},
				},
				{
					Type:       "umount",
					UseSyscall: false,
					Args:       []string{"umount", "mymachine" + ":/app/sites/oldhostname.test"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mymachine", "--", "sudo", "rm", "/etc/nginx/sites-enabled/oldhostname.test"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mymachine", "--", "sudo", "cp", "/opt/nitro/nginx/template.conf", "/etc/nginx/sites-available/newhostname.test"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mymachine", "--", "sudo", "sed", "-i", "s|CHANGEWEBROOTDIR|/nitro/sites/newhostname.test|g", "/etc/nginx/sites-available/newhostname.test"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mymachine", "--", "sudo", "sed", "-i", "s|CHANGESERVERNAME|newhostname.test|g", "/etc/nginx/sites-available/newhostname.test"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mymachine", "--", "sudo", "sed", "-i", "s|CHANGEPHPVERSION|7.4|g", "/etc/nginx/sites-available/newhostname.test"},
				},
				{
					Type:       "exec",
					UseSyscall: false,
					Args:       []string{"exec", "mymachine", "--", "sudo", "service", "nginx", "restart"},
				},
			},
			wantErr: false,
		},
		//{
		//	name: "if there is a relative mount it does not unmount",
		//	args: args{
		//		machine: "mymachine",
		//		php:     "7.4",
		//		oldSite: config.Site{
		//			Hostname: "oldhostname.test",
		//			Webroot:  "/nitro/sites/oldhostname.test",
		//		},
		//		newSite: config.Site{
		//			Hostname: "newhostname.test",
		//			Webroot:  "/nitro/sites/newhostname.test",
		//		},
		//		configFile: config.Config{
		//			Mounts: []config.Mount{
		//				{
		//					Source: "./testdata/example-source",
		//					Dest:   "/nitro/sites/",
		//				},
		//			},
		//		},
		//	},
		//	want: []nitro.Action{
		//		{
		//			Type:       "exec",
		//			UseSyscall: false,
		//			Args:       []string{"exec", "mymachine", "--", "sudo", "rm", "/etc/nginx/sites-enabled/oldhostname.test"},
		//		},
		//		{
		//			Type:       "exec",
		//			UseSyscall: false,
		//			Args:       []string{"exec", "mymachine", "--", "sudo", "cp", "/opt/nitro/nginx/template.conf", "/etc/nginx/sites-available/newhostname.test"},
		//		},
		//		{
		//			Type:       "exec",
		//			UseSyscall: false,
		//			Args:       []string{"exec", "mymachine", "--", "sudo", "sed", "-i", "s|CHANGEWEBROOTDIR|/nitro/sites/newhostname.test|g", "/etc/nginx/sites-available/newhostname.test"},
		//		},
		//		{
		//			Type:       "exec",
		//			UseSyscall: false,
		//			Args:       []string{"exec", "mymachine", "--", "sudo", "sed", "-i", "s|CHANGESERVERNAME|newhostname.test|g", "/etc/nginx/sites-available/newhostname.test"},
		//		},
		//		{
		//			Type:       "exec",
		//			UseSyscall: false,
		//			Args:       []string{"exec", "mymachine", "--", "sudo", "sed", "-i", "s|CHANGEPHPVERSION|7.4|g", "/etc/nginx/sites-available/newhostname.test"},
		//		},
		//		{
		//			Type:       "exec",
		//			UseSyscall: false,
		//			Args:       []string{"exec", "mymachine", "--", "sudo", "service", "nginx", "restart"},
		//		},
		//	},
		//	wantErr: false,
		//},
		//{
		//	name: "returns an action for mounts that do not have a direct mount",
		//	args: args{
		//		machine: "mymachine",
		//		php:     "7.4",
		//		oldSite: config.Site{
		//			Hostname: "oldhostname.test",
		//			Webroot:  "/nitro/sites/oldhostname.test",
		//		},
		//		newSite: config.Site{
		//			Hostname: "newhostname.test",
		//			Webroot:  "/nitro/sites/newhostname.test",
		//		},
		//		configFile: config.Config{
		//			Mounts: []config.Mount{
		//				{
		//					Source: "./testdata/example-source",
		//					Dest:   "/nitro/sites/",
		//				},
		//			},
		//		},
		//	},
		//	want: []nitro.Action{
		//		{
		//			Type:       "exec",
		//			UseSyscall: false,
		//			Args:       []string{"exec", "mymachine", "--", "sudo", "rm", "/etc/nginx/sites-enabled/oldhostname.test"},
		//		},
		//		{
		//			Type:       "exec",
		//			UseSyscall: false,
		//			Args:       []string{"exec", "mymachine", "--", "sudo", "cp", "/opt/nitro/nginx/template.conf", "/etc/nginx/sites-available/newhostname.test"},
		//		},
		//		{
		//			Type:       "exec",
		//			UseSyscall: false,
		//			Args:       []string{"exec", "mymachine", "--", "sudo", "sed", "-i", "s|CHANGEWEBROOTDIR|/nitro/sites/newhostname.test|g", "/etc/nginx/sites-available/newhostname.test"},
		//		},
		//		{
		//			Type:       "exec",
		//			UseSyscall: false,
		//			Args:       []string{"exec", "mymachine", "--", "sudo", "sed", "-i", "s|CHANGESERVERNAME|newhostname.test|g", "/etc/nginx/sites-available/newhostname.test"},
		//		},
		//		{
		//			Type:       "exec",
		//			UseSyscall: false,
		//			Args:       []string{"exec", "mymachine", "--", "sudo", "sed", "-i", "s|CHANGEPHPVERSION|7.4|g", "/etc/nginx/sites-available/newhostname.test"},
		//		},
		//		{
		//			Type:       "exec",
		//			UseSyscall: false,
		//			Args:       []string{"exec", "mymachine", "--", "sudo", "service", "nginx", "restart"},
		//		},
		//	},
		//	wantErr: false,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RenameSite(tt.args.machine, tt.args.php, tt.args.oldSite, tt.args.newSite, tt.args.configFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenameSite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("did not get the number of actions we wanted, got %d, wanted: %d", len(got), len(tt.want))
				return
			}

			for i, action := range tt.want {
				if !reflect.DeepEqual(action, got[i]) {
					t.Errorf("RenameSite() got = \n%v, \nwant \n%v", action, tt.want[i])
				}
			}
		})
	}
}
