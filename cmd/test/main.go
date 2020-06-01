package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/scripts"
)

func main() {
	machine := "nitro-dev"
	multipass, err := exec.LookPath("multipass")
	if err != nil {
		fmt.Println(err)
		return
	}

	script := scripts.New(multipass, machine)

	var confs []string
	if output, err := script.Run(false, `ls /etc/nginx/sites-enabled`); err == nil {
		sc := bufio.NewScanner(strings.NewReader(output))
		for sc.Scan() {
			if sc.Text() == "default" {
				continue
			}

			confs = append(confs, strings.TrimSpace(sc.Text()))
		}
	}

	var sites []config.Site
	for _, conf := range confs {
		s := config.Site{}
		// get the webroot
		if output, err := script.Run(false, fmt.Sprintf(scripts.FmtNginxSiteWebroot, conf)); err == nil {
			s.Webroot = output
		}

		// get the server_name
		if output, err := script.Run(false, fmt.Sprintf(`grep "server_name " /etc/nginx/sites-available/%s | while read -r line; do echo "$line"; done`, conf)); err == nil {
			sp := strings.Fields(output)
			if len(sp) >= 2 {
				s.Hostname = sp[1]
			}
		}

		// get the hostname
		if s.Webroot != "" && s.Hostname != "" {
			sites = append(sites, s)
		}
	}

	fmt.Println(sites)
}
