#!/bin/bash -e

WKDIR=`pwd`

cd ${WKDIR}/cmd/gdaxcli
CGO_ENABLED=0 
GOOS=linux GOARCH=amd64 vgo build -ldflags "-s -w"
GOOS=windows GOARCH=amd64 vgo build -ldflags "-s -w"

cd ${WKDIR}/cmd/orders
CGO_ENABLED=0 
GOOS=linux GOARCH=amd64 vgo build -ldflags "-s -w"
GOOS=windows GOARCH=amd64 vgo build -ldflags "-s -w"