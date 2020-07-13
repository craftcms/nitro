package api

type Response struct {
	Success bool     `json:"success"`
	Output  string   `json:"output"`
	Errors  []string `json:"errors,omitempty"`
}
