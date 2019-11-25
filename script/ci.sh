#!/bin/sh
set -e

# lint
echo "Checking lint"
golangci-lint run
echo "Lint success!"

# test
echo "Running tests"
go test `go list ./... | grep -v 'docs'` -mod=vendor -race -cover -tags=integration -coverprofile=coverage.txt -covermode=atomic
echo "Test execution completed!"