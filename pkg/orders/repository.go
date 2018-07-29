package orders

import (
	"context"
	"errors"

	"github.com/LAtanassov/godax/pkg/accessor"
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

const (
	inmem = "inmem"
	mysql = "mysql"
)

// NewRepository return a repository depending on driver
func NewRepository(dbDriver, dbURL, tableName string) (Repository, error) {

	switch dbDriver {
	case inmem:
		return newInMemRepository(), nil
	case mysql:
		accessor, err := accessor.New(dbDriver, dbURL, tableName)
		if err != nil {
			return nil, err
		}
		store, err := newMysqlStore(tableName, accessor)
		if err != nil {
			return nil, err
		}

		return newRepository(store), nil
	default:
		return nil, ErrUnsupportedDriver
	}
}

func newInMemRepository(observers ...func(event eventsource.Event)) Repository {
	return eventsource.New(&orderbook.Order{},
		eventsource.WithSerializer(serializer),
		eventsource.WithObservers(observers...),
	)
}

func newRepository(store eventsource.Store, observers ...func(event eventsource.Event)) Repository {
	return eventsource.New(&orderbook.Order{},
		eventsource.WithStore(store),
		eventsource.WithSerializer(serializer),
		eventsource.WithObservers(observers...),
	)
}

// NewMysqlStore return a repository with oberserves
func newMysqlStore(tableName string, accessor mysqlstore.Accessor) (eventsource.Store, error) {
	return mysqlstore.New(tableName, accessor)
}
