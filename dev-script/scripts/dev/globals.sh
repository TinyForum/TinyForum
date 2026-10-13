# ============================================
# Module: Global Variables
# ============================================

# ----- ������� -----
PROJECT_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../" && pwd)
PROJECT_NAME="tinyforum"
PROJECT_VERSION="0.0.1"
PROJECT_BACKEND="${PROJECT_ROOT}/backend"
PROJECT_FRONTEND="${PROJECT_ROOT}/frontend"
# logo
BANNER=$(cat "$PROJECT_ROOT/dev-script/scripts/dev/banner.txt" 2>/dev/null || true)

# ----- �����ļ�λ�� -----
RISK_CONTROL_CONFIG_PATH="${PROJECT_ROOT}/backend/config/risk_control.yml"
PRIVATE_CONFIG_PATH="${PROJECT_ROOT}/backend/config/private.yml"
BASIC_CONFIG_PATH="${PROJECT_ROOT}/backend/config/basic.yml"
POSTGRES_CONFIG_PATH="${PROJECT_ROOT}/backend/config/postgres.yml"
REDIS_CONFIG_PATH="${PROJECT_ROOT}/backend/config/redis.yml"
FRONTEND_CONFIG_PATH="${PROJECT_ROOT}/frontend/config.yml"

# ----- ϵͳ��Ϣ -----
OS_RAW=$(uname -s 2>/dev/null || echo "Windows_NT")
case "$OS_RAW" in
    MINGW*|MSYS*|CYGWIN*|Windows_NT*)
        OS="Windows_NT"
        ;;
    *)
        OS="$OS_RAW"
        ;;
esac
CURRENT_USER=$(whoami 2>/dev/null || cmd /c "echo %USERNAME%" 2>/dev/null || echo "user")

# ----- PostgreSQL ���� -----
PSQL_HOST="localhost"
PSQL_PORT="5432"
PSQL_DB_NAME="tiny_forum"
PSQL_USER="tinyforum"
PSQL_PASS="tf-password"

# ----- Redis ���� -----
REDIS_HOST="localhost"
REDIS_PORT="6379"
REDIS_USER="tinyforum"
REDIS_PASSWORD="tf-password"

# ----- Node.js ���� -----
NODE_ENV="development"
PACKAGE_MANAGER=""   # ���м�ʱ�Զ����

# ----- ���� IP（��� Linux / macOS / Windows��-----
_detect_local_ip() {
    local ip=""
    case "$(uname -s 2>/dev/null)" in
        Linux*)
            if command -v ip >/dev/null 2>&1; then
                ip=$(ip route get 1.1.1.1 2>/dev/null | awk '"'"'{for(i=1;i<=NF;i++) if($i=="src") print $(i+1)}'"'"')
            fi
            if [ -z "$ip" ] && command -v hostname >/dev/null 2>&1; then
                ip=$(hostname -I 2>/dev/null | awk '"'"'{print $1}'"'"')
            fi
            ;;
        Darwin*)
            local iface
            iface=$(route -n get default 2>/dev/null | awk '"'"'/interface:/{print $2}'"'"')
            if [ -n "$iface" ]; then
                ip=$(ipconfig getifaddr "$iface" 2>/dev/null)
            fi
            if [ -z "$ip" ]; then
                ip=$(ifconfig 2>/dev/null | awk '"'"'/inet /{print $2}'"'"' | grep -v '"'"'^127.'"'"' | head -n1)
            fi
            ;;
        MINGW*|MSYS*|CYGWIN*|Windows_NT*)
            # On Windows in Git Bash, try ipconfig
            if command -v ipconfig >/dev/null 2>&1; then
                ip=$(ipconfig 2>/dev/null | grep -A 5 -B 5 "IPv4 Address" | grep -m1 "192\.168\|10\." | awk -F: '"'"'{print $2}'"'"' | tr -d ' ')
            fi
            if [ -z "$ip" ] && command -v hostname >/dev/null 2>&1; then
                ip=$(hostname -I 2>/dev/null | awk '"'"'{print $1}'"'"')
            fi
            ;;
    esac
    echo "${ip:-localhost}"
}

LOCAL_IP=$(_detect_local_ip)
BACKEND_HOST="${LOCAL_IP}"
FRONTEND_HOST="${LOCAL_IP}"

# ----- ǰ�� / ��� URL -----
FRONTEND_PORT="3000"
BACKEND_PORT="8080"
FRONTEND_URL="http://${LOCAL_IP}:${FRONTEND_PORT}"
BACKEND_URL="http://${LOCAL_IP}:${BACKEND_PORT}"