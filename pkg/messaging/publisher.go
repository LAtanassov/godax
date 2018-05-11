package messaging

import "github.com/altairsix/eventsource"

// Publisher publishes a event on a work queue
type Publisher interface {
	Publish(event eventsource.Event) error
	Close() error
}

type publisher struct {
}

// NewPublisher returns a publisher
func NewPublisher() Publisher {
	return &publisher{}
}

func (p *publisher) Publish(event eventsource.Event) error {
	return nil
}

func (p *publisher) Close() error {
	return nil
}
