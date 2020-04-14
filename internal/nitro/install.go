package nitro

import (
	"strings"

	"github.com/craftcms/nitro/validate"
)

const (
	php70Packages = "php7.0 php7.0-mbstring php7.0-cli php7.0-curl php7.0-fpm php7.0-gd php7.0-intl php7.0-json php7.0-mysql php7.0-opcache php7.0-pgsql php7.0-zip php7.0-xml php-xdebug php-imagick blackfire-agent blackfire-php"
	php71Packages = "php7.1 php7.1-mbstring php7.1-cli php7.1-curl php7.1-fpm php7.1-gd php7.1-intl php7.1-json php7.1-mysql php7.1-opcache php7.1-pgsql php7.1-zip php7.1-xml php-xdebug php-imagick blackfire-agent blackfire-php"
	php72Packages = "php7.2 php7.2-mbstring php7.2-cli php7.2-curl php7.2-fpm php7.2-gd php7.2-intl php7.2-json php7.2-mysql php7.2-opcache php7.2-pgsql php7.2-zip php7.2-xml php-xdebug php-imagick blackfire-agent blackfire-php"
	php73Packages = "php7.3 php7.3-mbstring php7.3-cli php7.3-curl php7.3-fpm php7.3-gd php7.3-intl php7.3-json php7.3-mysql php7.3-opcache php7.3-pgsql php7.3-zip php7.3-xml php-xdebug php-imagick blackfire-agent blackfire-php"
	php74Packages = "php7.4 php7.4-mbstring php7.4-cli php7.4-curl php7.4-fpm php7.4-gd php7.4-intl php7.4-json php7.4-mysql php7.4-opcache php7.4-pgsql php7.4-zip php7.4-xml php-xdebug php-imagick blackfire-agent blackfire-php"
)

// InstallPackages is used to install the core PHP packages needed by the
// nitro machine to run.
func InstallPackages(name, php string) (*Action, error) {
	if err := validate.MachineName(name); err != nil {
		return nil, err
	}
	if err := validate.PHPVersion(php); err != nil {
		return nil, err
	}

	args := []string{"exec", name, "--", "sudo", "apt-get", "install", "-y"}

	switch php {
	case "7.0":
		args = append(args, strings.Split(php70Packages, " ")...)
	case "7.1":
		args = append(args, strings.Split(php71Packages, " ")...)
	case "7.2":
		args = append(args, strings.Split(php72Packages, " ")...)
	case "7.3":
		args = append(args, strings.Split(php73Packages, " ")...)
	default:
		args = append(args, strings.Split(php74Packages, " ")...)
	}

	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       args,
	}, nil
}
