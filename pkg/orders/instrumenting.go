package orders

import (
	"context"
	"time"

	"github.com/LAtanassov/godax/pkg/orderbook"

	"github.com/go-kit/kit/metrics"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

// NewInstrumentingMiddleware returns an instance of the instrumented middleware.
func NewInstrumentingMiddleware(counter metrics.Counter, latency metrics.Histogram) ServiceMiddleware {
	return func(next Service) Service {
		return &instrumentingService{
			requestCount:   counter,
			requestLatency: latency,
			Service:        next,
		}
	}
}

func (s *instrumentingService) CreateOrder(ctx context.Context, size, price float32,
	orderType orderbook.OrderType, orderSide orderbook.OrderSide, productID orderbook.ProductID) (id string, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "create").Add(1)
		s.requestLatency.With("method", "create").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.CreateOrder(ctx, size, price, orderType, orderSide, productID)
}

func (s *instrumentingService) GetOrder(ctx context.Context, id string) (order orderbook.Order, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "get").Add(1)
		s.requestLatency.With("method", "get").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetOrder(ctx, id)
}

func (s *instrumentingService) CancelOrder(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "load").Add(1)
		s.requestLatency.With("method", "load").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.CancelOrder(ctx, id)
}
