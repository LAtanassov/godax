package orders

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LAtanassov/godax/pkg/orderbook"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/gorilla/mux"
)

// MakeHandler returns a handler for the order service.
func MakeHandler(s Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	createOrderHandler := kithttp.NewServer(
		makeCreateOrderEndpoint(s),
		decodeCreateOrderRequest,
		encodeResponse,
		opts...,
	)

	getOrderHandler := kithttp.NewServer(
		makeGetOrderEndpoint(s),
		decodeGetOrderRequest,
		encodeResponse,
		opts...,
	)

	cancelOrderHandler := kithttp.NewServer(
		makeCancelOrderEndpoint(s),
		decodeCancelOrderRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/godax/v1/orders", createOrderHandler).Methods("POST")
	r.Handle("/godax/v1/orders/{id}", getOrderHandler).Methods("GET")
	r.Handle("/godax/v1//orders/{id}", cancelOrderHandler).Methods("DELETE")

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

func decodeCancelOrderRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	return getOrderRequest{ID: id}, nil
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
	orderbook.BTC_USD.String(): orderbook.BTC_USD,
}

const EPSILON float32 = 0.00000001

func floatEquals(a, b float32) bool {
	return (a-b) < EPSILON && (b-a) < EPSILON
}
