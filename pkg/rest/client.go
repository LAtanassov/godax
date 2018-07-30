// Package rest provides a generic rest client
// inspired byhttps://medium.com/@marcus.olsson/writing-a-go-client-for-your-restful-api-c193a2f4998c
package rest

import (
	"net/http"
	"net/url"
)

// Client provides basic http client methods
type Client interface {
	NewRequest(method, path string, body interface{}) (*http.Request, error)
	Do(req *http.Request, v interface{}) (*http.Response, error)
}

type client struct {
	BaseURL   *url.URL
	UserAgent string

	httpClient *http.Client
}

// NewClient return a json rest client
func NewClient(h *http.Client) (Client, error) {
	return nil, nil
}

func (c *client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	return nil, nil
}

func (c *client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	return nil, nil
}
