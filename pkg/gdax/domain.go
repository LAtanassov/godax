package gdax

import (
	"time"
)

// ProductID represents exchange product types from GDAX
type ProductID string

const (
	// EthUsd represent the exchange product from Ethereum to US Dollar
	EthUsd ProductID = "ETH-USD"
)

// The OrderEvent type represents an order event.
type OrderEvent struct {
	Type      string    `json:"type"`
	Time      time.Time `json:"time"`
	ProductID string    `json:"product_id"`
	Sequence  int       `json:"sequence"`
	OrderID   string    `json:"order_id"`
	Size      string    `json:"size"`
	Price     string    `json:"price"`
	Side      string    `json:"side"`
	OrderType string    `json:"order_type"`
}
