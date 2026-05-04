check_backend_config() {
  

    # 检查后端配置文件列表
    local configs=(
        "$POSTGRES_CONFIG_PATH"
        "$REDIS_CONFIG_PATH"
        "$BASIC_CONFIG_PATH"
        "$PRIVATE_CONFIG_PATH"
        "$RISK_CONTROL_CONFIG_PATH"
    )

    for config in "${configs[@]}"; do
        if [ ! -f "$config" ]; then
            echo -e "${RED}Backend config file not found at ${config}${NC}"
            exit 1
        fi
    done

 
}

check_frontend_config() {
      # 检查前端配置文件
    if [ ! -f "${FRONTEND_CONFIG_PATH}" ]; then
        echo -e "${RED}Frontend config file not found at ${FRONTEND_CONFIG_PATH}${NC}"
        exit 1
    fi
}

check_config() {
    check_backend_config
       case $? in
        0)
            echo -e "${GREEN}   ✅ Config checked${NC}"
            ;;
        1)
            echo -e "${RED}   ❌ Failed to check config${NC}"
            exit 1
            ;;

        esac
      echo -e "${GREEN}All configuration files exist.${NC}"
}