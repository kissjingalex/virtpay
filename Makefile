include .env

openapi_http:
	@./scripts/openapi-http.sh wallet

proto:
	@./scripts/proto.sh wallet

test:
	@./scripts/test.sh wallet .test.env

genModel:
	@./scripts/gen-db-model.sh wallet_info db_wallet wallet