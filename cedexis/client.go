package cedexis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

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

	zoneCache                map[int]*Zone
	privatePlatformListCache map[int]*PlatformInfo
	privatePlatformCache     map[int]*PlatformConfig
	appCache                 map[int]*Application
	countriesCache           map[int]*Country
}

// NewClient creates a new Cedexis API client
func NewClient(ctx context.Context, clientID string, clientSecret string) *Client {
	config := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     "https://portal.cedexis.com/api/oauth/token",
	}

	return &Client{
		httpClient:               config.Client(ctx),
		zoneCache:                map[int]*Zone{},
		privatePlatformListCache: map[int]*PlatformInfo{},
		privatePlatformCache:     map[int]*PlatformConfig{},
		appCache:                 map[int]*Application{},
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
	data := []byte{}
	var err error
	if send != nil {
		data, err = json.Marshal(send)
		if err != nil {
			return err
		}
	}

	resp, err := c.doHTTP(method, url, data)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if recv != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, &recv)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) doHTTP(method string, url string, toSend []byte) (*http.Response, error) {
	delay := time.Duration(1)

	for {
		req, err := http.NewRequest(method, url, bytes.NewReader(toSend))
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
		if resp.StatusCode == 429 {
			fmt.Printf("Rate Limited, sleeping for %v\n", delay*time.Second)
			time.Sleep(delay * time.Second)
			delay = delay << 1
			continue
		}
		if resp.StatusCode >= 400 {
			// dump, _ := httputil.DumpResponse(resp, true)
			// fmt.Printf("%v\n", string(dump))
			return nil, errorFromHTTPFailure(resp)
		}

		return resp, nil
	}
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
