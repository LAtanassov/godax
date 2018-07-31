package riskmonitor

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/LAtanassov/godax/pkg/orderbook"
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

func (s *loggingService) AcceptOrder(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "AcceptOrder",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.Service.AcceptOrder(ctx, id)
}

func (s *loggingService) RejectOrder(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "RejectOrder",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.Service.RejectOrder(ctx, id)
}

func (s *loggingService) GetPendingOrders() (orders []orderbook.Order, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "GetPendingOrders",
			"ordersLength", len(orders),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.Service.GetPendingOrders()
}
