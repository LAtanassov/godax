package orders

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/LAtanassov/godax/pkg/orderbook"
	"github.com/altairsix/eventsource"
)

func Test_service_CreateOrder(t *testing.T) {
	type fields struct {
		idGenerator Generator
		repository  Repository
	}
	tests := []struct {
		name    string
		fields  fields
		ctx     context.Context
		want    string
		wantErr bool
	}{
		{"should apply CreateOrder command to repository",
			fields{&mockIDGenerator{id: "AB-CD"}, &mockRepository{wantErr: false, err: nil,
				command: orderbook.CreateOrder{Size: 1.0, Price: 1.0, OrderType: orderbook.Limit, OrderSide: orderbook.Buy,
					ProductID: orderbook.BtcUsd, CommandModel: eventsource.CommandModel{ID: "AB-CD"}}}},
			context.Background(), "AB-CD", false},

		{"should return error when the repository returns so",
			fields{&mockIDGenerator{id: "AB-CD"}, &mockRepository{wantErr: true, err: errors.New("error")}}, context.Background(), "AB-CD", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.fields.idGenerator, tt.fields.repository)
			got, err := s.CreateOrder(tt.ctx, 1.0, 1.0, orderbook.Limit, orderbook.Buy, orderbook.BtcUsd)

			if tt.wantErr && err != nil {
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("service.CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("service.CreateOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_GetOrder(t *testing.T) {

	type fields struct {
		idGenerator Generator
		repository  Repository
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    orderbook.Order
		wantErr bool
	}{
		{"should get order by id",
			fields{idGenerator: nil, repository: &mockRepository{wantErr: false, err: nil, aggregate: &testOrder,
				command: orderbook.CreateOrder{Size: 1.0, Price: 1.0, OrderType: orderbook.Limit, OrderSide: orderbook.Buy,
					ProductID: orderbook.BtcUsd, CommandModel: eventsource.CommandModel{ID: "AB-CD"}}}},
			args{context.Background(), "AB-CD"}, testOrder, false},
		{"should return error when the repository returns so",
			fields{idGenerator: nil, repository: &mockRepository{wantErr: true, err: errors.New(""), aggregate: &testOrder}},
			args{context.Background(), "AB-CD"}, testOrder, true},
		{"should return error when the repository returns not a order",
			fields{idGenerator: nil, repository: &mockRepository{wantErr: false, err: errors.New(""), aggregate: &testAggregate}},
			args{context.Background(), "AB-CD"}, testOrder, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				idGenerator: tt.fields.idGenerator,
				repository:  tt.fields.repository,
			}
			got, err := s.GetOrder(tt.args.ctx, tt.args.id)

			if tt.wantErr && err != nil {
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_CancelOrder(t *testing.T) {
	type fields struct {
		idGenerator Generator
		repository  Repository
	}
	tests := []struct {
		name    string
		fields  fields
		ctx     context.Context
		want    string
		wantErr bool
	}{
		{"should apply CancelOrder command to repository",
			fields{nil, &mockRepository{wantErr: false, err: nil,
				command: orderbook.CreateOrder{Size: 1.0, Price: 1.0, OrderType: orderbook.Limit, OrderSide: orderbook.Buy,
					ProductID: orderbook.BtcUsd, CommandModel: eventsource.CommandModel{ID: "AB-CD"}}}},
			context.Background(), "AB-CD", false},

		{"should return error when the repository returns so",
			fields{nil, &mockRepository{wantErr: true, err: errors.New("error")}}, context.Background(), "AB-CD", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.fields.idGenerator, tt.fields.repository)
			err := s.CancelOrder(tt.ctx, "AB-CD")

			if tt.wantErr && err != nil {
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("service.CancelOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_service_AcceptOrder(t *testing.T) {
	type fields struct {
		idGenerator Generator
		repository  Repository
	}
	tests := []struct {
		name    string
		fields  fields
		ctx     context.Context
		want    string
		wantErr bool
	}{
		{"should apply AcceptOrder command to repository",
			fields{nil, &mockRepository{wantErr: false, err: nil,
				command: orderbook.AcceptOrder{CommandModel: eventsource.CommandModel{ID: "AB-CD"}}}},
			context.Background(), "AB-CD", false},

		{"should return error when the repository returns so",
			fields{nil, &mockRepository{wantErr: true, err: errors.New("error")}}, context.Background(), "AB-CD", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.fields.idGenerator, tt.fields.repository)
			err := s.AcceptOrder(tt.ctx, "AB-CD")

			if tt.wantErr && err != nil {
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("service.AcceptOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_service_PublishOrder(t *testing.T) {
	type fields struct {
		idGenerator Generator
		repository  Repository
	}
	tests := []struct {
		name    string
		fields  fields
		ctx     context.Context
		want    string
		wantErr bool
	}{
		{"should apply PublishOrder command to repository",
			fields{nil, &mockRepository{wantErr: false, err: nil,
				command: orderbook.PublishOrder{CommandModel: eventsource.CommandModel{ID: "AB-CD"}}}},
			context.Background(), "AB-CD", false},

		{"should return error when the repository returns so",
			fields{nil, &mockRepository{wantErr: true, err: errors.New("error")}}, context.Background(), "AB-CD", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.fields.idGenerator, tt.fields.repository)
			err := s.PublishOrder(tt.ctx, "AB-CD")

			if tt.wantErr && err != nil {
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("service.PublishOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_service_MatchOrder(t *testing.T) {
	type fields struct {
		idGenerator Generator
		repository  Repository
	}
	tests := []struct {
		name    string
		fields  fields
		ctx     context.Context
		want    string
		wantErr bool
	}{
		{"should apply MatchOrder command to repository",
			fields{nil, &mockRepository{wantErr: false, err: nil,
				command: orderbook.MatchOrder{CommandModel: eventsource.CommandModel{ID: "AB-CD"}}}},
			context.Background(), "AB-CD", false},

		{"should return error when the repository returns so",
			fields{nil, &mockRepository{wantErr: true, err: errors.New("error")}}, context.Background(), "AB-CD", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.fields.idGenerator, tt.fields.repository)
			err := s.MatchOrder(tt.ctx, "AB-CD")

			if tt.wantErr && err != nil {
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("service.MatchOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_service_ConfirmOrder(t *testing.T) {
	type fields struct {
		idGenerator Generator
		repository  Repository
	}
	tests := []struct {
		name    string
		fields  fields
		ctx     context.Context
		want    string
		wantErr bool
	}{
		{"should apply ConfirmOrder command to repository",
			fields{nil, &mockRepository{wantErr: false, err: nil,
				command: orderbook.ConfirmOrder{CommandModel: eventsource.CommandModel{ID: "AB-CD"}}}},
			context.Background(), "AB-CD", false},

		{"should return error when the repository returns so",
			fields{nil, &mockRepository{wantErr: true, err: errors.New("error")}}, context.Background(), "AB-CD", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.fields.idGenerator, tt.fields.repository)
			err := s.ConfirmOrder(tt.ctx, "AB-CD")

			if tt.wantErr && err != nil {
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("service.ConfirmOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_service_ClearOrder(t *testing.T) {
	type fields struct {
		idGenerator Generator
		repository  Repository
	}
	tests := []struct {
		name    string
		fields  fields
		ctx     context.Context
		want    string
		wantErr bool
	}{
		{"should apply ClearOrder command to repository",
			fields{nil, &mockRepository{wantErr: false, err: nil,
				command: orderbook.ClearOrder{CommandModel: eventsource.CommandModel{ID: "AB-CD"}}}},
			context.Background(), "AB-CD", false},

		{"should return error when the repository returns so",
			fields{nil, &mockRepository{wantErr: true, err: errors.New("error")}}, context.Background(), "AB-CD", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.fields.idGenerator, tt.fields.repository)
			err := s.ClearOrder(tt.ctx, "AB-CD")

			if tt.wantErr && err != nil {
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("service.ClearOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_service_SettleOrder(t *testing.T) {
	type fields struct {
		idGenerator Generator
		repository  Repository
	}
	tests := []struct {
		name    string
		fields  fields
		ctx     context.Context
		want    string
		wantErr bool
	}{
		{"should apply SettleOrder command to repository",
			fields{nil, &mockRepository{wantErr: false, err: nil,
				command: orderbook.SettleOrder{CommandModel: eventsource.CommandModel{ID: "AB-CD"}}}},
			context.Background(), "AB-CD", false},

		{"should return error when the repository returns so",
			fields{nil, &mockRepository{wantErr: true, err: errors.New("error")}}, context.Background(), "AB-CD", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.fields.idGenerator, tt.fields.repository)
			err := s.SettleOrder(tt.ctx, "AB-CD")

			if tt.wantErr && err != nil {
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("service.SettleOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

var testOrder = orderbook.Order{}
var testAggregate = mockAggregate{}

type mockAggregate struct {
}

func (a *mockAggregate) On(event eventsource.Event) error {
	return nil
}

type mockIDGenerator struct {
	id string
}

func (i *mockIDGenerator) Generate() string {
	return i.id
}

type mockRepository struct {
	err       error
	wantErr   bool
	aggregate eventsource.Aggregate
	command   eventsource.Command
}

func (m *mockRepository) Apply(ctx context.Context, command eventsource.Command) (int, error) {
	if m.wantErr {
		return 0, m.err
	}
	if m.command.AggregateID() == command.AggregateID() {
		return 123, nil
	}
	return 0, errors.New("unknown command")
}

func (m *mockRepository) Load(ctx context.Context, aggregateID string) (eventsource.Aggregate, error) {
	if m.wantErr {
		return nil, m.err
	}
	return m.aggregate, nil
}
