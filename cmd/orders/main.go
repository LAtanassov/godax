package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/LAtanassov/godax/pkg/accessor"
	"github.com/LAtanassov/godax/pkg/orders"
	"github.com/altairsix/eventsource"
	kitlog "github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	inmem     = "inmem"
	mysql     = "mysql"
	tableName = "orders"
)

func main() {

	var (
		envHTTPAddr  = envString("HTTP_ADDR", ":8080")
		envSQLDriver = envString("SQL_DRIVER", inmem)
		envSQLUser   = envString("SQL_USER", "")
		envSQLPwd    = envString("SQL_PASSWORD", "")
		envSQLHost   = envString("SQL_HOST", "")
		envSQLDbName = envString("SQL_DB_NAME", "godax")
		envTabName   = envString("SQL_DB_NAME", tableName)

		httpAddr = flag.String("http.addr", envHTTPAddr, "HTTP listen address")

		dbCon = orders.DatabaseConnection{
			Driver:   *flag.String("sql.driver", envSQLDriver, "SQL driver"),
			Username: *flag.String("sql.user", envSQLUser, "SQL user"),
			Password: *flag.String("sql.pwd", envSQLPwd, "SQL password"),
			Host:     *flag.String("sql.host", envSQLHost, "SQL host"),
			Database: *flag.String("sql.dbname", envSQLDbName, "SQL database name"),
		}
		tableName = *flag.String("sql.tabname", envTabName, "Table name")
	)

	flag.Parse()

	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)

	observers := []func(event eventsource.Event){}
	repo, err := createRepository(dbCon, tableName, observers...)
	if err != nil {
		log.Fatal("terminated: database connection failed ", err)
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

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	log.Fatal("terminated", <-errs)
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

func createRepository(dbCon orders.DatabaseConnection, tableName string, observers ...func(event eventsource.Event)) (orders.Repository, error) {

	switch dbCon.Driver {
	case inmem:
		return orders.NewInMemRepository(observers...), nil
	case mysql:
		accessor, err := accessor.New(dbCon.Driver, fmt.Sprintf("%s:%s@tcp(%s)/%s", dbCon.Username, dbCon.Password, dbCon.Host, dbCon.Database), tableName)
		if err != nil {
			return nil, err
		}

		store, err := orders.NewMysqlStore(tableName, accessor)
		if err != nil {
			return nil, err
		}

		return orders.NewRepository(store, observers...), nil
	default:
		return nil, errors.New("unsupported driver")
	}
}

func envString(env, fallback string) string {
	e, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	return e
}
