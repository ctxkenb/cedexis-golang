package cedexis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2/clientcredentials"
)

const baseURL = "https://portal.cedexis.com/api/v2"

type cedexisError struct {
	DeveloperMessage string `json:"developerMessage"`
	UserMessage      string `json:"userMessage"`
	Field            string `json:"field"`
	ErrorCode        string `json:"errorCode"`
	MoreInfo         string `json:"moreInfo"`
	RootCause        string `json:"rootCause"`
}

type cedexisErrorResponse struct {
	HTTPStatus   int            `json:"httpStatus"`
	ErrorDetails []cedexisError `json:"errorDetails"`
}

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

func (c *Client) delete(url string) error {
	_, err := c.doHTTP("DELETE", url, nil)
	return err
}

func (c *Client) getJSON(url string, recv interface{}) error {
	return c.doJSON("GET", url, nil, recv)
}

func (c *Client) postJSON(url string, send interface{}, recv interface{}) error {
	return c.doJSON("POST", url, send, recv)
}

func (c *Client) putJSON(url string, send interface{}, recv interface{}) error {
	return c.doJSON("PUT", url, send, recv)
}

func (c *Client) doJSON(method string, url string, send interface{}, recv interface{}) error {
	toSend := new(bytes.Buffer)
	if send != nil {
		err := json.NewEncoder(toSend).Encode(send)
		if err != nil {
			return err
		}
	}

	resp, err := c.doHTTP(method, url, toSend)
	if err != nil {
		return err
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

func (c *Client) doHTTP(method string, url string, toSend io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, toSend)
	if err != nil {
		return nil, err
	}

	if toSend != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	req.Header.Set("User-Agent", "github.com/ctxkenb/cedexis-golang")

	// dump, err := httputil.DumpRequestOut(req, true)
	// fmt.Printf("%v\n", string(dump))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		// dump, _ := httputil.DumpResponse(resp, true)
		// fmt.Printf("%v\n", string(dump))
		return nil, errorFromHTTPFailure(resp)
	}

	return resp, nil
}

func errorFromHTTPFailure(resp *http.Response) error {
	defer resp.Body.Close()
	body, errErr := ioutil.ReadAll(resp.Body)
	if errErr != nil {
		return fmt.Errorf("Call to Cedexis failed, error code %v", resp.StatusCode)
	}

	cedexisError := cedexisErrorResponse{}
	errErr = json.Unmarshal(body, &cedexisError)
	if errErr != nil {
		return fmt.Errorf("Call to Cedexis failed, error code %v", resp.StatusCode)
	}

	return fmt.Errorf("Call to Cedexis failed, because '%s'", cedexisError.ErrorDetails[0].UserMessage)
}
