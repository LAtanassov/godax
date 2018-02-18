package gdax

import (
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
				messageType, message, err := conn.ReadMessage()
				if err != nil {
					fmt.Println(err)
					if messageType != websocket.CloseNormalClosure {
						t.Fail()
					}
				}
				conn.WriteMessage(websocket.TextMessage, message)
			}
		}()

	}))

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
