package orders

import (
	"context"

	"github.com/LAtanassov/godax/pkg/orderbook"

	"github.com/go-kit/kit/endpoint"
)

type createOrderRequest struct {
	Size      float32
	Price     float32
	OrderType orderbook.OrderType
	OrderSide orderbook.OrderSide
	ProductID orderbook.ProductID
}

type createOrderResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

func (r createOrderResponse) error() error { return r.Err }

func makeCreateOrderEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(createOrderRequest)
		if !ok {
			return nil, ErrTypeCast
		}
		id, err := s.CreateOrder(ctx, req.Size, req.Price, req.OrderType, req.OrderSide, req.ProductID)
		return createOrderResponse{ID: id, Err: err}, nil
	}
}

type getOrderRequest struct {
	ID string `json:"id"`
}

type getOrderResponse struct {
	Order orderbook.Order `json:"order"`
	Err   error           `json:"error,omitempty"`
}

func (r getOrderResponse) error() error { return r.Err }

func makeGetOrderEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r, ok := request.(getOrderRequest)
		if !ok {
			return nil, ErrTypeCast
		}
		o, err := s.GetOrder(ctx, r.ID)
		return getOrderResponse{Order: o, Err: err}, nil
	}
}

type cancelOrderRequest struct {
	ID string `json:"id"`
}

type cancelOrderResponse struct {
	Err error `json:"error,omitempty"`
}

func (r cancelOrderResponse) error() error { return r.Err }

func makeCancelOrderEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(cancelOrderRequest)
		if !ok {
			return nil, ErrTypeCast
		}
		err := s.CancelOrder(ctx, req.ID)
		return cancelOrderResponse{Err: err}, nil
	}
}
