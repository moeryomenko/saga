#!/usr/bin/env bash

set -a

export DB_HOST=localhost
export DB_NAME=orders
export DB_USER=test
export DB_PASSWORD=pass
export GRACE_PERIOD=3s

go run cmd/$1/main.go

set +a
