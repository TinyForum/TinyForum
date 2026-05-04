# ============================================
# Module 7: Summary & Final Checks
# ============================================

print_summary() {
    # Bug fix: test_db_connection 需要三个参数；使用 setup_postgres 导出的变量
    local pg_user="${PG_FINAL_USER:-${PSQL_USER:-unknown}}"
    local pg_pass="${PG_FINAL_PASS:-${PSQL_PASS:-}}"
    local pg_db="${PG_FINAL_DB:-${PSQL_DB_NAME:-tiny_forum}}"
    local redis_user="${REDIS_FINAL_USER:-${REDIS_USER:-unknown}}"

    echo ""
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}${BOLD}  ✅  TinyForum Dev Environment Ready!${NC}"
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo "  📦 Services"
    echo "     PostgreSQL : localhost:${PSQL_PORT:-5432}  db=${pg_db}  user=${pg_user}"
    echo "     Redis      : localhost:${REDIS_PORT:-6379}  user=${redis_user}"
    echo ""
    echo "  🚀 Start Commands"
    echo "     Backend  : cd backend  && go run ./cmd/server/main.go"
    echo "     Frontend : cd frontend && ${PACKAGE_MANAGER:-npm} run dev"
    echo ""
    echo "  🌐 URLs"
    echo "     Frontend : ${FRONTEND_URL:-http://localhost:3000}"
    echo "     Backend  : ${BACKEND_URL:-http://localhost:8080}"
    echo "     API      : ${BACKEND_URL:-http://localhost:8080}/api/v1"
    echo ""
    echo "  📝 Generated Configs"
    echo "     backend/config/postgres.yml"
    echo "     backend/config/redis.yml"
    echo "     backend/config/basic.yml"
    echo "     backend/config/private.yml    ← ⚠️  sensitive, gitignored"
    echo "     backend/config/risk_control.yml"
    echo "     frontend/config.yml"
    echo ""

    # 最终数据库连通性验证（非阻断）
    echo "  🔍 Final connectivity check..."
    if test_db_connection "$pg_user" "$pg_pass" "$pg_db" 2>/dev/null; then
        echo -e "     PostgreSQL : ${GREEN}✅ OK${NC}"
    else
        echo -e "     PostgreSQL : ${YELLOW}⚠️  Cannot connect (check backend/config/postgres.yml)${NC}"
    fi

    if command -v redis-cli >/dev/null 2>&1 && \
       redis-cli --user "${REDIS_FINAL_USER:-tinyforum}" \
                 --pass "${REDIS_FINAL_PASS:-tf@password}" \
                 PING >/dev/null 2>&1; then
        echo -e "     Redis      : ${GREEN}✅ OK${NC}"
    else
        echo -e "     Redis      : ${YELLOW}⚠️  Cannot connect (check backend/config/redis.yml)${NC}"
    fi

    echo ""
    echo -e "${BOLD}  💡 Troubleshooting Tips${NC}"
    echo "  ┌─ CORS errors"
    echo "  │   backend/config/basic.yml   → allow_origins"
    echo "  │   frontend/config.yml        → allowed_dev_origins"
    echo "  ├─ DB auth errors"
    echo "  │   backend/config/private.yml → database section"
    echo "  │   Check pg_hba.conf allows md5/scram for localhost"
    echo "  └─ Redis auth errors"
    echo "      backend/config/private.yml → redis section"
    echo "      Run: redis-cli ACL LIST"
    echo  -e "${GREEN} to running: make backend ${NC}"
    echo ""
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}