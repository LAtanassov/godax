# go-gdax
go gdax - crypto currency trading platform - pet project

# Build & Deployment

[![Build Status](https://travis-ci.org/LAtanassov/godax.svg?branch=master)](https://travis-ci.org/LAtanassov/godax)

## Orders

```sh
# builds
$> cd ./cmd/orders
$> CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' .

$> docker build -t latanassov/orders .
$> docker tag latanassov/orders:latest latanassov/orders
$> docker push latanassov/orders

$> kubectl apply -f mysql-deployment.yaml
$> kubectl apply -f orders-deployment.yaml
```

## Risk Monitor

```sh
$> docker run -it -p 5672:5672 --hostname test-rabbitmq rabbitmq:3.7.4

```

## TODO

* SECURITY: validation

* FEATURE: not fault-tolerant yet (observers might miss a message)
* FEATURE: publisher/consumer - exponetial backoff reconnect, reuse channel

## Useful Links

value vs. pointer: https://www.ardanlabs.com/blog/2014/12/using-pointers-in-go.html