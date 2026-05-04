#!/usr/bin/env bash
# ============================================
# TinyForum Development Environment Setup
# ============================================
# 用法: bash dev-script/start-dev.sh
#       make init-dev
#
# 修复:
#   [BUG-1] set -e 与 source 的交互：source 的文件中 exit 1 会终止父 shell
#           修复: 各模块改用 return 1，主流程用 if ! func; then 判断
#   [BUG-2] psql_user 检测后未被使用（赋值到局部变量但 setup_postgres 用全局）
#           修复: 移除冗余变量，detect_default_db_user 结果仅用于日志展示
#   [BUG-3] set -e 会导致 read 命令在 EOF 时退出脚本
#           修复: 改为 set -eEo pipefail，并在交互式读取处局部关闭

set -eEo pipefail

# ── 路径解析 ──────────────────────────────────────────────────────────────────
DEV_SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPTS_DIR="${DEV_SCRIPT_DIR}"
DEV_DIR="${SCRIPTS_DIR}/dev"

# ── 加载模块 ──────────────────────────────────────────────────────────────────
# shellcheck source=lib/colors.sh
source "${DEV_DIR}/colors.sh"
# shellcheck source=lib/globals.sh
source "${DEV_DIR}/globals.sh"
# shellcheck source=lib/detect_os.sh
source "${DEV_DIR}/detect_os.sh"
# shellcheck source=lib/check_deps.sh
source "${DEV_DIR}/check_deps.sh"
# shellcheck source=lib/postgres.sh
source "${DEV_DIR}/postgres.sh"
# shellcheck source=lib/redis.sh
source "${DEV_DIR}/redis.sh"
# shellcheck source=lib/config.sh
source "${DEV_DIR}/config.sh"
# shellcheck source=lib/backend.sh
source "${DEV_DIR}/backend.sh"
# shellcheck source=lib/frontend.sh
source "${DEV_DIR}/frontend.sh"
# shellcheck source=lib/check_config.sh
source "${DEV_DIR}/check_config.sh"
# shellcheck source=lib/summary.sh
source "${DEV_DIR}/summary.sh"


echo "  Script Dir: ${SCRIPT_DIR}"
echo "  Lib Dir: ${DEV_DIR}"
echo "  Project Root: ${PROJECT_ROOT}"


# ── 全局错误处理 ───────────────────────────────────────────────────────────────
_on_error() {
    local lineno="$1" cmd="$2"
    echo ""
    echo -e "${RED}${BOLD}━━━ Unexpected Error ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${RED}  Script failed at line ${lineno}: ${cmd}${NC}"
    echo -e "${RED}  Run with 'bash -x ${BASH_SOURCE[0]}' for full trace.${NC}"
    echo -e "${RED}${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}
trap '_on_error $LINENO "$BASH_COMMAND"' ERR

# ── 主流程 ────────────────────────────────────────────────────────────────────
main() {
    echo "${BANNER}"
    echo ""
    echo -e "${BOLD}🚀 TinyForum — Development Environment Setup${NC}"
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo "  Project : ${PROJECT_ROOT}"
    echo "  Version : ${PROJECT_VERSION}"
    echo ""

    # ── Module 1: OS Detection ────────────────────────────────────────────
    detect_os

    # detect_default_db_user 仅用于信息展示；
    # setup_postgres 内部自行判断管理员用户，无需外部传入
    local detected_db_user
    detected_db_user=$(detect_default_db_user)
    echo "  Detected DB admin candidate: ${detected_db_user}"
    echo ""

    # ── Module 2: Dependency Check ────────────────────────────────────────
    if ! check_all_dependencies; then
        echo -e "${RED}❌ Dependency check failed. Fix missing tools and re-run.${NC}"
        exit 1
    fi

    # ── Module 3: PostgreSQL ──────────────────────────────────────────────
    # set -e 对交互式 read 不友好；临时关闭
    set +e
    setup_postgres
    local pg_rc=$?
    set -e
    if [ "$pg_rc" -ne 0 ]; then
        echo -e "${RED}❌ PostgreSQL setup failed (exit ${pg_rc}).${NC}"
        echo "   Run 'bash -x ${BASH_SOURCE[0]}' for details."
        exit 1
    fi

    # ── Module 4: Redis ───────────────────────────────────────────────────
    set +e
    setup_redis
    local redis_rc=$?
    set -e
    if [ "$redis_rc" -ne 0 ]; then
        echo -e "${RED}❌ Redis setup failed (exit ${redis_rc}).${NC}"
        exit 1
    fi

    # ── Module 4.5: Config files ──────────────────────────────────────────
    if ! setup_configurations; then
        echo -e "${RED}❌ Configuration setup failed.${NC}"
        exit 1
    fi

    # ── Module 5: Backend ─────────────────────────────────────────────────
    if ! setup_backend; then
        echo -e "${RED}❌ Backend setup failed.${NC}"
        exit 1
    fi

    # ── Module 6: Frontend ────────────────────────────────────────────────
    if ! setup_frontend; then
        echo -e "${RED}❌ Frontend setup failed.${NC}"
        exit 1
    fi
    # ── Module 7: Config Check  ────────────────────────────────────────────────
    if ! check_config; then
        echo -e "${RED}❌ Configuration check failed.${NC}"
        exit 1
    fi
    # ── Module 8: Summary ─────────────────────────────────────────────────
    print_summary
}

main "$@"