# ============================================
# Module 1: System Detection
# ============================================

detect_os() {
    echo -e "${BOLD}Detecting operating system...${NC}"
    case "$OS" in
        Darwin*)
            local macos_ver
            macos_ver=$(sw_vers -productVersion 2>/dev/null || echo "unknown")
            echo -e "${GREEN}✓ macOS detected (${macos_ver})${NC}"
            ;;
        Linux*)
            local distro=""
            if [ -f /etc/os-release ]; then
                # shellcheck source=/dev/null
                distro=$(. /etc/os-release && echo "${PRETTY_NAME:-$NAME}")
            fi
            echo -e "${GREEN}✓ Linux detected${distro:+ (${distro})}${NC}"
            ;;
        *)
            echo -e "${YELLOW}⚠ Unknown operating system: $OS${NC}"
            ;;
    esac
    echo -e "  Local IP  : ${LOCAL_IP}"
    echo -e "  System user: ${CURRENT_USER}"
}

# detect_default_db_user
# Bug fix: macOS 不应 sudo -u postgres（该账户通常不存在）
# 只做连接测试，不做任何写操作
detect_default_db_user() {
    case "$OS" in
        Darwin*)
            # macOS Homebrew: 超级用户是当前系统用户
            echo "$(whoami)"
            ;;
        Linux*)
            # Linux: 先用 peer auth 测试 postgres，再测试当前用户
            if sudo -n -u postgres psql -c "SELECT 1" >/dev/null 2>&1; then
                echo "postgres"
            else
                echo "$(whoami)"
            fi
            ;;
        *)
            echo "$(whoami)"
            ;;
    esac
}