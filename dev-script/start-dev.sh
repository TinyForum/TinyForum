#!/bin/bash
set -e

# ============================================
# Tiny Forum Development Setup Script
# ============================================

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Global variables
OS=$(uname -s)
CURRENT_USER=$(whoami)
DEFAULT_DB_USER=""
DB_USER=""
PACKAGE_MANAGER=""
NEW_DB_USER=""
NEW_DB_PASS=""

# ============================================
# Module 1: System Detection
# ============================================

detect_os() {
    echo -e "${GREEN}Detecting operating system...${NC}"
    case "$OS" in
        Darwin*)
            echo -e "${GREEN}✓ macOS detected${NC}"
            ;;
        Linux*)
            echo -e "${GREEN}✓ Linux detected${NC}"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            echo -e "${GREEN}✓ Windows detected${NC}"
            ;;
        *)
            echo -e "${YELLOW}⚠ Unknown operating system: $OS${NC}"
            ;;
    esac
}

detect_default_db_user() {
    case "$OS" in
        Darwin*)
            echo "$(whoami)"
            ;;
        Linux*)
            if sudo -u postgres psql -c "SELECT 1" >/dev/null 2>&1; then
                echo "postgres"
            else
                echo "$(whoami)"
            fi
            ;;
        MINGW*|MSYS*|CYGWIN*)
            echo "postgres"
            ;;
        *)
            echo "$(whoami)"
            ;;
    esac
}

# ============================================
# Module 2: Dependency Checking
# ============================================

check_go() {
    echo "Checking Go..."
    if command -v go >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Go found: $(go version)${NC}"
        return 0
    else
        echo -e "${RED}❌ Go is not installed. Visit https://go.dev/dl/${NC}"
        exit 1
    fi
}

check_nodejs() {
    echo "Checking Node.js..."
    if command -v node >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Node.js found: $(node --version)${NC}"
        return 0
    else
        echo -e "${RED}❌ Node.js is not installed. Visit https://nodejs.org/${NC}"
        exit 1
    fi
}

check_ollama() {
    echo "Checking Ollama..."
    if command -v ollama >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Ollama found: $(ollama --version)${NC}"
        return 0
    else
        echo -e "${RED}❌ Ollama is not installed. Visit https://ollama.com${NC}"
        exit 1
    fi
}

check_package_manager() {
    echo "Checking package manager..."
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
}

check_postgres_client() {
    echo "Checking PostgreSQL client..."
    if command -v psql >/dev/null 2>&1; then
        echo -e "${GREEN}✅ psql found: $(psql --version | head -n1)${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠️  psql not found - PostgreSQL client tools not installed${NC}"
        echo "   Install with:"
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

check_all_dependencies() {
    echo ""
    echo "Checking dependencies..."
    check_go
    check_nodejs
    check_ollama
    check_package_manager
    check_postgres_client
}

# ============================================
# Module 3: PostgreSQL Management
# ============================================

check_postgres_running() {
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

connect_postgres() {
    local user=$1
    local db=${2:-postgres}
    
    if [ "$user" = "postgres" ] && [ "$OS" = "Linux" ]; then
        sudo -u postgres psql -d "$db" -c "SELECT 1" >/dev/null 2>&1
    else
        psql -h localhost -U "$user" -d "$db" -c "SELECT 1" >/dev/null 2>&1
    fi
}

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

user_exists() {
    local username=$1
    if [ "$DB_USER" = "postgres" ] && [ "$OS" = "Linux" ]; then
        sudo -u postgres psql -tAc "SELECT 1 FROM pg_roles WHERE rolname='$username'"
    else
        psql -h localhost -U "$DB_USER" -tAc "SELECT 1 FROM pg_roles WHERE rolname='$username'"
    fi
}

create_database_user() {
    local username=$1
    local password=$2
    
    if [ "$DB_USER" = "postgres" ] && [ "$OS" = "Linux" ]; then
        sudo -u postgres psql << EOF
        CREATE USER $username WITH PASSWORD '$password';
        ALTER USER $username CREATEDB;
        GRANT ALL PRIVILEGES ON DATABASE tiny_forum TO $username;
EOF
    else
        psql -h localhost -U "$DB_USER" << EOF
        CREATE USER $username WITH PASSWORD '$password';
        ALTER USER $username CREATEDB;
        GRANT ALL PRIVILEGES ON DATABASE tiny_forum TO $username;
EOF
    fi
}

grant_schema_privileges() {
    local username=$1
    
    if [ "$DB_USER" = "postgres" ] && [ "$OS" = "Linux" ]; then
        sudo -u postgres psql -d tiny_forum << EOF
        GRANT ALL PRIVILEGES ON SCHEMA public TO $username;
        ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO $username;
EOF
    else
        psql -h localhost -U "$DB_USER" -d tiny_forum << EOF
        GRANT ALL PRIVILEGES ON SCHEMA public TO $username;
        ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO $username;
EOF
    fi
}

test_db_connection() {
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

setup_postgres() {
    echo ""
    echo "Step 1: Checking PostgreSQL..."
    
    if ! check_postgres_running; then
        echo -e "${RED}❌ PostgreSQL is not running${NC}"
        echo "   Please start PostgreSQL:"
        case "$OS" in
            Darwin*)
                echo "   - macOS: brew services start postgresql"
                ;;
            Linux*)
                echo "   - Linux: sudo systemctl start postgresql"
                ;;
            MINGW*|MSYS*|CYGWIN*)
                echo "   - Windows: net start postgresql-x64-15"
                ;;
        esac
        exit 1
    fi
    echo -e "${GREEN}✅ PostgreSQL is running${NC}"
    
    # Try to connect
    if connect_postgres "$DB_USER"; then
        echo -e "${GREEN}✅ Connected to PostgreSQL as user: $DB_USER${NC}"
    else
        echo -e "${YELLOW}⚠️  Cannot connect as $DB_USER, trying default users...${NC}"
        for alt_user in "postgres" "$CURRENT_USER"; do
            if [ "$alt_user" != "$DB_USER" ] && connect_postgres "$alt_user"; then
                echo -e "${GREEN}✅ Connected as $alt_user${NC}"
                DB_USER=$alt_user
                echo -e "${GREEN}   Switching to database user: $DB_USER${NC}"
                break
            fi
        done
        
        if ! connect_postgres "$DB_USER"; then
            echo -e "${RED}❌ Cannot connect to PostgreSQL${NC}"
            echo "   Possible solutions:"
            echo "   1. Set environment variable: export DB_USER=your_username"
            echo "   2. Create a database user: createuser -s $(whoami)"
            exit 1
        fi
    fi
    
    # Create database if not exists
    echo "  Checking database 'tiny_forum'..."
    if database_exists; then
        echo -e "${GREEN}  Database 'tiny_forum' already exists${NC}"
    else
        echo "  Creating database 'tiny_forum'..."
        if create_database; then
            echo -e "${GREEN}  Database 'tiny_forum' created${NC}"
        else
            echo -e "${RED}  Failed to create database${NC}"
            exit 1
        fi
    fi
    
    # Ask to create new user
    echo ""
    read -p "Create a new database user? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        read -p "Enter username: " NEW_DB_USER
        read -sp "Enter password: " NEW_DB_PASS
        echo
        
        if user_exists "$NEW_DB_USER" | grep -q 1; then
            echo -e "${YELLOW}⚠️  User $NEW_DB_USER already exists${NC}"
        else
            create_database_user "$NEW_DB_USER" "$NEW_DB_PASS"
            echo -e "${GREEN}✅ User $NEW_DB_USER created${NC}"
        fi
        
        grant_schema_privileges "$NEW_DB_USER"
        DB_USER=$NEW_DB_USER
        echo -e "${GREEN}✅ Now using database user: $DB_USER${NC}"
    fi
}

# ============================================
# Module 4: Configuration Management
# ============================================

create_private_config() {
    local config_file="config/private.yaml"
    
    if [ -f "$config_file" ]; then
        echo -e "${YELLOW}  Private config already exists: $config_file${NC}"
        # Update database user
        if [[ "$OSTYPE" == "darwin"* ]]; then
            sed -i '' "s/user:.*/user: $DB_USER/" "$config_file" 2>/dev/null || true
        else
            sed -i "s/user:.*/user: $DB_USER/" "$config_file" 2>/dev/null || true
        fi
        if [[ -n "$NEW_DB_PASS" ]]; then
            if [[ "$OSTYPE" == "darwin"* ]]; then
                sed -i '' "s/password:.*/password: $NEW_DB_PASS/" "$config_file" 2>/dev/null || true
            else
                sed -i "s/password:.*/password: $NEW_DB_PASS/" "$config_file" 2>/dev/null || true
            fi
        fi
        echo -e "${GREEN}  ✓ Updated database credentials${NC}"
        return 0
    fi
    
    mkdir -p config
    cat > "$config_file" << EOF
# Private Configuration
# 包含敏感信息，请勿提交到版本控制

# 邮件配置
email:
  host: smtp.gmail.com
  port: 587
  username: noreply@example.com
  password: your-email-password
  from: noreply@example.com
  from_name: Tiny Forum
  ssl: false
  tls: true
  pool_size: 5

# JWT配置
jwt:
  secret: "tiny-forum-secret-change-in-production-32chars"
  expire: 24h
  refresh_expire: 168h
  issuer: "tiny-forum"

# 数据库配置
database:
  host: localhost
  port: 5432
  user: $DB_USER
  password: "${NEW_DB_PASS:-}"
  dbname: tiny_forum
  sslmode: disable
  timezone: Asia/Shanghai

# 服务器配置
server:
  port: 8080
  mode: debug

# Redis配置
redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

# 管理员账户
admin:
  email: admin@test.com
  password: password
  username: admin
  role: super_admin
  score: 10000
EOF
    
    echo -e "${GREEN}  ✓ Created private config: $config_file${NC}"
    return 0
}

create_basic_config() {
    local config_file="config/basic.yaml"
    
    if [ -f "$config_file" ]; then
        echo -e "${YELLOW}  Basic config already exists: $config_file${NC}"
        return 0
    fi
    
    mkdir -p config
    cat > "$config_file" << 'EOF'
# Basic Configuration
# 基础配置，可提交到版本控制

# 服务器配置
server:
  port: 8080
  mode: debug
  read_timeout: 30s
  write_timeout: 30s
  max_header_bytes: 1048576

# API配置
api:
  protocol: http
  host: localhost
  port: 8080
  version: v1
  prefix: /api

# JWT配置（默认值，会被 private.yaml 覆盖）
jwt:
  expire: 24h
  refresh_expire: 168h
  issuer: "tiny-forum"

# 日志配置
log:
  level: info
  filename: ./logs/app.log
  max_size: 100
  max_backups: 10
  max_age: 30
  compress: true
  console: true
  json_format: false

# 限流配置
rate_limit:
  enabled: true
  requests: 100
  duration: 60
  burst: 50

# Ollama AI 配置
ollama:
  base_url: http://localhost:11434
  model: qwen3:0.6b
  num_predict: 256
  temperature: 0.7
  timeout: 60

# CORS 跨域配置
allow_origins:
  - http://localhost:3000
  - http://127.0.0.1:3000
  - http://localhost:8080
  - http://127.0.0.1:8080

# 上传配置
upload:
  max_size: 10485760
  allowed_types:
    - .jpg
    - .jpeg
    - .png
    - .gif
    - .webp
    - .pdf
    - .zip
  avatar:
    max_size: 2097152
    allowed_types:
      - .jpg
      - .jpeg
      - .png
      - .webp
    width: 512
    height: 512
  post_image:
    max_size: 5242880
    allowed_types:
      - .jpg
      - .jpeg
      - .png
      - .gif
      - .webp
  storage:
    type: local
    local_path: ./uploads
    url_prefix: /uploads/
EOF
    
    echo -e "${GREEN}  ✓ Created basic config: $config_file${NC}"
    return 0
}

create_risk_config() {
    local config_file="config/risk_control.yaml"
    
    if [ -f "$config_file" ]; then
        echo -e "${YELLOW}  Risk control config already exists: $config_file${NC}"
        return 0
    fi
    
    mkdir -p config
    cat > "$config_file" << 'EOF'
# Risk Control Configuration
# 风控配置，可提交到版本控制

# 不同风险等级的限流策略
rate_limit:
  risk_levels:
    normal:
      create_post:
        limit: 20
        window: 1h
      create_comment:
        limit: 60
        window: 1h
      send_report:
        limit: 10
        window: 1h
      update_profile:
        limit: 5
        window: 1h
      like_post:
        limit: 100
        window: 1h
      follow_user:
        limit: 50
        window: 1h
    
    observe:
      create_post:
        limit: 5
        window: 1h
      create_comment:
        limit: 20
        window: 1h
      send_report:
        limit: 5
        window: 1h
      update_profile:
        limit: 3
        window: 1h
      like_post:
        limit: 30
        window: 1h
      follow_user:
        limit: 15
        window: 1h
    
    restrict:
      create_post:
        limit: 2
        window: 1h
      create_comment:
        limit: 5
        window: 1h
      send_report:
        limit: 2
        window: 1h
      update_profile:
        limit: 1
        window: 1h
      like_post:
        limit: 10
        window: 1h
      follow_user:
        limit: 5
        window: 1h

# 内容过滤配置
content_filter:
  enabled: true
  sensitive_words:
    - "暴力"
    - "色情"
    - "赌博"
    - "毒品"
    - "诈骗"
  custom_patterns:
    - pattern: "\\b\\d{11}\\b"
      action: mask
    - pattern: "\\b[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}\\b"
      action: mask

# 反垃圾配置
anti_spam:
  enabled: true
  max_similar_posts: 3
  similarity_threshold: 0.85
  check_interval: 5m
  ban_duration: 24h

# IP 黑名单
ip_blacklist: []

# 用户黑名单
user_blacklist: []

# 风控日志
audit_log:
  enabled: true
  filename: ./logs/audit.log
  max_size: 100
  max_backups: 30
  max_age: 90
  compress: true
EOF
    
    echo -e "${GREEN}  ✓ Created risk control config: $config_file${NC}"
    return 0
}

create_config_example() {
    local example_file="config/private.example.yaml"
    
    if [ -f "$example_file" ]; then
        echo -e "${YELLOW}  Config example already exists: $example_file${NC}"
        return 0
    fi
    
    mkdir -p config
    cat > "$example_file" << 'EOF'
# Private Configuration Example
# Copy this to private.yaml and fill in your values

email:
  host: smtp.example.com
  port: 587
  username: your-email@example.com
  password: your-password
  from: your-email@example.com
  from_name: Tiny Forum

jwt:
  secret: "change-this-to-a-random-string-32chars"
  expire: 24h
  refresh_expire: 168h
  issuer: "tiny-forum"

database:
  host: localhost
  port: 5432
  user: your_db_user
  password: your_db_password
  dbname: tiny_forum
  sslmode: disable
  timezone: Asia/Shanghai

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

admin:
  email: admin@example.com
  password: change-this-password
  username: admin
  role: super_admin
  score: 10000
EOF
    
    echo -e "${GREEN}  ✓ Created config example: $example_file${NC}"
    return 0
}

create_config_gitignore() {
    local gitignore_file="config/.gitignore"
    
    if [ -f "$gitignore_file" ]; then
        echo -e "${YELLOW}  Config .gitignore already exists${NC}"
        return 0
    fi
    
    cat > "$gitignore_file" << EOF
# 忽略敏感配置
private.yaml
*.key
*.pem
*.crt

# 保留基础配置
!basic.yaml
!risk_control.yaml
!*.example.yaml
EOF
    
    echo -e "${GREEN}  ✓ Created config/.gitignore${NC}"
    return 0
}

setup_configurations() {
    echo ""
    echo "Step 2.5: Setting up configuration files..."
    
    create_basic_config
    create_risk_config
    create_private_config
    create_config_example
    create_config_gitignore
    
    echo -e "${GREEN}  ✓ All configuration files ready${NC}"
}

# ============================================
# Module 5: Backend Setup
# ============================================

setup_backend() {
    echo ""
    echo "Step 2: Setting up Backend..."
    
    cd backend
    
    if [ ! -f "go.mod" ]; then
        echo -e "${RED}❌ go.mod not found. Are you in the right directory?${NC}"
        exit 1
    fi
    
    echo "  Running go mod tidy..."
    go mod tidy
    echo -e "${GREEN}  Dependencies downloaded.${NC}"
    
    cd ..
}

# ============================================
# Module 6: Frontend Setup
# ============================================

setup_frontend() {
    echo ""
    echo "Step 3: Setting up Frontend..."
    
    cd frontend
    
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
}

# ============================================
# Module 7: Summary & Final Checks
# ============================================

print_summary() {
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
    if [[ -n "$NEW_DB_PASS" ]]; then
        echo "  Password: $NEW_DB_PASS"
    else
        echo "  Password: (empty - using trust authentication)"
    fi
    echo "  Database: tiny_forum"
    echo "=================================="
    
    echo ""
    echo "Testing database connection..."
    if test_db_connection; then
        echo -e "${GREEN}✅ Database connection successful${NC}"
    else
        echo -e "${YELLOW}⚠️  Could not connect to database, but setup completed${NC}"
    fi
    
    echo ""
    echo -e "${GREEN}All services are ready! 🎉${NC}"
}

# ============================================
# Main Execution
# ============================================

main() {
    echo "🚀 Tiny Forum Development Startup"
    echo "=================================="
    
    # Module 1: System Detection
    detect_os
    DEFAULT_DB_USER=$(detect_default_db_user)
    echo -e "${GREEN}System user: $CURRENT_USER${NC}"
    echo -e "${GREEN}Default database user: $DEFAULT_DB_USER${NC}"
    
    # Allow override via environment variable
    DB_USER=${DB_USER:-$DEFAULT_DB_USER}
    echo -e "${GREEN}Using database user: $DB_USER${NC}"
    
    # Module 2: Dependency Checking
    check_all_dependencies
    
    # Module 3: PostgreSQL Setup
    setup_postgres
    
    # Module 4: Configuration Setup
    setup_configurations
    
    # Module 5: Backend Setup
    setup_backend
    
    # Module 6: Frontend Setup
    setup_frontend
    
    # Module 7: Summary
    print_summary
}

# Run main function
main