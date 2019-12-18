package gitlab

import (
	"net/http"
)

// New creates a new instance of Gitlab's client with the required configuration
func New(host string, apiToken string) *Client {
	return &Client{
		Host: host,
		APIToken: apiToken,
		client: http.Client{},
	}
}

// Client is a gitlab client that can call gitlab APIs and fetch results or
// trigger pipelines. Use the New() function to create a client
type Client struct {
	Host string
	APIToken string
	client http.Client
}
