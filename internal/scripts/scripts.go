package scripts

func InstallMariaDB() []string {
	return []string{"sudo", "apt", "install", "-y", "mariadb-server"}
}

func InstallPHP(version string) []string {
	return []string{"sudo", "apt", "install", "-y", "php7.4", "php7.4-mbstring", "php7.4-cli", "php7.4-curl", "php7.4-fpm", "php7.4-gd", "php7.4-intl", "php7.4-json", "php7.4-mysql", "php7.4-opcache", "php7.4-pgsql", "php7.4-zip", "php7.4-xml"}
}
