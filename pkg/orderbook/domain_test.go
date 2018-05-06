package orderbook

import (
	"context"
	"testing"
	"time"

	"github.com/altairsix/eventsource"
)

func TestOrder_On(t *testing.T) {
	tests := []struct {
		name      string
		order     Order
		event     eventsource.Event
		wantState string
		wantErr   bool
	}{
		{"should set stateCreated", Order{}, &OrderCreated{
			Model: eventsource.Model{ID: "", Version: 1, At: time.Now()},
		}, stateCreated, false},
		{"should set stateAccepted", Order{}, &OrderAccepted{
			Model: eventsource.Model{ID: "", Version: 1, At: time.Now()},
		}, stateAccepted, false},
		{"should set stateCanceled", Order{}, &OrderCanceled{
			Model: eventsource.Model{ID: "", Version: 1, At: time.Now()},
		}, stateCanceled, false},
		{"should set statePublished", Order{}, &OrderPublished{
			Model: eventsource.Model{ID: "", Version: 1, At: time.Now()},
		}, statePublished, false},
		{"should set stateMatched", Order{}, &OrderMatched{
			Model: eventsource.Model{ID: "", Version: 1, At: time.Now()},
		}, stateMatched, false},
		{"should set stateConfirmed", Order{}, &OrderConfirmed{
			Model: eventsource.Model{ID: "", Version: 1, At: time.Now()},
		}, stateConfirmed, false},
		{"should set stateCleared", Order{}, &OrderCleared{
			Model: eventsource.Model{ID: "", Version: 1, At: time.Now()},
		}, stateCleared, false},
		{"should set stateSettled", Order{}, &OrderSettled{
			Model: eventsource.Model{ID: "", Version: 1, At: time.Now()},
		}, stateSettled, false},
		{"should return erro if unknown model", Order{}, &UnknownModel{
			Model: eventsource.Model{ID: "", Version: 1, At: time.Now()},
		}, stateSettled, true},
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

		{"should return OrderAccepted Event for AcceptOrder command", Order{version: 0, state: stateCreated},
			args{context.Background(), &AcceptOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{OrderAccepted{
				Model: eventsource.Model{ID: "1", Version: 1, At: time.Now()},
			}}, false, nil},
		{"should return ErrInvalidStateTransition for AcceptOrder command with an Order with not stateCreated", Order{version: 0, state: stateAccepted},
			args{context.Background(), &AcceptOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{}, true, ErrInvalidStateTransition},

		{"should return OrderCanceled Event for CancelOrder command", Order{version: 0, state: stateCreated},
			args{context.Background(), &CancelOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{OrderCanceled{
				Model: eventsource.Model{ID: "1", Version: 1, At: time.Now()},
			}}, false, nil},
		{"should return ErrInvalidStateTransition for CancelOrder command with an Order with not stateCreated", Order{version: 0, state: stateAccepted},
			args{context.Background(), &CancelOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{}, true, ErrInvalidStateTransition},

		{"should return OrderPublished Event for PublishOrder command", Order{version: 0, state: stateAccepted},
			args{context.Background(), &PublishOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{OrderPublished{
				Model: eventsource.Model{ID: "1", Version: 1, At: time.Now()},
			}}, false, nil},
		{"should return ErrInvalidStateTransition for PublishOrder command with an Order with not stateAccepted", Order{version: 0, state: stateCreated},
			args{context.Background(), &PublishOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{}, true, ErrInvalidStateTransition},

		{"should return OrderMatched Event for MatchOrder command", Order{version: 0, state: statePublished},
			args{context.Background(), &MatchOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{OrderMatched{
				Model: eventsource.Model{ID: "1", Version: 1, At: time.Now()},
			}}, false, nil},
		{"should return ErrInvalidStateTransition for MatchOrder command with an Order with not statePublished", Order{version: 0, state: stateAccepted},
			args{context.Background(), &MatchOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{}, true, ErrInvalidStateTransition},

		{"should return OrderConfirmed Event for ConfirmOrder command", Order{version: 0, state: stateMatched},
			args{context.Background(), &ConfirmOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{OrderConfirmed{
				Model: eventsource.Model{ID: "1", Version: 1, At: time.Now()},
			}}, false, nil},
		{"should return ErrInvalidStateTransition for ConfirmOrder command with an Order with not stateMatched", Order{version: 0, state: stateCreated},
			args{context.Background(), &ConfirmOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{}, true, ErrInvalidStateTransition},

		{"should return OrderCleared Event for ClearOrder command", Order{version: 0, state: stateConfirmed},
			args{context.Background(), &ClearOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{OrderCleared{
				Model: eventsource.Model{ID: "1", Version: 1, At: time.Now()},
			}}, false, nil},
		{"should return ErrInvalidStateTransition for ClearOrder command with an Order with not stateConfirmed", Order{version: 0, state: stateCreated},
			args{context.Background(), &ClearOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{}, true, ErrInvalidStateTransition},

		{"should return OrderSettled Event for SettleOrder command", Order{version: 0, state: stateCleared},
			args{context.Background(), &SettleOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{OrderSettled{
				Model: eventsource.Model{ID: "1", Version: 1, At: time.Now()},
			}}, false, nil},
		{"should return ErrInvalidStateTransition for SettleOrder command with an Order with not stateCleared", Order{version: 0, state: stateCreated},
			args{context.Background(), &SettleOrder{
				CommandModel: eventsource.CommandModel{ID: "1"},
			}},
			[]eventsource.Event{}, true, ErrInvalidStateTransition},

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
