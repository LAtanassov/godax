package orderbook

import (
	"context"
	"testing"
	"time"

	"github.com/altairsix/eventsource"
)

func TestOrder_On(t *testing.T) {
	type fields struct {
		Size      float32
		Price     float32
		OrderType OrderType
		OrderSide OrderSide
		ProductID ProductID
	}
	type args struct {
		event eventsource.Event
	}
	tests := []struct {
		name      string
		order     Order
		event     eventsource.Event
		wantState string
		wantErr   bool
	}{
		{"should set StateCreated", Order{}, &OrderCreated{
			Model: eventsource.Model{ID: "", Version: 1, At: time.Now()},
		}, StateCreated, false},
		{"should set StateAccepted", Order{}, &OrderAccepted{
			Model: eventsource.Model{ID: "", Version: 1, At: time.Now()},
		}, StateAccepted, false},
		{"should set StateCanceled", Order{}, &OrderCanceled{
			Model: eventsource.Model{ID: "", Version: 1, At: time.Now()},
		}, StateCanceled, false},
		{"should set StatePublished", Order{}, &OrderPublished{
			Model: eventsource.Model{ID: "", Version: 1, At: time.Now()},
		}, StatePublished, false},
		{"should set StateMatched", Order{}, &OrderMatched{
			Model: eventsource.Model{ID: "", Version: 1, At: time.Now()},
		}, StateMatched, false},
		{"should set StateConfirmed", Order{}, &OrderConfirmed{
			Model: eventsource.Model{ID: "", Version: 1, At: time.Now()},
		}, StateConfirmed, false},
		{"should set StateCleared", Order{}, &OrderCleared{
			Model: eventsource.Model{ID: "", Version: 1, At: time.Now()},
		}, StateCleared, false},
		{"should set StateSettled", Order{}, &OrderSettled{
			Model: eventsource.Model{ID: "", Version: 1, At: time.Now()},
		}, StateSettled, false},
		{"should return erro if unknown model", Order{}, &UnknownModel{
			Model: eventsource.Model{ID: "", Version: 1, At: time.Now()},
		}, StateSettled, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.order.On(tt.event)

			if err != nil && tt.wantErr {
				return
			}

			if err != nil {
				t.Errorf("Order.On() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.order.state != tt.wantState {
				t.Errorf("Order.On() state = %v, wantState %v", tt.order.state, tt.wantState)
			}
		})
	}
}

type UnknownModel struct {
	eventsource.Model
}

func TestOrder_Apply(t *testing.T) {
	type args struct {
		ctx     context.Context
		command eventsource.Command
	}
	tests := []struct {
		name    string
		order   Order
		args    args
		want    []eventsource.Event
		wantErr bool
		err     error
	}{
		{"should return OrderCreated Event for CreateOrder command", Order{version: 0},
			args{context.Background(), &CreateOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{OrderCreated{
				Model: eventsource.Model{ID: "1", Version: 1, At: time.Now()},
			}}, false, nil},

		{"should return OrderAccepted Event for AcceptOrder command", Order{version: 0, state: StateCreated},
			args{context.Background(), &AcceptOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{OrderAccepted{
				Model: eventsource.Model{ID: "1", Version: 1, At: time.Now()},
			}}, false, nil},
		{"should return ErrOrderNotCreated for AcceptOrder command with an Order with not StateCreated", Order{version: 0, state: StateAccepted},
			args{context.Background(), &AcceptOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{}, true, ErrOrderNotCreated},

		{"should return OrderCanceled Event for CancelOrder command", Order{version: 0, state: StateCreated},
			args{context.Background(), &CancelOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{OrderCanceled{
				Model: eventsource.Model{ID: "1", Version: 1, At: time.Now()},
			}}, false, nil},
		{"should return ErrOrderNotCreated for CancelOrder command with an Order with not StateCreated", Order{version: 0, state: StateAccepted},
			args{context.Background(), &CancelOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{}, true, ErrOrderNotCreated},

		{"should return OrderPublished Event for PublishOrder command", Order{version: 0, state: StateAccepted},
			args{context.Background(), &PublishOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{OrderPublished{
				Model: eventsource.Model{ID: "1", Version: 1, At: time.Now()},
			}}, false, nil},
		{"should return ErrOrderNotAccepted for PublishOrder command with an Order with not StateAccepted", Order{version: 0, state: StateCreated},
			args{context.Background(), &PublishOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{}, true, ErrOrderNotAccepted},

		{"should return OrderMatched Event for MatchOrder command", Order{version: 0, state: StatePublished},
			args{context.Background(), &MatchOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{OrderMatched{
				Model: eventsource.Model{ID: "1", Version: 1, At: time.Now()},
			}}, false, nil},
		{"should return ErrOrderNotPublished for MatchOrder command with an Order with not StatePublished", Order{version: 0, state: StateAccepted},
			args{context.Background(), &MatchOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{}, true, ErrOrderNotPublished},

		{"should return OrderConfirmed Event for ConfirmOrder command", Order{version: 0, state: StateMatched},
			args{context.Background(), &ConfirmOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{OrderConfirmed{
				Model: eventsource.Model{ID: "1", Version: 1, At: time.Now()},
			}}, false, nil},
		{"should return ErrOrderNotMatched for ConfirmOrder command with an Order with not StateMatched", Order{version: 0, state: StateCreated},
			args{context.Background(), &ConfirmOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{}, true, ErrOrderNotMatched},

		{"should return OrderCleared Event for ClearOrder command", Order{version: 0, state: StateConfirmed},
			args{context.Background(), &ClearOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{OrderCleared{
				Model: eventsource.Model{ID: "1", Version: 1, At: time.Now()},
			}}, false, nil},
		{"should return ErrOrderNotConfirmed for ClearOrder command with an Order with not StateConfirmed", Order{version: 0, state: StateCreated},
			args{context.Background(), &ClearOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{}, true, ErrOrderNotConfirmed},

		{"should return OrderSettled Event for SettleOrder command", Order{version: 0, state: StateCleared},
			args{context.Background(), &SettleOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{OrderSettled{
				Model: eventsource.Model{ID: "1", Version: 1, At: time.Now()},
			}}, false, nil},
		{"should return ErrOrderNotCleared for SettleOrder command with an Order with not StateCleared", Order{version: 0, state: StateCreated},
			args{context.Background(), &SettleOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{}, true, ErrOrderNotCleared},

		{"should return ErrUnknownCommand", Order{},
			args{context.Background(), eventsource.CommandModel{ID: "1"}},
			[]eventsource.Event{}, true, ErrUnknownCommand},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.order.Apply(tt.args.ctx, tt.args.command)
			if tt.wantErr && err == tt.err {
				return
			}

			if err != nil {
				t.Errorf("Order.Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) || got[0].EventVersion() != tt.want[0].EventVersion() {
				t.Errorf("Order.Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}
