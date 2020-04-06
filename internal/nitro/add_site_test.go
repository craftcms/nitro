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
