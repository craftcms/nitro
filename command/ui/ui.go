package ui

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # ui command
  nitro ui

  # start on a specific port
  nitro ui --port 8000`

func NewCommand(home string, ocker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ui",
		Short:   "Start the Nitro UI.",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// get the current working directory
			port := cmd.Flag("port").Value.String()

			http.HandleFunc("/v2/config", func(w http.ResponseWriter, r *http.Request) {
				cfg, err := config.Load(home)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				c, err := json.Marshal(cfg)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
				w.Header().Add("Content-Type", "application/json")
				w.Write(c)
			})

			http.HandleFunc("/v2/sites", func(w http.ResponseWriter, r *http.Request) {
				cfg, err := config.Load(home)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				c, err := json.Marshal(cfg.Sites)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
				w.Header().Add("Content-Type", "application/json")
				w.Write(c)
			})

			http.HandleFunc("/v2/databases", func(w http.ResponseWriter, r *http.Request) {
				cfg, err := config.Load(home)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				name := r.URL.Query().Get("name")
				if name != "" {
					for _, d := range cfg.Databases {
						if host, _ := d.GetHostname(); name == host {
							c, err := json.Marshal(d)
							if err != nil {
								http.Error(w, err.Error(), http.StatusInternalServerError)
								return
							}

							w.WriteHeader(http.StatusOK)
							w.Header().Add("Content-Type", "application/json")
							w.Write(c)
							return
						}
						http.Error(w, "not found", http.StatusNotFound)
						return
					}
				}

				c, err := json.Marshal(cfg.Databases)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
				w.Header().Add("Content-Type", "application/json")
				w.Write(c)
			})

			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

				w.WriteHeader(http.StatusOK)
				w.Header().Add("Content-Type", "application/json")
				fmt.Fprint(w, "UI goes here!")
			})

			fmt.Printf("Starting server at port %s\n", port)

			return http.ListenAndServe(":"+port, nil)
		},
	}

	// set flags for the command
	cmd.Flags().String("port", "8000", "port to expose the API for the UI")

	return cmd
}
