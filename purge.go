package fastly

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Purge struct {
	Status string `json:"status"`
	ID     string `json:"id"`
}

// Purge items tagged with a Surrogate Key. See https://docs.fastly.com/guides/purging/
func (c *Client) PurgeSurrogateKey(key string) (*Purge, error) {
	url := fmt.Sprintf("%s/service/%s/purge/%s", fastlyBaseApiUrl, c.serviceID, key)

	return handlePurgeResponse(c.request(http.MethodPost, url, nil))
}

func handlePurgeResponse(response *http.Response, err error) (*Purge, error) {
	if err != nil {
		return nil, err
	}

	body, err := getResponseBody(response)
	if err != nil {
		return nil, err
	}

	var purge Purge
	err = json.Unmarshal(body, &purge)
	if err != nil {
		return nil, err
	}

	if purge.Status != "ok" {
		err = fmt.Errorf("cdn purge received a non-ok status[%s]", purge.Status)
	}
	return &purge, err
}
