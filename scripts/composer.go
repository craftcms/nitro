package scripts

// InstallComposer returns a script that will install the latest composer. It might be a good idea to provide a specific version?
// to prevent scenarios like this: https://github.com/composer/composer/issues/8710
func InstallComposer() string {
	return `php -r "readfile('http://getcomposer.org/installer')" | sudo php -- --install-dir=/usr/bin/ --filename=composer`
}
