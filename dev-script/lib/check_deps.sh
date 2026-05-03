
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
        echo -e "${YELLOW}⚠️ Ollama is not installed. Visit https://ollama.com${NC}"
        return 1
    fi
}

check_llamacpp() {
    echo "Checking llama.cpp..."
    if command -v llama-cli >/dev/null 2>&1 || command -v llamacpp >/dev/null 2>&1; then
        # 获取版本信息（可选）
        echo -e "${GREEN}✅ llama.cpp found${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠️ llama.cpp is not installed. Visit https://github.com/ggml-org/llama.cpp${NC}"
        return 1
    fi
}

check_llm_runtime() {
    echo "Checking LLM runtime (Ollama or llama.cpp)..."
    local ollama_ok=0
    local llamacpp_ok=0
    
    if command -v ollama >/dev/null 2>&1; then
        echo -e "${GREEN}  ✅ Ollama found: $(ollama --version 2>/dev/null || echo 'installed')${NC}"
        ollama_ok=1
    fi
    
    if command -v llama-cli >/dev/null 2>&1 || command -v llamacpp >/dev/null 2>&1; then
        echo -e "${GREEN}  ✅ llama.cpp found${NC}"
        llamacpp_ok=1
    fi
    
    if [ $ollama_ok -eq 0 ] && [ $llamacpp_ok -eq 0 ]; then
        echo -e "${RED}❌ Neither Ollama nor llama.cpp is installed.${NC}"
        echo "   Please install at least one:"
        echo "   - Ollama: https://ollama.com"
        echo "   - llama.cpp: https://github.com/ggml-org/llama.cpp"
        return 1
    fi
    
    echo -e "${GREEN}✅ LLM runtime available${NC}"
    return 0
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
    local ERRORS=0
    echo ""
    echo "Checking dependencies..."
    
    check_go || ((ERRORS++))
    check_nodejs || ((ERRORS++))
    check_llm_runtime || ((ERRORS++))       
    check_package_manager || ((ERRORS++))
    check_postgres_client || ((ERRORS++))
    check_redis_client || ((ERRORS++))
    
    if [ $ERRORS -gt 0 ]; then
        echo -e "${RED}Missing $ERRORS dependencies. Please install them and try again.${NC}"
        exit 1
    fi
    echo -e "${GREEN}All dependencies satisfied.${NC}"
}