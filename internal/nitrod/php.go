package nitrod

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type PhpFpmOptions struct {
	Version string
	Action  string
}

func (c *Client) ServicePhpFpm(ctx context.Context, options *PhpFpmOptions) (*SuccessResponse, error) {
	version := "7.4"
	action := "restart"
	if options != nil {
		version = options.Version
		action = options.Action
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/services/php-fpm", c.BaseURL),
		strings.NewReader(fmt.Sprintf(`{"action":"%s","version":"%s"}`, action, version)),
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

		return nil, errors.New("unable to " + action + " php-fpm")
	}

	success := SuccessResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&success); err != nil {
		return nil, err
	}

	return &success, nil
}
