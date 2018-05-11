package riskmonitor

import (
	"context"
)

// Service accepts or rejects orders either by automation or by a decision of a risk analyst.
type Service interface {
	// AcceptOrder accepts an existing Order
	AcceptOrder(ctx context.Context, id string) error
	// RejectOrder rejects an existing Order
	RejectOrder(ctx context.Context, id string) error
}
