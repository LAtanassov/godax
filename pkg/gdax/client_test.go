package gdax

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
)

func TestClientLifecycle(t *testing.T) {
	s := newTestWebsocketServer()
	defer s.Close()

	u, err := url.Parse(s.URL)
	if err != nil {
		t.Fail()
	}
	u.Scheme = "ws"

	c := NewClient(websocket.DefaultDialer)
	err = c.Connect(u)
	if err != nil {
		t.Fail()
	}
	_, err = c.Subscribe([]ProductID{EthUsd})
	if err != nil {
		t.Fail()
	}
	err = c.Disconnect()
	if err != nil {
		t.Fail()
	}
}

func TestSubscribeBeforeConnect(t *testing.T) {
	cli := NewClient(websocket.DefaultDialer)
	_, err := cli.Subscribe([]ProductID{EthUsd})
	if err != ErrNotConnected {
		t.Fail()
	}
}

func TestDisconnectBeforeConnect(t *testing.T) {
	c := NewClient(websocket.DefaultDialer)
	err := c.Disconnect()
	if err != ErrNotConnected {
		t.Fail()
	}
}

func TestSnapshot(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b := BookSnapshot{
			Sequence: "1",
			Bids:     []Order{},
			Asks:     []Order{},
		}

		m, err := json.Marshal(b)
		if err != nil {
			t.Fail()
		}
		w.Write(m)
	}))
	defer srv.Close()

	u, _ := url.Parse(srv.URL)

	_, err := Snapshot(u)
	if err != nil {
		t.Fail()
	}
}

func newTestWebsocketServer() (s *httptest.Server) {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, _ := u.Upgrade(w, r, nil)
		go func() {
			defer func() {
				conn.Close()
			}()
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					fmt.Println(err)
					return
				}
				conn.WriteMessage(websocket.TextMessage, message)
			}
		}()
	}))
}
