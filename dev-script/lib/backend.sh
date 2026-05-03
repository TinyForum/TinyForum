

# ============================================
# Module 5: Backend Setup
# ============================================

setup_backend() {
    echo ""
    echo "Step 2: Setting up Backend..."
    
    cd backend
    
    if [ ! -f "go.mod" ]; then
        echo -e "${RED}❌ go.mod not found. Are you in the right directory?${NC}"
        exit 1
    fi
    
    echo "  Running go mod tidy..."
    go mod tidy
    echo -e "${GREEN}  Dependencies downloaded.${NC}"
    
    cd ..
}