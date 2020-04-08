package scripts

import (
	"reflect"
	"testing"
)

func TestInstallPHP(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    []string
	}{
		{
			name:    "returns PHP 7.0 packages",
			version: "7.0",
			want: []string{
				"php7.0",
				"php7.0-mbstring",
				"php7.0-cli",
				"php7.0-curl",
				"php7.0-fpm",
				"php7.0-gd",
				"php7.0-intl",
				"php7.0-json",
				"php7.0-mysql",
				"php7.0-opcache",
				"php7.0-pgsql",
				"php7.0-zip",
				"php7.0-xml",
				"php-xdebug",
				"php-imagick",
				"blackfire-agent",
				"blackfire-php",
			},
		},
		{
			name:    "returns PHP 7.1 packages",
			version: "7.1",
			want: []string{
				"php7.1",
				"php7.1-mbstring",
				"php7.1-cli",
				"php7.1-curl",
				"php7.1-fpm",
				"php7.1-gd",
				"php7.1-intl",
				"php7.1-json",
				"php7.1-mysql",
				"php7.1-opcache",
				"php7.1-pgsql",
				"php7.1-zip",
				"php7.1-xml",
				"php-xdebug",
				"php-imagick",
				"blackfire-agent",
				"blackfire-php",
			},
		},
		{
			name:    "returns PHP 7.2 packages",
			version: "7.2",
			want: []string{
				"php7.2",
				"php7.2-mbstring",
				"php7.2-cli",
				"php7.2-curl",
				"php7.2-fpm",
				"php7.2-gd",
				"php7.2-intl",
				"php7.2-json",
				"php7.2-mysql",
				"php7.2-opcache",
				"php7.2-pgsql",
				"php7.2-zip",
				"php7.2-xml",
				"php-xdebug",
				"php-imagick",
				"blackfire-agent",
				"blackfire-php",
			},
		},
		{
			name:    "returns PHP 7.3 packages",
			version: "7.3",
			want: []string{
				"php7.3",
				"php7.3-mbstring",
				"php7.3-cli",
				"php7.3-curl",
				"php7.3-fpm",
				"php7.3-gd",
				"php7.3-intl",
				"php7.3-json",
				"php7.3-mysql",
				"php7.3-opcache",
				"php7.3-pgsql",
				"php7.3-zip",
				"php7.3-xml",
				"php-xdebug",
				"php-imagick",
				"blackfire-agent",
				"blackfire-php",
			},
		},
		{
			name:    "returns PHP 7.4 by default",
			version: "",
			want: []string{
				"php7.4",
				"php7.4-mbstring",
				"php7.4-cli",
				"php7.4-curl",
				"php7.4-fpm",
				"php7.4-gd",
				"php7.4-intl",
				"php7.4-json",
				"php7.4-mysql",
				"php7.4-opcache",
				"php7.4-pgsql",
				"php7.4-zip",
				"php7.4-xml",
				"php-xdebug",
				"php-imagick",
				"blackfire-agent",
				"blackfire-php",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InstallPHP(tt.version); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InstallPHP() = %v, want %v", got, tt.want)
			}
		})
	}
}
