#!/bin/bash

# Скрип для подсчета coverage

# запишет данные о покрытии в merge.out
go test -count=1 -coverpkg=./... -coverprofile=merge.out ./... 2>&1 > /dev/null

# посчитает и выдаст суммарное покрытие
go tool cover -func ./merge.out | tail -1 | tr -s '\t'

