

# ============================================
# Module 3: PostgreSQL Management
# ============================================

check_postgres_running() {
    case "$OS" in
        Darwin*)
            if brew services list | grep -q "postgresql.*started" || pg_isready >/dev/null 2>&1; then
                return 0
            fi
            ;;
        Linux*)
            if systemctl is-active --quiet postgresql 2>/dev/null || \
               systemctl is-active --quiet postgresql@*-main 2>/dev/null || \
               pg_isready >/dev/null 2>&1; then
                return 0
            fi
            ;;
        MINGW*|MSYS*|CYGWIN*)
            if net start | grep -qi "postgres" 2>/dev/null || pg_isready >/dev/null 2>&1; then
                return 0
            fi
            ;;
        *)
            if pg_isready >/dev/null 2>&1; then
                return 0
            fi
            ;;
    esac
    return 1
}

connect_postgres() {
    local user=$1
    local db=${2:-postgres}
    
    if [ "$user" = "postgres" ] && [ "$OS" = "Linux" ]; then
        sudo -u postgres psql -d "$db" -c "SELECT 1" >/dev/null 2>&1
    else
        psql -h localhost -U "$user" -d "$db" -c "SELECT 1" >/dev/null 2>&1
    fi
}

database_exists() {
    if [ "$DB_USER" = "postgres" ] && [ "$OS" = "Linux" ]; then
        sudo -u postgres psql -lqt | cut -d \| -f 1 | grep -qw "tiny_forum"
    else
        psql -h localhost -U "$DB_USER" -lqt | cut -d \| -f 1 | grep -qw "tiny_forum"
    fi
}

create_database() {
    if [ "$DB_USER" = "postgres" ] && [ "$OS" = "Linux" ]; then
        sudo -u postgres createdb tiny_forum 2>/dev/null || \
        sudo -u postgres psql -c "CREATE DATABASE tiny_forum;" 2>/dev/null
    else
        createdb -h localhost -U "$DB_USER" tiny_forum 2>/dev/null || \
        psql -h localhost -U "$DB_USER" -c "CREATE DATABASE tiny_forum;" 2>/dev/null
    fi
}

user_exists() {
    local username=$1
    if [ "$DB_USER" = "postgres" ] && [ "$OS" = "Linux" ]; then
        sudo -u postgres psql -tAc "SELECT 1 FROM pg_roles WHERE rolname='$username'"
    else
        psql -h localhost -U "$DB_USER" -tAc "SELECT 1 FROM pg_roles WHERE rolname='$username'"
    fi
}

create_database_user() {
    local username=$1
    local password=$2
    
    if [ "$DB_USER" = "postgres" ] && [ "$OS" = "Linux" ]; then
        sudo -u postgres psql << EOF
        CREATE USER $username WITH PASSWORD '$password';
        ALTER USER $username CREATEDB;
        GRANT ALL PRIVILEGES ON DATABASE tiny_forum TO $username;
EOF
    else
        psql -h localhost -U "$DB_USER" << EOF
        CREATE USER $username WITH PASSWORD '$password';
        ALTER USER $username CREATEDB;
        GRANT ALL PRIVILEGES ON DATABASE tiny_forum TO $username;
EOF
    fi
}

grant_schema_privileges() {
    local username=$1
    
    if [ "$DB_USER" = "postgres" ] && [ "$OS" = "Linux" ]; then
        sudo -u postgres psql -d tiny_forum << EOF
        GRANT ALL PRIVILEGES ON SCHEMA public TO $username;
        ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO $username;
EOF
    else
        psql -h localhost -U "$DB_USER" -d tiny_forum << EOF
        GRANT ALL PRIVILEGES ON SCHEMA public TO $username;
        ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO $username;
EOF
    fi
}

test_db_connection() {
    if [ -n "$NEW_DB_PASS" ]; then
        PGPASSWORD=$NEW_DB_PASS psql -h localhost -U "$DB_USER" -d tiny_forum -c "SELECT 'Database connected successfully!' as message;" >/dev/null 2>&1
    else
        if [ "$DB_USER" = "postgres" ] && [ "$OS" = "Linux" ]; then
            sudo -u postgres psql -d tiny_forum -c "SELECT 'Database connected successfully!' as message;" >/dev/null 2>&1
        else
            psql -h localhost -U "$DB_USER" -d tiny_forum -c "SELECT 'Database connected successfully!' as message;" >/dev/null 2>&1
        fi
    fi
}

setup_postgres() {
    echo ""
    echo "Step 1: Checking PostgreSQL..."
    
    if ! check_postgres_running; then
        echo -e "${RED}❌ PostgreSQL is not running${NC}"
        echo "   Please start PostgreSQL:"
        case "$OS" in
            Darwin*)
                echo "   - macOS: brew services start postgresql"
                ;;
            Linux*)
                echo "   - Linux: sudo systemctl start postgresql"
                ;;
            MINGW*|MSYS*|CYGWIN*)
                echo "   - Windows: net start postgresql-x64-15"
                ;;
        esac
        exit 1
    fi
    echo -e "${GREEN}✅ PostgreSQL is running${NC}"
    
    # Try to connect
    if connect_postgres "$DB_USER"; then
        echo -e "${GREEN}✅ Connected to PostgreSQL as user: $DB_USER${NC}"
    else
        echo -e "${YELLOW}⚠️  Cannot connect as $DB_USER, trying default users...${NC}"
        for alt_user in "postgres" "$CURRENT_USER"; do
            if [ "$alt_user" != "$DB_USER" ] && connect_postgres "$alt_user"; then
                echo -e "${GREEN}✅ Connected as $alt_user${NC}"
                DB_USER=$alt_user
                echo -e "${GREEN}   Switching to database user: $DB_USER${NC}"
                break
            fi
        done
        
        if ! connect_postgres "$DB_USER"; then
            echo -e "${RED}❌ Cannot connect to PostgreSQL${NC}"
            echo "   Possible solutions:"
            echo "   1. Set environment variable: export DB_USER=your_username"
            echo "   2. Create a database user: createuser -s $(whoami)"
            exit 1
        fi
    fi
    
    # Create database if not exists
    echo "  Checking database 'tiny_forum'..."
    if database_exists; then
        echo -e "${GREEN}  Database 'tiny_forum' already exists${NC}"
    else
        echo "  Creating database 'tiny_forum'..."
        if create_database; then
            echo -e "${GREEN}  Database 'tiny_forum' created${NC}"
        else
            echo -e "${RED}  Failed to create database${NC}"
            exit 1
        fi
    fi
    
    # Ask to create new user
    echo ""
    read -p "Create a new database user? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        read -p "Enter username: " NEW_DB_USER
        read -sp "Enter password: " NEW_DB_PASS
        echo
        
        if user_exists "$NEW_DB_USER" | grep -q 1; then
            echo -e "${YELLOW}⚠️  User $NEW_DB_USER already exists${NC}"
        else
            create_database_user "$NEW_DB_USER" "$NEW_DB_PASS"
            echo -e "${GREEN}✅ User $NEW_DB_USER created${NC}"
        fi
        
        grant_schema_privileges "$NEW_DB_USER"
        DB_USER=$NEW_DB_USER
        echo -e "${GREEN}✅ Now using database user: $DB_USER${NC}"
    fi
}