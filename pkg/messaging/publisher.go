package messaging

import (
	"io"
	"io/ioutil"

	kitlog "github.com/go-kit/kit/log"

	"github.com/streadway/amqp"
)

// Publisher publishes the content of a reader
type Publisher interface {
	Publish(r io.Reader) error

	Open(url, queue string) error
	Close() error
}

type publisher struct {
	url     string
	queue   string
	conn    *amqp.Connection
	ch      *amqp.Channel
	q       *amqp.Queue
	closeCh chan *amqp.Error

	logger kitlog.Logger
}

func (p *publisher) Publish(r io.Reader) error {
	if err := p.needReconnect(); err != nil {
		return err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return p.ch.Publish(
		"",
		p.q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        b,
		},
	)
}

func (p *publisher) Open(url, queue string) error {

	p.url = url
	p.queue = queue

	return p.needReconnect()
}

func (p *publisher) Close() error {
	p.logger.Log("close publisher")
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

func (p *publisher) needReconnect() error {
	select {
	case <-p.closeCh:
		return p.connect()
	default:
	}
	return nil
}

func (p *publisher) connect() error {
	p.logger.Log("connect publisher")
	conn, err := amqp.Dial(p.url)
	if err != nil {
		return err
	}
	p.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	p.ch = ch

	p.closeCh = make(chan *amqp.Error)
	conn.NotifyClose(p.closeCh)

	_, err = ch.QueueDeclare(
		p.queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// NewPublisher create a publisher and initialize it with a closed 'close' channel to trigger reconnect
func NewPublisher(opts ...Option) Publisher {

	ch := make(chan *amqp.Error)
	close(ch)

	p := &publisher{"", "", nil, nil, nil, ch, kitlog.NewNopLogger()}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Option provides functional configuration for a *Repository
type Option func(*publisher)

// WithLogger set a logger
func WithLogger(logger kitlog.Logger) Option {
	return func(p *publisher) {
		p.logger = logger
	}
}
