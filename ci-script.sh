#!/bin/bash -e
# used for development

GO_FILES=$(find . -iname '*.go' -type f | grep -v /vendor/) # All the .go files, excluding vendor/
# go get golang.org/x/vgo

vgo test -v -race ./...                  # Run all the tests with the race detector enabled
vgo vet ./...                            # go vet is the official Go static analyzer
test -z $(vgo fmt ./...)

#
go get github.com/golang/lint/golint     # Linter
go get honnef.co/go/tools/cmd/megacheck  # Badass static analyzer/linter
go get github.com/fzipp/gocyclo

vgo mod -vendor
megacheck ./...                          # "go vet on steroids" + linter
gocyclo -over 19 $GO_FILES               # forbid code with huge functions
golint -set_exit_status $(go list ./...) # one last linter