package nitro

import (
	"reflect"
	"testing"
)

func TestSQL(t *testing.T) {
	type args struct {
		name    string
		engine  string
		version string
		root    bool
	}
	tests := []struct {
		name string
		args args
		want []Command
	}{
		{
			name: "can get postgres shell",
			args: args{
				name:    "machine-name",
				engine:  "postgres",
				version: "11.5",
				root:    false,
			},
			want: []Command{
				{
					Machine:   "machine-name",
					Type:      "exec",
					Chainable: false,
					Args:      []string{"machine-name", "--", "docker", "exec", "-it", "nitro_postgres_11.5", "psql", "-U", "nitro"},
				},
			},
		},
		{
			name: "can get mysql root shell",
			args: args{
				name:    "machine-name",
				engine:  "mysql",
				version: "5.7",
				root:    false,
			},
			want: []Command{
				{
					Machine:   "machine-name",
					Type:      "exec",
					Chainable: false,
					Args:      []string{"machine-name", "--", "docker", "exec", "-it", "nitro_mysql_5.7", "mysql", "-u", "nitro", "-pnitro"},
				},
			},
		},
		{
			name: "can get mysql root shell",
			args: args{
				name:    "machine-name",
				engine:  "mysql",
				version: "5.7",
				root:    true,
			},
			want: []Command{
				{
					Machine:   "machine-name",
					Type:      "exec",
					Chainable: false,
					Args:      []string{"machine-name", "--", "docker", "exec", "-it", "nitro_mysql_5.7", "mysql", "-u", "root", "-pnitro"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SQL(tt.args.name, tt.args.engine, tt.args.version, tt.args.root); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQL() = \n%v, \nwant \n%v", got, tt.want)
			}
		})
	}
}
