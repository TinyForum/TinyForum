
#!/usr/bin/env bash
# config_manager.sh
# 用于管理 TinyForum 配置的工具脚本
# 用法：
#   source config_manager.sh           # 加载函数和默认配置
#   load_config_from_env [.env文件]    # 从 .env 更新配置
#   get_config <配置键>                # 获取配置值，如 get_config BASIC_SERVER_PORT
#   check_config_format                # 检查所有配置值是否有效（非空）

# ----------------------------- 默认配置定义 -----------------------------
# 所有键名格式：FILENAME_KEY_SUBKEY（全大写）
# 值被提取自提供的 YAML 文件，数组类型转为逗号分隔字符串

# basic.yml
TINYFORUM_VERSION="0.0.1"
BASIC_SERVER_HOST="localhost"
BASIC_SERVER_PORT="8080"
BASIC_SERVER_MODE="debug"
BASIC_SERVER_READ_TIMEOUT="30s"
BASIC_SERVER_WRITE_TIMEOUT="30s"
BASIC_SERVER_MAX_HEADER_BYTES="1048576"
BASIC_FRONTEND_PROTOCOL="http"
BASIC_FRONTEND_HOST="localhost"
BASIC_FRONTEND_PORT="3000"
BASIC_API_PROTOCOL="http"
BASIC_API_HOST="localhost"
BASIC_API_PORT="8080"
BASIC_API_VERSION="v1"
BASIC_API_PREFIX="/api"
BASIC_LOG_LEVEL="info"
BASIC_LOG_FILENAME="./logs/app.log"
BASIC_LOG_MAX_SIZE="100"
BASIC_LOG_MAX_BACKUPS="10"
BASIC_LOG_MAX_AGE="30"
BASIC_LOG_COMPRESS="true"
BASIC_LOG_CONSOLE="true"
BASIC_LOG_JSON_FORMAT="false"
BASIC_LOG_DB_DSN="./logs/log.db"
BASIC_LOG_DB_MAX_BUFFER="1024"
BASIC_LOG_DB_BATCH_SIZE="100"
BASIC_LOG_DB_FLUSH_EVERY="5s"
BASIC_LOG_DB_RETENTION="7"
BASIC_OLLAMA_BASE_URL="http://localhost:11434"
BASIC_OLLAMA_MODEL="qwen3:0.6b"
BASIC_OLLAMA_NUM_PREDICT="256"
BASIC_OLLAMA_TEMPERATURE="0.7"
BASIC_OLLAMA_TIMEOUT="60"
BASIC_LLAMACPP_BASE_URL="http://localhost:8080"
BASIC_LLAMACPP_MODEL="llama.cpp"
BASIC_LLAMACPP_NUM_PREDICT="256"
BASIC_LLAMACPP_TEMPERATURE="0.7"
BASIC_LLAMACPP_TIMEOUT="60"
BASIC_ALLOW_ORIGINS="http://localhost:3000,http://127.0.0.1:3000,http://localhost:8080,http://127.0.0.1:8080"
BASIC_UPLOAD_UPLOAD_DIR="uploads"
BASIC_UPLOAD_URL_PREFIX="/upload"
BASIC_UPLOAD_ALLOWED_EXT="png,jpg,jpeg,gif,mp4,webm,mp3,avi,mkv"
BASIC_UPLOAD_MAX_SIZE="10485760"
BASIC_UPLOAD_ALLOWED_TYPES=".jpg,.jpeg,.png,.gif,.webp,.pdf,.zip"
BASIC_UPLOAD_AVATAR_MAX_SIZE="2097152"
BASIC_UPLOAD_AVATAR_ALLOWED_TYPES=".jpg,.jpeg,.png,.webp"
BASIC_UPLOAD_AVATAR_WIDTH="512"
BASIC_UPLOAD_AVATAR_HEIGHT="512"
BASIC_UPLOAD_POST_IMAGE_MAX_SIZE="5242880"
BASIC_UPLOAD_POST_IMAGE_ALLOWED_TYPES=".jpg,.jpeg,.png,.gif,.webp"
BASIC_UPLOAD_STORAGE_TYPE="local"
BASIC_UPLOAD_STORAGE_LOCAL_PATH="./uploads"

# private.yml
PRIVATE_EMAIL_HOST="smtp.163.com"
PRIVATE_EMAIL_PORT="587"
PRIVATE_EMAIL_USERNAME="xxxxx@xxxx.com"
PRIVATE_EMAIL_PASSWORD="password"
PRIVATE_EMAIL_FROM="xxxxx@xxxx.com"
PRIVATE_EMAIL_FROM_NAME="TinyForum"
PRIVATE_EMAIL_SSL="false"
PRIVATE_EMAIL_TLS="true"
PRIVATE_EMAIL_POOL_SIZE="5"
PRIVATE_JWT_SECRET="tiny-forum-secret-change-in-production-32chars"
PRIVATE_JWT_EXPIRE="24h"
PRIVATE_JWT_REFRESH_EXPIRE="168h"
PRIVATE_JWT_ISSUER="tiny-forum"
PRIVATE_DATABASE_HOST="localhost"
PRIVATE_DATABASE_PORT="5432"
PRIVATE_DATABASE_USER="caoyang"
PRIVATE_DATABASE_PASSWORD="tf@password"
PRIVATE_DATABASE_DBNAME="tiny_forum"
PRIVATE_DATABASE_SSLMODE="disable"
PRIVATE_DATABASE_TIMEZONE="Asia/Shanghai"
PRIVATE_REDIS_HOST="localhost"
PRIVATE_REDIS_PORT="6379"
PRIVATE_REDIS_USER="tinyforum"
PRIVATE_REDIS_PASSWORD="tf@password"
PRIVATE_REDIS_DB="0"
PRIVATE_ADMIN_EMAIL="admin@test.com"
PRIVATE_ADMIN_PASSWORD="password"
PRIVATE_ADMIN_USERNAME="admin"
PRIVATE_ADMIN_ROLE="super_admin"
PRIVATE_ADMIN_SCORE="10000"

# redis.yml (单独配置块)
REDIS_HOST="localhost"
REDIS_PORT="6379"
REDIS_USER="tinyforum"
REDIS_PASSWORD="tf@password"
REDIS_DB="0"
REDIS_POOL_SIZE="10"
REDIS_MIN_IDLE_CONNS="2"
REDIS_DIAL_TIMEOUT="5s"
REDIS_READ_TIMEOUT="3s"
REDIS_WRITE_TIMEOUT="3s"

# risk_control.yml
RISK_CONTROL_RATE_LIMIT_ENABLED="true"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_NORMAL_CREATE_POST_LIMIT="20"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_NORMAL_CREATE_POST_WINDOW="1h"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_NORMAL_CREATE_COMMENT_LIMIT="60"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_NORMAL_CREATE_COMMENT_WINDOW="1h"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_NORMAL_SEND_REPORT_LIMIT="10"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_NORMAL_SEND_REPORT_WINDOW="1h"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_NORMAL_UPDATE_PROFILE_LIMIT="5"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_NORMAL_UPDATE_PROFILE_WINDOW="1h"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_OBSERVE_CREATE_POST_LIMIT="5"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_OBSERVE_CREATE_POST_WINDOW="1h"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_OBSERVE_CREATE_COMMENT_LIMIT="20"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_OBSERVE_CREATE_COMMENT_WINDOW="1h"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_OBSERVE_SEND_REPORT_LIMIT="5"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_OBSERVE_SEND_REPORT_WINDOW="1h"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_OBSERVE_UPDATE_PROFILE_LIMIT="3"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_OBSERVE_UPDATE_PROFILE_WINDOW="1h"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_RESTRICT_CREATE_POST_LIMIT="2"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_RESTRICT_CREATE_POST_WINDOW="1h"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_RESTRICT_CREATE_COMMENT_LIMIT="5"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_RESTRICT_CREATE_COMMENT_WINDOW="1h"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_RESTRICT_SEND_REPORT_LIMIT="2"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_RESTRICT_SEND_REPORT_WINDOW="1h"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_RESTRICT_UPDATE_PROFILE_LIMIT="1"
RISK_CONTROL_RATE_LIMIT_RISK_CONTROL_LEVELS_RESTRICT_UPDATE_PROFILE_WINDOW="1h"
RISK_CONTROL_IP_WHITELIST="127.0.0.1,::1,localhost,10.0.0.0/8,172.16.0.0/12,192.168.5.180"

# ----------------------------- 函数定义 -----------------------------

# 从 .env 文件更新配置值
# 用法: load_config_from_env [env_file]
# 环境变量格式: KEY=value (KEY 必须与默认配置中的键名一致)
load_config_from_env() {
    local env_file="${1:-.env}"
    if [[ ! -f "$env_file" ]]; then
        echo "Warning: $env_file not found, skip loading." >&2
        return 0
    fi

    while IFS= read -r line || [[ -n "$line" ]]; do
        # 跳过注释和空行
        [[ -z "$line" || "$line" =~ ^[[:space:]]*# ]] && continue
        # 匹配 KEY=value 格式
        if [[ "$line" =~ ^([A-Z_][A-Z0-9_]*)=(.*)$ ]]; then
            key="${BASH_REMATCH[1]}"
            value="${BASH_REMATCH[2]}"
            # 检查该键是否在默认配置中定义
            if declare -p "$key" &>/dev/null; then
                # 使用 eval 进行安全赋值（值可能包含空格或特殊字符）
                printf -v "$key" "%s" "$value"
                echo "Updated $key = $value" >&2
            else
                echo "Warning: unknown key '$key' in $env_file, ignored." >&2
            fi
        else
            echo "Warning: invalid line in $env_file: $line" >&2
        fi
    done < "$env_file"
}

# 获取配置值
# 用法: get_config <KEY>
# 返回: 输出配置值，若键不存在则返回非0退出码
get_config() {
    local key="$1"
    if [[ -z "$key" ]]; then
        echo "Error: missing key argument" >&2
        return 1
    fi
    if declare -p "$key" &>/dev/null; then
        echo "${!key}"
        return 0
    else
        echo "Error: unknown config key '$key'" >&2
        return 2
    fi
}

# 检查所有配置值的格式（非空检查）
# 可选：对特定键做更严格的格式校验
check_config_format() {
    local all_keys
    # 收集所有以 FILENAME_ 开头的变量名（即配置键）
    all_keys=$(compgen -A variable | grep -E '^(BASIC_|PRIVATE_|REDIS_|RISK_CONTROL_)')
    local failed=0

    echo "Checking configuration format..."
    for key in $all_keys; do
        value="${!key}"
        if [[ -z "$value" ]]; then
            echo "❌ $key is empty!"
            failed=1
        else
            # 简单示例：对某些数值字段检查是否数字
            case "$key" in
                BASIC_SERVER_PORT|BASIC_FRONTEND_PORT|BASIC_API_PORT)
                    if [[ ! "$value" =~ ^[0-9]+$ ]]; then
                        echo "❌ $key should be a number, got '$value'"
                        failed=1
                    fi
                    ;;
                PRIVATE_EMAIL_SSL|PRIVATE_EMAIL_TLS|BASIC_LOG_COMPRESS|BASIC_LOG_CONSOLE|BASIC_LOG_JSON_FORMAT)
                    if [[ ! "$value" =~ ^(true|false)$ ]]; then
                        echo "❌ $key should be true/false, got '$value'"
                        failed=1
                    fi
                    ;;
            esac
        fi
    done

    if [[ $failed -eq 0 ]]; then
        echo "✅ All configuration keys are valid."
        return 0
    else
        echo "❌ Some configuration keys have invalid values."
        return 1
    fi
}

# ----------------------------- 辅助函数（可选）-----------------------------
# 导出所有配置变量，供子进程使用（如果不希望污染环境，可以注释）
export_config() {
    local all_keys
    all_keys=$(compgen -A variable | grep -E '^(BASIC_|PRIVATE_|REDIS_|RISK_CONTROL_)')
    for key in $all_keys; do
        export "$key"
    done
}

# 演示：打印所有配置
print_all_config() {
    local all_keys
    all_keys=$(compgen -A variable | grep -E '^(BASIC_|PRIVATE_|REDIS_|RISK_CONTROL_)' | sort)
    for key in $all_keys; do
        echo "$key = ${!key}"
    done
}

# ----------------------------- 主入口（如果直接执行脚本）-----------------------------
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    echo "TinyForum Config Manager"
    echo "Usage: source $0   # to load functions and default config"
    echo "Then use: load_config_from_env, get_config, check_config_format"
    echo ""
    echo "Current configuration (default values):"
    print_all_config
fi