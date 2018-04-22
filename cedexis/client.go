package cedexis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2/clientcredentials"
)

const baseURL = "https://portal.cedexis.com/api/v2"

// Client implements a client for the Cedexis API
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new Cedexis API client
func NewClient(ctx context.Context, clientID string, clientSecret string) *Client {
	config := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     "https://portal.cedexis.com/api/oauth/token",
	}

	return &Client{
		httpClient: config.Client(ctx),
	}
}

func (c *Client) getJSON(url string, recv interface{}) error {
	return c.doJSON("GET", url, nil, recv)
}

func (c *Client) postJSON(url string, send interface{}, recv interface{}) error {
	return c.doJSON("POST", url, send, recv)
}

func (c *Client) doJSON(method string, url string, send interface{}, recv interface{}) error {
	toSend := new(bytes.Buffer)
	if send != nil {
		err := json.NewEncoder(toSend).Encode(send)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, url, toSend)
	if err != nil {
		return err
	}

	if send != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	req.Header.Set("User-Agent", "github.com/ctxkenb/cedexis-golang")

	//dump, err := httputil.DumpRequestOut(req, true)
	//fmt.Printf("%v\n", string(dump))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		//dump, err = httputil.DumpResponse(resp, true)
		//fmt.Printf("%v\n", string(dump))
		return fmt.Errorf("Call to Cedexis failed, error code %v", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &recv)
	if err != nil {
		return err
	}

	return nil
}
