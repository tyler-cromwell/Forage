#!/bin/sh

source .env
go test -v -coverpkg=./... -coverprofile coverage.out -failfast -tags test ./...
go tool cover -html=coverage.out -o coverage.html
