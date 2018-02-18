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

var feedURI = flag.String("FEED_URI", "wss://ws-feed.gdax.com", "GDAX websocket feed")
var snapshotURI = flag.String("SNAPSHOT_URI", "https://api.gdax.com/products/ETH-USD/book?level=3", "GDAX book snapshot")

func main() {

	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	c := gdax.NewClient(websocket.DefaultDialer)
	c.WithLogger(logger)

	u, err := url.Parse(*feedURI)
	if err != nil {
		logger.Log("could not parse feed uri", err)
		return
	}

	err = c.Connect(u)
	if err != nil {
		logger.Log("could not connect to gdax.", err)
		return
	}
	defer c.Disconnect()

	oc, err := c.Subscribe([]gdax.ProductID{gdax.EthUsd})
	if err != nil {
		logger.Log("could not subscribe for Eth-Usd", err)
		return
	}

	s, err := url.Parse(*snapshotURI)
	if err != nil {
		logger.Log("could not parse feed uri", err)
		return
	}

	b, err := gdax.Snapshot(s)
	if err != nil {
		logger.Log("could not parse feed uri", err)
		return
	}
	logger.Log(b)

	go func() {
		for o := range oc {
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
			return
		}
	}
}
