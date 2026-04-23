#!/bin/bash
set -e

echo "🚀 Tiny Forum Development Startup"
echo "=================================="

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get current system user (macOS default PostgreSQL user)
DB_USER=$(whoami)
echo -e "${GREEN}Using database user: $DB_USER${NC}"

# Check dependencies
echo "Checking dependencies..."

# Check Go
command -v go >/dev/null 2>&1 || { echo -e "${RED}❌ Go is not installed. Visit https://go.dev/dl/${NC}"; exit 1; }
echo -e "${GREEN}✅ Go found: $(go version)${NC}"

# Check Node.js and npm/pnpm
command -v node >/dev/null 2>&1 || { echo -e "${RED}❌ Node.js is not installed. Visit https://nodejs.org/${NC}"; exit 1; }
echo -e "${GREEN}✅ Node.js found: $(node --version)${NC}"

# Chack Ollama
command -v ollama >/dev/null 2>&1 || { echo -e "${RED}❌ Ollama is not installed. Visit https://ollama.com${NC}"; exit 1; }
echo -e "${GREEN}✅ Ollama found: $(ollama --version)${NC}"

# Check package manager (pnpm preferred, fallback to npm)
if command -v pnpm >/dev/null 2>&1; then
    PACKAGE_MANAGER="pnpm"
    echo -e "${GREEN}✅ pnpm found: $(pnpm --version)${NC}"
elif command -v npm >/dev/null 2>&1; then
    PACKAGE_MANAGER="npm"
    echo -e "${YELLOW}⚠️  pnpm not found, using npm instead${NC}"
    echo -e "${GREEN}✅ npm found: $(npm --version)${NC}"
else
    echo -e "${RED}❌ No package manager (npm/pnpm) found. Visit https://nodejs.org/${NC}"
    exit 1
fi

# Check PostgreSQL client
if command -v psql >/dev/null 2>&1; then
    echo -e "${GREEN}✅ psql found: $(psql --version)${NC}"
else
    echo -e "${YELLOW}⚠️  psql not found - PostgreSQL client tools not installed${NC}"
    echo "   Install with: brew install libpq (macOS) or apt-get install postgresql-client (Ubuntu)"
fi

echo ""
echo "Step 1: Checking PostgreSQL..."
# Check if PostgreSQL is running
if pg_isready -h localhost -p 5432 >/dev/null 2>&1; then
    echo -e "${GREEN}✅ PostgreSQL is running${NC}"
else
    echo -e "${RED}❌ PostgreSQL is not running${NC}"
    echo "   Please start PostgreSQL:"
    echo "   - macOS: brew services start postgresql"
    echo "   - Linux: sudo systemctl start postgresql"
    exit 1
fi

# Check if we can connect with current user
if psql -h localhost -U $DB_USER -d postgres -c "SELECT 1" >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Connected to PostgreSQL as user: $DB_USER${NC}"
else
    echo -e "${RED}❌ Cannot connect to PostgreSQL as user: $DB_USER${NC}"
    echo "   Please ensure PostgreSQL is installed and configured correctly"
    exit 1
fi

# Check if database exists, create if not
echo "  Checking database 'tiny_forum'..."
if psql -h localhost -U $DB_USER -lqt | cut -d \| -f 1 | grep -qw "tiny_forum"; then
    echo -e "${GREEN}  Database 'tiny_forum' already exists${NC}"
else
    echo "  Creating database 'tiny_forum'..."
    createdb -h localhost -U $DB_USER tiny_forum 2>/dev/null || \
    psql -h localhost -U $DB_USER -c "CREATE DATABASE tiny_forum;" 2>/dev/null || {
        echo -e "${RED}  Failed to create database. Please check permissions${NC}"
        echo "  You can manually create the database with:"
        echo "  createdb tiny_forum"
        exit 1
    }
    echo -e "${GREEN}  Database 'tiny_forum' created${NC}"
fi

# Ask whether to create a new database user (moved after database creation)
echo ""
read -p "Create a new database user? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    read -p "Enter username: " NEW_DB_USER
    read -sp "Enter password: " NEW_DB_PASS
    echo
    
    # Check if user already exists
    if psql -h localhost -U $DB_USER -d postgres -tAc "SELECT 1 FROM pg_roles WHERE rolname='$NEW_DB_USER'" | grep -q 1; then
        echo -e "${YELLOW}⚠️  User $NEW_DB_USER already exists${NC}"
    else
        psql -h localhost -U $DB_USER -d postgres << EOF
        CREATE USER $NEW_DB_USER WITH PASSWORD '$NEW_DB_PASS';
        ALTER USER $NEW_DB_USER CREATEDB;
        GRANT ALL PRIVILEGES ON DATABASE tiny_forum TO $NEW_DB_USER;
EOF
        echo -e "${GREEN}✅ User $NEW_DB_USER created${NC}"
    fi
    
    # Grant schema privileges
    psql -h localhost -U $DB_USER -d tiny_forum << EOF
        GRANT ALL PRIVILEGES ON SCHEMA public TO $NEW_DB_USER;
        ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO $NEW_DB_USER;
EOF
    
    DB_USER=$NEW_DB_USER
    echo -e "${GREEN}✅ Now using database user: $DB_USER${NC}"
fi

echo ""
echo "Step 2: Setting up Backend..."
cd backend

# Check if go.mod exists
if [ ! -f "go.mod" ]; then
    echo -e "${RED}❌ go.mod not found. Are you in the right directory?${NC}"
    exit 1
fi

echo "  Running go mod tidy..."
go mod tidy
echo -e "${GREEN}  Dependencies downloaded.${NC}"

# Create/update config file with correct database settings
echo "  Configuring database settings..."
CONFIG_FILE="config/config.yaml"
if [ ! -f "$CONFIG_FILE" ]; then
    mkdir -p config
    cat > "$CONFIG_FILE" << EOF
database:
  host: localhost
  port: 5432
  user: $DB_USER
  password: ""
  dbname: tiny_forum
  sslmode: disable
  timezone: Asia/Shanghai

server:
  port: 8080
  mode: debug

jwt:
  secret: your-secret-key-change-this-in-production
  expire_hours: 24

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
EOF
    echo -e "${GREEN}  Created config file at $CONFIG_FILE${NC}"
else
    echo -e "${GREEN}  Config file already exists${NC}"
    # Update the database user in existing config (macOS sed)
    sed -i '' "s/user:.*/user: $DB_USER/" "$CONFIG_FILE" 2>/dev/null || \
    sed -i "s/user:.*/user: $DB_USER/" "$CONFIG_FILE" 2>/dev/null || true
    
    # Update password if a new user was created with password
    if [[ $REPLY =~ ^[Yy]$ ]] && [ -n "$NEW_DB_PASS" ]; then
        sed -i '' "s/password:.*/password: $NEW_DB_PASS/" "$CONFIG_FILE" 2>/dev/null || \
        sed -i "s/password:.*/password: $NEW_DB_PASS/" "$CONFIG_FILE" 2>/dev/null || true
    fi
fi

cd ..

echo ""
echo "Step 3: Setting up Frontend..."
cd frontend

# Check if package.json exists
if [ ! -f "package.json" ]; then
    echo -e "${RED}❌ package.json not found. Are you in the right directory?${NC}"
    exit 1
fi

echo "  Installing dependencies with $PACKAGE_MANAGER..."
if [ "$PACKAGE_MANAGER" = "pnpm" ]; then
    pnpm install
else
    npm install
fi
echo -e "${GREEN}  Frontend dependencies installed.${NC}"

cd ..

echo ""
echo "=================================="
echo -e "${GREEN}✅ Setup complete!${NC}"
echo ""
echo "To start the backend:"
echo "  cd backend && go run ./cmd/server/main.go"
echo ""
echo "To start the frontend:"
echo "  cd frontend && ${PACKAGE_MANAGER} run dev"
echo ""
echo "Database connection info:"
echo "  Host: localhost:5432"
echo "  User: $DB_USER"
if [[ $REPLY =~ ^[Yy]$ ]] && [ -n "$NEW_DB_PASS" ]; then
    echo "  Password: $NEW_DB_PASS"
else
    echo "  Password: (empty - using trust authentication)"
fi
echo "  Database: tiny_forum"
echo "=================================="

# Test database connection
echo ""
echo "Testing database connection..."
if [ -n "$NEW_DB_PASS" ]; then
    PGPASSWORD=$NEW_DB_PASS psql -h localhost -U $DB_USER -d tiny_forum -c "SELECT 'Database connected successfully!' as message;" >/dev/null 2>&1 && \
    echo -e "${GREEN}✅ Database connection successful${NC}" || \
    echo -e "${YELLOW}⚠️  Could not connect to database, but setup completed${NC}"
else
    psql -h localhost -U $DB_USER -d tiny_forum -c "SELECT 'Database connected successfully!' as message;" >/dev/null 2>&1 && \
    echo -e "${GREEN}✅ Database connection successful${NC}" || \
    echo -e "${YELLOW}⚠️  Could not connect to database, but setup completed${NC}"
fi

echo ""
echo -e "${GREEN}All services are ready! 🎉${NC}"