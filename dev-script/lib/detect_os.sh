# ============================================
# Module 1: System Detection
# ============================================

detect_os() {
    echo -e "${GREEN}Detecting operating system...${NC}"
    case "$OS" in
        Darwin*)
            echo -e "${GREEN}✓ macOS detected${NC}"
            ;;
        Linux*)
            echo -e "${GREEN}✓ Linux detected${NC}"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            echo -e "${GREEN}✓ Windows detected${NC}"
            ;;
        *)
            echo -e "${YELLOW}⚠ Unknown operating system: $OS${NC}"
            ;;
    esac
}

detect_default_db_user() {
    case "$OS" in
        Darwin*)
            echo "$(whoami)"
            ;;
        Linux*)
            if sudo -u postgres psql -c "SELECT 1" >/dev/null 2>&1; then
                echo "postgres"
            else
                echo "$(whoami)"
            fi
            ;;
        MINGW*|MSYS*|CYGWIN*)
            echo "postgres"
            ;;
        *)
            echo "$(whoami)"
            ;;
    esac
}