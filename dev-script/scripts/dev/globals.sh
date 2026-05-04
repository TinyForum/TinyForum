# ============================================
# Module: Global Variables
# ============================================

# ----- 项目配置 -----
PROJECT_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../" && pwd)
PROJECT_NAME="tinyforum"
PROJECT_VERSION="0.0.1"
PROJECT_BACKEND="${PROJECT_ROOT}/backend"
PROJECT_FRONTEND="${PROJECT_ROOT}/frontend"
# logo
BANNER=$(cat $PROJECT_ROOT/dev-script/scripts/dev/banner.txt)

# ----- 配置文件位置 -----
RISK_CONTROL_CONFIG_PATH="${PROJECT_ROOT}/backend/config/risk_control.yml"
PRIVATE_CONFIG_PATH="${PROJECT_ROOT}/backend/config/private.yml"
BASIC_CONFIG_PATH="${PROJECT_ROOT}/backend/config/basic.yml"
POSTGRES_CONFIG_PATH="${PROJECT_ROOT}/backend/config/postgres.yml"
REDIS_CONFIG_PATH="${PROJECT_ROOT}/backend/config/redis.yml"
FRONTEND_CONFIG_PATH="${PROJECT_ROOT}/frontend/config.yml"


# ----- 系统信息 -----
OS=$(uname -s)
CURRENT_USER=$(whoami)

# ----- PostgreSQL 配置 -----
PSQL_HOST="localhost"
PSQL_PORT="5432"
PSQL_DB_NAME="tiny_forum"
PSQL_USER="tinyforum"
PSQL_PASS="tf@password"

# ----- Redis 配置 -----
REDIS_HOST="localhost"
REDIS_PORT="6379"
REDIS_USER="tinyforum"
REDIS_PASSWORD="tf@password"

# ----- Node.js 配置 -----
NODE_ENV="development"
PACKAGE_MANAGER=""   # 运行时自动检测

# ----- 本机 IP（兼容 Linux / macOS）-----
# Bug fix: macOS 无 hostname -I，改用 ipconfig getifaddr / ifconfig 兼容方案
_detect_local_ip() {
    local ip=""
    case "$(uname -s)" in
        Linux*)
            # 优先 ip route（最可靠）
            if command -v ip >/dev/null 2>&1; then
                ip=$(ip route get 1.1.1.1 2>/dev/null | awk '{for(i=1;i<=NF;i++) if($i=="src") print $(i+1)}')
            fi
            # fallback: hostname -I
            if [ -z "$ip" ] && command -v hostname >/dev/null 2>&1; then
                ip=$(hostname -I 2>/dev/null | awk '{print $1}')
            fi
            ;;
        Darwin*)
            # macOS：通过 route 找出口网卡，再用 ipconfig 取地址
            local iface
            iface=$(route -n get default 2>/dev/null | awk '/interface:/{print $2}')
            if [ -n "$iface" ]; then
                ip=$(ipconfig getifaddr "$iface" 2>/dev/null)
            fi
            # fallback: ifconfig
            if [ -z "$ip" ]; then
                ip=$(ifconfig 2>/dev/null \
                    | awk '/inet /{print $2}' \
                    | grep -v '^127\.' \
                    | head -n1)
            fi
            ;;
    esac
    # 最终 fallback
    echo "${ip:-localhost}"
}

LOCAL_IP=$(_detect_local_ip)

# ----- 前端 / 后端 URL -----
FRONTEND_PORT="3000"
BACKEND_PORT="8080"
FRONTEND_URL="http://${LOCAL_IP}:${FRONTEND_PORT}"
BACKEND_URL="http://${LOCAL_IP}:${BACKEND_PORT}"
