#!/bin/bash

go build -o ./test/srv cmd/server/main.go

go test -v -count=1 -json ./... | jq -c 'select(.Action=="fail")'

