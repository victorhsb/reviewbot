tools:
	@which migrate > /dev/null || (echo "migrate is not installed, installing..." && go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest)
	@which air > /dev/null || (echo "air is not installed, installing..." && go install github.com/cosmtrek/air@latest)
	@which golangci-lint > /dev/null || (echo "golangci-lint is not installed, installing..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@which mockery > /dev/null || (echo "mockery is not installed, installing..." && go install github.com/vektra/mockery/v2@v2.42.3)
	@which sqlc > /dev/null || (echo "sqlc is not installed, installing..." && go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest)

test:
	@go test -v ./...

gen: tools
	@sqlc generate
	@mockery

deps:
	@docker compose up -d db

run: tools
	@docker compose up -d db
	@REVIEWBOT_DATABASE_URL="postgresql://user:password@localhost:5432/product-reviews?sslmode=disable" REVIEWBOT_PORT=":8080" air

