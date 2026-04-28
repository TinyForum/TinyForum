#!/bin/bash
set -e

echo "🚀 Tiny Forum Development Startup"
echo "=================================="

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Detect OS and set default database user
detect_default_db_user() {
    OS=$(uname -s)
    case "$OS" in
        Darwin*)
            # macOS - default user is current user
            echo "$(whoami)"
            ;;
        Linux*)
            # Linux - default is postgres user (or current user if it exists)
            if sudo -u postgres psql -c "SELECT 1" >/dev/null 2>&1; then
                echo "postgres"
            else
                echo "$(whoami)"
            fi
            ;;
        MINGW*|MSYS*|CYGWIN*)
            # Windows (Git Bash, Cygwin) - default is postgres
            echo "postgres"
            ;;
        *)
            # Unknown OS - fallback to current user
            echo "$(whoami)"
            ;;
    esac
}

# Get current system user and platform-specific default DB user
CURRENT_USER=$(whoami)
DEFAULT_DB_USER=$(detect_default_db_user)
echo -e "${GREEN}System user: $CURRENT_USER${NC}"
echo -e "${GREEN}Default database user: $DEFAULT_DB_USER${NC}"

# Allow override via environment variable
DB_USER=${DB_USER:-$DEFAULT_DB_USER}
echo -e "${GREEN}Using database user: $DB_USER${NC}"

# Check dependencies
echo ""
echo "Checking dependencies..."

# Check Go
command -v go >/dev/null 2>&1 || { echo -e "${RED}❌ Go is not installed. Visit https://go.dev/dl/${NC}"; exit 1; }
echo -e "${GREEN}✅ Go found: $(go version)${NC}"

# Check Node.js and npm/pnpm
command -v node >/dev/null 2>&1 || { echo -e "${RED}❌ Node.js is not installed. Visit https://nodejs.org/${NC}"; exit 1; }
echo -e "${GREEN}✅ Node.js found: $(node --version)${NC}"

# Check Ollama
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
check_postgres_client() {
    if command -v psql >/dev/null 2>&1; then
        echo -e "${GREEN}✅ psql found: $(psql --version | head -n1)${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠️  psql not found - PostgreSQL client tools not installed${NC}"
        echo "   Install with:"
        
        OS=$(uname -s)
        case "$OS" in
            Darwin*)
                echo "      brew install libpq (macOS - client only)"
                echo "      brew install postgresql (macOS - full installation)"
                ;;
            Linux*)
                echo "      sudo apt update && sudo apt install postgresql-client (Ubuntu/Debian)"
                echo "      sudo yum install postgresql (RHEL/CentOS/Fedora)"
                ;;
            MINGW*|MSYS*|CYGWIN*)
                echo "      Download from: https://www.postgresql.org/download/windows/"
                echo "      Or use: winget install PostgreSQL.PostgreSQL"
                ;;
            *)
                echo "      Visit: https://www.postgresql.org/download/"
                ;;
        esac
        return 1
    fi
}

check_postgres_client

echo ""
echo "Step 1: Checking PostgreSQL..."

# Check if PostgreSQL is running
check_postgres_running() {
    OS=$(uname -s)
    
    case "$OS" in
        Darwin*)
            if brew services list | grep -q "postgresql.*started" || pg_isready >/dev/null 2>&1; then
                return 0
            fi
            ;;
        Linux*)
            if systemctl is-active --quiet postgresql 2>/dev/null || \
               systemctl is-active --quiet postgresql@*-main 2>/dev/null || \
               pg_isready >/dev/null 2>&1; then
                return 0
            fi
            ;;
        MINGW*|MSYS*|CYGWIN*)
            if net start | grep -qi "postgres" 2>/dev/null || pg_isready >/dev/null 2>&1; then
                return 0
            fi
            ;;
        *)
            if pg_isready >/dev/null 2>&1; then
                return 0
            fi
            ;;
    esac
    return 1
}

if check_postgres_running; then
    echo -e "${GREEN}✅ PostgreSQL is running${NC}"
else
    echo -e "${RED}❌ PostgreSQL is not running${NC}"
    echo "   Please start PostgreSQL:"
    OS=$(uname -s)
    case "$OS" in
        Darwin*)
            echo "   - macOS: brew services start postgresql"
            echo "   - macOS: pg_ctl -D /usr/local/var/postgres start"
            ;;
        Linux*)
            echo "   - Ubuntu/Debian: sudo systemctl start postgresql"
            echo "   - RHEL/CentOS:   sudo systemctl start postgresql"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            echo "   - Windows: net start postgresql-x64-15"
            echo "   - Or start service from Services panel"
            ;;
    esac
    exit 1
fi

# Check if we can connect with current user
connect_postgres() {
    local user=$1
    local db=${2:-postgres}
    
    if [ "$user" = "postgres" ] && [ "$OS" = "Linux" ]; then
        # On Linux, try with sudo for postgres user
        sudo -u postgres psql -d "$db" -c "SELECT 1" >/dev/null 2>&1
    else
        psql -h localhost -U "$user" -d "$db" -c "SELECT 1" >/dev/null 2>&1
    fi
}

if connect_postgres "$DB_USER"; then
    echo -e "${GREEN}✅ Connected to PostgreSQL as user: $DB_USER${NC}"
else
    echo -e "${YELLOW}⚠️  Cannot connect as $DB_USER, trying default users...${NC}"
    
    # Try alternative default users
    for alt_user in "postgres" "$CURRENT_USER"; do
        if [ "$alt_user" != "$DB_USER" ] && connect_postgres "$alt_user"; then
            echo -e "${GREEN}✅ Connected as $alt_user${NC}"
            DB_USER=$alt_user
            echo -e "${GREEN}   Switching to database user: $DB_USER${NC}"
            break
        fi
    done
    
    # If still can't connect, provide helpful error
    if ! connect_postgres "$DB_USER"; then
        echo -e "${RED}❌ Cannot connect to PostgreSQL${NC}"
        echo "   Possible solutions:"
        echo "   1. Set environment variable: export DB_USER=your_username"
        echo "   2. Create a database user:"
        echo "      - macOS/Linux: createuser -s $(whoami) (as postgres user)"
        echo "      - Linux: sudo -u postgres createuser -s $(whoami)"
        echo "   3. Use postgres user: PGUSER=postgres ./script.sh"
        exit 1
    fi
fi

# Check if database exists, create if not
echo "  Checking database 'tiny_forum'..."
database_exists() {
    if [ "$DB_USER" = "postgres" ] && [ "$OS" = "Linux" ]; then
        sudo -u postgres psql -lqt | cut -d \| -f 1 | grep -qw "tiny_forum"
    else
        psql -h localhost -U "$DB_USER" -lqt | cut -d \| -f 1 | grep -qw "tiny_forum"
    fi
}

create_database() {
    if [ "$DB_USER" = "postgres" ] && [ "$OS" = "Linux" ]; then
        sudo -u postgres createdb tiny_forum 2>/dev/null || \
        sudo -u postgres psql -c "CREATE DATABASE tiny_forum;" 2>/dev/null
    else
        createdb -h localhost -U "$DB_USER" tiny_forum 2>/dev/null || \
        psql -h localhost -U "$DB_USER" -c "CREATE DATABASE tiny_forum;" 2>/dev/null
    fi
}

if database_exists; then
    echo -e "${GREEN}  Database 'tiny_forum' already exists${NC}"
else
    echo "  Creating database 'tiny_forum'..."
    if create_database; then
        echo -e "${GREEN}  Database 'tiny_forum' created${NC}"
    else
        echo -e "${RED}  Failed to create database. Please check permissions${NC}"
        echo "  You can manually create the database with:"
        if [ "$OS" = "Linux" ] && [ "$DB_USER" = "postgres" ]; then
            echo "  sudo -u postgres createdb tiny_forum"
        else
            echo "  createdb tiny_forum"
        fi
        exit 1
    fi
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
    user_exists() {
        if [ "$DB_USER" = "postgres" ] && [ "$OS" = "Linux" ]; then
            sudo -u postgres psql -tAc "SELECT 1 FROM pg_roles WHERE rolname='$NEW_DB_USER'"
        else
            psql -h localhost -U "$DB_USER" -tAc "SELECT 1 FROM pg_roles WHERE rolname='$NEW_DB_USER'"
        fi
    }
    
    if user_exists | grep -q 1; then
        echo -e "${YELLOW}⚠️  User $NEW_DB_USER already exists${NC}"
    else
        create_user_sql() {
            if [ "$DB_USER" = "postgres" ] && [ "$OS" = "Linux" ]; then
                sudo -u postgres psql << EOF
                CREATE USER $NEW_DB_USER WITH PASSWORD '$NEW_DB_PASS';
                ALTER USER $NEW_DB_USER CREATEDB;
                GRANT ALL PRIVILEGES ON DATABASE tiny_forum TO $NEW_DB_USER;
EOF
            else
                psql -h localhost -U "$DB_USER" << EOF
                CREATE USER $NEW_DB_USER WITH PASSWORD '$NEW_DB_PASS';
                ALTER USER $NEW_DB_USER CREATEDB;
                GRANT ALL PRIVILEGES ON DATABASE tiny_forum TO $NEW_DB_USER;
EOF
            fi
        }
        
        create_user_sql
        echo -e "${GREEN}✅ User $NEW_DB_USER created${NC}"
    fi
    
    # Grant schema privileges
    grant_schema_privs() {
        if [ "$DB_USER" = "postgres" ] && [ "$OS" = "Linux" ]; then
            sudo -u postgres psql -d tiny_forum << EOF
                GRANT ALL PRIVILEGES ON SCHEMA public TO $NEW_DB_USER;
                ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO $NEW_DB_USER;
EOF
        else
            psql -h localhost -U "$DB_USER" -d tiny_forum << EOF
                GRANT ALL PRIVILEGES ON SCHEMA public TO $NEW_DB_USER;
                ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO $NEW_DB_USER;
EOF
        fi
    }
    
    grant_schema_privs
    
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
    # Update the database user in existing config (cross-platform sed)
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' "s/user:.*/user: $DB_USER/" "$CONFIG_FILE" 2>/dev/null || true
    else
        sed -i "s/user:.*/user: $DB_USER/" "$CONFIG_FILE" 2>/dev/null || true
    fi
    
    # Update password if a new user was created with password
    if [[ $REPLY =~ ^[Yy]$ ]] && [ -n "$NEW_DB_PASS" ]; then
        if [[ "$OSTYPE" == "darwin"* ]]; then
            sed -i '' "s/password:.*/password: $NEW_DB_PASS/" "$CONFIG_FILE" 2>/dev/null || true
        else
            sed -i "s/password:.*/password: $NEW_DB_PASS/" "$CONFIG_FILE" 2>/dev/null || true
        fi
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
test_connection() {
    if [ -n "$NEW_DB_PASS" ]; then
        PGPASSWORD=$NEW_DB_PASS psql -h localhost -U "$DB_USER" -d tiny_forum -c "SELECT 'Database connected successfully!' as message;" >/dev/null 2>&1
    else
        if [ "$DB_USER" = "postgres" ] && [ "$OS" = "Linux" ]; then
            sudo -u postgres psql -d tiny_forum -c "SELECT 'Database connected successfully!' as message;" >/dev/null 2>&1
        else
            psql -h localhost -U "$DB_USER" -d tiny_forum -c "SELECT 'Database connected successfully!' as message;" >/dev/null 2>&1
        fi
    fi
}

if test_connection; then
    echo -e "${GREEN}✅ Database connection successful${NC}"
else
    echo -e "${YELLOW}⚠️  Could not connect to database, but setup completed${NC}"
fi

echo ""
echo -e "${GREEN}All services are ready! 🎉${NC}"