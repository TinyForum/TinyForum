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

    check_redis_client
    
    # Module 4: Configuration Setup
    setup_configurations
    
    # Module 5: Backend Setup
    setup_backend
    
    # Module 6: Frontend Setup
    setup_frontend
    
    # Module 7: Summary
    print_summary
    echo "如果无法访问后端，请检查 cors 配置"
    echo "1. 检查 backend/config/basic.yaml 中的 allow_origins 配置"
    echo "2. 检查 frontend/config.yaml 中的 allowedDevOrigins 配置"
    echo "如果无法访问数据库，请检查用户名是否正确，或者尝试新建用户，"
    echo "1. 检查 backend/config/private.yaml 中的 database 配置"
}

# Run main function
main