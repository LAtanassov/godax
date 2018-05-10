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

	// AcceptOrder accepts an existing Order
	AcceptOrder(ctx context.Context, id string) error
	// PublishOrder publish an existing Order
	PublishOrder(ctx context.Context, id string) error
	// MatchOrder matches an existing Order
	MatchOrder(ctx context.Context, id string) error
	// ConfirmOrder confirms an existing Order
	ConfirmOrder(ctx context.Context, id string) error
	// ClearOrder clears an existing Order
	ClearOrder(ctx context.Context, id string) error
	// SettleOrder settles an existing Order
	SettleOrder(ctx context.Context, id string) error
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

// AcceptOrder creates a AcceptOrder command and apply it on the Order.
func (s *service) AcceptOrder(ctx context.Context, id string) error {

	acceptOrder := &orderbook.AcceptOrder{
		CommandModel: eventsource.CommandModel{ID: id},
	}

	_, err := s.repository.Apply(ctx, acceptOrder)
	if err != nil {
		return err
	}
	return nil
}

// PublishOrder creates a PublishOrder command and apply it on the Order.
func (s *service) PublishOrder(ctx context.Context, id string) error {

	publishOrder := &orderbook.PublishOrder{
		CommandModel: eventsource.CommandModel{ID: id},
	}

	_, err := s.repository.Apply(ctx, publishOrder)
	if err != nil {
		return err
	}
	return nil
}

// MatchOrder creates a MatchOrder command and apply it on the Order.
func (s *service) MatchOrder(ctx context.Context, id string) error {

	matchOrder := &orderbook.MatchOrder{
		CommandModel: eventsource.CommandModel{ID: id},
	}

	_, err := s.repository.Apply(ctx, matchOrder)
	if err != nil {
		return err
	}
	return nil
}

// ConfirmOrder creates a ConfirmOrder command and apply it on the Order.
func (s *service) ConfirmOrder(ctx context.Context, id string) error {

	confirmOrder := &orderbook.ConfirmOrder{
		CommandModel: eventsource.CommandModel{ID: id},
	}

	_, err := s.repository.Apply(ctx, confirmOrder)
	if err != nil {
		return err
	}
	return nil
}

// ClearOrder creates a ClearOrder command and apply it on the Order.
func (s *service) ClearOrder(ctx context.Context, id string) error {

	clearOrder := &orderbook.ClearOrder{
		CommandModel: eventsource.CommandModel{ID: id},
	}

	_, err := s.repository.Apply(ctx, clearOrder)
	if err != nil {
		return err
	}
	return nil
}

// SettleOrder creates a SettleOrder command and apply it on the Order.
func (s *service) SettleOrder(ctx context.Context, id string) error {

	settleOrder := &orderbook.SettleOrder{
		CommandModel: eventsource.CommandModel{ID: id},
	}

	_, err := s.repository.Apply(ctx, settleOrder)
	if err != nil {
		return err
	}
	return nil
}
