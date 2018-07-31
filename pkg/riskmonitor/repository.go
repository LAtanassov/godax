package riskmonitor

import "github.com/LAtanassov/godax/pkg/orderbook"

// Repository abstracts database
type Repository interface {
	// GetPendingOrders returns them sorted (oldest first) and limited to 50
	GetPendingOrders() ([]orderbook.Order, error)
}
