#!/usr/bin/env bash
# ============================================
# Nginx 配置部署模块（函数库）
# ============================================

# 颜色常量
export RED='\033[0;31m'
export GREEN='\033[0;32m'
export YELLOW='\033[1;33m'
export BLUE='\033[0;34m'
export NC='\033[0m'

NGINX_CONFIG=""
CONF_FILE=""
ENABLED_LINK=""
OS=""

detect_os() {
    OS=$(uname -s)
    echo -e "${BLUE}🔍 检测操作系统: $OS${NC}"
}

check_nginx_installed() {
    if ! command -v nginx &> /dev/null; then
        echo -e "${RED}❌ Nginx 未安装${NC}"
        if [[ "$OS" == "Darwin" ]]; then
            echo "请运行: brew install nginx"
        else
            echo "请运行: sudo apt update && sudo apt install nginx -y"
        fi
        exit 1
    fi
    echo -e "${GREEN}✅ Nginx 已安装: $(nginx -v 2>&1)${NC}"
}

get_nginx_paths() {
    if [[ "$OS" == "Darwin" ]]; then
        NGINX_PREFIX=$(brew --prefix nginx 2>/dev/null || echo "/usr/local")
        NGINX_CONF_DIR="$NGINX_PREFIX/etc/nginx/conf.d"
        NGINX_SITES_AVAILABLE="$NGINX_PREFIX/etc/nginx/sites-available"
        NGINX_SITES_ENABLED="$NGINX_PREFIX/etc/nginx/sites-enabled"
        # 主配置文件通常在 /opt/homebrew/etc/nginx/nginx.conf (符号链接)
        if [[ -f "/opt/homebrew/etc/nginx/nginx.conf" ]]; then
            NGINX_MAIN_CONF="/opt/homebrew/etc/nginx/nginx.conf"
        else
            NGINX_MAIN_CONF="$NGINX_PREFIX/etc/nginx/nginx.conf"
        fi
    else
        NGINX_CONF_DIR="/etc/nginx/conf.d"
        NGINX_SITES_AVAILABLE="/etc/nginx/sites-available"
        NGINX_SITES_ENABLED="/etc/nginx/sites-enabled"
        NGINX_MAIN_CONF="/etc/nginx/nginx.conf"
    fi
    echo -e "${BLUE}📁 Nginx 配置目录: $NGINX_CONF_DIR${NC}"
}

determine_config_path() {
    local use_sites="$1"
    if [[ "$use_sites" == true ]] && [[ -d "$NGINX_SITES_AVAILABLE" ]]; then
        CONF_FILE="$NGINX_SITES_AVAILABLE/tinyforum"
        ENABLED_LINK="$NGINX_SITES_ENABLED/tinyforum"
        mkdir -p "$NGINX_SITES_ENABLED"
        echo -e "${GREEN}📝 使用 sites-available 模式，配置文件: $CONF_FILE${NC}"
    else
        CONF_FILE="$NGINX_CONF_DIR/tinyforum.conf"
        mkdir -p "$NGINX_CONF_DIR"
        echo -e "${GREEN}📝 使用 conf.d 模式，配置文件: $CONF_FILE${NC}"
    fi
}

define_nginx_config() {
    read -r -d '' NGINX_CONFIG << 'EOF' || true
server {
    listen 80;
    server_name _;

    location /api/ {
        proxy_pass         http://127.0.0.1:8080/api/;
        proxy_http_version 1.1;
        proxy_set_header   Host              $host;
        proxy_set_header   X-Real-IP         $remote_addr;
        proxy_set_header   X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto $scheme;
        proxy_read_timeout 120s;
        proxy_send_timeout 120s;
    }

    location /uploads/ {
        proxy_pass http://127.0.0.1:8080/uploads/;
        expires 7d;
        add_header Cache-Control "public, immutable";
    }

    location / {
        proxy_pass         http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header   Upgrade           $http_upgrade;
        proxy_set_header   Connection        "upgrade";
        proxy_set_header   Host              $host;
        proxy_set_header   X-Real-IP         $remote_addr;
        proxy_set_header   X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
EOF
}

write_config() {
    echo -e "${GREEN}📝 写入 Nginx 配置文件: $CONF_FILE${NC}"
    sudo tee "$CONF_FILE" > /dev/null <<< "$NGINX_CONFIG"
    if [[ -n "$ENABLED_LINK" ]] && [[ ! -L "$ENABLED_LINK" ]]; then
        echo -e "${GREEN}🔗 启用站点配置 (ln -sf)${NC}"
        sudo ln -sf "$CONF_FILE" "$ENABLED_LINK"
    fi
}

enhance_nginx_conf() {
    local general_conf="$NGINX_CONF_DIR/00-general.conf"
    if [[ -f "$general_conf" ]] && grep -q "client_max_body_size 5m" "$general_conf"; then
        echo -e "${GREEN}✅ 通用配置已存在: $general_conf${NC}"
        return
    fi

    echo -e "${GREEN}📝 写入通用配置文件: $general_conf${NC}"
    sudo tee "$general_conf" > /dev/null << 'EOF'
client_max_body_size 5m;
gzip on;
gzip_types text/plain text/css application/json application/javascript text/xml application/xml image/svg+xml;
gzip_min_length 1024;
EOF
}

reload_nginx() {
    local should_reload="${1:-false}"
    echo -e "${BLUE}🔍 测试 Nginx 配置语法...${NC}"
    if sudo nginx -t; then
        echo -e "${GREEN}✅ 配置语法正确${NC}"
        if [[ "$should_reload" == true ]]; then
            echo -e "${GREEN}🔄 重新加载 Nginx${NC}"
            # 尝试重载，如果失败则启动
            if sudo nginx -s reload 2>/dev/null; then
                echo -e "${GREEN}✅ Nginx 重载成功${NC}"
            else
                echo -e "${YELLOW}⚠️  Nginx 未运行或重载失败，尝试启动...${NC}"
                if sudo nginx; then
                    sleep 1
                    echo -e "${GREEN}✅ Nginx 启动成功${NC}"
                else
                    echo -e "${RED}❌ Nginx 启动失败，请手动检查${NC}"
                    exit 1
                fi
            fi
        fi
    else
        echo -e "${RED}❌ 配置语法错误，请检查 $CONF_FILE${NC}"
        exit 1
    fi
}

preview_config() {
    echo -e "${YELLOW}🔍 预览模式：将写入以下内容到 $CONF_FILE${NC}"
    echo "----------------------------------------"
    echo "$NGINX_CONFIG"
    echo "----------------------------------------"
}

main_setup() {
    local use_sites="${1:-false}"
    local do_reload="${2:-true}"
    local dry_run="${3:-false}"

    detect_os
    check_nginx_installed
    get_nginx_paths
    define_nginx_config
    determine_config_path "$use_sites"

    if [[ "$dry_run" == true ]]; then
        preview_config
        return 0
    fi

    write_config
    enhance_nginx_conf
    reload_nginx "$do_reload"

    echo -e "${GREEN}🎉 Nginx 配置部署完成${NC}"
}