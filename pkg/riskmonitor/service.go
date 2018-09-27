package riskmonitor

import (
	"context"

	"github.com/LAtanassov/godax/pkg/orderbook"
	"github.com/LAtanassov/godax/pkg/orders"
)

// Service accepts or rejects orders either by automation or by a decision of a risk analyst.
type Service interface {
	// AcceptOrder accepts an existing Order
	AcceptOrder(ctx context.Context, id string) error
	// RejectOrder rejects an existing Order
	RejectOrder(ctx context.Context, id string) error
	// GetPendingOrders returns them sorted (oldest first) and limited to 50
	GetPendingOrders() ([]orderbook.Order, error)
}

// ServiceMiddleware is a chainable behavior modifier for Service.
type ServiceMiddleware func(Service) Service

type service struct {
	client     orders.Client
	repository Repository
}

// NewService creates a booking service with necessary dependencies.
func NewService(c orders.Client, r Repository) Service {
	return &service{
		client:     c,
		repository: r,
	}
}

func (s *service) AcceptOrder(ctx context.Context, id string) error {
	return s.client.AcceptOrder(ctx, id)
}

func (s *service) RejectOrder(ctx context.Context, id string) error {
	return s.client.RejectOrder(ctx, id)
}

func (s *service) GetPendingOrders() ([]orderbook.Order, error) {
	return s.repository.GetPendingOrders()
}
