package cedexis

import "errors"

const pingPath = "/meta/system.json/ping"

type pingResponse struct {
	Result string `json:"result"`
}

// Ping validates connectivity to Cedexis API - returns error on failure
func (c *Client) Ping() error {
	var resp pingResponse
	err := c.getJSON(baseURL+pingPath, &resp)

	if err != nil {
		return err
	}

	if resp.Result == "pong" {
		return nil
	}

	return errors.New("Invalid ping response")
}
