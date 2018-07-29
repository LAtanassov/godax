package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LAtanassov/godax/pkg/orders"
	kitlog "github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func main() {

	var ( // configuration
		envHTTPAddr  = envString("HTTP_ADDR", ":8080")
		envDbDriver  = envString("DB_DRIVER", "inmem")
		envDbURL     = envString("DB_URL", "")
		envTableName = envString("DB_TABLE_NAME", "orders")

		httpAddr  = *flag.String("http.addr", envHTTPAddr, "HTTP listen address")
		dbDriver  = *flag.String("db.driver", envDbDriver, "database driver")
		dbURL     = *flag.String("db.url", envDbURL, "database connection url")
		tableName = *flag.String("sql.tabname", envTableName, "Table name")
	)
	flag.Parse()

	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)

	repo, err := orders.NewRepository(dbDriver, dbURL, tableName)
	if err != nil {
		log.Fatal("terminated", err)
	}

	idg := orders.NewIDGenerator()

	fieldKeys := []string{"method"}
	o := orders.NewService(idg, repo)
	o = orders.NewLoggingMiddleware(kitlog.With(logger, "component", "orders"))(o)
	o = orders.NewInstrumentingMiddleware(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "orders_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "orders_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys))(o)

	httpLogger := kitlog.With(logger, "component", "http")

	mux := http.NewServeMux()
	mux.Handle("/godax/v1/", orders.MakeHandler(o, httpLogger))

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
	err = <-errs
	log.Println("recv. signal", err)
	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Log(err)
	}
	log.Println("terminate in 5 sec")
	time.Sleep(time.Duration(5) * time.Second)
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

func readinessHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
}

func livenessHandler() http.Handler {
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
