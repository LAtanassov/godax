#!/bin/sh -e

GO_FILES=$(find . -iname '*.go' -type f | grep -v /vendor/) # All the .go files, excluding vendor/
test -z $(gofmt -s -l $GO_FILES)         # Fail if a .go file hasn't been formatted with gofmt
go test -v -race ./...                   # Run all the tests with the race detector enabled
go vet ./...                             # go vet is the official Go static analyzer
megacheck ./...                          # "go vet on steroids" + linter
gocyclo -over 19 $GO_FILES               # forbid code with huge functions
golint -set_exit_status $(go list ./...) # one last linter