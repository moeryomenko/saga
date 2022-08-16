#!/usr/bin/env bash

case $1 in
	"order")
		DB_HOST=localhost DB_NAME=orders DB_USER=test DB_PASSWORD=pass GRACE_PERIOD=3s go run cmd/order/main.go
		;;
esac
