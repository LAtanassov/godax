package orders

import (
	"context"
	"errors"

	"github.com/LAtanassov/godax/pkg/orderbook"
	"github.com/altairsix/eventsource"
)

var (
	// ErrTypeCast represents a unexpected type cast error
	ErrTypeCast = errors.New("type cast failed")
)

// Service specifies methods for Order API.
type Service interface {
	// CreateNewOrder create a new order
	CreateOrder(ctx context.Context, size, price float32,
		orderType orderbook.OrderType, side orderbook.OrderSide, productID orderbook.ProductID) (string, error)
	// CreateNewOrder create a new order
	GetOrder(ctx context.Context, id string) (orderbook.Order, error)
	// CancelOrder cancels an existing Order
	CancelOrder(ctx context.Context, id string) error
}

type service struct {
	idGenerator Generator
	repository  Repository
}

// NewService creates a booking service with necessary dependencies.
func NewService(idGenerator Generator, repository Repository) Service {
	return &service{
		idGenerator: idGenerator,
		repository:  repository,
	}
}

// CreateOrder creates a CreateOrder command and apply it on the Order.
func (s *service) CreateOrder(ctx context.Context, size, price float32,
	orderType orderbook.OrderType, orderSide orderbook.OrderSide, productID orderbook.ProductID) (string, error) {

	id := s.idGenerator.Generate()
	createOrder := &orderbook.CreateOrder{
		Size:      size,
		Price:     price,
		OrderType: orderType,
		OrderSide: orderSide,
		ProductID: productID,

		CommandModel: eventsource.CommandModel{ID: id},
	}

	_, err := s.repository.Apply(ctx, createOrder)
	if err != nil {
		return "", err
	}
	return id, nil
}

// GetOrder loads and returns the order from the repository
func (s *service) GetOrder(ctx context.Context, id string) (orderbook.Order, error) {

	v, err := s.repository.Load(ctx, id)
	if err != nil {
		return orderbook.Order{}, err
	}

	o, ok := v.(*orderbook.Order)
	if !ok {
		return orderbook.Order{}, ErrTypeCast
	}
	return *o, nil
}

// CancelOrder creates a CancelOrder command and apply it on the Order.
func (s *service) CancelOrder(ctx context.Context, id string) error {

	cancelOrder := &orderbook.CancelOrder{
		CommandModel: eventsource.CommandModel{ID: id},
	}

	_, err := s.repository.Apply(ctx, cancelOrder)
	if err != nil {
		return err
	}
	return nil
}
