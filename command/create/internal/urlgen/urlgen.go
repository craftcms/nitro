package urlgen

import (
	"fmt"
	"net/url"
	"strings"
)

const download = "https://github.com/craftcms/craft/archive/HEAD.zip"

// Generate is a helper that is used to build the Github download link
// of a repository (e.g. https://github.com/craftcms/craft/archive/HEAD.zip).
// It supports short hand urls such as `craftcms/craft`. If no address is
// provided or the case does not match, it defaults to the craft repo.
func Generate(addr string) (*url.URL, error) {
	// split the address by /
	sp := strings.Split(addr, "/")

	// check the length
	switch len(sp) {
	case 5:
		// parse a conplete url for github and append the download
		u, err := url.Parse(fmt.Sprintf("%s/archive/HEAD.zip", addr))
		if err != nil {
			return nil, err
		}

		return u, nil
	// if this is using the short hand address
	case 2:
		u, err := url.Parse(fmt.Sprintf("https://github.com/%s/%s/archive/HEAD.zip", sp[0], sp[1]))
		if err != nil {
			return nil, err
		}

		return u, nil
	}

	// setup the default download url
	return url.Parse(download)
}
