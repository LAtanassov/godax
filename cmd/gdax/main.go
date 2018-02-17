package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

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
		for {
			_, m, err := c.ReadMessage()
			if err != nil {
				log.Println("could not read message from websocket", err)
				return
			}
			log.Println("recv:", m) // debug
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
