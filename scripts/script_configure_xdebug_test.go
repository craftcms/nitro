package scripts

import (
	"reflect"
	"testing"
)

func TestConfigureXdebug(t *testing.T) {
	type args struct {
		php string
	}
	tests := []struct {
		name string
		args args
		want Script
		err  bool
	}{
		{
			name: "PHP 7.4 xdebug config",
			args: args{php: "7.4"},
			want: Script{
				Name: "configure xdebug remote connection for PHP 7.4",
				Args: []string{"sudo", "sed", "-i", "-e", `"\$axdebug.remote_enable=1\nxdebug.remote_connect_back=0\nxdebug.remote_host=CHANGEMEIP\nxdebug.remote_port=9000\nxdebug.remote_log=/var/log/nitro/xdebug.log"`, "/etc/php/7.4/mods-available/xdebug.ini"},
			},
			err: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConfigureXdebug(tt.args.php)
			if tt.err {
				if err != nil {
					t.Error("wanted an error, got nil instead")
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConfigureXdebug() = %v, \nwant\n %v", got, tt.want)
			}
		})
	}
}
