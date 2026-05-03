#!/bin/bash
set -e

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LIB_DIR="$SCRIPT_DIR/lib"

# 引入所有模块
source "$LIB_DIR/colors.sh"
source "$LIB_DIR/globals.sh"
source "$LIB_DIR/detect_os.sh"
source "$LIB_DIR/check_deps.sh"
source "$LIB_DIR/postgres.sh"
source "$LIB_DIR/redis.sh"
source "$LIB_DIR/config.sh"
source "$LIB_DIR/backend.sh"
source "$LIB_DIR/frontend.sh"
source "$LIB_DIR/summary.sh"


# ============================================
# Main Execution
# ============================================

main() {
    echo "🚀 Tiny Forum Development Startup"
    echo "=================================="
    echo "Project root: ${PROJECT_ROOT}"
    local psql_user
   local psql_password
   local redis_user
   local redis_password
    
    # Module 1: System Detection
    detect_os >/dev/null   # 输出操作系统信息（如果需要显示可移除重定向）
    psql_user=$(detect_default_db_user) # 获取默认数据库用户
    echo -e "${GREEN}System user: $CURRENT_USER${NC}" # 显示当前用户
    echo -e "${GREEN}Default database user: $psql_user${NC}" # 显示默认数据库用户
    
    # 允许通过环境变量覆盖
    psql_user=${psql_user:-$psql_user} # 如果未设置psql_user，则使用默认值
    echo -e "${GREEN}Using database user: $psql_user${NC}"
    
    # Module 2: Dependency Checking
    if ! check_all_dependencies; then
        echo -e "${RED}Dependency check failed. Exiting.${NC}"
        exit 1
    fi
    
    # Module 3: PostgreSQL & Redis Setup
    if ! setup_postgres; then
        echo -e "${RED}PostgreSQL setup failed. Exiting.${NC}"
        exit 1
    fi
    if ! setup_redis; then
        echo -e "${RED}Redis setup failed. Exiting.${NC}"
        exit 1
    fi
    
    # Module 4: Configuration Setup
    if ! setup_configurations; then
        echo -e "${RED}Configuration setup failed. Exiting.${NC}"
        exit 1
    fi
    
    # Module 5: Backend Setup
    if ! setup_backend; then
        echo -e "${RED}Backend setup failed. Exiting.${NC}"
        exit 1
    fi
    
    # Module 6: Frontend Setup
    if ! setup_frontend; then
        echo -e "${RED}Frontend setup failed. Exiting.${NC}"
        exit 1
    fi
    
    # Module 7: Summary
    print_summary
    
    echo ""
    echo -e "${YELLOW}⚠️  Troubleshooting Tips:${NC}"
    echo "如果无法访问后端，请检查 CORS 配置："
    echo "  1. backend/config/basic.yaml 中的 allow_origins 配置"
    echo "  2. frontend/config.yaml 中的 allowed_dev_origins 配置"
    echo "如果无法访问数据库，请检查用户名/密码："
    echo "  1. backend/config/private.yaml 中的 database 配置"
}

# Run main function
main