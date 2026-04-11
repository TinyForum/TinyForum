.PHONY: help frontend backend docs install clean build dev all

init-config:
	@echo "🚀 Initializing project..."
	@bash ./start-dev.sh
# Default target
help:
	@echo "Available targets:"
	@echo "  make frontend  - Start frontend development server"
	@echo "  make backend   - Start backend development server"
	@echo "  make docs      - Start docsify documentation server"
	@echo "  make dev       - Run both frontend and backend concurrently"
	@echo "  make install   - Install all dependencies"
	@echo "  make build     - Build production versions"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make all       - Install and build everything"

# Frontend development server
frontend:
	@echo "🚀 Starting frontend development server..."
	@cd frontend && pnpm run dev

# Backend development server
backend:
	@echo "🚀 Starting backend development server..."
	@cd backend && go run ./cmd/server/main.go

# Documentation server
docs:
	@echo "📚 Starting documentation server..."
	@cd docs && docsify serve 
api:
	@echo "📚 Starting API documentation..."
	@cd backend && swag init -g ./cmd/server/main.go --output ./docs
	@echo "Open Swagger API documentation at http://localhost:8080/swagger/index.html"


# Run both frontend and backend concurrently (requires 'concurrently' pnpm package or similar)
dev:
	@echo "🚀 Starting both frontend and backend..."
	@if command -v concurrently >/dev/null 2>&1; then \
		concurrently "make frontend" "make backend"; \
	else \
		echo "⚠️  'concurrently' not found. Installing..."; \
		pnpm install -g concurrently && concurrently "make frontend" "make backend"; \
	fi

# Install dependencies
install:
	@echo "📦 Installing dependencies..."
	@cd frontend && pnpm install
	@cd backend && go mod download
	@echo "✅ Dependencies installed"

# Build production versions
build:
	@echo "🔨 Building frontend..."
	@cd frontend && pnpm run build
	@echo "🔨 Building backend..."
	@cd backend && go build -o bin/server ./cmd/server/main.go
	@echo "✅ Build complete"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@cd frontend && rm -rf node_modules .next dist build
	@cd backend && rm -rf bin
	@echo "✅ Clean complete"

# Install and build everything
all: install build
	@echo "✅ All tasks completed"

# Run tests (if you add tests later)
test:
	@echo "🧪 Running tests..."
	@cd frontend && pnpm test || true
	@cd backend && go test ./... || true