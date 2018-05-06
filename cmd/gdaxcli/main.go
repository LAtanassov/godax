package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/LAtanassov/godax/pkg/gdax"
	kitlog "github.com/go-kit/kit/log"
	"github.com/gorilla/websocket"
)

var feedURI = flag.String("FEED_URI", "wss://ws-feed.gdax.com", "GDAX websocket feed")
var snapshotURI = flag.String("SNAPSHOT_URI", "https://api.gdax.com/products/ETH-USD/book?level=3", "GDAX book snapshot")

func main() {

	flag.Parse()

	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)

	u, err := url.Parse(*feedURI)
	if err != nil {
		logger.Log("could not parse feed uri", err)
		return
	}

	c := gdax.NewClient(websocket.DefaultDialer)
	c.WithLogger(logger)
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

	errs := make(chan error, 1)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	log.Fatal("terminated", <-errs)
}
