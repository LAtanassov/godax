// Package rest provides a generic rest client
// inspired by https://medium.com/@marcus.olsson/writing-a-go-client-for-your-restful-api-c193a2f4998c
package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client represents RestClient
type Client interface {
	NewRequest(method string, path *url.URL, body interface{}) (*http.Request, error)
	Do(req *http.Request, v interface{}) (interface{}, error)
}

type client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

// NewClient return a json rest client
func NewClient(c *http.Client, b *url.URL) Client {
	if c == nil {
		c = http.DefaultClient
	}
	return &client{
		baseURL:    b,
		httpClient: c}
}

// NewRequest return a request with headers set
func (c *client) NewRequest(method string, path *url.URL, body interface{}) (*http.Request, error) {
	u := c.baseURL.ResolveReference(path)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	fmt.Println(c.baseURL.String())
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	return req, nil
}

// Do execute the request and domain object
func (c *client) Do(req *http.Request, v interface{}) (interface{}, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return v, err
}
