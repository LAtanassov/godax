package messaging

import (
	"github.com/streadway/amqp"
)

// Consumer provides an interface to consume messages
type Consumer interface {
	Consume() (<-chan amqp.Delivery, error)
	Close()
}

type consumer struct {
	url   string
	queue string

	conn *amqp.Connection
	ch   *amqp.Channel
}

// NewConsumer returns a simple consumer
func NewConsumer(url, queue string) Consumer {
	return &consumer{url, queue, nil, nil}
}

func (c *consumer) Consume() (<-chan amqp.Delivery, error) {
	c.Close()
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (c *consumer) Close() {
	if c.ch != nil {
		c.ch.Close()
	}

	if c.conn != nil {
		c.conn.Close()
	}
}
