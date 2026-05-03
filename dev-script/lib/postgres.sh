#!/bin/bash
# ============================================
# Module 3: PostgreSQL Management
# ============================================
# 支持: Linux (systemd + peer auth), macOS (Homebrew)
#
# 修复清单:
#   [BUG-1] create_database_user: 多条 SQL 合并为单个 -c 调用，
#           任一语句失败（如用户已存在）导致整批失败且错误被吞。
#           修复: 拆分为独立 psql 调用，单独捕获每步错误。
#   [BUG-2] 进入"创建用户"流程前未检查 default_user 是否已存在。
#           修复: 在选项1/2创建前先 user_exists 检查，已存在则跳过 CREATE。
#   [BUG-3] _ADMIN_USER 在 setup_postgres 中声明为 local，
#           _determine_admin_user 通过 bash 动态作用域修改它（可工作，
#           但令人困惑）。修复: 改用 stdout 返回值模式，更清晰。
#   [BUG-4] >/dev/null 2>&1 吞掉所有 psql 错误，无法诊断。
#           修复: 错误时输出 psql stderr 到日志。
#   [BUG-5] grant_schema_privileges: owner 问题。PostgreSQL 15+ 默认
#           撤销 public schema 的 CREATE 权限，需额外 GRANT。
#           修复: 补充 GRANT CREATE ON SCHEMA public。

# ── 颜色（未定义时补全）────────────────────────────────────────────────────
: "${RED:=\033[0;31m}"
: "${GREEN:=\033[0;32m}"
: "${YELLOW:=\033[1;33m}"
: "${BLUE:=\033[0;34m}"
: "${BOLD:=\033[1m}"
: "${NC:=\033[0m}"

# ============================================================
# _pg_log_error: 执行 psql 并在失败时打印错误（不吞掉 stderr）
# 用法: _pg_log_error <描述> <psql命令数组>
# ============================================================
_pg_log_error() {
    local desc="$1"; shift
    local output
    # 捕获 stderr；stdout 丢弃（不污染主输出）
    if ! output=$("$@" 2>&1 >/dev/null); then
        echo -e "${RED}   ✗ ${desc} failed${NC}"
        # 只打印第一行错误，避免刷屏
        echo -e "${YELLOW}     psql: $(echo "$output" | head -n1)${NC}"
        return 1
    fi
    return 0
}

# ============================================================
# sql_escape: 转义单引号防止 SQL 注入
# ============================================================
sql_escape() {
    printf '%s' "$1" | sed "s/'/''/g"
}

# ============================================================
# _psql_as: 以指定管理员身份执行 psql（统一入口）
#   _psql_as <admin_user> <db> [psql_args...]
#
# Linux postgres 用户: 走 sudo -u postgres（peer auth，无需密码）
# 其他情况: 走 psql -h localhost -U <user>
# ============================================================
_psql_as() {
    local admin_user="$1"; shift
    local db="$1"; shift
    # 剩余参数直接透传给 psql

    if [ "$admin_user" = "postgres" ] && [ "$(uname -s)" = "Linux" ]; then
        sudo -u postgres psql -d "$db" "$@"
    else
        psql -h localhost -U "$admin_user" -d "$db" "$@"
    fi
}

# 查询版本（用于判断 PG15+ 行为）
_pg_major_version() {
    local admin_user="$1"
    _psql_as "$admin_user" "postgres" -tAc "SHOW server_version_num;" 2>/dev/null \
        | cut -c1-2
}

# ============================================================
# check_postgres_running
# ============================================================
check_postgres_running() {
    case "$(uname -s)" in
        Darwin*)
            brew services list 2>/dev/null | grep -qE "postgresql.*started" && return 0
            pg_isready -q 2>/dev/null && return 0
            ;;
        Linux*)
            systemctl is-active --quiet postgresql 2>/dev/null && return 0
            # 兼容 postgresql@14-main 等实例化服务
            systemctl list-units --type=service --state=active 2>/dev/null \
                | grep -q "postgresql" && return 0
            pg_isready -q 2>/dev/null && return 0
            ;;
        *)
            pg_isready -q 2>/dev/null && return 0
            ;;
    esac
    return 1
}

# ============================================================
# connect_postgres <user> [db]
# 返回 0 = 可连接，1 = 不可连接
# ============================================================
connect_postgres() {
    local user="$1" db="${2:-postgres}"
    _psql_as "$user" "$db" -c "SELECT 1" >/dev/null 2>&1
}

# ============================================================
# _determine_admin_user <system_user>
# stdout 输出可用的管理员用户名；失败返回 1
# Bug fix: 改用 stdout 返回值，避免 local 变量作用域混淆
# ============================================================
_determine_admin_user() {
    local system_user="$1"

    if [ "$(uname -s)" = "Linux" ]; then
        # Linux: 优先 postgres peer auth
        if connect_postgres "postgres"; then
            echo "postgres"; return 0
        elif connect_postgres "$system_user"; then
            echo "$system_user"; return 0
        fi
    else
        # macOS: Homebrew 超级用户通常是当前系统用户
        if connect_postgres "$system_user"; then
            echo "$system_user"; return 0
        elif connect_postgres "postgres"; then
            echo "postgres"; return 0
        fi
    fi
    return 1
}

# ============================================================
# database_exists <admin_user> <db_name>
# ============================================================
database_exists() {
    local admin_user="$1" db_name="$2"
    local escaped; escaped=$(sql_escape "$db_name")
    local result
    result=$(_psql_as "$admin_user" "postgres" \
        -tAc "SELECT 1 FROM pg_database WHERE datname='${escaped}';" 2>/dev/null)
    [ "$result" = "1" ]
}

# ============================================================
# create_database <admin_user> <db_name>
# ============================================================
create_database() {
    local admin_user="$1" db_name="$2"
    _pg_log_error "CREATE DATABASE \"${db_name}\"" \
        _psql_as "$admin_user" "postgres" -c "CREATE DATABASE \"${db_name}\";"
}

# ============================================================
# user_exists <admin_user> <username>
# ============================================================
user_exists() {
    local admin_user="$1" username="$2"
    local escaped; escaped=$(sql_escape "$username")
    local result
    result=$(_psql_as "$admin_user" "postgres" \
        -tAc "SELECT 1 FROM pg_roles WHERE rolname='${escaped}';" 2>/dev/null)
    [ "$result" = "1" ]
}

# ============================================================
# create_pg_user <admin_user> <db_name> <new_user> <password>
# Bug fix: 拆分为独立的 psql 调用，失败时明确提示哪步出错
# Bug fix: 先检查用户是否已存在，避免 CREATE USER 报 already exists
# ============================================================
create_pg_user() {
    local admin_user="$1" db_name="$2" new_user="$3" password="$4"
    local escaped_pass; escaped_pass=$(sql_escape "$password")

    # 1. 创建用户（如已存在则跳过）
    if user_exists "$admin_user" "$new_user"; then
        echo -e "${YELLOW}   ℹ️  User '${new_user}' already exists, skipping CREATE USER${NC}"
    else
        echo -e "   Creating user '${new_user}'..."
        _pg_log_error "CREATE USER '${new_user}'" \
            _psql_as "$admin_user" "postgres" \
            -c "CREATE USER \"${new_user}\" WITH PASSWORD '${escaped_pass}';" \
            || return 1
    fi

    # 2. 允许该用户创建数据库（开发环境常用）
    _pg_log_error "ALTER USER CREATEDB" \
        _psql_as "$admin_user" "postgres" \
        -c "ALTER USER \"${new_user}\" CREATEDB;" \
        || return 1

    # 3. 授予数据库权限
    _pg_log_error "GRANT DATABASE" \
        _psql_as "$admin_user" "postgres" \
        -c "GRANT ALL PRIVILEGES ON DATABASE \"${db_name}\" TO \"${new_user}\";" \
        || return 1

    return 0
}

# ============================================================
# grant_schema_privileges <admin_user> <db_name> <target_user>
# Bug fix: PostgreSQL 15+ 收紧了 public schema 权限，需显式授予
# ============================================================
grant_schema_privileges() {
    local admin_user="$1" db_name="$2" target_user="$3"

    # GRANT USAGE + CREATE on public schema
    _pg_log_error "GRANT SCHEMA public" \
        _psql_as "$admin_user" "$db_name" \
        -c "GRANT USAGE, CREATE ON SCHEMA public TO \"${target_user}\";" \
        || return 1

    # 当前已有表的权限
    _pg_log_error "GRANT existing tables" \
        _psql_as "$admin_user" "$db_name" \
        -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO \"${target_user}\";" \
        || return 1

    # 未来新建表的默认权限
    _pg_log_error "ALTER DEFAULT PRIVILEGES" \
        _psql_as "$admin_user" "$db_name" \
        -c "ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO \"${target_user}\";" \
        || return 1

    return 0
}

# ============================================================
# set_user_password <admin_user> <target_user> <password>
# ============================================================
set_user_password() {
    local admin_user="$1" target_user="$2" password="$3"
    local escaped_pass; escaped_pass=$(sql_escape "$password")
    _pg_log_error "ALTER USER password" \
        _psql_as "$admin_user" "postgres" \
        -c "ALTER USER \"${target_user}\" WITH PASSWORD '${escaped_pass}';"
}

# ============================================================
# test_db_connection <user> <password> <db_name>
# ============================================================
test_db_connection() {
    local user="$1" password="$2" db_name="$3"

    if [ -n "$password" ]; then
        PGPASSWORD="$password" psql -h localhost -U "$user" -d "$db_name" \
            -c "SELECT 1" >/dev/null 2>&1
    elif [ "$user" = "postgres" ] && [ "$(uname -s)" = "Linux" ]; then
        sudo -u postgres psql -d "$db_name" -c "SELECT 1" >/dev/null 2>&1
    else
        psql -h localhost -U "$user" -d "$db_name" -c "SELECT 1" >/dev/null 2>&1
    fi
}

# ============================================================
# write_postgres_yaml
# ============================================================
write_postgres_yaml() {
    local host="$1" port="$2" db_name="$3" user="$4" password="$5"
    local config_dir="backend/config"
    local config_file="${config_dir}/postgres.yml"

    mkdir -p "$config_dir" || {
        echo -e "${RED}❌ Cannot create config directory: ${config_dir}${NC}"
        return 1
    }

    cat > "$config_file" << EOF
# PostgreSQL Configuration
# Generated by TinyForum setup script
# DO NOT EDIT MANUALLY — run 'make init-dev' to regenerate
host: ${host}
port: ${port}
dbname: ${db_name}
user: ${user}
password: ${password}
sslmode: disable
max_idle_conns: 10
max_open_conns: 100
conn_max_lifetime: 5m0s
logger:
  slow_threshold: 200ms
  log_level: info
  ignore_record_not_found_error: true
  colorful: true
EOF
    echo -e "${GREEN}   📝 Saved: ${config_file}${NC}"
}

# ============================================================
# _create_or_reuse_user
#   内部辅助：根据用户选择创建或复用数据库用户
#   stdout 输出 "final_user:final_password"
# ============================================================
_create_or_reuse_user() {
    local admin_user="$1" db_name="$2" default_user="$3" default_password="$4" system_user="$5"

    if user_exists "$admin_user" "$system_user"; then
        echo -e "${GREEN}   ✅ PostgreSQL user '${system_user}' already exists${NC}" >&2
        echo -e "   Updating password..." >&2
        set_user_password "$admin_user" "$system_user" "$default_password" || return 1
        echo "${system_user}:${default_password}"
        return 0
    fi

    echo -e "${YELLOW}   ℹ️  System user '${system_user}' not found in PostgreSQL${NC}" >&2
    echo "   How would you like to set up the database user?" >&2
    echo "   1) Use default user  : ${default_user} / ${default_password}" >&2
    echo "   2) Create custom user" >&2
    echo "   3) Exit setup" >&2
    echo "" >&2

    while true; do
        read -rp "   Enter choice [1-3]: " choice
        case "$choice" in
            1)
                echo -e "   Setting up default user '${default_user}'..." >&2
                create_pg_user "$admin_user" "$db_name" "$default_user" "$default_password" || return 1
                echo "${default_user}:${default_password}"
                return 0
                ;;
            2)
                local custom_user custom_pass
                while true; do
                    read -rp "   Username: " custom_user
                    [ -n "$custom_user" ] && break
                    echo -e "${RED}   Username cannot be empty.${NC}" >&2
                done
                while true; do
                    read -rsp "   Password: " custom_pass; echo >&2
                    [ -n "$custom_pass" ] && break
                    echo -e "${RED}   Password cannot be empty.${NC}" >&2
                done
                create_pg_user "$admin_user" "$db_name" "$custom_user" "$custom_pass" || return 1
                echo "${custom_user}:${custom_pass}"
                return 0
                ;;
            3|[qQ])
                echo -e "${YELLOW}   Cancelled by user.${NC}" >&2
                return 2
                ;;
            *)
                echo -e "${RED}   Invalid choice, enter 1, 2, or 3.${NC}" >&2
                ;;
        esac
    done
}

# ============================================================
# setup_postgres （主入口）
# 依赖全局变量: PSQL_DB_NAME, PSQL_HOST, PSQL_PORT,
#               PSQL_USER, PSQL_PASS, CURRENT_USER
# ============================================================
setup_postgres() {
    local db_name="$PSQL_DB_NAME"
    local db_host="$PSQL_HOST"
    local db_port="$PSQL_PORT"
    local default_user="$PSQL_USER"
    local default_password="$PSQL_PASS"
    local system_user="$CURRENT_USER"
    local admin_user final_user final_password pair

    echo ""
    echo -e "${BOLD}━━━ PostgreSQL Setup ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

    # ── Step 1: 确认 PostgreSQL 正在运行 ──────────────────────────────────
    echo "[1/5] Checking PostgreSQL service..."
    if ! check_postgres_running; then
        echo -e "${RED}❌ PostgreSQL is not running.${NC}"
        case "$(uname -s)" in
            Darwin*) echo "     → brew services start postgresql@<version>" ;;
            Linux*)  echo "     → sudo systemctl start postgresql" ;;
        esac
        return 1
    fi
    echo -e "${GREEN}   ✅ PostgreSQL is running${NC}"

    # ── Step 2: 确定可用管理员 ───────────────────────────────────────────
    echo "[2/5] Determining admin connection..."
    if ! admin_user=$(_determine_admin_user "$system_user"); then
        echo -e "${RED}❌ Cannot connect to PostgreSQL as 'postgres' or '${system_user}'.${NC}"
        echo "   Troubleshooting:"
        echo "     Linux  → sudo -u postgres psql  (peer auth)"
        echo "     macOS  → psql -h localhost -U \$(whoami)  (Homebrew default)"
        return 1
    fi
    echo -e "${GREEN}   ✅ Admin: '${admin_user}'${NC}"

    # ── Step 3: 确保数据库存在 ────────────────────────────────────────────
    echo "[3/5] Checking database '${db_name}'..."
    if database_exists "$admin_user" "$db_name"; then
        echo -e "${GREEN}   ✅ Database '${db_name}' already exists${NC}"
    else
        echo "   Creating database '${db_name}'..."
        if create_database "$admin_user" "$db_name"; then
            echo -e "${GREEN}   ✅ Database '${db_name}' created${NC}"
        else
            echo -e "${RED}❌ Failed to create database. See error above.${NC}"
            return 1
        fi
    fi

    # ── Step 4: 用户管理 ──────────────────────────────────────────────────
    echo "[4/5] Database user setup..."
    pair=$(_create_or_reuse_user \
        "$admin_user" "$db_name" "$default_user" "$default_password" "$system_user")
    local user_rc=$?

    if [ "$user_rc" -eq 2 ]; then
        # 用户主动取消
        return 0
    elif [ "$user_rc" -ne 0 ]; then
        echo -e "${RED}❌ User setup failed. See error above.${NC}"
        return 1
    fi

    # 解析 "user:password"（密码可能含冒号，只切第一个冒号）
    final_user="${pair%%:*}"
    final_password="${pair#*:}"

    # 授予 schema 权限
    echo "   Granting schema privileges to '${final_user}'..."
    if grant_schema_privileges "$admin_user" "$db_name" "$final_user"; then
        echo -e "${GREEN}   ✅ Schema privileges granted${NC}"
    else
        echo -e "${YELLOW}   ⚠️  Could not grant all schema privileges (may already be set)${NC}"
    fi

    # ── Step 5: 连接验证 ──────────────────────────────────────────────────
    echo "[5/5] Verifying connection as '${final_user}'..."
    if test_db_connection "$final_user" "$final_password" "$db_name"; then
        echo -e "${GREEN}   ✅ Connection verified${NC}"
    else
        echo -e "${RED}❌ Cannot connect as '${final_user}'.${NC}"
        echo "   Possible causes:"
        echo "     • pg_hba.conf requires md5/scram but password was not set correctly"
        echo "     • PostgreSQL requires 'host' entry for 127.0.0.1"
        echo "   Try: sudo -u postgres psql -c \"ALTER USER \\\"${final_user}\\\" WITH PASSWORD '${final_password}';\""
        return 1
    fi

    # ── 写入配置 ──────────────────────────────────────────────────────────
    write_postgres_yaml "$db_host" "$db_port" "$db_name" "$final_user" "$final_password" \
        || return 1

    # ── 可选：额外用户 ────────────────────────────────────────────────────
    echo ""
    read -rp "   Create an additional PostgreSQL user? (y/N): " add_extra
    if [[ "$add_extra" =~ ^[Yy]$ ]]; then
        local extra_user extra_pass
        read -rp "   Extra username: " extra_user
        if [ -z "$extra_user" ]; then
            echo -e "${YELLOW}   Skipped (empty username).${NC}"
        elif user_exists "$admin_user" "$extra_user"; then
            echo -e "${YELLOW}   ⚠️  User '${extra_user}' already exists.${NC}"
        else
            read -rsp "   Password for '${extra_user}': " extra_pass; echo
            if create_pg_user "$admin_user" "$db_name" "$extra_user" "$extra_pass"; then
                grant_schema_privileges "$admin_user" "$db_name" "$extra_user"
                echo -e "${GREEN}   ✅ Extra user '${extra_user}' created${NC}"
            else
                echo -e "${RED}   ❌ Failed to create '${extra_user}'${NC}"
            fi
        fi
    fi

    # ── 完成摘要 ──────────────────────────────────────────────────────────
    echo ""
    echo -e "${GREEN}🎉 PostgreSQL setup completed!${NC}"
    echo -e "   Database : ${BOLD}${db_name}${NC}"
    echo -e "   Host     : ${db_host}:${db_port}"
    echo -e "   User     : ${BOLD}${final_user}${NC}"
    echo -e "   Config   : backend/config/postgres.yml"
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

    # 导出供 summary 使用
    PG_FINAL_USER="$final_user"
    PG_FINAL_PASS="$final_password"
    PG_FINAL_DB="$db_name"
}