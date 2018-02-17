package gdax

import (
	"time"
)

// The Subscribe type represents a request to which product ids you want to subscibe.
type Subscribe struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
}

// The Subscription type represents a subscription.
type Subscription struct {
	Type string `json:"type"`
}

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
