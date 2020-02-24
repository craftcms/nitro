package install

import "errors"

// PHP takes a version and returns the commands that
// are used to install that version of PHP. If an
// unknown version is passed, it will return an error
func PHP(v string) (string, error) {
	switch v {
	case "7.4":
		return "php7.4 php7.4-mbstring php7.4-cli php7.4-curl php7.4-fpm php7.4-gd php7.4-intl php7.4-json php7.4-mysql php7.4-opcache php7.4-pgsql php7.4-zip php7.4-xml", nil
	case "7.3":
		return "php7.3 php7.3-mbstring php7.3-cli php7.3-curl php7.3-fpm php7.3-gd php7.3-intl php7.3-json php7.3-mysql php7.3-opcache php7.3-pgsql php7.3-zip php7.3-xml", nil
	}

	return "", errors.New("unsupported version of PHP provided")
}
