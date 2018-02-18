package main

import (
	"flag"
	"net/url"
	"os"
	"os/signal"

	"github.com/LAtanassov/godax/pkg/gdax"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/websocket"
)

var feedWs = flag.String("FEED_WS", "wss://ws-feed.gdax.com", "GDAX websocket feed")

func main() {

	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	done := make(chan bool)
	c := gdax.NewClient(websocket.DefaultDialer)
	c.WithLogger(logger)

	u, err := url.Parse(*feedWs)
	if err != nil {
		logger.Log("could not parse feed uri", err)
		done <- true
	}

	err = c.Connect(u)
	if err != nil {
		logger.Log("could not connect to gdax.", err)
		done <- true
	}

	oc, err := c.Subscribe([]gdax.ProductID{gdax.EthUsd})
	if err != nil {
		logger.Log("could not subscribe for Eth-Usd", err)
		done <- true
	}

	go func() {
		for {
			o := <-oc
			logger.Log("recv:", o.OrderType)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case <-interrupt:
			err := c.Disconnect()
			if err != nil {
				logger.Log("could not disconnect.", err)
			}
			<-done
			return
		case <-done:
			return
		}
	}
}
