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

## TODO

* SECURITY: validation