package bridge

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	ipRegex, _ = regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
)

const exampleText = `  # bridge your network ip to share nitro on the local network
  nitro bridge`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bridge",
		Short:   "Share sites on your local network",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// get all of the interfaces
			ifaces, err := net.Interfaces()
			if err != nil {
				log.Fatal(err)
			}

			// create the options for the prompt
			var ifaceOptions []string

			// get the ip addresses
			for _, i := range ifaces {
				addrs, err := i.Addrs()
				if err != nil {
					log.Fatal(err)
				}

				for _, addr := range addrs {
					var ip net.IP
					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP
					case *net.IPAddr:
						ip = v.IP
					}

					// ignore ipv6
					if isIpv4(ip.String()) {
						if !strings.Contains(ip.String(), "127.0.0.1") {
							ifaceOptions = append(ifaceOptions, ip.String())
						}
					}
				}
			}

			// prompt for which interface to use
			selected, err := output.Select(cmd.InOrStdin(), "Which IP address should we use for the bridge? ", ifaceOptions)
			if err != nil {
				return err
			}

			// get the ip for the bridge
			ip := ifaceOptions[selected]

			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", labels.Nitro)

			// get a context aware list of sites
			sites := cfg.ListOfSitesByDirectory(home, wd)

			// create the options for the sites
			var options []string
			for _, s := range sites {
				options = append(options, s.Hostname)
			}

			// if there are found sites we want to show or connect to the first one, otherwise prompt for
			// which site to connect to.
			var site config.Site
			switch len(sites) {
			case 0:
				// prompt for the site to ssh into
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
				if err != nil {
					return err
				}

				// add the label to get the site
				filter.Add("label", labels.Host+"="+sites[selected].Hostname)

				site = sites[selected]
			case 1:
				output.Info("connecting to", sites[0].Hostname)

				// add the label to get the site
				filter.Add("label", labels.Host+"="+sites[0].Hostname)
				site = sites[0]
			default:
				// prompt for the site to ssh into
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
				if err != nil {
					return err
				}

				// add the label to get the site
				filter.Add("label", labels.Host+"="+sites[selected].Hostname)
				site = sites[selected]
			}

			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter, All: true})
			if err != nil {
				return err
			}

			if len(containers) == 0 {
				return fmt.Errorf("no containers found")
			}

			// start the containers
			for _, c := range containers {
				// start the containers if not running
				if c.State != "running" {
					for _, command := range cmd.Root().Commands() {
						if command.Use == "start" {
							if err := command.RunE(cmd, []string{}); err != nil {
								return err
							}
						}
					}
				}
			}

			target, err := url.Parse(fmt.Sprintf("http://%s", site.Hostname))
			if err != nil {
				return err
			}

			// set the port
			port, err := cmd.Flags().GetString("port")
			if err != nil {
				port = "8000"
			}

			logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

			// create the handle
			http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
				proxy := httputil.NewSingleHostReverseProxy(target)

				r.URL.Host = target.Host
				r.URL.Scheme = target.Scheme
				r.Header.Set("X-Forwarded-Host", target.Host)
				r.Host = target.Host

				logger.Println(r.RequestURI)

				proxy.ServeHTTP(rw, r)
			})

			output.Info(fmt.Sprintf("bridge server listening on http://%s:%s", ip, port))

			log.Fatal(http.ListenAndServe(ip+":"+port, nil))

			return nil
		},
	}

	cmd.Flags().String("port", "8000", "which port to use for the bridge")

	return cmd
}

func isIpv4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")
	return ipRegex.MatchString(ipAddress)
}
