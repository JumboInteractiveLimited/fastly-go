# fastly-go [![Build Status](https://travis-ci.org/JumboInteractiveLimited/fastly-go.svg?branch=master)](https://travis-ci.org/JumboInteractiveLimited/fastly-go) [![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg)](https://godoc.org/github.com/JumboInteractiveLimited/fastly-go)

A Golang interface for the [Fastly API](https://docs.fastly.com/api).

# Usage

Get the library:

    $ go get -v github.com/JumboInteractiveLimited/fastly-go

Purge items tagged with a Surrogate Key. See https://docs.fastly.com/guides/purging/

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/JumboInteractiveLimited/fastly-go"
)

func main() {
	err := purgeSurrogateKey("surrogate-key")
	if err != nil {
		fmt.Println(err)
	}
}

func purgeSurrogateKey(surrogatekey string) error {
	// Create a Fastly client. The client can be reused for multiple calls to Fastly.
	// Adding an ApiKey is optional, since some calls to Fastly do not require one.
	// ServiceID and HttpClient must be set.
	// HttpClient is a go interface that implements the http.Client's 'Do' method.
	fastlyClient, err := fastly.NewClient(fastly.Config{
		ApiKey:     "your-api-key",
		ServiceID:  "your-service-id",
		HttpClient: &http.Client{},
	})
	if err != nil {
		return err
	}

	// PurgeSurrogateKey returns an error if the call fails, or if a non-ok status is received.
	// The Purge object is returned if the caller can make use of the Purge ID
	_, err = fastlyClient.PurgeSurrogateKey(surrogatekey)
	return err
}
```
