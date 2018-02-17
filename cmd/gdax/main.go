package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/LAtanassov/godax/pkg/gdax"
	"github.com/gorilla/websocket"
)

var feedWs = flag.String("FEED_WS", "wss://ws-feed.gdax.com", "GDAX websocket feed")

func main() {

	flag.Parse()

	c, _, err := websocket.DefaultDialer.Dial(*feedWs, nil)
	if err != nil {
		log.Println("could not connect to websocket.", err)
		return
	}
	defer c.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)

		c.WriteJSON(gdax.Subscribe{Type: "subscribe", ProductIds: []string{"ETH-USD"}})

		var s gdax.Subscription
		err = c.ReadJSON(&s)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("recv:", s)

		for {
			var o gdax.OrderEvent
			err = c.ReadJSON(&o)
			if err != nil {
				log.Println(err)
				return
			}

			log.Println("recv:", o)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case <-interrupt:
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("could not close websocket conneciton.", err)
				return
			}
			<-done
			return
		case <-done:
			return
		}
	}
}
