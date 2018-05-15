package orders

import (
	"context"
	"errors"

	"github.com/LAtanassov/godax/pkg/orderbook"
	"github.com/altairsix/eventsource"
	"github.com/altairsix/eventsource/mysqlstore"
)

var (
	// ErrUnsupportedDriver is returned when NewRepository is called with an unsupported driver
	ErrUnsupportedDriver = errors.New("unsupported driver")
)

// Repository executes the command specified and returns the current version of the aggregate
type Repository interface {
	// Apply executes the command specified and returns the current version of the aggregate
	Apply(ctx context.Context, command eventsource.Command) (int, error)
	// Load retrieves the specified aggregate from the underlying store
	Load(ctx context.Context, aggregateID string) (eventsource.Aggregate, error)
}

// DatabaseConnection contains all fields to establish a database connection
type DatabaseConnection struct {
	Driver   string
	Username string
	Password string
	Host     string
	Database string
}

var serializer = eventsource.NewJSONSerializer(
	orderbook.OrderAccepted{},
	orderbook.OrderCanceled{},
	orderbook.OrderCleared{},
	orderbook.OrderConfirmed{},
	orderbook.OrderCreated{},
	orderbook.OrderMatched{},
	orderbook.OrderPublished{},
	orderbook.OrderSettled{},
)

// NewInMemRepository return a in memory repository with oberserves
func NewInMemRepository(observers ...func(event eventsource.Event)) Repository {
	return eventsource.New(&orderbook.Order{},
		eventsource.WithSerializer(serializer),
		eventsource.WithObservers(observers...),
	)
}

// NewRepository return a repository with oberserves
func NewRepository(store eventsource.Store, observers ...func(event eventsource.Event)) Repository {
	return eventsource.New(&orderbook.Order{},
		eventsource.WithStore(store),
		eventsource.WithSerializer(serializer),
		eventsource.WithObservers(observers...),
	)
}

// NewMysqlStore return a repository with oberserves
func NewMysqlStore(tableName string, accessor mysqlstore.Accessor) (eventsource.Store, error) {
	return mysqlstore.New(tableName, accessor)
}
