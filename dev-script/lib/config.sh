
# ============================================
# Module 4: Configuration Management
# ============================================

create_private_config() {
    local config_file="backend/config/private.yaml"
    
    if [ -f "$config_file" ]; then
        echo -e "${YELLOW}  Private config already exists: $config_file${NC}"
        # Update database user
        if [[ "$OSTYPE" == "darwin"* ]]; then
            sed -i '' "s/user:.*/user: $DB_USER/" "$config_file" 2>/dev/null || true
        else
            sed -i "s/user:.*/user: $DB_USER/" "$config_file" 2>/dev/null || true
        fi
        if [[ -n "$NEW_DB_PASS" ]]; then
            if [[ "$OSTYPE" == "darwin"* ]]; then
                sed -i '' "s/password:.*/password: $NEW_DB_PASS/" "$config_file" 2>/dev/null || true
            else
                sed -i "s/password:.*/password: $NEW_DB_PASS/" "$config_file" 2>/dev/null || true
            fi
        fi
        echo -e "${GREEN}  ✓ Updated database credentials${NC}"
        return 0
    fi
    
    mkdir -p backend/config
    cat > "$config_file" << EOF
# Private Configuration
# 包含敏感信息，请勿提交到版本控制

# 邮件配置
email:
  host: smtp.gmail.com
  port: 587
  username: noreply@example.com
  password: your-email-password
  from: noreply@example.com
  from_name: Tiny Forum
  ssl: false
  tls: true
  pool_size: 5

# JWT配置
jwt:
  secret: "tiny-forum-secret-change-in-production-32chars"
  expire: 24h
  refresh_expire: 168h
  issuer: "tiny-forum"

# 数据库配置
database:
  host: localhost
  port: 5432
  user: $DB_USER
  password: "${NEW_DB_PASS:-}"
  dbname: tiny_forum
  sslmode: disable
  timezone: Asia/Shanghai

# 服务器配置
server:
  port: 8080
  mode: debug

# Redis配置
redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

# 管理员账户
admin:
  email: admin@test.com
  password: password
  username: admin
  role: super_admin
  score: 10000
EOF
    
    echo -e "${GREEN}  ✓ Created private config: $config_file${NC}"
    return 0
}

create_basic_config() {
    local config_file="backend/config/basic.yaml"
    
    if [ -f "$config_file" ]; then
        echo -e "${YELLOW}  Basic config already exists: $config_file${NC}"
        return 0
    fi
    
    mkdir -p backend/config
    cat > "$config_file" << EOF
# Basic Configuration
# 基础配置，可提交到版本控制

# 服务器配置
server:
  protocol: http
  host: $LOCAL_IP
  port: 8080
  mode: debug
  read_timeout: 30s
  write_timeout: 30s
  max_header_bytes: 1048576

# API配置
api:
  protocol: http
  host: $LOCAL_IP
  port: 8080
  version: v1
  prefix: /api

# JWT配置（默认值，会被 private.yaml 覆盖）
jwt:
  expire: 24h
  refresh_expire: 168h
  issuer: "tiny-forum"

# 日志配置
log:
  level: info
  filename: ./logs/app.log
  max_size: 100
  max_backups: 10
  max_age: 30
  compress: true
  console: true
  json_format: false

# 限流配置
rate_limit:
  enabled: true
  requests: 100
  duration: 60
  burst: 50

# Ollama AI 配置
ollama:
  base_url: http://localhost:11434
  model: qwen3:0.6b
  num_predict: 256
  temperature: 0.7
  timeout: 60

# CORS 跨域配置
allow_origins:
  - http://localhost:3000
  - http://127.0.0.1:3000
  - http://localhost:8080
  - http://127.0.0.1:8080
  - http://$LOCAL_IP:3000
  - http://$LOCAL_IP:8080

# 上传配置
upload:
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
    allowed_types:
      - .jpg
      - .jpeg
      - .png
      - .webp
    width: 512
    height: 512
  post_image:
    max_size: 5242880
    allowed_types:
      - .jpg
      - .jpeg
      - .png
      - .gif
      - .webp
  storage:
    type: local
    local_path: ./uploads
    url_prefix: /uploads/
EOF
    
    echo -e "${GREEN}  ✓ Created basic config: $config_file${NC}"
    return 0
}

create_risk_config() {
    local config_file="backend/config/risk_control.yaml"
    
    if [ -f "$config_file" ]; then
        echo -e "${YELLOW}  Risk control config already exists: $config_file${NC}"
        return 0
    fi
    
    mkdir -p backend/config
    cat > "$config_file" << 'EOF'
# Risk Control Configuration
# 风控配置，可提交到版本控制

# 不同风险等级的限流策略
rate_limit:
  risk_levels:
    normal:
      create_post:
        limit: 20
        window: 1h
      create_comment:
        limit: 60
        window: 1h
      send_report:
        limit: 10
        window: 1h
      update_profile:
        limit: 5
        window: 1h
      like_post:
        limit: 100
        window: 1h
      follow_user:
        limit: 50
        window: 1h
    
    observe:
      create_post:
        limit: 5
        window: 1h
      create_comment:
        limit: 20
        window: 1h
      send_report:
        limit: 5
        window: 1h
      update_profile:
        limit: 3
        window: 1h
      like_post:
        limit: 30
        window: 1h
      follow_user:
        limit: 15
        window: 1h
    
    restrict:
      create_post:
        limit: 2
        window: 1h
      create_comment:
        limit: 5
        window: 1h
      send_report:
        limit: 2
        window: 1h
      update_profile:
        limit: 1
        window: 1h
      like_post:
        limit: 10
        window: 1h
      follow_user:
        limit: 5
        window: 1h

# 内容过滤配置
content_filter:
  enabled: true
  sensitive_words:
    - "暴力"
    - "色情"
    - "赌博"
    - "毒品"
    - "诈骗"
  custom_patterns:
    - pattern: "\\b\\d{11}\\b"
      action: mask
    - pattern: "\\b[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}\\b"
      action: mask

# 反垃圾配置
anti_spam:
  enabled: true
  max_similar_posts: 3
  similarity_threshold: 0.85
  check_interval: 5m
  ban_duration: 24h

# IP 黑名单
ip_blacklist: []

# 用户黑名单
user_blacklist: []

# 风控日志
audit_log:
  enabled: true
  filename: ./logs/audit.log
  max_size: 100
  max_backups: 30
  max_age: 90
  compress: true
EOF
    
    echo -e "${GREEN}  ✓ Created risk control config: $config_file${NC}"
    return 0
}

create_config_example() {
    local example_file="backend/config/private.example.yaml"
    
    if [ -f "$example_file" ]; then
        echo -e "${YELLOW}  Config example already exists: $example_file${NC}"
        return 0
    fi
    
    mkdir -p config
    cat > "$example_file" << 'EOF'
# Private Configuration Example
# Copy this to private.yaml and fill in your values

email:
  host: smtp.example.com
  port: 587
  username: your-email@example.com
  password: your-password
  from: your-email@example.com
  from_name: Tiny Forum

jwt:
  secret: "change-this-to-a-random-string-32chars"
  expire: 24h
  refresh_expire: 168h
  issuer: "tiny-forum"

database:
  host: localhost
  port: 5432
  user: your_db_user
  password: your_db_password
  dbname: tiny_forum
  sslmode: disable
  timezone: Asia/Shanghai

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

admin:
  email: admin@example.com
  password: change-this-password
  username: admin
  role: super_admin
  score: 10000
EOF
    
    echo -e "${GREEN}  ✓ Created config example: $example_file${NC}"
    return 0
}

create_config_gitignore() {
    local gitignore_file="backend/config/.gitignore"
    
    if [ -f "$gitignore_file" ]; then
        echo -e "${YELLOW}  Config .gitignore already exists${NC}"
        return 0
    fi
    
    cat > "$gitignore_file" << EOF
# 忽略敏感配置
private.yaml
*.key
*.pem
*.crt

# 保留基础配置
!basic.yaml
!risk_control.yaml
!*.example.yaml
EOF
    
    echo -e "${GREEN}  ✓ Created config/.gitignore${NC}"
    return 0
}

setup_configurations() {
    echo ""
    echo "Step 2.5: Setting up configuration files..."
    
    create_basic_config
    create_risk_config
    create_private_config
    create_config_example
    create_config_gitignore
    
    echo -e "${GREEN}  ✓ All configuration files ready${NC}"
}
