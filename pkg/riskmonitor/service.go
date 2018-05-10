package riskmonitor

import (
	"context"

	"github.com/LAtanassov/godax/pkg/orderbook"
)

type Service interface {
	GetOrderCreated(ctx context.Context) ([]orderbook.OrderCreated, error)
	// AcceptOrder accepts an existing Order
	AcceptOrder(ctx context.Context, id string) error
	// RejectOrder rejects an existing Order
	RejectOrder(ctx context.Context, id string) error
}
