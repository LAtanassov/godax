package orderbook

import (
	"context"
	"errors"
	"time"

	"github.com/altairsix/eventsource"
)

type OrderType int

const (
	Limit OrderType = iota
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

type OrderSide int

const (
	Sell OrderSide = iota
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

type ProductID int

const (
	BTC_USD ProductID = iota
)

func (p ProductID) String() string {
	switch p {
	case BTC_USD:
		return "BTC-USD"
	}

	return ""
}

const (
	StateCreated   = "created"
	StateAccepted  = "accepted"
	StatePublished = "published"
	StateCanceled  = "canceled"
	StateMatched   = "matched"
	StateConfirmed = "confirmed"
	StateCleared   = "cleared"
	StateSettled   = "settled"
)

var (
	ErrUnknownEvent = errors.New("unknown event")

	ErrOrderNotCreated   = errors.New("order not created")
	ErrOrderNotAccepted  = errors.New("order not accepted")
	ErrOrderNotPublished = errors.New("order not published")
	ErrOrderNotMatched   = errors.New("order not matched")
	ErrOrderNotConfirmed = errors.New("order not confirmed")
	ErrOrderNotCleared   = errors.New("order not cleared")
	ErrUnknownCommand    = errors.New("unknown command")
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
		o.state = StateCreated
	case *OrderAccepted:
		o.state = StateAccepted
	case *OrderCanceled:
		o.state = StateCanceled
	case *OrderPublished:
		o.state = StatePublished
	case *OrderMatched:
		o.state = StateMatched
	case *OrderConfirmed:
		o.state = StateConfirmed
	case *OrderCleared:
		o.state = StateCleared
	case *OrderSettled:
		o.state = StateSettled
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
		if o.state != StateCreated {
			return nil, ErrOrderNotCreated
		}
		orderAccepted := &OrderAccepted{
			Model: eventsource.Model{ID: v.AggregateID(), Version: o.version + 1, At: time.Now()},
		}
		return []eventsource.Event{orderAccepted}, nil
	case *CancelOrder:
		if o.state != StateCreated {
			return nil, ErrOrderNotCreated
		}
		orderCanceled := &OrderCanceled{
			Model: eventsource.Model{ID: v.AggregateID(), Version: o.version + 1, At: time.Now()},
		}
		return []eventsource.Event{orderCanceled}, nil
	case *PublishOrder:
		if o.state != StateAccepted {
			return nil, ErrOrderNotAccepted
		}
		orderPublished := &OrderPublished{
			Model: eventsource.Model{ID: v.AggregateID(), Version: o.version + 1, At: time.Now()},
		}
		return []eventsource.Event{orderPublished}, nil
	case *MatchOrder:
		if o.state != StatePublished {
			return nil, ErrOrderNotPublished
		}
		orderMatched := &OrderMatched{
			Model: eventsource.Model{ID: v.AggregateID(), Version: o.version + 1, At: time.Now()},
		}
		return []eventsource.Event{orderMatched}, nil
	case *ConfirmOrder:
		if o.state != StateMatched {
			return nil, ErrOrderNotMatched
		}
		orderConfirmed := &OrderConfirmed{
			Model: eventsource.Model{ID: v.AggregateID(), Version: o.version + 1, At: time.Now()},
		}
		return []eventsource.Event{orderConfirmed}, nil
	case *ClearOrder:
		if o.state != StateConfirmed {
			return nil, ErrOrderNotConfirmed
		}
		orderCleared := &OrderCleared{
			Model: eventsource.Model{ID: v.AggregateID(), Version: o.version + 1, At: time.Now()},
		}
		return []eventsource.Event{orderCleared}, nil
	case *SettleOrder:
		if o.state != StateCleared {
			return nil, ErrOrderNotCleared
		}
		orderSettled := &OrderSettled{
			Model: eventsource.Model{ID: v.AggregateID(), Version: o.version + 1, At: time.Now()},
		}
		return []eventsource.Event{orderSettled}, nil

	default:
		return nil, ErrUnknownCommand
	}
}
