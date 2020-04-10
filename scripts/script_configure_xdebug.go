package scripts

import (
	"fmt"

	"github.com/craftcms/nitro/validate"
)

func ConfigureXdebug(php string) (Script, error) {
	if err := validate.PHPVersion(php); err != nil {
		return Script{}, err
	}

	cmd := `"\$axdebug.remote_enable=1\nxdebug.remote_connect_back=0\nxdebug.remote_host=CHANGEMEIP\nxdebug.remote_port=9000\nxdebug.remote_log=/var/log/nitro/xdebug.log"`
	path := fmt.Sprintf("/etc/php/%v/mods-available/xdebug.ini", php)

	return Script{
		Name: "configure xdebug remote connection for PHP 7.4",
		Args: []string{"sudo", "sed", "-i", "-e", cmd, path},
	}, nil
}
