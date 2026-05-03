# ============================================
# Module 5: Backend Setup
# ============================================

# ── 内部工具 ─────────────────────────────────────────────────────────────────


_write_config() {
    local filename="$1"
    local content="$2"
    local config_dir="backend/config"
    mkdir -p "$config_dir"
    printf '%s\n' "$content" > "${config_dir}/${filename}"
    echo -e "${GREEN}   📝 Generated ${config_dir}/${filename}${NC}"
}

_get_default() {
    case "$1" in
        server_host)        echo "${BACKEND_HOST:-localhost}" ;;
        server_port)        echo "${BACKEND_PORT:-8080}" ;;
        server_mode)        echo "${BACKEND_MODE:-debug}" ;;
        frontend_proto)     echo "${FRONTEND_PROTOCOL:-http}" ;;
        frontend_host)      echo "${FRONTEND_HOST:-localhost}" ;;
        frontend_port)      echo "${FRONTEND_PORT:-3000}" ;;
        api_proto)          echo "${API_PROTOCOL:-http}" ;;
        api_host)           echo "${API_HOST:-localhost}" ;;
        api_port)           echo "${API_PORT:-8080}" ;;
        log_level)          echo "${LOG_LEVEL:-info}" ;;
        ollama_url)         echo "${OLLAMA_BASE_URL:-http://localhost:11434}" ;;
        ollama_model)       echo "${OLLAMA_MODEL:-qwen3:0.6b}" ;;
        llamacpp_url)       echo "${LLAMACPP_BASE_URL:-http://localhost:8080}" ;;
        llamacpp_model)     echo "${LLAMACPP_MODEL:-llama.cpp}" ;;
        email_host)         echo "${EMAIL_HOST:-smtp.163.com}" ;;
        email_port)         echo "${EMAIL_PORT:-465}" ;;
        email_user)         echo "${EMAIL_USERNAME:-xxxxx@xxxx.com}" ;;
        email_pass)         echo "${EMAIL_PASSWORD:-password}" ;;
        email_from)         echo "${EMAIL_FROM:-${EMAIL_USERNAME:-xxxxx@xxxx.com}}" ;;
        email_from_name)    echo "${EMAIL_FROM_NAME:-TinyForum}" ;;
        email_ssl)          echo "${EMAIL_SSL:-false}" ;;
        email_tls)          echo "${EMAIL_TLS:-true}" ;;
        jwt_secret)         echo "${JWT_SECRET:-tiny-forum-secret-change-in-production-32chars}" ;;
        jwt_expire)         echo "${JWT_EXPIRE:-24h}" ;;
        jwt_refresh)        echo "${JWT_REFRESH_EXPIRE:-168h}" ;;
        jwt_issuer)         echo "${JWT_ISSUER:-tiny-forum}" ;;
        db_host)            echo "${PSQL_HOST:-localhost}" ;;
        db_port)            echo "${PSQL_PORT:-5432}" ;;
        db_user)            echo "${PG_FINAL_USER:-${PSQL_USER:-tinyforum}}" ;;
        db_pass)            echo "${PG_FINAL_PASS:-${PSQL_PASS:-tf@password}}" ;;
        db_name)            echo "${PSQL_DB_NAME:-tiny_forum}" ;;
        db_sslmode)         echo "${DB_SSLMODE:-disable}" ;;
        db_timezone)        echo "${DB_TIMEZONE:-Asia/Shanghai}" ;;
        redis_host)         echo "${REDIS_HOST:-localhost}" ;;
        redis_port)         echo "${REDIS_PORT:-6379}" ;;
        redis_user)         echo "${REDIS_FINAL_USER:-${REDIS_USER:-tinyforum}}" ;;
        redis_pass)         echo "${REDIS_FINAL_PASS:-${REDIS_PASSWORD:-tf@password}}" ;;
        redis_db)           echo "${REDIS_DB:-0}" ;;
        admin_email)        echo "${ADMIN_EMAIL:-admin@test.com}" ;;
        admin_pass)         echo "${ADMIN_PASSWORD:-password}" ;;
        admin_user)         echo "${ADMIN_USERNAME:-admin}" ;;
        admin_role)         echo "${ADMIN_ROLE:-super_admin}" ;;
        admin_score)        echo "${ADMIN_SCORE:-10000}" ;;
        upload_dir)         echo "${UPLOAD_DIR:-uploads}" ;;
        upload_url_prefix)  echo "${UPLOAD_URL_PREFIX:-upload}" ;;
        upload_allowed_ext) echo "${UPLOAD_ALLOWED_EXT:-png,jpg,jpeg,gif,mp4,webm,mp3,avi,mkv}" ;;
        config_version)     echo "${CONFIG_VERSION:-1.0.0}" ;;
    esac
}

# ── risk_control.yml ─────────────────────────────────────────────────────────
write_risk_control_yml() {
    local config_content
    config_content=$(cat <<EOF
# risk_control.yml
rate_limit:
  enabled: true
  risk_control_levels:
    normal:
      create_post:    {limit: 20, window: 1h} # 每小时最多创建 20 个帖子
      create_comment: {limit: 60, window: 1h} # 每小时最多创建 60 个评论
      send_report:    {limit: 10, window: 1h} # 每小时最多发送 10 个举报
      update_profile: {limit:  5, window: 1h} # 每小时最多更新 5 次个人资料
    observe: 
      create_post:    {limit:  5, window: 1h} # 每小时最多创建 5 个帖子
      create_comment: {limit: 20, window: 1h} # 每小时最多创建 20 个评论
      send_report:    {limit:  5, window: 1h} # 每小时最多发送 5 个举报
      update_profile: {limit:  3, window: 1h} # 每小时最多更新 3 次个人资料
    restrict:
      create_post:    {limit:  2, window: 1h} # 每小时最多创建 2 个帖子
      create_comment: {limit:  5, window: 1h} # 每小时最多创建 5 个评论
      send_report:    {limit:  2, window: 1h} # 每小时最多发送 2 个举报
      update_profile: {limit:  1, window: 1h} # 每小时最多更新 1 次个人资料

ip_whitelist:
  - "127.0.0.1" # 本地回环地址
  - "::1" # IPv6 本地回环地址
  - "localhost" # 本地主机名
  - "10.0.0.0/8" # 私有网络地址
  - "172.16.0.0/12" # 私有网络地址
  - "${LOCAL_IP}" # 本地 IP 地址
EOF
)
    _write_config "risk_control.yml" "$config_content"
}


# ── basic.yml ────────────────────────────────────────────────────────────────
write_basic_yaml() {
    local server_host; server_host=$(_get_default server_host)
    local server_port; server_port=$(_get_default server_port)
    local server_mode; server_mode=$(_get_default server_mode)
    local frontend_proto; frontend_proto=$(_get_default frontend_proto)
    local frontend_host; frontend_host=$(_get_default frontend_host)
    local frontend_port; frontend_port=$(_get_default frontend_port)
    local api_proto; api_proto=$(_get_default api_proto)
    local api_host; api_host=$(_get_default api_host)
    local api_port; api_port=$(_get_default api_port)
    local log_level; log_level=$(_get_default log_level)
    local ollama_url; ollama_url=$(_get_default ollama_url)
    local ollama_model; ollama_model=$(_get_default ollama_model)
    local llamacpp_url; llamacpp_url=$(_get_default llamacpp_url)
    local llamacpp_model; llamacpp_model=$(_get_default llamacpp_model)
    local upload_dir; upload_dir=$(_get_default upload_dir)
    local upload_url_prefix; upload_url_prefix=$(_get_default upload_url_prefix)
    local upload_allowed_ext; upload_allowed_ext=$(_get_default upload_allowed_ext)
    local config_version; config_version=$(_get_default config_version)

    # allow_origins
    local allow_origins_yaml
    if [ -n "${ALLOW_ORIGINS:-}" ]; then
        allow_origins_yaml=""
        IFS=',' read -ra _ORIGINS <<< "$ALLOW_ORIGINS"
        for _origin in "${_ORIGINS[@]}"; do
            allow_origins_yaml+="  - \"${_origin}\"\n"
        done
    else
       allow_origins_yaml=$(cat <<EOF
  - "http://localhost:3000"
  - "http://127.0.0.1:3000"
  - "http://localhost:8080"
  - "http://127.0.0.1:8080"
  - "${BACKEND_URL}"
  - "${FRONTEND_URL}"
EOF
       )
    fi

    # Note: 使用 printf 而非 heredoc 避免变量展开歧义
    local content
    content="# basic.yml — non-sensitive configuration
# Generated by TinyForum setup script
# DO NOT EDIT MANUALLY — run 'make init-dev' to regenerate
version: ${config_version}

server:
  host: ${server_host}
  port: ${server_port}
  mode: ${server_mode}
  read_timeout: 30s
  write_timeout: 30s
  max_header_bytes: 1048576

frontend:
  protocol: ${frontend_proto}
  host: ${frontend_host}
  port: ${frontend_port}

api:
  protocol: ${api_proto}
  host: ${api_host}
  port: ${api_port}
  version: v1
  prefix: /api

log:
  level: ${log_level}
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
  base_url: ${ollama_url}
  model: ${ollama_model}
  num_predict: 256
  temperature: 0.7
  timeout: 60

llamacpp:
  base_url: ${llamacpp_url}
  model: ${llamacpp_model}
  num_predict: 256
  temperature: 0.7
  timeout: 60

allow_origins:
${allow_origins_yaml}

upload:
  upload_dir: ${upload_dir}
  url_prefix: /${upload_url_prefix}
  allowed_ext: ${upload_allowed_ext}
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
    local_path: ./uploads"

    _write_config "basic.yml" "$content"
}

# ── private.yml ──────────────────────────────────────────────────────────────
write_private_yaml() {
    local email_host; email_host=$(_get_default email_host)
    local email_port; email_port=$(_get_default email_port)
    local email_user; email_user=$(_get_default email_user)
    local email_pass; email_pass=$(_get_default email_pass)
    local email_from; email_from=$(_get_default email_from)
    local email_from_name; email_from_name=$(_get_default email_from_name)
    local email_ssl; email_ssl=$(_get_default email_ssl)
    local email_tls; email_tls=$(_get_default email_tls)
    local jwt_secret; jwt_secret=$(_get_default jwt_secret)
    local jwt_expire; jwt_expire=$(_get_default jwt_expire)
    local jwt_refresh; jwt_refresh=$(_get_default jwt_refresh)
    local jwt_issuer; jwt_issuer=$(_get_default jwt_issuer)
    local db_host; db_host=$(_get_default db_host)
    local db_port; db_port=$(_get_default db_port)
    local db_user; db_user=$(_get_default db_user)
    local db_pass; db_pass=$(_get_default db_pass)
    local db_name; db_name=$(_get_default db_name)
    local db_sslmode; db_sslmode=$(_get_default db_sslmode)
    local db_timezone; db_timezone=$(_get_default db_timezone)
    local redis_host; redis_host=$(_get_default redis_host)
    local redis_port; redis_port=$(_get_default redis_port)
    local redis_user; redis_user=$(_get_default redis_user)
    local redis_pass; redis_pass=$(_get_default redis_pass)
    local redis_db; redis_db=$(_get_default redis_db)
    local admin_email; admin_email=$(_get_default admin_email)
    local admin_pass; admin_pass=$(_get_default admin_pass)
    local admin_user; admin_user=$(_get_default admin_user)
    local admin_role; admin_role=$(_get_default admin_role)
    local admin_score; admin_score=$(_get_default admin_score)

    local content
    content="# private.yml — SENSITIVE configuration
# Generated by TinyForum setup script
# DO NOT EDIT MANUALLY — run 'make init-dev' to regenerate
# ⚠️  This file is in .gitignore and should NEVER be committed.

email:
  host: \"${email_host}\"
  port: ${email_port}
  username: \"${email_user}\"
  password: \"${email_pass}\"
  from: \"${email_from}\"
  from_name: \"${email_from_name}\"
  ssl: ${email_ssl}
  tls: ${email_tls}
  pool_size: 5

jwt:
  secret: \"${jwt_secret}\"
  expire: \"${jwt_expire}\"
  refresh_expire: \"${jwt_refresh}\"
  issuer: \"${jwt_issuer}\"

database:
  host: ${db_host}
  port: ${db_port}
  user: \"${db_user}\"
  password: \"${db_pass}\"
  dbname: ${db_name}
  sslmode: ${db_sslmode}
  timezone: ${db_timezone}

redis:
  host: ${redis_host}
  port: ${redis_port}
  user: \"${redis_user}\"
  password: \"${redis_pass}\"
  db: ${redis_db}

admin:
  email: ${admin_email}
  password: ${admin_pass}
  username: ${admin_user}
  role: ${admin_role}
  score: ${admin_score}"

    _write_config "private.yml" "$content"
}

# ── 主函数 ───────────────────────────────────────────────────────────────────

setup_backend_config() {
    echo ""
    echo -e "${BOLD}━━━ Backend Config Generation ━━━━━━━━━━━━━━━━━━━━━${NC}"
    write_risk_control_yml
    write_basic_yaml
    write_private_yaml
    echo -e "${GREEN}✅ Backend configuration files generated.${NC}"
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

setup_backend() {
    echo ""
    echo -e "${BOLD}━━━ Backend Dependencies ━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

    # Bug fix: 确保始终在 PROJECT_ROOT 下操作，用子 shell 隔离 cd
    (
      echo -e "${YELLOW}   cd to ${PROJECT_ROOT}...${NC}"
      
        cd $PROJECT_ROOT || { echo -e "${RED}❌ Cannot cd to PROJECT_ROOT${NC}"; exit 1; }

        setup_backend_config
        pwd 

        cd backend || { echo -e "${RED}❌ 'backend' directory not found${NC}"; exit 1; }

        if [ ! -f "go.mod" ]; then
            echo -e "${RED}❌ go.mod not found in ${PROJECT_BACKEND}. Are you in the right directory?${NC}"
            exit 1
        fi

        echo "   Running go mod tidy..."
        if go mod tidy; then
            echo -e "${GREEN}   ✅ Go dependencies resolved.${NC}"
        else
            echo -e "${RED}❌ go mod tidy failed.${NC}"
            exit 1
        fi
    ) || return 1

    echo -e "${GREEN}✅ Backend setup completed.${NC}"
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}