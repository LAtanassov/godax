// riskmonitor listens to incoming OrderCreated and either
// - accepts them automatically by some predefined rules or manually by risk analyst or
// - reject them manually by a risk analysts.
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	var ( // configuration
		envHTTPAddr = envString("HTTP_ADDR", ":8080")
		httpAddr    = *flag.String("http.addr", envHTTPAddr, "HTTP listen address")
	)
	flag.Parse()

	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)

	mux := http.NewServeMux()

	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/_status/liveness", livenessHandler())
	http.Handle("/_status/readiness", readinessHandler())
	srv := http.Server{
		Addr:    httpAddr,
		Handler: nil,
	}

	errs := make(chan error, 2)

	go func() {
		logger.Log("transport", "http", "address", httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(httpAddr, nil)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGABRT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	// gracefully shutdown
	err := <-errs
	logger.Log("shutdown", "http_server", "signal_recv", err)
	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Log("shutdown", "http_server", "err", err)
	}
	logger.Log("shutdown", "cooldown_5_sec")
	time.Sleep(time.Duration(5) * time.Second)
	logger.Log("shutdown", "byebye")
	os.Exit(1)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func livenessHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
}

func readinessHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
}

func envString(env, fallback string) string {
	e, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	return e
}
