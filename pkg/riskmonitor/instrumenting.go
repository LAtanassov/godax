package riskmonitor

import (
	"context"
	"time"

	"github.com/LAtanassov/godax/pkg/orderbook"
	"github.com/go-kit/kit/metrics"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	ordersLength   metrics.Histogram
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

func (s *instrumentingService) AcceptOrder(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "AcceptOrder").Add(1)
		s.requestLatency.With("method", "AcceptOrder").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.AcceptOrder(ctx, id)
}

func (s *instrumentingService) RejectOrder(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "RejectOrder").Add(1)
		s.requestLatency.With("method", "RejectOrder").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.RejectOrder(ctx, id)
}

func (s *instrumentingService) GetPendingOrders() (orders []orderbook.Order, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "GetPendingOrders").Add(1)
		s.requestLatency.With("method", "GetPendingOrders").Observe(time.Since(begin).Seconds())
		s.ordersLength.With("method", "GetPendingOrders").Observe(float64(len(orders)))
	}(time.Now())

	return s.Service.GetPendingOrders()
}
