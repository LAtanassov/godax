package orders

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/LAtanassov/godax/pkg/orderbook"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	kithttp "github.com/go-kit/kit/transport/http"
	"golang.org/x/time/rate"

	"github.com/gorilla/mux"
)

// circuit breaker and rate limit does not belong here
// rate limit per IP adress would be better
func newCircuitBreakerMiddleware(commandName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return circuitbreaker.Hystrix(commandName)(next)
	}
}

func newRatelimitMiddleware(limit ratelimit.Allower) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return ratelimit.NewErroringLimiter(limit)(next)
	}
}

// MakeHandler returns a handler for the order service.
func MakeHandler(s Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	c := makeCreateOrderEndpoint(s)
	c = newCircuitBreakerMiddleware("create order")(c)
	c = newRatelimitMiddleware(rate.NewLimiter(rate.Every(time.Second), 100))(c)
	createOrderHandler := kithttp.NewServer(
		c,
		decodeCreateOrderRequest,
		encodeResponse,
		opts...,
	)

	g := makeGetOrderEndpoint(s)
	g = newCircuitBreakerMiddleware("get order")(g)
	g = newRatelimitMiddleware(rate.NewLimiter(rate.Every(time.Second), 100))(g)
	getOrderHandler := kithttp.NewServer(
		g,
		decodeGetOrderRequest,
		encodeResponse,
		opts...,
	)

	cancelOrderHandler := kithttp.NewServer(
		makeCancelOrderEndpoint(s),
		decodeCommonOrderRequest,
		encodeResponse,
		opts...,
	)

	acceptOrderHandler := kithttp.NewServer(
		makeAcceptOrderEndpoint(s),
		decodeCommonOrderRequest,
		encodeResponse,
		opts...,
	)

	publishOrderHandler := kithttp.NewServer(
		makePublishOrderEndpoint(s),
		decodeCommonOrderRequest,
		encodeResponse,
		opts...,
	)

	matchOrderHandler := kithttp.NewServer(
		makeMatchOrderEndpoint(s),
		decodeCommonOrderRequest,
		encodeResponse,
		opts...,
	)

	confirmOrderHandler := kithttp.NewServer(
		makeConfirmOrderEndpoint(s),
		decodeCommonOrderRequest,
		encodeResponse,
		opts...,
	)

	clearOrderHandler := kithttp.NewServer(
		makeClearOrderEndpoint(s),
		decodeCommonOrderRequest,
		encodeResponse,
		opts...,
	)

	settleOrderHandler := kithttp.NewServer(
		makeSettleOrderEndpoint(s),
		decodeCommonOrderRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/godax/v1/orders", createOrderHandler).Methods("POST")
	r.Handle("/godax/v1/orders/{id}", getOrderHandler).Methods("GET")
	r.Handle("/godax/v1/orders/{id}", cancelOrderHandler).Methods("DELETE")

	r.Handle("/godax/v1/orders/{id}/accept", acceptOrderHandler).Methods("PUT")
	r.Handle("/godax/v1/orders/{id}/publish", publishOrderHandler).Methods("PUT")
	r.Handle("/godax/v1/orders/{id}/match", matchOrderHandler).Methods("PUT")
	r.Handle("/godax/v1/orders/{id}/confirm", confirmOrderHandler).Methods("PUT")
	r.Handle("/godax/v1/orders/{id}/clear", clearOrderHandler).Methods("PUT")
	r.Handle("/godax/v1/orders/{id}/settle", settleOrderHandler).Methods("PUT")

	return r
}

var errBadRoute = errors.New("bad route")
var errIllegalArgument = errors.New("illegal argument")

func decodeCreateOrderRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Size      float32 `json:"size"`
		Price     float32 `json:"price"`
		OrderType string  `json:"type"`
		OrderSide string  `json:"side"`
		ProductID string  `json:"product_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	defer r.Body.Close()

	// validation
	if floatEquals(body.Price, 0.0) || floatEquals(body.Size, 0.0) {
		return nil, errIllegalArgument
	}

	orderType, ok := orderTypes[body.OrderType]
	if !ok {
		return nil, errIllegalArgument
	}

	orderSide, ok := orderSides[body.OrderSide]
	if !ok {
		return nil, errIllegalArgument
	}

	productID, ok := productIDs[body.ProductID]
	if !ok {
		return nil, errIllegalArgument
	}

	return createOrderRequest{
		Size:      body.Size,
		Price:     body.Price,
		OrderType: orderType,
		OrderSide: orderSide,
		ProductID: productID,
	}, nil
}

func decodeGetOrderRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	return getOrderRequest{ID: id}, nil
}

func decodeCommonOrderRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	return commonOrderRequest{ID: id}, nil
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case errBadRoute:
		w.WriteHeader(http.StatusNotFound)
	case errIllegalArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

var orderTypes = map[string]orderbook.OrderType{
	orderbook.Limit.String():  orderbook.Limit,
	orderbook.Market.String(): orderbook.Market,
}

var orderSides = map[string]orderbook.OrderSide{
	orderbook.Sell.String(): orderbook.Sell,
	orderbook.Buy.String():  orderbook.Sell,
}

var productIDs = map[string]orderbook.ProductID{
	orderbook.BtcUsd.String(): orderbook.BtcUsd,
}

// EPSILON is a small float number
const EPSILON float32 = 0.00000001

func floatEquals(a, b float32) bool {
	return (a-b) < EPSILON && (b-a) < EPSILON
}
