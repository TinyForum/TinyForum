#!/usr/bin/env bash
# =============================================================================
# lib/schema.sh — 环境变量 Schema 注册表
#
# 格式：schema::define <KEY> <default> <description> [validator_fn]
# 所有键名、默认值、说明、校验规则集中在此，其余模块按需调用。
# =============================================================================
[[ -n "${_LIB_SCHEMA_LOADED:-}" ]] && return 0
readonly _LIB_SCHEMA_LOADED=1

# shellcheck source=lib/core.sh
source "$(dirname "${BASH_SOURCE[0]}")/core.sh"

# ── 内部存储（关联数组）──────────────────────────────────────────────────────
declare -gA _SCHEMA_DEFAULT=()    # KEY → default value
declare -gA _SCHEMA_DESC=()       # KEY → human description
declare -gA _SCHEMA_VALIDATOR=()  # KEY → validator function name
declare -gA _SCHEMA_REQUIRED=()   # KEY → "1" if required
declare -gA _SCHEMA_SECRET=()     # KEY → "1" if sensitive (mask in output)
declare -ga _SCHEMA_ORDER=()      # insertion order

# ── 注册函数 ─────────────────────────────────────────────────────────────────
schema::define() {
  # schema::define KEY default "description" [validator] [required=0] [secret=0]
  local key="$1" default="$2" desc="$3"
  local validator="${4:-}"
  local required="${5:-0}"
  local secret="${6:-0}"

  _SCHEMA_DEFAULT["$key"]="$default"
  _SCHEMA_DESC["$key"]="$desc"
  _SCHEMA_VALIDATOR["$key"]="$validator"
  _SCHEMA_REQUIRED["$key"]="$required"
  _SCHEMA_SECRET["$key"]="$secret"
  _SCHEMA_ORDER+=("$key")
}

schema::keys()   { printf '%s\n' "${_SCHEMA_ORDER[@]}"; }
schema::default(){ echo "${_SCHEMA_DEFAULT[${1}]:-}"; }
schema::desc()   { echo "${_SCHEMA_DESC[${1}]:-}"; }
schema::secret() { [[ "${_SCHEMA_SECRET[${1}]:-0}" == "1" ]]; }
schema::required(){ [[ "${_SCHEMA_REQUIRED[${1}]:-0}" == "1" ]]; }
schema::validator(){ echo "${_SCHEMA_VALIDATOR[${1}]:-}"; }

# ── 校验函数库 ────────────────────────────────────────────────────────────────
validate::nonempty()  { [[ -n "$1" ]]; }
validate::port()      { [[ "$1" =~ ^[0-9]{1,5}$ ]] && (( $1 >= 1 && $1 <= 65535 )); }
validate::bool()      { [[ "$1" == "true" || "$1" == "false" ]]; }
validate::duration()  { [[ "$1" =~ ^[0-9]+(s|m|h|d)$ ]]; }
validate::url()       { [[ "$1" =~ ^https?:// ]]; }
validate::email()     { [[ "$1" =~ ^[^@]+@[^@]+\.[^@]+$ ]]; }
validate::loglevel()  { [[ "$1" =~ ^(debug|info|warn|error|fatal)$ ]]; }
validate::sslmode()   { [[ "$1" =~ ^(disable|require|verify-ca|verify-full)$ ]]; }
validate::mode()      { [[ "$1" =~ ^(debug|release|test)$ ]]; }
validate::bytes()     { [[ "$1" =~ ^[0-9]+$ ]] && (( $1 > 0 )); }
validate::secret32()  {
  # JWT secret 至少 32 字符且不是默认值
  [[ ${#1} -ge 32 ]] && [[ "$1" != "tiny-forum-secret-change-in-production-32chars" ]]
}

# =============================================================================
# Schema 定义区 — 分组，与 yml 文件结构对应
# =============================================================================

# ── [tiny forum] ─────────────────────────────────────────────────────────────────
schema::define TINYFORUM_VERSION "0.0.1" "TinyForum 应用版本号" "" 0 0

# ── [server] ─────────────────────────────────────────────────────────────────
schema::define SERVER_HOST          "localhost"  "后端监听主机"         "" 0 0
schema::define SERVER_PORT          "8080"       "后端 HTTP 端口"       "validate::port" 1 0
schema::define SERVER_MODE          "debug"      "运行模式 debug/release/test" "validate::mode" 0 0
schema::define SERVER_READ_TIMEOUT  "30s"        "读超时"               "validate::duration" 0 0
schema::define SERVER_WRITE_TIMEOUT "30s"        "写超时"               "validate::duration" 0 0

# ── [frontend] ────────────────────────────────────────────────────────────────
schema::define FRONTEND_PROTOCOL    "http"       "前端协议"             "" 0 0
schema::define FRONTEND_HOST        "localhost"  "前端主机"             "" 0 0
schema::define FRONTEND_PORT        "3000"       "前端端口"             "validate::port" 0 0

# ── [database] ───────────────────────────────────────────────────────────────
schema::define DB_HOST              "localhost"  "数据库主机"           "" 1 0
schema::define DB_PORT              "5432"       "数据库端口"           "validate::port" 1 0
schema::define DB_USER              "postgres"   "数据库用户名"         "validate::nonempty" 1 0
schema::define DB_PASSWORD          ""           "数据库密码"           "validate::nonempty" 1 1
schema::define DB_NAME              "tiny_forum" "数据库名"             "validate::nonempty" 1 0
schema::define DB_SSLMODE           "disable"    "SSL 模式"             "validate::sslmode" 0 0
schema::define DB_TIMEZONE          "Asia/Shanghai" "时区"              "" 0 0
schema::define DB_MAX_IDLE_CONNS    "10"         "最大空闲连接数"       "validate::bytes" 0 0
schema::define DB_MAX_OPEN_CONNS    "100"        "最大打开连接数"       "validate::bytes" 0 0

# ── [redis] ──────────────────────────────────────────────────────────────────
schema::define REDIS_HOST           "localhost"  "Redis 主机"           "" 1 0
schema::define REDIS_PORT           "6379"       "Redis 端口"           "validate::port" 1 0
schema::define REDIS_USER           "tinyforum"  "Redis 用户名"         "" 0 0
schema::define REDIS_PASSWORD       ""           "Redis 密码"           "" 0 1
schema::define REDIS_DB             "0"          "Redis DB 编号"        "" 0 0
schema::define REDIS_POOL_SIZE      "10"         "Redis 连接池大小"     "validate::bytes" 0 0
schema::define REDIS_MIN_IDLE_CONNS "2"          "Redis 最小空闲连接"   "" 0 0
schema::define REDIS_DIAL_TIMEOUT   "5s"         "Redis 连接超时"       "validate::duration" 0 0
schema::define REDIS_READ_TIMEOUT   "3s"         "Redis 读超时"         "validate::duration" 0 0
schema::define REDIS_WRITE_TIMEOUT  "3s"         "Redis 写超时"         "validate::duration" 0 0

# ── [jwt] ────────────────────────────────────────────────────────────────────
schema::define JWT_SECRET           "tiny-forum-secret-change-in-production-32chars" \
                                    "JWT 签名密钥 (≥32字符, 生产环境必须修改)" \
                                    "validate::nonempty" 1 1
schema::define JWT_EXPIRE           "24h"        "JWT 有效期"           "validate::duration" 0 0
schema::define JWT_REFRESH_EXPIRE   "168h"       "JWT Refresh 有效期"   "validate::duration" 0 0
schema::define JWT_ISSUER           "tiny-forum" "JWT 签发者"           "" 0 0

# ── [email] ──────────────────────────────────────────────────────────────────
schema::define EMAIL_HOST           ""           "SMTP 主机"            "" 0 0
schema::define EMAIL_PORT           "587"        "SMTP 端口"            "validate::port" 0 0
schema::define EMAIL_USERNAME       ""           "SMTP 用户名/邮箱"     "" 0 0
schema::define EMAIL_PASSWORD       ""           "SMTP 密码"            "" 0 1
schema::define EMAIL_FROM           ""           "发件人地址"           "" 0 0
schema::define EMAIL_FROM_NAME      "TinyForum"  "发件人名称"           "" 0 0
schema::define EMAIL_SSL            "false"      "启用 SSL"             "validate::bool" 0 0
schema::define EMAIL_TLS            "true"       "启用 TLS"             "validate::bool" 0 0

# ── [log] ────────────────────────────────────────────────────────────────────
schema::define LOG_LEVEL            "info"       "日志级别"             "validate::loglevel" 0 0
schema::define LOG_FILENAME         "./logs/app.log" "日志文件路径"     "" 0 0
schema::define LOG_MAX_SIZE         "100"        "日志最大大小 (MB)"    "validate::bytes" 0 0
schema::define LOG_MAX_BACKUPS      "10"         "日志最大备份数"       "" 0 0
schema::define LOG_MAX_AGE          "30"         "日志最大保留天数"     "" 0 0
schema::define LOG_COMPRESS         "true"       "压缩旧日志"           "validate::bool" 0 0
schema::define LOG_CONSOLE          "true"       "输出到控制台"         "validate::bool" 0 0
schema::define LOG_JSON_FORMAT      "false"      "JSON 格式日志"        "validate::bool" 0 0

# ── [ai / ollama] ─────────────────────────────────────────────────────────────
schema::define OLLAMA_BASE_URL      "http://localhost:11434" "Ollama API 地址" "validate::url" 0 0
schema::define OLLAMA_MODEL         "qwen3:0.6b" "Ollama 默认模型"      "" 0 0
schema::define OLLAMA_NUM_PREDICT   "256"        "Ollama 最大 token 数" "validate::bytes" 0 0
schema::define OLLAMA_TEMPERATURE   "0.7"        "Ollama 温度"          "" 0 0

# ── [llamacpp] ────────────────────────────────────────────────────────────────
schema::define LLAMACPP_BASE_URL    "http://localhost:8080" "llama.cpp API 地址" "validate::url" 0 0
schema::define LLAMACPP_MODEL       "llama.cpp"  "llama.cpp 模型名"     "" 0 0

# ── [upload] ─────────────────────────────────────────────────────────────────
schema::define UPLOAD_DIR           "uploads"    "上传目录"             "" 0 0
schema::define UPLOAD_URL_PREFIX    "/upload"    "上传 URL 前缀"        "" 0 0
schema::define UPLOAD_MAX_SIZE      "10485760"   "最大上传大小 (bytes)" "validate::bytes" 0 0

# ── [admin] ──────────────────────────────────────────────────────────────────
schema::define ADMIN_EMAIL          "admin@test.com" "管理员邮箱"       "" 0 0
schema::define ADMIN_USERNAME       "admin"      "管理员用户名"         "" 0 0
schema::define ADMIN_PASSWORD       "password"   "管理员密码"           "validate::nonempty" 0 1
schema::define ADMIN_ROLE           "super_admin" "管理员角色"          "" 0 0

# ── [proxy / frontend config] ─────────────────────────────────────────────────
schema::define PROXY_ENABLED        "true"       "启用前端代理"         "validate::bool" 0 0
schema::define PROXY_BACKEND_URL    "http://localhost:8080" "代理后端地址" "validate::url" 0 0

# ── [cdn] ────────────────────────────────────────────────────────────────────
schema::define CDN_DOMAIN           ""           "CDN 域名 (可选)"      "" 0 0

# ── [risk control] ───────────────────────────────────────────────────────────
schema::define RISK_CONTROL_RATE_LIMIT_ENABLED "true" "是否启用全局限流" "" 0 0

# Normal 等级限流配置
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_NORMAL_CREATE_POST_LIMIT     "20"  "Normal用户：发帖限流次数"      "" 0 0
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_NORMAL_CREATE_POST_WINDOW     "1h"  "Normal用户：发帖限流窗口"      "validate::duration" 0 0
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_NORMAL_CREATE_COMMENT_LIMIT   "60"  "Normal用户：评论限流次数"      "" 0 0
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_NORMAL_CREATE_COMMENT_WINDOW  "1h"  "Normal用户：评论限流窗口"      "validate::duration" 0 0
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_NORMAL_SEND_REPORT_LIMIT      "10"  "Normal用户：举报限流次数"      "" 0 0
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_NORMAL_SEND_REPORT_WINDOW     "1h"  "Normal用户：举报限流窗口"      "validate::duration" 0 0
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_NORMAL_UPDATE_PROFILE_LIMIT   "5"   "Normal用户：更新资料限流次数"  "" 0 0
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_NORMAL_UPDATE_PROFILE_WINDOW  "1h"  "Normal用户：更新资料限流窗口"  "validate::duration" 0 0

# Observe 等级限流配置
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_OBSERVE_CREATE_POST_LIMIT     "5"   "Observe用户：发帖限流次数"     "" 0 0
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_OBSERVE_CREATE_POST_WINDOW    "1h"  "Observe用户：发帖限流窗口"     "validate::duration" 0 0
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_OBSERVE_CREATE_COMMENT_LIMIT  "20"  "Observe用户：评论限流次数"     "" 0 0
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_OBSERVE_CREATE_COMMENT_WINDOW "1h"  "Observe用户：评论限流窗口"     "validate::duration" 0 0
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_OBSERVE_SEND_REPORT_LIMIT     "5"   "Observe用户：举报限流次数"     "" 0 0
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_OBSERVE_SEND_REPORT_WINDOW    "1h"  "Observe用户：举报限流窗口"     "validate::duration" 0 0
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_OBSERVE_UPDATE_PROFILE_LIMIT  "3"   "Observe用户：更新资料限流次数" "" 0 0
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_OBSERVE_UPDATE_PROFILE_WINDOW "1h"  "Observe用户：更新资料限流窗口" "validate::duration" 0 0

# Restrict 等级限流配置
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_RESTRICT_CREATE_POST_LIMIT     "2"   "Restrict用户：发帖限流次数"    "" 0 0
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_RESTRICT_CREATE_POST_WINDOW    "1h"  "Restrict用户：发帖限流窗口"    "validate::duration" 0 0
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_RESTRICT_CREATE_COMMENT_LIMIT  "5"   "Restrict用户：评论限流次数"    "" 0 0
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_RESTRICT_CREATE_COMMENT_WINDOW "1h"  "Restrict用户：评论限流窗口"    "validate::duration" 0 0
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_RESTRICT_SEND_REPORT_LIMIT     "2"   "Restrict用户：举报限流次数"    "" 0 0
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_RESTRICT_SEND_REPORT_WINDOW    "1h"  "Restrict用户：举报限流窗口"    "validate::duration" 0 0
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_RESTRICT_UPDATE_PROFILE_LIMIT  "1"   "Restrict用户：更新资料限流次数" "" 0 0
schema::define RISK_CONTROL_RATE_LIMIT_LEVELS_RESTRICT_UPDATE_PROFILE_WINDOW "1h"  "Restrict用户：更新资料限流窗口" "validate::duration" 0 0

# IP 白名单
schema::define RISK_CONTROL_IP_WHITELIST "127.0.0.1,::1,localhost,10.0.0.0/8,172.16.0.0/12" "IP 白名单列表（逗号分隔）" "" 0 0