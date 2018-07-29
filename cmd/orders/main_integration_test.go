// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func Test_main_create(t *testing.T) {
	tests := []struct {
		Name       string
		Method     string
		URL        string
		Order      orderRequest
		StatusCode int
	}{
		{Name: "should create an order", Method: "POST", URL: "http://localhost:8080/godax/v1/orders",
			Order:      orderRequest{Size: 1.34, Price: 13.34, OrderType: "limit", OrderSide: "sell", ProductID: "BTC-USD"},
			StatusCode: http.StatusOK},
		{Name: "should return Bad Request (400) for invalid body", Method: "POST", URL: "http://localhost:8080/godax/v1/orders",
			Order:      orderRequest{},
			StatusCode: http.StatusBadRequest},

		// should get an order
		// should return Not Found (404) if order does not exists
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			b, err := json.Marshal(tt.Order)
			if err != nil {
				t.Errorf("could not marshal resume %v", err)
			}
			req, err := http.NewRequest(tt.Method, tt.URL, bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/pdf")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("could do http request %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != tt.StatusCode {
				t.Errorf("http.Client{}.Do(request) status = %v, wantStatus %v", resp.StatusCode, tt.StatusCode)
			}
		})
	}
}

type orderRequest struct {
	Size      float32 `json:"size,omitempty"`
	Price     float32 `json:"price,omitempty"`
	OrderType string  `json:"type,omitempty"`
	OrderSide string  `json:"side,omitempty"`
	ProductID string  `json:"product_id,omitempty"`
}
