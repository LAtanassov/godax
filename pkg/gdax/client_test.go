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

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, err := u.Upgrade(w, r, nil)
		if err != nil {
			t.Fail()
		}

		go func() {
			defer func() {
				conn.Close()
			}()
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					fmt.Println(err)
				}
				conn.WriteMessage(websocket.TextMessage, message)
			}
		}()

	}))
	defer srv.Close()

	u, _ := url.Parse(srv.URL)
	u.Scheme = "ws"

	cli := NewClient(websocket.DefaultDialer)
	err := cli.Connect(u)
	if err != nil {
		t.Fail()
	}
	_, err = cli.Subscribe([]ProductID{EthUsd})
	if err != nil {
		t.Fail()
	}
	err = cli.Disconnect()
	if err != nil {
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
