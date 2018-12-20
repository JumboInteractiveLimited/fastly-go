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
		body []byte
		err  error
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
			name: "Valid Purge",
			key:  "surrogate-key",
			mock: mock{
				body: []byte(`{"status":"ok","id":"test-id-1"}`),
				err:  nil,
			},
			want: want{
				purge: &Purge{
					Status: "ok",
					ID:     "test-id-1",
				},
				err: nil,
			},
		},
		{
			name: "Bad Purge Response",
			key:  "surrogate-key",
			mock: mock{
				body: []byte(`{"status":"bad","id":"test-id-2"}`),
				err:  nil,
			},
			want: want{
				purge: &Purge{
					Status: "bad",
					ID:     "test-id-2",
				},
				err: errors.New("cdn purge received a non-ok status[bad]"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(st *testing.T) {
			recorder := httptest.NewRecorder()
			recorder.Write(test.body)

			client := &Client{
				apiKey:    "test-api-key",
				serviceID: "test-service-id",
				httpClient: TestClient{
					test:          st,
					expectMethod:  http.MethodPost,
					expectURL:     "https://api.fastly.com/service/test-service-id/purge/" + test.key,
					expectBody:    nil,
					expectHeaders: map[string]string{fastlyKeyHeader: "test-api-key"},
					response:      recorder.Result(),
					err:           test.mock.err,
				},
			}

			purge, err := client.PurgeSurrogateKey(test.key)

			assert.Equal(st, test.purge, purge)
			assert.Equal(st, test.want.err, err)
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
