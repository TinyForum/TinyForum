# ============================================
# Module 2: Dependency Checking
# ============================================

check_go() {
    echo -n "  Checking Go... "
    if command -v go >/dev/null 2>&1; then
        echo -e "${GREEN}✅ $(go version)${NC}"
        return 0
    else
        echo -e "${RED}❌ not found${NC}"
        echo "     → Install: https://go.dev/dl/"
        return 1
    fi
}

check_nodejs() {
    echo -n "  Checking Node.js... "
    if command -v node >/dev/null 2>&1; then
        echo -e "${GREEN}✅ $(node --version)${NC}"
        return 0
    else
        echo -e "${RED}❌ not found${NC}"
        echo "     → Install: https://nodejs.org/"
        return 1
    fi
}

check_llm_runtime() {
    echo "  Checking LLM runtime (Ollama or llama.cpp)..."
    local ollama_ok=0 llamacpp_ok=0

    if command -v ollama >/dev/null 2>&1; then
        local ver
        ver=$(ollama --version 2>/dev/null || echo "installed")
        echo -e "     ${GREEN}✅ Ollama: ${ver}${NC}"
        ollama_ok=1
    fi

    if command -v llama-cli >/dev/null 2>&1 || command -v llamacpp >/dev/null 2>&1; then
        echo -e "     ${GREEN}✅ llama.cpp found${NC}"
        llamacpp_ok=1
    fi

    if [ "$ollama_ok" -eq 0 ] && [ "$llamacpp_ok" -eq 0 ]; then
        echo -e "     ${YELLOW}⚠️  Neither Ollama nor llama.cpp found (optional)${NC}"
        echo "        → Ollama  : https://ollama.com"
        echo "        → llama.cpp: https://github.com/ggml-org/llama.cpp"
        return 1
    fi
    return 0
}

check_package_manager() {
    echo -n "  Checking JS package manager... "
    if command -v pnpm >/dev/null 2>&1; then
        PACKAGE_MANAGER="pnpm"
        echo -e "${GREEN}✅ pnpm $(pnpm --version)${NC}"
    elif command -v npm >/dev/null 2>&1; then
        PACKAGE_MANAGER="npm"
        echo -e "${YELLOW}⚠️  pnpm not found, using npm $(npm --version)${NC}"
    else
        echo -e "${RED}❌ No package manager (npm/pnpm) found${NC}"
        echo "     → Install Node.js: https://nodejs.org/"
        return 1
    fi
}

check_postgres_client() {
    echo -n "  Checking PostgreSQL client (psql)... "
    if command -v psql >/dev/null 2>&1; then
        echo -e "${GREEN}✅ $(psql --version | head -n1)${NC}"
        return 0
    fi
    echo -e "${RED}❌ not found${NC}"
    case "$OS" in
        Darwin*)
            echo "     → brew install libpq && brew link --force libpq"
            echo "       (or: brew install postgresql for full install)"
            ;;
        Linux*)
            if [ -f /etc/os-release ]; then
                # shellcheck source=/dev/null
                local id; id=$(. /etc/os-release && echo "${ID}")
                case "$id" in
                    ubuntu|debian)
                        echo "     → sudo apt update && sudo apt install postgresql-client"
                        ;;
                    rhel|centos|fedora|rocky|alma)
                        echo "     → sudo dnf install postgresql"
                        ;;
                    *)
                        echo "     → Install postgresql-client via your package manager"
                        ;;
                esac
            fi
            ;;
    esac
    return 1
}

check_redis_client() {
    echo -n "  Checking Redis client (redis-cli)... "
    if command -v redis-cli >/dev/null 2>&1; then
        echo -e "${GREEN}✅ $(redis-cli --version | head -n1)${NC}"
        return 0
    fi
    echo -e "${RED}❌ not found${NC}"
    case "$OS" in
        Darwin*)
            echo "     → brew install redis"
            ;;
        Linux*)
            if [ -f /etc/os-release ]; then
                # shellcheck source=/dev/null
                local id; id=$(. /etc/os-release && echo "${ID}")
                case "$id" in
                    ubuntu|debian)
                        echo "     → sudo apt update && sudo apt install redis-tools"
                        ;;
                    rhel|centos|fedora|rocky|alma)
                        echo "     → sudo dnf install redis"
                        ;;
                    *)
                        echo "     → Install redis-tools via your package manager"
                        ;;
                esac
            fi
            ;;
    esac
    return 1
}

check_all_dependencies() {
    echo ""
    echo -e "${BOLD}━━━ Dependency Check ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    local errors=0

    check_go             || ((errors++))
    check_nodejs         || ((errors++))
    check_llm_runtime    || true   # LLM runtime 是可选的，不阻断流程
    check_package_manager || ((errors++))
    check_postgres_client || ((errors++))
    check_redis_client   || ((errors++))

    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    if [ "$errors" -gt 0 ]; then
        echo -e "${RED}❌ $errors required dependency/dependencies missing. Please install them and retry.${NC}"
        return 1
    fi
    echo -e "${GREEN}✅ All required dependencies satisfied.${NC}"
    return 0
}