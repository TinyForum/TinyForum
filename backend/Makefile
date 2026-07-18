.PHONY: run build tidy

run:
	@echo "Running server..."
	@go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

tidy:
	go mod tidy

wire:
	wire ./internal/wire/...

docs-api:
	@swag init -g ./cmd/server --output ./docs


