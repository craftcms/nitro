package nitrod

import (
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Output  string `json:"output"`
}

type ErrorResponse struct {
	Errors struct {
		Version []string `json:"version"`
		Action  []string `json:"action"`
	} `json:"errors"`
}

func NewClient(ip string) *Client {
	return &Client{
		BaseURL: fmt.Sprintf("http://%s:9999/v1", ip),
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}
