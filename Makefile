.PHONY: vendor
vendor:
	@go mod vendor

.PHONY: run/api
run/api:
	@go run ./cmd/api -db-dsn=${DB_DSN} -jwt-secret=${JWT_SECRET}