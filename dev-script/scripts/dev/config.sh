# ============================================
# Module: Configuration File Setup
# ============================================

create_config_gitignore() {
    local gitignore_file="backend/config/.gitignore"

    mkdir -p "backend/config"

    if [ -f "$gitignore_file" ]; then
        echo -e "${YELLOW}   ℹ️  Config .gitignore already exists, skipping${NC}"
        return 0
    fi

    cat > "$gitignore_file" << 'EOF'
# 忽略敏感配置文件（由 make init-dev 生成，不提交到版本控制）
private.yml
private.yaml
postgres.yml
postgres.yaml
redis.yml
redis.yaml
*.key
*.pem
*.crt
EOF

    echo -e "${GREEN}   ✅ Created backend/config/.gitignore${NC}"
}

setup_configurations() {
    echo ""
    echo -e "${BOLD}━━━ Configuration Files ━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    create_config_gitignore
    echo -e "${GREEN}✅ Configuration setup done.${NC}"
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}