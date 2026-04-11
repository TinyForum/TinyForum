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

docs-:
	@swag init -g ./cmd/server/main.go --output ./docs
	@go get -u github.com/swaggo/gin-swagger
	@go get -u github.com/swaggo/files
	@go get -u github.com/swaggo/swag
