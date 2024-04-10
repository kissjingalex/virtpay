#!/bin/bash
set -e

readonly service="$1"
#readonly output_dir="$2"
#readonly package="$3"

oapi-codegen -generate types -o "internal/$service/ports/openapi_types.gen.go" -package "ports" "api/openapi/$service.yml"
oapi-codegen -generate chi-server -o "internal/$service/ports/openapi_api.gen.go" -package "ports" "api/openapi/$service.yml"
oapi-codegen -generate types -o "internal/common/client/$service/openapi_types.gen.go" -package "$service" "api/openapi/$service.yml"
oapi-codegen -generate client -o "internal/common/client/$service/openapi_client_gen.go" -package "$service" "api/openapi/$service.yml"
