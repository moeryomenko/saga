#!/usr/bin/env bash

set -a

export DB_HOST=localhost
export DB_USER=test
export DB_PASSWORD=pass
export GRACE_PERIOD=3s

SVC=$1

air --build.cmd "go build -o bin/api cmd/${SVC}/main.go" \
	--build.bin "./bin/api" \
	--build.exclude_dir "api,hack,migrations"

set +a
