package action

import (
	"reflect"
	"testing"
)

func TestInstallCorePackages(t *testing.T) {
	type args struct {
		name string
		php  string
	}
	tests := []struct {
		name    string
		args    args
		want    *Action
		wantErr bool
	}{
		{
			name: "can get commands to install PHP 7.4",
			args: args{
				name: "somename",
				php:  "7.4",
			},
			want: &Action{
				Type:       "exec",
				UseSyscall: false,
				Args:       []string{"exec", "somename", "--", "sudo", "apt-get", "install", "-y", "php7.4", "php7.4-mbstring", "php7.4-cli", "php7.4-curl", "php7.4-fpm", "php7.4-gd", "php7.4-intl", "php7.4-json", "php7.4-mysql", "php7.4-opcache", "php7.4-pgsql", "php7.4-zip", "php7.4-xml", "php-xdebug", "php-imagick", "blackfire-agent", "blackfire-php"},
			},
			wantErr: false,
		},
		{
			name: "wrong version of php returns an error",
			args: args{
				name: "somename",
				php:  "7.9",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InstallPackages(tt.args.name, tt.args.php)
			if (err != nil) != tt.wantErr {
				t.Errorf("InstallPackages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InstallPackages() got = %v, want %v", got, tt.want)
			}
		})
	}
}
