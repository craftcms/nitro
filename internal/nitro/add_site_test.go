package nitro

import (
	"reflect"
	"testing"
)

func Test_changePathInNginxPath(t *testing.T) {
	type args struct {
		name     string
		site     string
		variable string
		actual   string
	}
	tests := []struct {
		name string
		args args
		want Command
	}{
		{
			name: "returns the correct sed site",
			args: args{
				name:     "machine-name",
				site:     "example.test",
				variable: "CHANGEPATH",
				actual:   "example.test",
			},
			want: Command{
				Machine:   "machine-name",
				Type:      "exec",
				Chainable: true,
				Args:      []string{"machine-name", "--", "sudo", "sed", "-i", "s|CHANGEPATH|example.test|g", "/etc/nginx/sites-available/example.test"},
			},
		},
		{
			name: "returns the correct sed php version",
			args: args{
				name:     "machine-name",
				site:     "another-site.test",
				variable: "CHANGEPHPVERSION",
				actual:   "7.4",
			},
			want: Command{
				Machine:   "machine-name",
				Type:      "exec",
				Chainable: true,
				Args:      []string{"machine-name", "--", "sudo", "sed", "-i", "s|CHANGEPHPVERSION|7.4|g", "/etc/nginx/sites-available/another-site.test"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := changeVariableInNginxTemplate(tt.args.name, tt.args.site, tt.args.variable, tt.args.actual); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("changeVariableInNginxTemplate() = \n%v, \nwant:\n %v", got, tt.want)
			}
		})
	}
}

func TestCreateNewDirectoryForSite(t *testing.T) {
	type args struct {
		name string
		site string
	}
	tests := []struct {
		name string
		args args
		want Command
	}{
		{
			name: "creates directory",
			args: args{
				name: "machine-name",
				site: "example.tes",
			},
			want: Command{
				Machine:   "machine-name",
				Type:      "exec",
				Chainable: true,
				Args:      []string{"machine-name", "--", "mkdir", "-p", "/app/sites/example.tes"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateNewDirectoryForSite(tt.args.name, tt.args.site); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateNewDirectoryForSite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopyNginxTemplate(t *testing.T) {
	type args struct {
		name string
		site string
	}
	tests := []struct {
		name string
		args args
		want Command
	}{
		{
			name: "copies template",
			args: args{
				name: "machine-name",
				site: "example.tes",
			},
			want: Command{
				Machine:   "machine-name",
				Type:      "exec",
				Chainable: true,
				Args:      []string{"machine-name", "--", "sudo", "cp", "/opt/nitro/nginx/template.conf", "/etc/nginx/sites-available/example.tes"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CopyNginxTemplate(tt.args.name, tt.args.site); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CopyNginxTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_changeVariableInNginxTemplate(t *testing.T) {
	type args struct {
		name     string
		site     string
		variable string
		actual   string
	}
	tests := []struct {
		name string
		args args
		want Command
	}{
		{
			name: "can replace variables in template",
			args: args{
				name:     "machine-name",
				site:     "example.test",
				variable: "CHANGEME",
				actual:   "changed",
			},
			want: Command{
				Machine:   "machine-name",
				Type:      "exec",
				Chainable: true,
				Args:      []string{"machine-name", "--", "sudo", "sed", "-i", "s|CHANGEME|changed|g", "/etc/nginx/sites-available/example.test"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := changeVariableInNginxTemplate(tt.args.name, tt.args.site, tt.args.variable, tt.args.actual); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("changeVariableInNginxTemplate() = \n%v, \nwant \n%v", got, tt.want)
			}
		})
	}
}

func TestLinkNginxSite(t *testing.T) {
	type args struct {
		name string
		site string
	}
	tests := []struct {
		name string
		args args
		want Command
	}{
		{
			name: "can create symlink for the nginx site",
			args: args{
				name: "machine-test",
				site: "example.test",
			},
			want: Command{
				Machine:   "machine-test",
				Type:      "exec",
				Chainable: true,
				Args:      []string{"machine-test", "--", "sudo", "ln", "-s", "/etc/nginx/sites-available/example.test", "/etc/nginx/sites-enabled/"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LinkNginxSite(tt.args.name, tt.args.site); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LinkNginxSite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReloadNginx(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want Command
	}{
		{
			name: "can restart nginx",
			args: args{
				name: "machine-name",
			},
			want: Command{
				Machine:   "machine-name",
				Type:      "exec",
				Chainable: true,
				Args:      []string{"machine-name", "--", "sudo", "service", "nginx", "restart"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReloadNginx(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReloadNginx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChangeVariablesInTemplate(t *testing.T) {
	type args struct {
		name   string
		domain string
		dir    string
		php    string
	}
	tests := []struct {
		name string
		args args
		want []Command
	}{
		{
			name: "can get all of the commands to change all the variables in the template",
			args: args{
				name:   "machine-name",
				domain: "example.test",
				dir:    "somedir",
				php:    "7.4",
			},
			want: []Command{
				{
					Machine:   "machine-name",
					Type:      "exec",
					Chainable: true,
					Args:      []string{"machine-name", "--", "sudo", "sed", "-i", "s|CHANGEPATH|example.test|g", "/etc/nginx/sites-available/example.test"},
				},
				{
					Machine:   "machine-name",
					Type:      "exec",
					Chainable: true,
					Args:      []string{"machine-name", "--", "sudo", "sed", "-i", "s|CHANGESERVERNAME|example.test|g", "/etc/nginx/sites-available/example.test"},
				},
				{
					Machine:   "machine-name",
					Type:      "exec",
					Chainable: true,
					Args:      []string{"machine-name", "--", "sudo", "sed", "-i", "s|CHANGEPUBLICDIR|somedir|g", "/etc/nginx/sites-available/example.test"},
				},
				{
					Machine:   "machine-name",
					Type:      "exec",
					Chainable: true,
					Args:      []string{"machine-name", "--", "sudo", "sed", "-i", "s|CHANGEPHPVERSION|7.4|g", "/etc/nginx/sites-available/example.test"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ChangeVariablesInTemplate(tt.args.name, tt.args.domain, tt.args.dir, tt.args.php)

			for i, command := range got {
				if !reflect.DeepEqual(tt.want[i], command) {
					t.Errorf("ChangeVariablesInTemplate() = \n%v, \nwant \n%v", got, tt.want)
				}
			}
		})
	}
}
