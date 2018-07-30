package orders

import (
	"context"
	"net/http"

	"github.com/LAtanassov/godax/pkg/rest"
)

type Client interface {
	AcceptOrder(ctx context.Context, id string) error
	RejectOrder(ctx context.Context, id string) error
}

type client struct {
	client *rest.Client
}

// NewClient return an orders api client
func NewClient(h *http.Client) (Client, error) {
	return &client{}, nil
}

func (c *client) AcceptOrder(ctx context.Context, id string) error {
	return nil
}

func (c *client) RejectOrder(ctx context.Context, id string) error {
	return nil
}
