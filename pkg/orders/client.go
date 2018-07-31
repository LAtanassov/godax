package orders

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/LAtanassov/godax/pkg/rest"
)

// Client apply action on orders api
type Client interface {
	AcceptOrder(ctx context.Context, id string) error
	RejectOrder(ctx context.Context, id string) error
}

type client struct {
	restClient rest.Client
}

// NewClient return an orders api client
func NewClient(h *http.Client, u *url.URL) Client {
	return &client{restClient: rest.NewClient(h, u)}
}

func (c *client) AcceptOrder(ctx context.Context, id string) error {
	p := fmt.Sprintf("/godax/v1/orders/%s/accept", id)
	r, err := c.restClient.NewRequest("PUT", &url.URL{Path: p}, nil)
	if err != nil {
		return err
	}
	_, err = c.restClient.Do(r, nil)
	return err
}

func (c *client) RejectOrder(ctx context.Context, id string) error {
	p := fmt.Sprintf("/godax/v1/orders/%s/reject", id)
	r, err := c.restClient.NewRequest("PUT", &url.URL{Path: p}, nil)
	if err != nil {
		return err
	}
	_, err = c.restClient.Do(r, nil)
	return err
}
