#!/bin/bash

go build -o ./test/srv cmd/server/main.go
if [[ R_VAL -ne "0" ]] ; then
  echo "!!! TESTS FAILED !!!"
  exit 99
fi

#BIN_TEST=TRUE go test -v -race -count=1 -json ./... | jq -c 'select(.Action=="fail")'
BIN_TEST=TRUE go test -v -race -count=1 ./...
