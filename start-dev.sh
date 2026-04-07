#!/bin/bash
set -e

echo "🚀 BBS Forum Development Startup"
echo "=================================="

# Check dependencies
command -v go >/dev/null 2>&1 || { echo "❌ Go is not installed. Visit https://go.dev/dl/"; exit 1; }
command -v node >/dev/null 2>&1 || { echo "❌ Node.js is not installed. Visit https://nodejs.org/"; exit 1; }
command -v psql >/dev/null 2>&1 || echo "⚠️  psql not found locally - make sure PostgreSQL is running"

echo ""
echo "Step 1: Setting up PostgreSQL with Docker..."
docker run -d \
  --name bbs_postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=bbs_forum \
  -p 5432:5432 \
  postgres:16-alpine 2>/dev/null || echo "  (Container may already exist, continuing...)"

echo "  Waiting for PostgreSQL to be ready..."
sleep 3

echo ""
echo "Step 2: Setting up Backend..."
cd backend
go mod tidy
echo "  Dependencies downloaded."

echo ""
echo "Step 3: Setting up Frontend..."
cd ../frontend
npm install
echo "  Dependencies installed."

echo ""
echo "=================================="
echo "✅ Setup complete!"
echo ""
echo "To start the backend:   cd backend && go run ./cmd/server/main.go"
echo "To start the frontend:  cd frontend && npm run dev"
echo ""
echo "Or use Docker Compose:  docker compose up -d"
echo "=================================="
