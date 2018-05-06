package orderbook

import (
	"context"
	"errors"
	"time"

	"github.com/altairsix/eventsource"
)

// OrderType represents an enum of order types
type OrderType int

const (
	// Limit order type lets you set your own price
	Limit OrderType = iota
	// Market order type will be executed immediately at the current market price
	Market
)

func (o OrderType) String() string {
	switch o {
	case Limit:
		return "limit"
	case Market:
		return "market"
	}

	return ""
}

// OrderSide represents an enum of order sides
type OrderSide int

const (
	// Sell a specific product
	Sell OrderSide = iota
	// Buy a specific product
	Buy
)

func (o OrderSide) String() string {
	switch o {
	case Sell:
		return "sell"
	case Buy:
		return "buy"
	}

	return ""
}

// ProductID represents an enum of product ids
type ProductID int

const (
	// BtcUsd product id represents the market of Bitcoin and US dollar
	BtcUsd ProductID = iota
)

func (p ProductID) String() string {
	switch p {
	case BtcUsd:
		return "BTC-USD"
	}

	return ""
}

const (
	stateCreated   = "created"
	stateAccepted  = "accepted"
	statePublished = "published"
	stateCanceled  = "canceled"
	stateMatched   = "matched"
	stateConfirmed = "confirmed"
	stateCleared   = "cleared"
	stateSettled   = "settled"
)

var (
	// ErrUnknownEvent represents an unknown event
	ErrUnknownEvent = errors.New("unknown event")

	// ErrUnknownCommand represents an unknown command
	ErrUnknownCommand = errors.New("unknown command")
	// ErrInvalidStateTransition is returned when preconditions are not met
	ErrInvalidStateTransition = errors.New("invalid state transition")
)

// Events --------------

// OrderCreated Event - created by the system
type OrderCreated struct {
	Size      float32
	Price     float32
	OrderType OrderType
	OrderSide OrderSide
	ProductID ProductID
	eventsource.Model
}

// OrderAccepted Event - accepted by a risk analyst
type OrderAccepted struct {
	eventsource.Model
}

// OrderCanceled Event - an order was canceled
type OrderCanceled struct {
	eventsource.Model
}

// OrderPublished Event - published on the exchange
type OrderPublished struct {
	eventsource.Model
}

// OrderMatched Event - sell and buy order was matched
type OrderMatched struct {
	eventsource.Model
}

// OrderConfirmed - both clients confirmed the trade
type OrderConfirmed struct {
	eventsource.Model
}

// OrderCleared - all calculations and obligations are fullfilled
type OrderCleared struct {
	eventsource.Model
}

// OrderSettled - money and security settled
type OrderSettled struct {
	eventsource.Model
}

// Commands --------------

// CreateOrder Command
type CreateOrder struct {
	Size      float32
	Price     float32
	OrderType OrderType
	OrderSide OrderSide
	ProductID ProductID

	eventsource.CommandModel
}

// AcceptOrder Command
type AcceptOrder struct {
	eventsource.CommandModel
}

// CancelOrder Command
type CancelOrder struct {
	eventsource.CommandModel
}

// PublishOrder Command
type PublishOrder struct {
	eventsource.CommandModel
}

// MatchOrder Command
type MatchOrder struct {
	eventsource.CommandModel
}

// ConfirmOrder Command
type ConfirmOrder struct {
	eventsource.CommandModel
}

// ClearOrder Command
type ClearOrder struct {
	eventsource.CommandModel
}

// SettleOrder Command
type SettleOrder struct {
	eventsource.CommandModel
}

// Aggregates --------------

// Order is an Aggregate which apply Events
type Order struct {
	Size      float32
	Price     float32
	OrderType OrderType
	OrderSide OrderSide
	ProductID ProductID

	id        string
	version   int
	createdAt time.Time
	updatedAt time.Time
	state     string
}

// On an incoming event apply updates to the order (aggregate).
// After all events were applied the order represents the latest state.
func (o *Order) On(event eventsource.Event) error {
	switch v := event.(type) {
	case *OrderCreated:
		o.Size = v.Size
		o.Price = v.Price
		o.OrderType = v.OrderType
		o.ProductID = v.ProductID
		o.OrderSide = v.OrderSide

		o.createdAt = v.At
		o.state = stateCreated
	case *OrderAccepted:
		o.state = stateAccepted
	case *OrderCanceled:
		o.state = stateCanceled
	case *OrderPublished:
		o.state = statePublished
	case *OrderMatched:
		o.state = stateMatched
	case *OrderConfirmed:
		o.state = stateConfirmed
	case *OrderCleared:
		o.state = stateCleared
	case *OrderSettled:
		o.state = stateSettled
	default:
		return ErrUnknownEvent
	}

	o.id = event.AggregateID()
	o.version = event.EventVersion()
	o.updatedAt = event.EventAt()

	return nil
}

// Apply generates events from a command
func (o *Order) Apply(ctx context.Context, command eventsource.Command) ([]eventsource.Event, error) {
	switch v := command.(type) {
	case *CreateOrder:
		orderCreated := &OrderCreated{
			Model: eventsource.Model{ID: v.AggregateID(), Version: o.version + 1, At: time.Now()},
		}
		return []eventsource.Event{orderCreated}, nil
	case *AcceptOrder:
		if o.state != stateCreated {
			return nil, ErrInvalidStateTransition
		}
		orderAccepted := &OrderAccepted{
			Model: eventsource.Model{ID: v.AggregateID(), Version: o.version + 1, At: time.Now()},
		}
		return []eventsource.Event{orderAccepted}, nil
	case *CancelOrder:
		if o.state != stateCreated {
			return nil, ErrInvalidStateTransition
		}
		orderCanceled := &OrderCanceled{
			Model: eventsource.Model{ID: v.AggregateID(), Version: o.version + 1, At: time.Now()},
		}
		return []eventsource.Event{orderCanceled}, nil
	case *PublishOrder:
		if o.state != stateAccepted {
			return nil, ErrInvalidStateTransition
		}
		orderPublished := &OrderPublished{
			Model: eventsource.Model{ID: v.AggregateID(), Version: o.version + 1, At: time.Now()},
		}
		return []eventsource.Event{orderPublished}, nil
	case *MatchOrder:
		if o.state != statePublished {
			return nil, ErrInvalidStateTransition
		}
		orderMatched := &OrderMatched{
			Model: eventsource.Model{ID: v.AggregateID(), Version: o.version + 1, At: time.Now()},
		}
		return []eventsource.Event{orderMatched}, nil
	case *ConfirmOrder:
		if o.state != stateMatched {
			return nil, ErrInvalidStateTransition
		}
		orderConfirmed := &OrderConfirmed{
			Model: eventsource.Model{ID: v.AggregateID(), Version: o.version + 1, At: time.Now()},
		}
		return []eventsource.Event{orderConfirmed}, nil
	case *ClearOrder:
		if o.state != stateConfirmed {
			return nil, ErrInvalidStateTransition
		}
		orderCleared := &OrderCleared{
			Model: eventsource.Model{ID: v.AggregateID(), Version: o.version + 1, At: time.Now()},
		}
		return []eventsource.Event{orderCleared}, nil
	case *SettleOrder:
		if o.state != stateCleared {
			return nil, ErrInvalidStateTransition
		}
		orderSettled := &OrderSettled{
			Model: eventsource.Model{ID: v.AggregateID(), Version: o.version + 1, At: time.Now()},
		}
		return []eventsource.Event{orderSettled}, nil

	default:
		return nil, ErrUnknownCommand
	}
}
