package wmata

import (
	"net/http"
	"time"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewWMATADefaultClient returns a new client to make requests to the WMATA API
// This creates a default http.Client with a 30 second timeout
func NewWMATADefaultClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// NewWMATAClient returns a new client to make requests to the WMATA API
func NewWMATAClient(apiKey string, httpClient http.Client) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &httpClient,
	}
}
