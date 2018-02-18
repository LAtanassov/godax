package gdax

import (
	"net/url"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/websocket"
)

// The Subscribe type represents a request to which product ids you want to subscibe.
type Subscribe struct {
	Type       string      `json:"type"`
	ProductIds []ProductID `json:"product_ids"`
}

// The Subscription type represents a subscription.
type Subscription struct {
	Type string `json:"type"`
}

// Client interface to consume GDAX API.
type Client interface {
	Connect(u *url.URL) error
	Disconnect() error
	Subscribe(p []ProductID) (<-chan OrderEvent, error)

	WithLogger(l log.Logger)
}

type client struct {
	uri        *url.URL
	productIDs []ProductID

	dialer *websocket.Dialer
	logger log.Logger
	conn   *websocket.Conn
}

// NewClient return a new GDAX client for real time order events
func NewClient(d *websocket.Dialer) Client {
	return &client{
		dialer: d,
		logger: log.NewNopLogger(),
	}
}

func (c *client) WithLogger(l log.Logger) {
	c.logger = l
}

// Connect to GDAX websocket api and returns an error if the connection attempt failed.
func (c *client) Connect(u *url.URL) error {
	c.uri = u

	conn, _, err := c.dialer.Dial(c.uri.String(), nil)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

// Disconnect by sending a close message via websocket and afterwards closing the websocket connection.
func (c *client) Disconnect() error {
	defer c.conn.Close()

	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}
	return nil
}

// Subscribe sents a subscribe request and return order event as a channel.
// In case the subscribe request fail Subscribe will return an error.
// An unexpected message will not cancel the subscription or connection and will be logged.
func (c *client) Subscribe(p []ProductID) (<-chan OrderEvent, error) {
	c.productIDs = p
	oc := make(chan OrderEvent, 2048)

	c.conn.WriteJSON(Subscribe{Type: "subscribe", ProductIds: []ProductID{EthUsd}})

	var s Subscription
	err := c.conn.ReadJSON(&s)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			var o OrderEvent
			err := c.conn.ReadJSON(&o)
			if err != nil {
				c.logger.Log(err)
				return
			}
			oc <- o
		}
	}()
	return oc, nil
}
