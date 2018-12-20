package fastly

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPurgeSurrogateKey(t *testing.T) {
	type mock struct {
		url      string
		response *http.Response
		respErr  error
	}
	type want struct {
		purge *Purge
		err   error
	}
	tests := []struct {
		name   string
		key    string
		mock
		want
	}{
		{
			name: "",
			key: "",
			want: want{
				purge: &Purge{

				},
				err: nil,
			},
		},
	}

	apiKey := "test-api-key"
	serviceID := "test-service-id"

	for _, test := range tests {
		t.Run(test.name, func(st *testing.T) {
			client := &Client{
				apiKey:    apiKey,
				serviceID: serviceID,
				httpClient: TestClient{
					test: st,
					expectMethod:  http.MethodPost,
					expectURL:     test.url,
					expectBody:    nil,
					expectHeaders: map[string]string{
						fastlyKeyHeader:apiKey,
					},
					response: test.response,
					err:      test.respErr,
				},
			}
			purge, err := client.PurgeSurrogateKey(test.key)

			assert.Equal(st, test.purge, purge)
			assert.Equal(st, test.err, err)
		})
	}
}

func TestHandlePurgeResponse(t *testing.T) {
	type args struct {
		body []byte
		err  error
	}
	type want struct {
		purge  *Purge
		hasErr bool
	}
	tests := []struct {
		name string
		args
		want
	}{
		{
			name: "Pre-error",
			args: args{
				body: nil,
				err:  errors.New("http error"),
			},
			want: want{
				purge:  nil,
				hasErr: true,
			},
		},
		{
			name: "Valid response",
			args: args{
				body: []byte(`{"status":"ok","id":"some_id_1"}`),
				err:  nil,
			},
			want: want{
				purge: &Purge{
					Status: "ok",
					ID:     "some_id_1",
				},
				hasErr: false,
			},
		},
		{
			name: "Invalid json",
			args: args{
				body: []byte("status:ok"),
				err:  nil,
			},
			want: want{
				purge:  nil,
				hasErr: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(st *testing.T) {
			recorder := httptest.NewRecorder()
			recorder.Write(test.body)

			purge, err := handlePurgeResponse(recorder.Result(), test.err)

			assert.Equal(st, test.purge, purge)
			if test.hasErr {
				assert.NotNil(st, err)
			} else {
				assert.Nil(st, err)
			}
		})
	}
}
