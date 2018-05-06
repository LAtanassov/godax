package orders

import (
	"context"
	"time"

	"github.com/LAtanassov/godax/pkg/orderbook"

	"github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log.Logger
	Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) CreateOrder(ctx context.Context, size, price float32,
	orderType orderbook.OrderType, orderSide orderbook.OrderSide, productID orderbook.ProductID) (id string, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "create",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.Service.CreateOrder(ctx, size, price, orderType, orderSide, productID)
}

func (s *loggingService) GetOrder(ctx context.Context, id string) (order orderbook.Order, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "get",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.Service.GetOrder(ctx, id)
}

func (s *loggingService) CancelOrder(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "create",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.Service.CancelOrder(ctx, id)
}
