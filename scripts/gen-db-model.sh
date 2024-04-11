#!/bin/bash
set -e

readonly tableName="$1"
readonly dbName="$2"
readonly service="$3"

env $(cat ".env" | grep -Ev '^#' | xargs) go run ./tools/gen-db-model/main.go -tb=$tableName -db=$dbName -o=./internal/$service/adapter/pg/model