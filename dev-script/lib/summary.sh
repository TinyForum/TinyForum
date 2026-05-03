

# ============================================
# Module 7: Summary & Final Checks
# ============================================

print_summary() {
    echo ""
    echo "=================================="
    echo -e "${GREEN}✅ Setup complete!${NC}"
    echo ""
    echo "To start the backend:"
    echo "  cd backend && go run ./cmd/server/main.go"
    echo ""
    echo "To start the frontend:"
    echo "  cd frontend && ${PACKAGE_MANAGER} run dev"
    echo ""
    echo "Database connection info:"
    echo "  Host: localhost:5432"
    echo "  User: $DB_USER"
    if [[ -n "$NEW_DB_PASS" ]]; then
        echo "  Password: $NEW_DB_PASS"
    else
        echo "  Password: (empty - using trust authentication)"
    fi
    echo "  Database: tiny_forum"
    echo "=================================="
    
    echo ""
    echo "Testing database connection..."
    if test_db_connection; then
        echo -e "${GREEN}✅ Database connection successful${NC}"
    else
        echo -e "${YELLOW}⚠️  Could not connect to database, but setup completed${NC}"
    fi
    
    echo ""
    echo -e "${GREEN}All services are ready! 🎉${NC}"
}
