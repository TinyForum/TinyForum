# ============================================
# Global variables
# ============================================

# ----- 项目配置 -----
PROJECT_ROOT=$(cd "$(dirname "$0")/../" && pwd)  # 项目根目录
PROJECT_NAME="tinyforum"                        # 项目名称
PROJECT_VERSION="0.0.1"                          # 项目版本

# ----- 系统信息 -----
OS=$(uname -s)                           # 操作系统
CURRENT_USER=$(whoami)                   # 当前用户

# ----- PostgreSQL 配置 -----
PSQL_HOST="localhost"                    # 数据库地址
PSQL_PORT="5432"                         # 数据库端口
PSQL_DB_NAME="tiny_forum"                # 默认创建数据库名
PSQL_USER="tinyforum"                    # 默认创建用户名
PSQL_PASS="tf@password"                  # 默认密码

# ----- Redis 配置 -----
REDIS_HOST="localhost"                   # 数据库地址
REDIS_PORT="6379"                        # 数据库端口
REDIS_USER="tinyforum"                   # 默认创建用户名
REDIS_PASSWORD="tf@password"             # 默认创建用户密码

# ----- Node.js 配置 -----
NODE_ENV="development"
PACKAGE_MANAGER=""                       # 运行时自动检测（npm/pnpm/yarn）

# ----- 网络地址（用于生成前端/后端访问 URL）-----
LOCAL_IP="localhost"                    # 本机 IP

# 获取本机 IP（兼容 Linux / macOS）
if command -v hostname >/dev/null 2>&1 && hostname -I 2>/dev/null | grep -q '[0-9]'; then
    LOCAL_IP=$(hostname -I | awk '{print $1}')
elif command -v ip >/dev/null 2>&1; then
    LOCAL_IP=$(ip route get 1 | awk '{print $7; exit}' 2>/dev/null)
else
    # macOS 或其它：使用 ifconfig 获取第一个活跃的 IPv4 地址
    LOCAL_IP=$(ifconfig | grep -E 'inet (addr:)?([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)' | grep -v '127.0.0.1' | head -n1 | awk '{print $2}' | cut -d: -f2)
fi

# 前端配置
FRONTEND_PORT="3000"                     # 前端运行端口
FRONTEND_URL="http://$LOCAL_IP:$FRONTEND_PORT"  # 前端访问 URL

# 后端配置
BACKEND_PORT="8080"                      # 后端运行端口
BACKEND_URL="http://$LOCAL_IP:$BACKEND_PORT"    # 后端访问 URL