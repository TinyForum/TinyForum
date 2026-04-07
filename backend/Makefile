.PHONY: run build tidy

run:
	go run ./cmd/server/main.go

build:
	go build -o bin/server ./cmd/server/main.go

tidy:
	go mod tidy

wire:
	wire ./internal/wire/...
