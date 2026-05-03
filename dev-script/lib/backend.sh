# ============================================
# Module 5: Backend Setup (Optimized)
# ============================================

# 辅助函数：写入配置文件
_write_config() {
    local filename=$1
    local content=$2
    local config_dir="backend/config"
    mkdir -p "$config_dir"
    echo "$content" > "$config_dir/$filename"
    echo -e "${GREEN}   📝 Generated $config_dir/$filename${NC}"
}

# 获取配置默认值（集中管理）
_get_default() {
    case "$1" in
        server_host)     echo "${BACKEND_HOST:-localhost}" ;;
        server_port)     echo "${BACKEND_PORT:-8080}" ;;
        server_mode)     echo "${BACKEND_MODE:-debug}" ;;
        frontend_proto)  echo "${FRONTEND_PROTOCOL:-http}" ;;
        frontend_host)   echo "${FRONTEND_HOST:-localhost}" ;;
        frontend_port)   echo "${FRONTEND_PORT:-3000}" ;;
        api_proto)       echo "${API_PROTOCOL:-http}" ;;
        api_host)        echo "${API_HOST:-localhost}" ;;
        api_port)        echo "${API_PORT:-8080}" ;;
        log_level)       echo "${LOG_LEVEL:-info}" ;;
        ollama_url)      echo "${OLLAMA_BASE_URL:-http://localhost:11434}" ;;
        ollama_model)    echo "${OLLAMA_MODEL:-qwen3:0.6b}" ;;
        llamacpp_url)    echo "${LLAMACPP_BASE_URL:-http://localhost:8080}" ;;
        llamacpp_model)  echo "${LLAMACPP_MODEL:-llama.cpp}" ;;
        email_host)      echo "${EMAIL_HOST:-smtp.163.com}" ;;
        email_port)      echo "${EMAIL_PORT:-465}" ;;
        email_user)      echo "${EMAIL_USERNAME:-xxxxx@xxxx.com}" ;;
        email_pass)      echo "${EMAIL_PASSWORD:-password}" ;;
        email_from)      echo "${EMAIL_FROM:-$(_get_default email_user)}" ;;
        email_ssl)       echo "${EMAIL_SSL:-false}" ;;
        email_tls)       echo "${EMAIL_TLS:-true}" ;;
        jwt_secret)      echo "${JWT_SECRET:-tiny-forum-secret-change-in-production-32chars}" ;;
        jwt_expire)      echo "${JWT_EXPIRE:-24h}" ;;
        jwt_refresh)     echo "${JWT_REFRESH_EXPIRE:-168h}" ;;
        jwt_issuer)      echo "${JWT_ISSUER:-tiny-forum}" ;;
        db_host)         echo "${PSQL_HOST:-localhost}" ;;
        db_port)         echo "${PSQL_PORT:-5432}" ;;
        db_user)         echo "${PSQL_USER:-caoyang}" ;;
        db_pass)         echo "${PSQL_PASS:-your_database_password}" ;;
        db_name)         echo "${PSQL_DB_NAME:-tiny_forum}" ;;
        db_sslmode)      echo "${DB_SSLMODE:-disable}" ;;
        db_timezone)     echo "${DB_TIMEZONE:-Asia/Shanghai}" ;;
        redis_host)      echo "${REDIS_HOST:-localhost}" ;;
        redis_port)      echo "${REDIS_PORT:-6379}" ;;
        redis_pass)      echo "${REDIS_PASSWORD:-}" ;;
        redis_db)        echo "${REDIS_DB:-0}" ;;
        admin_email)     echo "${ADMIN_EMAIL:-admin@test.com}" ;;
        admin_pass)      echo "${ADMIN_PASSWORD:-password}" ;;
        admin_user)      echo "${ADMIN_USERNAME:-admin}" ;;
        admin_role)      echo "${ADMIN_ROLE:-super_admin}" ;;
        admin_score)     echo "${ADMIN_SCORE:-10000}" ;;
        upload_dir)      echo "${UPLOAD_DIR:-uploads}" ;;
        upload_url_prefix)     echo "${UPLOAD_URL_PREFIX:-upload}" ;;
        upload_allowed_ext) echo "${UPLOAD_ALLOWED_EXT:-png,jpg,jpeg,gif,mp4,webm,mp3,avi,mkv}" ;;
        config_version) echo "${CONFIG_VERSION:-1.0.0}" ;;
        email_from_name) echo "${EMAIL_FROM_NAME:-TinyForum}" ;;
    esac
}

# 1. 生成 risk_control.yml（风控配置）—— 静态内容，无动态变量
# MARK: Risk Control
write_risk_control_yml() {
    _write_config "risk_control.yml" '# risk_control_control
rate_limit:
  enabled: true
  risk_control_levels:
    normal:
      create_post: {limit: 20, window: 1h}
      create_comment: {limit: 60, window: 1h}
      send_report: {limit: 10, window: 1h}
      update_profile: {limit: 5, window: 1h}
    observe:
      create_post: {limit: 5, window: 1h}
      create_comment: {limit: 20, window: 1h}
      send_report: {limit: 5, window: 1h}
      update_profile: {limit: 3, window: 1h}
    restrict:
      create_post: {limit: 2, window: 1h}
      create_comment: {limit: 5, window: 1h}
      send_report: {limit: 2, window: 1h}
      update_profile: {limit: 1, window: 1h}

# IP 白名单
ip_whitelist:
  - "127.0.0.1"
  - "::1"
  - "localhost"
  - "192.168.1.100"
  - "10.0.0.0/8"
  - "172.16.0.0/12"'
}

# 2. 生成 basic.yml（非敏感配置）
# MARK: Basic
write_basic_yaml() {
    local server_host=$(_get_default server_host)
    local server_port=$(_get_default server_port)
    local server_mode=$(_get_default server_mode)
    local frontend_proto=$(_get_default frontend_proto)
    local frontend_host=$(_get_default frontend_host)
    local frontend_port=$(_get_default frontend_port)
    local api_proto=$(_get_default api_proto)
    local api_host=$(_get_default api_host)
    local api_port=$(_get_default api_port)
    local log_level=$(_get_default log_level)
    local ollama_url=$(_get_default ollama_url)
    local ollama_model=$(_get_default ollama_model)
    local llamacpp_url=$(_get_default llamacpp_url)
    local llamacpp_model=$(_get_default llamacpp_model)
    local upload_dir=$(_get_default upload_dir)
    local upload_url_prefix=$(_get_default upload_url_prefix)
    local upload_allowed_ext=$(_get_default upload_allowed_ext)
    local config_version=$(_get_default config_version)

    # 处理 allow_origins
    local allow_origins_str
    if [ -n "$ALLOW_ORIGINS" ]; then
        allow_origins_str=""
        IFS=',' read -ra ORIGINS <<< "$ALLOW_ORIGINS"
        for origin in "${ORIGINS[@]}"; do
            allow_origins_str+="  - \"$origin\"\n"
        done
    else
        allow_origins_str='  - "http://localhost:3000"
  - "http://127.0.0.1:3000"
  - "http://localhost:8080"
  - "http://127.0.0.1:8080"'
    fi

    local content="
# basic
version: $config_version # 配置版本
server:
  host: $server_host
  port: $server_port
  mode: $server_mode
  read_timeout: 30s
  write_timeout: 30s
  max_header_bytes: 1048576

frontend:
  protocol: $frontend_proto
  host: $frontend_host
  port: $frontend_port

api:
  protocol: $api_proto
  host: $api_host
  port: $api_port
  version: v1
  prefix: /api

log:
  level: $log_level
  filename: ./logs/app.log
  max_size: 100
  max_backups: 10
  max_age: 30
  compress: true
  console: true
  json_format: false
  db:
    dsn: ./logs/log.db
    max_buffer: 1024
    batch_size: 100
    flush_every: 5s
    retention: 7

ollama:
  base_url: $ollama_url
  model: $ollama_model
  num_predict: 256
  temperature: 0.7
  timeout: 60

llamacpp:
  base_url: $llamacpp_url
  model: $llamacpp_model
  num_predict: 256
  temperature: 0.7
  timeout: 60

allow_origins:
$allow_origins_str

upload:
  upload_dir: $upload_dir
  url_prefix: /$upload_url_prefix
  allowed_ext: $upload_allowed_ext
  max_size: 10485760
  allowed_types:
    - .jpg
    - .jpeg
    - .png
    - .gif
    - .webp
    - .pdf
    - .zip
  avatar:
    max_size: 2097152
    allowed_types: [.jpg, .jpeg, .png, .webp]
    width: 512
    height: 512
  post_image:
    max_size: 5242880
    allowed_types: [.jpg, .jpeg, .png, .gif, .webp]
  storage:
    type: local
    local_path: ./uploads
    # url_prefix: /uploads/
"
    _write_config "basic.yml" "$content"
}

# 3. 生成 private.yml（敏感配置）
# MARK: Private
write_private_yaml() {
    local email_host=$(_get_default email_host)
    local email_port=$(_get_default email_port)
    local email_user=$(_get_default email_user)
    local email_pass=$(_get_default email_pass)
    local email_from=$(_get_default email_from)
    local email_ssl=$(_get_default email_ssl)
    local email_tls=$(_get_default email_tls)

    local jwt_secret=$(_get_default jwt_secret)
    local jwt_expire=$(_get_default jwt_expire)
    local jwt_refresh=$(_get_default jwt_refresh)
    local jwt_issuer=$(_get_default jwt_issuer)

    local db_host=$(_get_default db_host)
    local db_port=$(_get_default db_port)
    local db_user=$(_get_default db_user)
    local db_pass=$(_get_default db_pass)
    local db_name=$(_get_default db_name)
    local db_sslmode=$(_get_default db_sslmode)
    local db_timezone=$(_get_default db_timezone)

    local redis_host=$(_get_default redis_host)
    local redis_port=$(_get_default redis_port)
    local redis_pass=$(_get_default redis_pass)
    local redis_db=$(_get_default redis_db)

    local admin_email=$(_get_default admin_email)
    local admin_pass=$(_get_default admin_pass)
    local admin_user=$(_get_default admin_user)
    local admin_role=$(_get_default admin_role)
    local admin_score=$(_get_default admin_score)
    local email_from_name=$(_get_default email_from_name)

    local content="
# private
email:
  host: \"$email_host\"
  port: $email_port
  username: \"$email_user\"
  password: \"$email_pass\"
  from: \"$email_from\"
  from_name: \"$email_from_name\"
  ssl: $email_ssl
  tls: $email_tls
  pool_size: 5

jwt:
  secret: \"$jwt_secret\"
  expire: \"$jwt_expire\"
  refresh_expire: \"$jwt_refresh\"
  issuer: \"$jwt_issuer\"

database:
  host: $db_host
  port: $db_port
  user: $db_user
  password: \"$db_pass\"
  dbname: $db_name
  sslmode: $db_sslmode
  timezone: $db_timezone

redis:
  host: $redis_host
  port: $redis_port
  password: \"$redis_pass\"
  db: $redis_db

admin:
  email: $admin_email
  password: $admin_pass
  username: $admin_user
  role: $admin_role
  score: $admin_score
"
    _write_config "private.yml" "$content"
}

# 主函数：生成所有后端配置
setup_backend_config() {
    echo ""
    echo "Step: Generating backend configuration files..."
    write_risk_control_yml
    write_basic_yaml
    write_private_yaml
    echo -e "${GREEN}✅ Backend configuration files generated successfully${NC}"
}

# 后端环境搭建（依赖安装）
setup_backend() {
    echo ""
    echo "Step 2: Setting up Backend..."
    cd "$PROJECT_ROOT" || exit 1
    setup_backend_config
    cd backend || exit 1
    if [ ! -f "go.mod" ]; then
        echo -e "${RED}❌ go.mod not found. Are you in the right directory?${NC}"
        exit 1
    fi
    echo "  Running go mod tidy..."
    go mod tidy
    echo -e "${GREEN}  Dependencies downloaded.${NC}"
    cd "$PROJECT_ROOT" || exit 1
}