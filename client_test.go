package fastly

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockClient struct {}

func (client MockClient) Do(req *http.Request) (*http.Response, error) {
	return nil, nil
}

type TestClient struct {
	test          *testing.T
	expectMethod  string
	expectURL     string
	expectBody    []byte
	expectHeaders map[string]string

	// return values
	response *http.Response
	err      error
}

func (client TestClient) Do(req *http.Request) (*http.Response, error) {
	assert.Equal(client.test, client.expectMethod, req.Method)
	assert.Equal(client.test, client.expectURL, req.URL.String())

	if client.expectBody == nil {
		assert.Nil(client.test, req.Body)
	} else {
		body, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		assert.Nil(client.test, err)
		assert.Equal(client.test, client.expectBody, body)
	}

	for key, header := range client.expectHeaders {
		assert.Equal(client.test, header, req.Header.Get(key))
	}

	return client.response, client.err
}

func TestNewClient(t *testing.T) {
	type want struct {
		client *Client
		err    error
	}
	tests := []struct {
		name   string
		config Config
		want
	}{
		{
			name: "Valid Config",
			config: Config{
				ApiKey:     "test-api-key",
				ServiceID:  "test-service-id",
				HttpClient: MockClient{},
			},
			want: want{
				client: &Client{
					apiKey:     "test-api-key",
					serviceID:  "test-service-id",
					httpClient: MockClient{},
				},
				err: nil,
			},
		},
		{
			name: "No Api Key",
			config: Config{
				ServiceID:  "test-service-id",
				HttpClient: MockClient{},
			},
			want: want{
				client: &Client{
					apiKey:     "",
					serviceID:  "test-service-id",
					httpClient: MockClient{},
				},
				err: nil,
			},
		},
		{
			name: "No Service ID",
			config: Config{
				ApiKey:     "test-api-key",
				HttpClient: MockClient{},
			},
			want: want{
				client: nil,
				err:    ErrNoServiceID,
			},
		},
		{
			name: "No Http Client",
			config: Config{
				ApiKey:    "test-api-key",
				ServiceID: "test-service-id",
			},
			want: want{
				client: nil,
				err:    ErrNoHttpClient,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(st *testing.T) {
			client, err := NewClient(test.config)

			assert.Equal(st, test.client, client)
			assert.Equal(st, test.err, err)
		})
	}
}

func TestRequest(t *testing.T) {
	type args struct {
		method string
		url    string
		body   []byte
	}
	tests := []struct {
		name   string
		apiKey string
		args
	}{
		{
			name:   "GET|no api key header",
			apiKey: "",
			args: args{
				method: http.MethodGet,
				url:    "https://test-api.fastly.com/test",
				body:   nil,
			},
		},
		{
			name:   "POST|no api key header|has body",
			apiKey: "",
			args: args{
				method: http.MethodPost,
				url:    "https://test-api.fastly.com/test",
				body:   []byte("test body"),
			},
		},
		{
			name:   "GET|has api key header",
			apiKey: "test-api-key",
			args: args{
				method: http.MethodGet,
				url:    "https://test-api.fastly.com/test",
				body:   nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(st *testing.T) {
			var expectBody []byte
			if test.body == nil {
				expectBody = []byte{}
			} else {
				expectBody = test.body
			}

			headers := map[string]string{}
			if len(test.apiKey) > 0 {
				headers[fastlyKeyHeader] = test.apiKey
			}

			client := &Client{
				apiKey: test.apiKey,
				httpClient: TestClient{
					test: st,
					expectMethod:  test.method,
					expectURL:     test.url,
					expectBody:    expectBody,
					expectHeaders: headers,
				},
			}
			client.request(test.method, test.url, bytes.NewBuffer(test.body))
		})
	}
}

func TestGetResponseBody(t *testing.T) {
	type want struct {
		body []byte
		err  error
	}
	tests := []struct {
		name string
		body []byte
		want
	}{
		{
			name: "Valid Body",
			body: []byte("test body"),
			want: want{
				body: []byte("test body"),
				err:  nil,
			},
		},
		{
			name: "Empty Body",
			body: nil,
			want: want{
				body: []byte{},
				err:  nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(st *testing.T) {
			recorder := httptest.NewRecorder()
			recorder.Write(test.body)

			body, err := getResponseBody(recorder.Result())

			assert.Equal(st, test.want.body, body)
			assert.Equal(st, test.err, err)
		})
	}
}
