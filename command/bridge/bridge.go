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

	"github.com/craftcms/nitro/pkg/appaware"
	"github.com/craftcms/nitro/pkg/flags"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
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
		Short:   "Shares apps on your local network.",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// get all interfaces
			ifaces, err := net.Interfaces()
			if err != nil {
				return err
			}

			// create the options for the prompt
			var interfaces []string

			// get the ip addresses
			for _, i := range ifaces {
				addrs, err := i.Addrs()
				if err != nil {
					return err
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
							interfaces = append(interfaces, ip.String())
						}
					}
				}
			}

			// prompt for which interface to use
			selected, err := output.Select(cmd.InOrStdin(), "Which IP address should we use for the bridge? ", interfaces)
			if err != nil {
				return err
			}

			// get the ip for the bridge
			ip := interfaces[selected]

			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			var appName string
			switch flags.AppName != "" {
			case false:
				// get the current working directory
				wd, err := os.Getwd()
				if err != nil {
					return err
				}

				// get a context aware list of sites
				appName, err = appaware.Detect(*cfg, wd)
				if err != nil {
					return err
				}
			default:
				appName = flags.AppName
			}

			// find the app by hostname
			app, err := cfg.FindAppByHostname(appName)
			if err != nil {
				return err
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

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

			// set the port
			port, err := cmd.Flags().GetString("port")
			if err != nil {
				port = "8000"
			}

			logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

			target, err := url.Parse(fmt.Sprintf("http://%s", app.Hostname))
			if err != nil {
				return err
			}

			output.Info("connecting to",app.Hostname)

			// create the handle
			http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
				proxy := httputil.NewSingleHostReverseProxy(target)

				original := r.Host

				r.Host = target.Host
				r.URL.Host = target.Host
				r.URL.Scheme = target.Scheme
				r.Header.Set("X-Forwarded-Host", original)

				logger.Println(r.Header.Get("Host"), r.RequestURI)

				proxy.ServeHTTP(rw, r)
			})

			output.Info(fmt.Sprintf("Bridge server listening on http://%s:%s", ip, port))

			return http.ListenAndServe(ip+":"+port, nil)
		},
	}

	cmd.Flags().String("port", "8000", "which port to use for the bridge")

	return cmd
}

func isIpv4(ipAddress string) bool {
	return ipRegex.MatchString(strings.Trim(ipAddress, " "))
}
