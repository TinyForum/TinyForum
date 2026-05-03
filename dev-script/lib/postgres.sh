#!/bin/bash
# ============================================
# Module 3: PostgreSQL Management
# ============================================
# 支持平台: macOS (Homebrew), Linux (systemd)
# 修复:
#   - user_exists / database_exists 逻辑错误（检查输出而非退出码）
#   - PSQL_USER 全局变量污染（helper 改为接受参数）
#   - grant_schema_privileges 返回值检测错误
#   - $? 跨语句失效问题
#   - 密码含单引号时的 SQL 注入风险
#   - macOS 下管理员用户判断逻辑错误

# ── 颜色（未定义时设默认值）──────────────────────────────────────────────────
: "${RED:=\033[0;31m}"
: "${GREEN:=\033[0;32m}"
: "${YELLOW:=\033[0;33m}"
: "${NC:=\033[0m}"

# ============================================================
# 内部工具：转义 SQL 单引号（防止注入）
# 用法: escaped=$(sql_escape "O'Brien")
# ============================================================
sql_escape() {
    # 将单引号替换为两个单引号（标准 SQL 转义）
    printf '%s' "$1" | sed "s/'/''/g"
}

# ============================================================
# 内部工具：统一执行 psql 命令
#   _psql_exec <admin_user> <db> <sql>
#   在 Linux 上若 admin_user=postgres，则走 sudo -u postgres
#   否则走 psql -h localhost
# ============================================================
_psql_exec() {
    local admin_user="$1"
    local db="$2"
    local sql="$3"

    if [ "$admin_user" = "postgres" ] && [ "$(uname -s)" = "Linux" ]; then
        sudo -u postgres psql -d "$db" -c "$sql"
    else
        psql -h localhost -U "$admin_user" -d "$db" -c "$sql"
    fi
}

# 同上，但返回单行文本结果（-tAc）
_psql_query() {
    local admin_user="$1"
    local db="$2"
    local sql="$3"

    if [ "$admin_user" = "postgres" ] && [ "$(uname -s)" = "Linux" ]; then
        sudo -u postgres psql -d "$db" -tAc "$sql" 2>/dev/null
    else
        psql -h localhost -U "$admin_user" -d "$db" -tAc "$sql" 2>/dev/null
    fi
}

# ============================================================
# check_postgres_running
# ============================================================
check_postgres_running() {
    local os
    os="$(uname -s)"
    case "$os" in
        Darwin*)
            # Homebrew 服务 或 pg_isready 任一成功即可
            if brew services list 2>/dev/null | grep -qE "postgresql.*started"; then
                return 0
            fi
            pg_isready -q 2>/dev/null && return 0
            ;;
        Linux*)
            # systemd 普通服务名 或 实例化服务名（postgresql@14-main 等）
            if systemctl is-active --quiet postgresql 2>/dev/null; then
                return 0
            fi
            # 兼容带版本号的服务（如 postgresql@14-main）
            if systemctl list-units --type=service --state=active 2>/dev/null \
                    | grep -q "postgresql"; then
                return 0
            fi
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
# 仅测试是否能连接，不执行任何修改
# ============================================================
connect_postgres() {
    local user="$1"
    local db="${2:-postgres}"
    local os
    os="$(uname -s)"

    if [ "$user" = "postgres" ] && [ "$os" = "Linux" ]; then
        sudo -u postgres psql -d "$db" -c "SELECT 1" >/dev/null 2>&1
    else
        psql -h localhost -U "$user" -d "$db" -c "SELECT 1" >/dev/null 2>&1
    fi
}

# ============================================================
# database_exists <admin_user> <db_name>
# 修复：通过查询 pg_database 而非解析 \l 输出，更可靠
# ============================================================
database_exists() {
    local admin_user="$1"
    local db_name="$2"
    local escaped_name
    escaped_name="$(sql_escape "$db_name")"

    local result
    result="$(_psql_query "$admin_user" "postgres" \
        "SELECT 1 FROM pg_database WHERE datname='${escaped_name}'")"
    [ "$result" = "1" ]
}

# ============================================================
# create_database <admin_user> <db_name>
# ============================================================
create_database() {
    local admin_user="$1"
    local db_name="$2"
    _psql_exec "$admin_user" "postgres" "CREATE DATABASE \"${db_name}\";" >/dev/null 2>&1
}

# ============================================================
# user_exists <admin_user> <target_user>
# 修复：检查查询输出是否为 "1"，而非退出码
# ============================================================
user_exists() {
    local admin_user="$1"
    local target_user="$2"
    local escaped_user
    escaped_user="$(sql_escape "$target_user")"

    local result
    result="$(_psql_query "$admin_user" "postgres" \
        "SELECT 1 FROM pg_roles WHERE rolname='${escaped_user}'")"
    [ "$result" = "1" ]
}

# ============================================================
# create_database_user <admin_user> <db_name> <new_user> <password>
# ============================================================
create_database_user() {
    local admin_user="$1"
    local db_name="$2"
    local new_user="$3"
    local password="$4"
    local escaped_pass
    escaped_pass="$(sql_escape "$password")"

    local sql
    sql="CREATE USER \"${new_user}\" WITH PASSWORD '${escaped_pass}';
         ALTER USER \"${new_user}\" CREATEDB;
         GRANT ALL PRIVILEGES ON DATABASE \"${db_name}\" TO \"${new_user}\";"

    _psql_exec "$admin_user" "postgres" "$sql" >/dev/null 2>&1
}

# ============================================================
# grant_schema_privileges <admin_user> <db_name> <target_user>
# 修复：移除对 $? 的延迟检查，改为直接捕获返回值
# ============================================================
grant_schema_privileges() {
    local admin_user="$1"
    local db_name="$2"
    local target_user="$3"

    local sql
    sql="GRANT ALL PRIVILEGES ON SCHEMA public TO \"${target_user}\";
         ALTER DEFAULT PRIVILEGES IN SCHEMA public
           GRANT ALL ON TABLES TO \"${target_user}\";"

    _psql_exec "$admin_user" "$db_name" "$sql" >/dev/null 2>&1
    # 直接返回 _psql_exec 的退出码，调用方可直接 if grant_schema_privileges ...
}

# ============================================================
# set_user_password <admin_user> <target_user> <new_password>
# ============================================================
set_user_password() {
    local admin_user="$1"
    local target_user="$2"
    local new_password="$3"
    local escaped_pass
    escaped_pass="$(sql_escape "$new_password")"

    _psql_exec "$admin_user" "postgres" \
        "ALTER USER \"${target_user}\" WITH PASSWORD '${escaped_pass}';" >/dev/null 2>&1
}

# ============================================================
# test_db_connection <user> <password> <db_name>
# ============================================================
test_db_connection() {
    local user="$1"
    local password="$2"
    local db_name="$3"
    local os
    os="$(uname -s)"

    if [ -n "$password" ]; then
        PGPASSWORD="$password" psql -h localhost -U "$user" -d "$db_name" \
            -c "SELECT 1" >/dev/null 2>&1
    elif [ "$user" = "postgres" ] && [ "$os" = "Linux" ]; then
        sudo -u postgres psql -d "$db_name" -c "SELECT 1" >/dev/null 2>&1
    else
        psql -h localhost -U "$user" -d "$db_name" -c "SELECT 1" >/dev/null 2>&1
    fi
}

# ============================================================
# write_postgres_yaml <host> <port> <db_name> <user> <password>
# ============================================================
write_postgres_yaml() {
    local host="$1"
    local port="$2"
    local db_name="$3"
    local user="$4"
    local password="$5"
    local config_dir="backend/config"
    local config_file="${config_dir}/postgres.yml"

    mkdir -p "$config_dir" || {
        echo -e "${RED}❌ Cannot create config directory: $config_dir${NC}"
        return 1
    }

    cat > "$config_file" << EOF
# PostgreSQL Configuration
# Generated by TinyForum setup script
# DO NOT EDIT MANUALLY - Run 'make init-dev' to regenerate
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
  log_level: info          # 可选: silent, error, warn, info
  ignore_record_not_found_error: true
  colorful: true
EOF

    echo -e "${GREEN}   📝 Configuration saved to ${config_file}${NC}"
}

# ============================================================
# _determine_admin_user
# 输出：将可用管理员用户名写入全局 _ADMIN_USER
# 修复：macOS 上优先尝试当前系统用户（Homebrew 默认超级用户），
#       再尝试 postgres，避免 macOS 上 postgres 用户通常不存在的问题
# ============================================================
_determine_admin_user() {
    local os system_user="$CURRENT_USER"
    os="$(uname -s)"

    if [ "$os" = "Linux" ]; then
        # Linux：优先走 postgres（peer auth）
        if connect_postgres "postgres"; then
            _ADMIN_USER="postgres"
        elif connect_postgres "$system_user"; then
            _ADMIN_USER="$system_user"
        else
            return 1
        fi
    else
        # macOS：Homebrew 默认超级用户是当前系统用户
        if connect_postgres "$system_user"; then
            _ADMIN_USER="$system_user"
        elif connect_postgres "postgres"; then
            _ADMIN_USER="postgres"
        else
            return 1
        fi
    fi
    return 0
}

# ============================================================
# setup_postgres  （主入口）
# 依赖外部变量: PSQL_DB_NAME, PSQL_HOST, PSQL_PORT,
#               PSQL_USER, PSQL_PASS, CURRENT_USER
# ============================================================
setup_postgres() {
    local db_name="$PSQL_DB_NAME"
    local db_host="$PSQL_HOST"
    local db_port="$PSQL_PORT"
    local default_user="$PSQL_USER"
    local default_password="$PSQL_PASS"
    local system_user="$CURRENT_USER"
    local _ADMIN_USER=""   # 由 _determine_admin_user 填充
    local final_user=""
    local final_password=""

    # ── Step 1: 检查 PostgreSQL 是否运行 ───────────────────────────────────
    echo ""
    echo "Step 1: Checking PostgreSQL..."

    if ! check_postgres_running; then
        echo -e "${RED}❌ PostgreSQL is not running${NC}"
        echo "   Please start PostgreSQL:"
        case "$(uname -s)" in
            Darwin*) echo "   macOS: brew services start postgresql@<version>" ;;
            Linux*)  echo "   Linux: sudo systemctl start postgresql" ;;
            *)       echo "   Start PostgreSQL manually." ;;
        esac
        return 1
    fi
    echo -e "${GREEN}✅ PostgreSQL is running${NC}"

    # ── 确定管理员用户 ─────────────────────────────────────────────────────
    echo "   Determining admin user..."
    if ! _determine_admin_user; then
        echo -e "${RED}❌ Cannot connect as 'postgres' or '$system_user'${NC}"
        echo "   Please ensure PostgreSQL is accessible."
        return 1
    fi
    echo -e "${GREEN}   Using admin: ${_ADMIN_USER}${NC}"

    # ── Step 2: 确保数据库存在 ────────────────────────────────────────────
    echo ""
    echo "Step 2: Checking database '${db_name}'..."

    if database_exists "$_ADMIN_USER" "$db_name"; then
        echo -e "${GREEN}   Database '${db_name}' already exists${NC}"
    else
        echo "   Creating database '${db_name}'..."
        if create_database "$_ADMIN_USER" "$db_name"; then
            echo -e "${GREEN}   Database '${db_name}' created${NC}"
        else
            echo -e "${RED}❌ Failed to create database '${db_name}'${NC}"
            return 1
        fi
    fi

    # ── Step 3: 数据库用户管理 ────────────────────────────────────────────
    echo ""
    echo "Step 3: Database user setup..."

    if user_exists "$_ADMIN_USER" "$system_user"; then
        # 系统用户已存在于 PostgreSQL，更新密码即可
        echo -e "${GREEN}   PostgreSQL user '${system_user}' already exists${NC}"
        echo "   Updating password for '${system_user}'..."
        if set_user_password "$_ADMIN_USER" "$system_user" "$default_password"; then
            echo -e "${GREEN}   Password updated for '${system_user}'${NC}"
            final_user="$system_user"
            final_password="$default_password"
        else
            echo -e "${RED}❌ Failed to update password for '${system_user}'${NC}"
            return 1
        fi
    else
        # 系统用户不存在，提示选择
        echo -e "${YELLOW}   PostgreSQL user '${system_user}' not found${NC}"
        echo "   Choose how to create a database user:"
        echo "   1) Create default user (${default_user} / ${default_password})"
        echo "   2) Create custom user"
        echo "   3) Exit"

        while true; do
            read -rp "   Enter your choice (1-3): " choice
            case "$choice" in
                1)
                    echo "   Creating default user '${default_user}'..."
                    if create_database_user "$_ADMIN_USER" "$db_name" \
                            "$default_user" "$default_password"; then
                        echo -e "${GREEN}   User '${default_user}' created${NC}"
                        final_user="$default_user"
                        final_password="$default_password"
                    else
                        echo -e "${RED}❌ Failed to create default user${NC}"
                        return 1
                    fi
                    break
                    ;;
                2)
                    read -rp "   Enter username: " custom_user
                    if [ -z "$custom_user" ]; then
                        echo -e "${RED}   Username cannot be empty.${NC}"
                        continue
                    fi
                    read -rsp "   Enter password: " custom_pass
                    echo
                    if [ -z "$custom_pass" ]; then
                        echo -e "${RED}   Password cannot be empty.${NC}"
                        continue
                    fi
                    if create_database_user "$_ADMIN_USER" "$db_name" \
                            "$custom_user" "$custom_pass"; then
                        echo -e "${GREEN}   User '${custom_user}' created${NC}"
                        final_user="$custom_user"
                        final_password="$custom_pass"
                    else
                        echo -e "${RED}❌ Failed to create user '${custom_user}'${NC}"
                        return 1
                    fi
                    break
                    ;;
                3|[qQ])
                    echo -e "${YELLOW}   Exiting. No user was created.${NC}"
                    return 0
                    ;;
                *)
                    echo -e "${RED}   Invalid choice, please enter 1, 2, or 3.${NC}"
                    ;;
            esac
        done
    fi

    # ── Step 4: 授予 schema 权限 ──────────────────────────────────────────
    echo ""
    echo "Step 4: Granting privileges to '${final_user}'..."
    if grant_schema_privileges "$_ADMIN_USER" "$db_name" "$final_user"; then
        echo -e "${GREEN}   Privileges granted successfully${NC}"
    else
        echo -e "${YELLOW}   ⚠️  Could not grant schema privileges (may already be set)${NC}"
    fi

    # ── Step 5: 验证连接 ──────────────────────────────────────────────────
    echo ""
    echo "Step 5: Verifying connection as '${final_user}'..."
    if test_db_connection "$final_user" "$final_password" "$db_name"; then
        echo -e "${GREEN}✅ Successfully connected to '${db_name}' as '${final_user}'${NC}"
    else
        echo -e "${RED}❌ Connection failed for user '${final_user}'${NC}"
        return 1
    fi

    # ── Step 6: 写入 YAML 配置 ────────────────────────────────────────────
    write_postgres_yaml "$db_host" "$db_port" "$db_name" "$final_user" "$final_password" \
        || return 1

    # ── 可选：创建额外用户 ────────────────────────────────────────────────
    echo ""
    read -rp "Create an additional PostgreSQL user? (y/n): " add_extra
    if [[ "$add_extra" =~ ^[Yy]$ ]]; then
        read -rp "   Enter additional username: " extra_user
        if [ -z "$extra_user" ]; then
            echo -e "${YELLOW}   Username empty, skipping.${NC}"
        elif user_exists "$_ADMIN_USER" "$extra_user"; then
            echo -e "${YELLOW}   ⚠️  User '${extra_user}' already exists, skipping creation.${NC}"
        else
            read -rsp "   Enter password for '${extra_user}': " extra_pass
            echo
            if create_database_user "$_ADMIN_USER" "$db_name" "$extra_user" "$extra_pass"; then
                echo -e "${GREEN}✅ Additional user '${extra_user}' created${NC}"
                if grant_schema_privileges "$_ADMIN_USER" "$db_name" "$extra_user"; then
                    echo -e "${GREEN}   Privileges granted to '${extra_user}'${NC}"
                else
                    echo -e "${YELLOW}   ⚠️  Could not grant privileges to '${extra_user}'${NC}"
                fi
            else
                echo -e "${RED}❌ Failed to create additional user '${extra_user}'${NC}"
            fi
        fi
    fi

    # ── 完成 ──────────────────────────────────────────────────────────────
    echo ""
    echo -e "${GREEN}🎉 PostgreSQL setup completed!${NC}"
    echo "   Database : ${db_name}"
    echo "   Host     : ${db_host}:${db_port}"
    echo "   User     : ${final_user}"
    echo "   Config   : backend/config/postgres.yml"
}