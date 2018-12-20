package fastly

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	fastlyKeyHeader  = "Fastly-Key"
	fastlyBaseApiUrl = "https://api.fastly.com" // The main entry point for fastly API
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// A reusable client for making calls to Fastly.
type Client struct {
	apiKey     string
	serviceID  string
	httpClient HttpClient
}

// Initial configuration for the client. ApiKey is optional as some calls to Fastly do not require one.
type Config struct {
	ApiKey     string
	ServiceID  string
	HttpClient HttpClient
}

var (
	ErrNoServiceID  = errors.New("no service ID")
	ErrNoHttpClient = errors.New("no http client")
)

// Creates a reusable client.
func NewClient(config Config) (*Client, error) {
	if len(config.ServiceID) == 0 {
		return nil, ErrNoServiceID
	}
	if config.HttpClient == nil {
		return nil, ErrNoHttpClient
	}

	return &Client{
		apiKey:     config.ApiKey,
		serviceID:  config.ServiceID,
		httpClient: config.HttpClient,
	}, nil
}

func (c *Client) request(method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if len(c.apiKey) > 0 {
		req.Header.Set(fastlyKeyHeader, c.apiKey)
	}

	return c.httpClient.Do(req)
}

func getResponseBody(response *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(response.Body)
	defer func() {
		err = response.Body.Close()
	}()

	return body, err
}
