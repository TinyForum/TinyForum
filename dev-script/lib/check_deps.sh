
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

check_redis_client() {
    echo "Checking Redis client..."
    if command -v redis-cli >/dev/null 2>&1; then
        echo -e "${GREEN}✅ redis-cli found: $(redis-cli --version | head -n1)${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠️  redis-cli not found - Redis client tools not installed${NC}"
        echo "   Install with:"
        case "$OS" in
            Darwin*)
                echo "      brew install redis                  # macOS (Homebrew)"
                echo "      # 安装后 redis-cli 即可使用"
                ;;
            Linux*)
                # 检测具体的 Linux 发行版（如果 $ID 可用，否则给通用提示）
                if [ -f /etc/os-release ]; then
                    . /etc/os-release
                    case "$ID" in
                        ubuntu|debian)
                            echo "      sudo apt update && sudo apt install redis-tools   # Ubuntu/Debian (client only)"
                            echo "      # 若需要完整 Redis 服务端，请运行: sudo apt install redis-server"
                            ;;
                        rhel|centos|fedora)
                            echo "      sudo yum install redis      # RHEL/CentOS/Fedora (includes both client and server)"
                            ;;
                        *)
                            echo "      sudo apt install redis-tools   # Debian/Ubuntu"
                            echo "      sudo yum install redis         # RHEL/CentOS/Fedora"
                            ;;
                    esac
                else
                    echo "      sudo apt install redis-tools   # Debian/Ubuntu"
                    echo "      sudo yum install redis         # RHEL/CentOS/Fedora"
                fi
                ;;
            MINGW*|MSYS*|CYGWIN*)
                echo "      winget install redis                # Windows (winget)"
                echo "      # 或从 https://github.com/microsoftarchive/redis/releases 下载安装"
                ;;
            *)
                echo "      Visit: https://redis.io/download/   # 其他系统"
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