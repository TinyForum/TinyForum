

create_config_gitignore() {
    local gitignore_file="backend/config/.gitignore"
    
    if [ -f "$gitignore_file" ]; then
        echo -e "${YELLOW}  Config .gitignore already exists${NC}"
        return 0
    fi
    
    cat > "$gitignore_file" << EOF
# 忽略敏感配置
private.*
*.key
*.pem
*.crt

EOF
    
    echo -e "${GREEN}  ✓ Created config/.gitignore${NC}"
    return 0
}

setup_configurations() {
    echo ""
    echo "Step 2.5: Setting up configuration files..."
    
   
    create_config_gitignore
    
    echo -e "${GREEN}  ✓ All configuration files ready${NC}"
}
