// +build integration

package messaging

import (
	"io"
	"strings"
	"testing"
)

func Test_publisher_Publish(t *testing.T) {
	tests := []struct {
		name    string
		p       Publisher
		r       io.Reader
		wantErr bool
	}{
		{"should publish an event", NewSimplePublisher("amqp://guest:guest@localhost:5672/", "gotest"), strings.NewReader("test"), false},
		{"should return an error if the event could not been publised", NewSimplePublisher("amqp://guest:guest@localhost:9999/", "gotest"), strings.NewReader("test"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.Publish(tt.r); (err != nil) != tt.wantErr {
				t.Errorf("publisher.Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
