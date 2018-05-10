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

type commonOrderRequest struct {
	ID string `json:"id"`
}

type commonOrderResponse struct {
	Err error `json:"error,omitempty"`
}

func (r commonOrderResponse) error() error { return r.Err }

func makeCancelOrderEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(commonOrderRequest)
		if !ok {
			return nil, ErrTypeCast
		}
		err := s.CancelOrder(ctx, req.ID)
		return commonOrderResponse{Err: err}, nil
	}
}

func makeAcceptOrderEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(commonOrderRequest)
		if !ok {
			return nil, ErrTypeCast
		}
		err := s.AcceptOrder(ctx, req.ID)
		return commonOrderResponse{Err: err}, nil
	}
}

func makePublishOrderEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(commonOrderRequest)
		if !ok {
			return nil, ErrTypeCast
		}
		err := s.PublishOrder(ctx, req.ID)
		return commonOrderResponse{Err: err}, nil
	}
}

func makeMatchOrderEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(commonOrderRequest)
		if !ok {
			return nil, ErrTypeCast
		}
		err := s.MatchOrder(ctx, req.ID)
		return commonOrderResponse{Err: err}, nil
	}
}

func makeConfirmOrderEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(commonOrderRequest)
		if !ok {
			return nil, ErrTypeCast
		}
		err := s.ConfirmOrder(ctx, req.ID)
		return commonOrderResponse{Err: err}, nil
	}
}

func makeClearOrderEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(commonOrderRequest)
		if !ok {
			return nil, ErrTypeCast
		}
		err := s.ClearOrder(ctx, req.ID)
		return commonOrderResponse{Err: err}, nil
	}
}

func makeSettleOrderEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(commonOrderRequest)
		if !ok {
			return nil, ErrTypeCast
		}
		err := s.SettleOrder(ctx, req.ID)
		return commonOrderResponse{Err: err}, nil
	}
}
