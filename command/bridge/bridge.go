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
		Short:   "Share sites on your local network",
		Example: exampleText,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			cfg, err := config.Load(home)
			if err != nil {
				return nil, cobra.ShellCompDirectiveDefault
			}

			var options []string
			for _, s := range cfg.Sites {
				options = append(options, s.Hostname)
			}

			return options, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// get all of the interfaces
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

			var site string
			if len(args) > 0 {
				site = strings.TrimSpace(args[0])
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

			// get a context aware list of sites
			sites := cfg.ListOfSitesByDirectory(home, wd)

			// create the options for the sites
			var options []string
			for _, s := range sites {
				options = append(options, s.Hostname)
			}

			var target *url.URL
			switch site == "" {
			case false:
				for k, v := range options {
					if site == v {
						target, err = url.Parse(fmt.Sprintf("http://%s", sites[k].Hostname))
						if err != nil {
							return err
						}

						break
					}
				}
			default:
				switch len(sites) {
				case 0:
					// prompt for the site to ssh into
					selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
					if err != nil {
						return err
					}

					// add the label to get the site
					filter.Add("label", containerlabels.Host+"="+sites[selected].Hostname)

					target, err = url.Parse(fmt.Sprintf("http://%s", sites[selected].Hostname))
					if err != nil {
						return err
					}
				case 1:
					output.Info("connecting to", sites[0].Hostname)

					// add the label to get the site
					filter.Add("label", containerlabels.Host+"="+sites[0].Hostname)

					target, err = url.Parse(fmt.Sprintf("http://%s", sites[0].Hostname))
					if err != nil {
						return err
					}
				default:
					// prompt for the site to ssh into
					selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
					if err != nil {
						return err
					}

					// add the label to get the site
					filter.Add("label", containerlabels.Host+"="+sites[selected].Hostname)

					target, err = url.Parse(fmt.Sprintf("http://%s", sites[selected].Hostname))
					if err != nil {
						return err
					}
				}
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

			// set the port
			port, err := cmd.Flags().GetString("port")
			if err != nil {
				port = "8000"
			}

			logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

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
