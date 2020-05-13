package scripts

import (
	"errors"
	"fmt"

	"github.com/craftcms/nitro/config"
)

var (
	FmtSiteAvailable = `if test -f '/etc/nginx/sites-available/%s'; then echo 'exists'; fi`
)

// SiteIsAvailable takes a site and returns the commands used
// to see if the site if available in NGINX.
func SiteIsAvailable(s config.Site) (string, error) {
	if s.Hostname == "" {
		return "", errors.New("site hostname is empty")
	}

	return fmt.Sprintf(`if test -f '/etc/nginx/sites-available/%s'; then echo 'exists'; fi`, s.Hostname), nil
}

// SiteIsEnabled takes a site and returns the commands used
// to see if the site if enabled in NGINX.
func SiteIsEnabled(s config.Site) (string, error) {
	if s.Hostname == "" {
		return "", errors.New("site hostname is empty")
	}

	return fmt.Sprintf(`if test -f '/etc/nginx/sites-enabled/%s'; then echo 'exists'; fi`, s.Hostname), nil
}

// SiteWebroot is used to return the root of the NGINX
// site.
func SiteWebroot(s config.Site) (string, error) {
	if s.Hostname == "" {
		return "", errors.New("site hostname is empty")
	}

	return fmt.Sprintf(`grep "root " %s | while read -r line; do echo "$line"; done`, "/etc/nginx/sites-available/"+s.Hostname), nil
}
