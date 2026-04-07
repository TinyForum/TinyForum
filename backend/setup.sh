#!/bin/bash
# Run this script once to download Go dependencies
set -e
echo "📦 Downloading Go dependencies..."
go mod tidy
echo "✅ Dependencies downloaded! Now run: go run ./cmd/server/main.go"
