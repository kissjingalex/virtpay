include .env

openapi_http:
	@./scripts/openapi-http.sh wallet

proto:
	@./scripts/proto.sh wallet
