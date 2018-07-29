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

// NewLoggingMiddleware returns a new instance of a logging middleware.
func NewLoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return &loggingService{logger, next}
	}
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
