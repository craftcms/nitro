package nitrod

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type NginxServiceOptions struct {
	Action string
}

func (c *Client) ServiceNginx(ctx context.Context, options *NginxServiceOptions) (*SuccessResponse, error) {
	action := "restart"
	if options != nil {
		action = options.Action
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/services/nginx", c.BaseURL),
		strings.NewReader(fmt.Sprintf(`{"action":"%s"}`, action)),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// check for validation error
	if resp.StatusCode == http.StatusUnprocessableEntity {
		errResp := ErrorResponse{}

		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, err
		}

		// action errors
		for _, e := range errResp.Errors.Action {
			if e != "" {
				fmt.Println(e)
			}
		}

		// version errors
		for _, e := range errResp.Errors.Version {
			if e != "" {
				fmt.Println(e)
			}
		}

		return nil, errors.New("unable to " + action + " nginx")
	}

	success := SuccessResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&success); err != nil {
		return nil, err
	}

	return &success, nil
}
