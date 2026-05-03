

# ============================================
# Module 6: Frontend Setup
# ============================================

setup_frontend() {
    echo ""
    echo "Step 3: Setting up Frontend..."
    
    cd frontend
    
    if [ ! -f "package.json" ]; then
        echo -e "${RED}❌ package.json not found. Are you in the right directory?${NC}"
        exit 1
    fi
    
    echo "  Installing dependencies with $PACKAGE_MANAGER..."
    if [ "$PACKAGE_MANAGER" = "pnpm" ]; then
        pnpm install
    else
        npm install
    fi
    echo -e "${GREEN}  Frontend dependencies installed.${NC}"
    
    cd ..
}