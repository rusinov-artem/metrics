#!/bin/bash

go build -o ./test/srv cmd/server/main.go

BIN_TEST=TRUE go test -v -race -count=1 -json ./... | jq -c 'select(.Action=="fail")'
