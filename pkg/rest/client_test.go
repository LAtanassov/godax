package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

type TestData struct {
	ID int `json:"id"`
}

func Test_client_NewRequest(t *testing.T) {

	type args struct {
		method string
		path   *url.URL
		body   interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *TestData
		wantErr bool
	}{
		{"should GET TestData", args{"GET", &url.URL{Path: "/test"}, nil}, &TestData{ID: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				json.NewEncoder(w).Encode(&TestData{ID: 1})
			}))
			defer ts.Close()
			u, err := url.Parse(ts.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.NewRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			c := NewClient(nil, u)

			r, err := c.NewRequest(tt.args.method, tt.args.path, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.NewRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := c.Do(r, &TestData{})
			if (err != nil) != tt.wantErr {
				t.Errorf("client.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("client.Do() = %v, want %v", got, tt.want)
			}
		})
	}
}
