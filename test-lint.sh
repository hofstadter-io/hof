#!/bin/bash
set -euo pipefail

echo "ci-tests"
echo "===================================="
GO_FILES=$(find . -iname '*.go' -type f | grep -v /vendor/) # All the .go files, excluding vendor/

echo ""
echo "CHECKING gofmt - fail for not being a gopher"
echo "===================================="
test -z $(gofmt -s -l $GO_FILES)         # Fail if a .go file hasn't been formatted with gofmt

echo ""
echo "CHECKING go test with race conditions"
echo "===================================="
go test -v -race ./...                   # Run all the tests with the race detector enabled

echo ""
echo "CHECKING go vet - offical static analyzer"
echo "===================================="
go vet ./...                             # go vet is the official Go static analyzer

echo ""
echo "CHECKING megacheck - vet on steroids"
echo "===================================="
megacheck ./...                          # "go vet on steroids" + linter

echo ""
echo "CHECKING gocyclo - forbid code with huge functions"
echo "===================================="
gocyclo -over 19 $GO_FILES               # forbid code with huge functions

echo ""
echo "CHECKING golint"
echo "===================================="
golint -set_exit_status $(go list ./...) # one last linter


