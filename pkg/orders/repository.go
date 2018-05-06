package orders

import (
	"context"
	"errors"
	"fmt"

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

// NewRepository return a Repository depending on driver (inmem, mysql (defaild))
func NewRepository(sqlDriver, sqlHost, sqlDbName, sqlUser, sqlPwd, tabName string) (Repository, error) {
	serializer := eventsource.NewJSONSerializer(
		orderbook.OrderCreated{},
		orderbook.OrderCanceled{},
	)

	if sqlDriver == "inmem" {
		return eventsource.New(&orderbook.Order{},
			eventsource.WithSerializer(serializer),
		), nil
	}

	if sqlDriver != "mysql" {
		return nil, ErrUnsupportedDriver
	}

	acc, err := accessor.New(sqlDriver, fmt.Sprintf("%s:%s@tcp(%s)/%s", sqlUser, sqlPwd, sqlHost, sqlDbName), tabName)
	if err != nil {
		return nil, err
	}

	store, err := mysqlstore.New(tabName, acc)
	if err != nil {
		return nil, err
	}

	repo := eventsource.New(&orderbook.Order{},
		eventsource.WithStore(store),
		eventsource.WithSerializer(serializer),
	)
	return repo, nil
}
