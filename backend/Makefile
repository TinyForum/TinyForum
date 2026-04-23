.PHONY: run build tidy

run:
	@echo "Running server..."
	@go run ./cmd/server/main.go

build:
	go build -o bin/server ./cmd/server/main.go

tidy:
	go mod tidy

wire:
	wire ./internal/wire/...

docs-api:
	@swag init -g ./cmd/server/main.go --output ./docs


