package gdax

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/websocket"
)

var (
	// ErrNotConnected will be returned when a method was called that required a the client to be connected first
	ErrNotConnected = errors.New("need to connect first")
)

// Client interface to consume GDAX API.
type Client interface {
	Connect(u *url.URL) error
	Disconnect() error
	Subscribe(p []ProductID) (<-chan OrderEvent, error)

	WithLogger(l log.Logger)
}

// NewClient return a new GDAX client to subscribe for real time order events
func NewClient(d *websocket.Dialer) Client {
	return &client{
		dialer: d,
		logger: log.NewNopLogger(),
	}
}

// WithLogger assigns a Logger
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
	if c.conn == nil {
		return ErrNotConnected
	}

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
	if c.conn == nil {
		return nil, ErrNotConnected
	}

	c.productIDs = p
	oc := make(chan OrderEvent, 2048)
	c.conn.WriteJSON(subscribe{Type: "subscribe", ProductIds: []ProductID{EthUsd}})

	var s subscription
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
				c.logger.Log("closing subscription")
				close(oc)
				return
			}
			oc <- o
		}
	}()
	return oc, nil
}

// Snapshot returns a BookSnapshot.
func Snapshot(u *url.URL) (*BookSnapshot, error) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s BookSnapshot
	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return nil, err
	}

	return &s, nil
}

// subscribe represents a request with the ProductIds that you want to subscibe to.
type subscribe struct {
	Type       string      `json:"type"`
	ProductIds []ProductID `json:"product_ids"`
}

// subscription represents a subscription.
type subscription struct {
	Type string `json:"type"`
}

type client struct {
	uri        *url.URL
	productIDs []ProductID

	dialer *websocket.Dialer
	logger log.Logger
	conn   *websocket.Conn
}
