.PHONY: help frontend backend docs install clean build dev all init-config
.PHONY: start stop restart status logs stop-all start-all
.PHONY: start-backend stop-backend start-frontend stop-frontend
.PHONY: start-docs stop-docs

# PID 文件存储目录
PID_DIR := .pids
BACKEND_PID_FILE := $(PID_DIR)/backend.pid
FRONTEND_PID_FILE := $(PID_DIR)/frontend.pid
DOCS_PID_FILE := $(PID_DIR)/docs.pid

# 日志目录
LOG_DIR := logs
BACKEND_LOG := $(LOG_DIR)/backend.log
FRONTEND_LOG := $(LOG_DIR)/frontend.log
DOCS_LOG := $(LOG_DIR)/docs.log

init-config:
	@echo "🚀 Initializing project..."
	@bash ./start-dev.sh

# Default target
help:
	@echo "Available targets:"
	@echo ""
	@echo "Development (foreground):"
	@echo "  make frontend  - Start frontend development server (foreground)"
	@echo "  make backend   - Start backend development server (foreground)"
	@echo "  make docs      - Start docsify documentation server (foreground)"
	@echo "  make dev       - Run both frontend and backend concurrently"
	@echo ""
	@echo "Daemon management (background):"
	@echo "  make start      - Start all services in background"
	@echo "  make start-backend  - Start backend only in background"
	@echo "  make start-frontend - Start frontend only in background"
	@echo "  make start-docs     - Start docs only in background"
	@echo "  make stop       - Stop all background services"
	@echo "  make stop-backend   - Stop backend only"
	@echo "  make stop-frontend  - Stop frontend only"
	@echo "  make stop-docs      - Stop docs only"
	@echo "  make restart    - Restart all services"
	@echo "  make status     - Check status of all services"
	@echo "  make logs       - Tail all logs"
	@echo "  make logs-backend   - Tail backend logs"
	@echo "  make logs-frontend  - Tail frontend logs"
	@echo "  make logs-docs      - Tail docs logs"
	@echo ""
	@echo "Setup & Build:"
	@echo "  make install   - Install all dependencies"
	@echo "  make build     - Build production versions"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make all       - Install and build everything"
	@echo "  make api       - Generate Swagger API documentation"
	@echo "  make test      - Run tests"

# 创建必要的目录
$(PID_DIR):
	@mkdir -p $(PID_DIR)

$(LOG_DIR):
	@mkdir -p $(LOG_DIR)

# ============================================
# 前台运行模式 (Foreground)
# ============================================

# Frontend development server
frontend:
	@echo "🚀 Starting frontend development server..."
	@cd frontend && pnpm run dev

# Backend development server
backend:
	@echo "🚀 Starting backend development server..."
	@cd backend && go run ./cmd/server/main.go

# Documentation server
docs:
	@echo "📚 Starting documentation server..."
	@cd docs && docsify serve

api:
	@echo "📚 Starting API documentation..."
	@cd backend && swag init -g ./cmd/server/main.go --output ./docs
	@echo "Open Swagger API documentation at http://localhost:8080/swagger/index.html"

# ============================================
# 后台运行模式 (Daemon/Background)
# ============================================

# 启动所有服务
start: | $(PID_DIR) $(LOG_DIR) start-backend start-frontend start-docs
	@echo "✅ All services started"
	@$(MAKE) status

# 启动后端服务（后台）
start-backend: | $(PID_DIR) $(LOG_DIR)
	@if [ -f $(BACKEND_PID_FILE) ] && kill -0 $$(cat $(BACKEND_PID_FILE)) 2>/dev/null; then \
		echo "⚠️  Backend is already running (PID: $$(cat $(BACKEND_PID_FILE)))"; \
	else \
		echo "🚀 Starting backend server in background..."; \
		cd backend && nohup go run ./cmd/server/main.go > ../$(BACKEND_LOG) 2>&1 & \
		echo $$! > $(BACKEND_PID_FILE); \
		sleep 1; \
		if kill -0 $$(cat $(BACKEND_PID_FILE)) 2>/dev/null; then \
			echo "✅ Backend started successfully (PID: $$(cat $(BACKEND_PID_FILE)))"; \
			echo "📝 Logs: $(BACKEND_LOG)"; \
		else \
			echo "❌ Backend failed to start. Check logs: $(BACKEND_LOG)"; \
			rm -f $(BACKEND_PID_FILE); \
			exit 1; \
		fi \
	fi

# 启动前端服务（后台）
start-frontend: | $(PID_DIR) $(LOG_DIR)
	@if [ -f $(FRONTEND_PID_FILE) ] && kill -0 $$(cat $(FRONTEND_PID_FILE)) 2>/dev/null; then \
		echo "⚠️  Frontend is already running (PID: $$(cat $(FRONTEND_PID_FILE)))"; \
	else \
		echo "🚀 Starting frontend server in background..."; \
		cd frontend && nohup pnpm run dev > ../$(FRONTEND_LOG) 2>&1 & \
		echo $$! > $(FRONTEND_PID_FILE); \
		sleep 2; \
		if kill -0 $$(cat $(FRONTEND_PID_FILE)) 2>/dev/null; then \
			echo "✅ Frontend started successfully (PID: $$(cat $(FRONTEND_PID_FILE)))"; \
			echo "📝 Logs: $(FRONTEND_LOG)"; \
		else \
			echo "❌ Frontend failed to start. Check logs: $(FRONTEND_LOG)"; \
			rm -f $(FRONTEND_PID_FILE); \
			exit 1; \
		fi \
	fi

# 启动文档服务（后台）
start-docs: | $(PID_DIR) $(LOG_DIR)
	@if [ -f $(DOCS_PID_FILE) ] && kill -0 $$(cat $(DOCS_PID_FILE)) 2>/dev/null; then \
		echo "⚠️  Docs server is already running (PID: $$(cat $(DOCS_PID_FILE)))"; \
	else \
		echo "📚 Starting docs server in background..."; \
		cd docs && nohup docsify serve > ../$(DOCS_LOG) 2>&1 & \
		echo $$! > $(DOCS_PID_FILE); \
		sleep 1; \
		if kill -0 $$(cat $(DOCS_PID_FILE)) 2>/dev/null; then \
			echo "✅ Docs server started successfully (PID: $$(cat $(DOCS_PID_FILE)))"; \
			echo "📝 Logs: $(DOCS_LOG)"; \
		else \
			echo "❌ Docs server failed to start. Check logs: $(DOCS_LOG)"; \
			rm -f $(DOCS_PID_FILE); \
			exit 1; \
		fi \
	fi

# ============================================
# 停止服务
# ============================================

# 停止所有服务
stop: stop-backend stop-frontend stop-docs
	@echo "✅ All services stopped"

# 停止后端服务（优雅停止）
stop-backend:
	@if [ -f $(BACKEND_PID_FILE) ]; then \
		PID=$$(cat $(BACKEND_PID_FILE)); \
		if kill -0 $$PID 2>/dev/null; then \
			echo "🛑 Stopping backend server (PID: $$PID)..."; \
			kill -TERM $$PID; \
			sleep 2; \
			if kill -0 $$PID 2>/dev/null; then \
				echo "⚠️  Process didn't stop, forcing..."; \
				kill -KILL $$PID; \
			fi; \
			echo "✅ Backend stopped"; \
		else \
			echo "⚠️  Backend not running (stale PID file)"; \
		fi; \
		rm -f $(BACKEND_PID_FILE); \
	else \
		echo "⚠️  Backend PID file not found (service not running)"; \
	fi

# 停止前端服务
stop-frontend:
	@if [ -f $(FRONTEND_PID_FILE) ]; then \
		PID=$$(cat $(FRONTEND_PID_FILE)); \
		if kill -0 $$PID 2>/dev/null; then \
			echo "🛑 Stopping frontend server (PID: $$PID)..."; \
			kill -TERM $$PID; \
			sleep 2; \
			if kill -0 $$PID 2>/dev/null; then \
				echo "⚠️  Process didn't stop, forcing..."; \
				kill -KILL $$PID; \
			fi; \
			echo "✅ Frontend stopped"; \
		else \
			echo "⚠️  Frontend not running (stale PID file)"; \
		fi; \
		rm -f $(FRONTEND_PID_FILE); \
	else \
		echo "⚠️  Frontend PID file not found (service not running)"; \
	fi

# 停止文档服务
stop-docs:
	@if [ -f $(DOCS_PID_FILE) ]; then \
		PID=$$(cat $(DOCS_PID_FILE)); \
		if kill -0 $$PID 2>/dev/null; then \
			echo "🛑 Stopping docs server (PID: $$PID)..."; \
			kill -TERM $$PID; \
			sleep 1; \
			if kill -0 $$PID 2>/dev/null; then \
				echo "⚠️  Process didn't stop, forcing..."; \
				kill -KILL $$PID; \
			fi; \
			echo "✅ Docs server stopped"; \
		else \
			echo "⚠️  Docs server not running (stale PID file)"; \
		fi; \
		rm -f $(DOCS_PID_FILE); \
	else \
		echo "⚠️  Docs PID file not found (service not running)"; \
	fi

# ============================================
# 强制停止 (立即 kill)
# ============================================

force-stop: force-stop-backend force-stop-frontend force-stop-docs
	@echo "✅ All services force-stopped"

force-stop-backend:
	@if [ -f $(BACKEND_PID_FILE) ]; then \
		PID=$$(cat $(BACKEND_PID_FILE)); \
		if kill -0 $$PID 2>/dev/null; then \
			echo "💀 Force stopping backend (PID: $$PID)..."; \
			kill -KILL $$PID; \
			sleep 1; \
		fi; \
		rm -f $(BACKEND_PID_FILE); \
		echo "✅ Backend force-stopped"; \
	else \
		echo "⚠️  Backend PID file not found"; \
	fi

force-stop-frontend:
	@if [ -f $(FRONTEND_PID_FILE) ]; then \
		PID=$$(cat $(FRONTEND_PID_FILE)); \
		if kill -0 $$PID 2>/dev/null; then \
			echo "💀 Force stopping frontend (PID: $$PID)..."; \
			kill -KILL $$PID; \
			sleep 1; \
		fi; \
		rm -f $(FRONTEND_PID_FILE); \
		echo "✅ Frontend force-stopped"; \
	else \
		echo "⚠️  Frontend PID file not found"; \
	fi

force-stop-docs:
	@if [ -f $(DOCS_PID_FILE) ]; then \
		PID=$$(cat $(DOCS_PID_FILE)); \
		if kill -0 $$PID 2>/dev/null; then \
			echo "💀 Force stopping docs server (PID: $$PID)..."; \
			kill -KILL $$PID; \
			sleep 1; \
		fi; \
		rm -f $(DOCS_PID_FILE); \
		echo "✅ Docs server force-stopped"; \
	else \
		echo "⚠️  Docs PID file not found"; \
	fi

# ============================================
# 重启服务
# ============================================

restart: stop start
	@echo "🔄 All services restarted"

restart-backend: stop-backend start-backend
	@echo "🔄 Backend restarted"

restart-frontend: stop-frontend start-frontend
	@echo "🔄 Frontend restarted"

restart-docs: stop-docs start-docs
	@echo "🔄 Docs server restarted"

# ============================================
# 状态检查
# ============================================

status:
	@echo "=========================================="
	@echo "📊 Service Status"
	@echo "=========================================="
	@echo ""
	@# Backend status
	@if [ -f $(BACKEND_PID_FILE) ]; then \
		PID=$$(cat $(BACKEND_PID_FILE)); \
		if kill -0 $$PID 2>/dev/null; then \
			echo "✅ Backend:     RUNNING (PID: $$PID)"; \
			ps -p $$PID -o pid,vsz,rss,pcpu,comm --no-headers 2>/dev/null || true; \
		else \
			echo "❌ Backend:     STOPPED (stale PID file)"; \
			rm -f $(BACKEND_PID_FILE); \
		fi \
	else \
		echo "❌ Backend:     STOPPED"; \
	fi
	@echo ""
	@# Frontend status
	@if [ -f $(FRONTEND_PID_FILE) ]; then \
		PID=$$(cat $(FRONTEND_PID_FILE)); \
		if kill -0 $$PID 2>/dev/null; then \
			echo "✅ Frontend:    RUNNING (PID: $$PID)"; \
			ps -p $$PID -o pid,vsz,rss,pcpu,comm --no-headers 2>/dev/null || true; \
		else \
			echo "❌ Frontend:    STOPPED (stale PID file)"; \
			rm -f $(FRONTEND_PID_FILE); \
		fi \
	else \
		echo "❌ Frontend:    STOPPED"; \
	fi
	@echo ""
	@# Docs status
	@if [ -f $(DOCS_PID_FILE) ]; then \
		PID=$$(cat $(DOCS_PID_FILE)); \
		if kill -0 $$PID 2>/dev/null; then \
			echo "✅ Docs:        RUNNING (PID: $$PID)"; \
			ps -p $$PID -o pid,vsz,rss,pcpu,comm --no-headers 2>/dev/null || true; \
		else \
			echo "❌ Docs:        STOPPED (stale PID file)"; \
			rm -f $(DOCS_PID_FILE); \
		fi \
	else \
		echo "❌ Docs:        STOPPED"; \
	fi
	@echo ""
	@echo "=========================================="
	@echo "📁 PID directory: $(PID_DIR)"
	@echo "📝 Log directory: $(LOG_DIR)"
	@echo "=========================================="

# ============================================
# 日志查看
# ============================================

logs:
	@echo "📝 Tailing all logs (Ctrl+C to exit)..."
	@tail -f $(BACKEND_LOG) $(FRONTEND_LOG) $(DOCS_LOG) 2>/dev/null || echo "No log files found. Run 'make start' first."

logs-backend:
	@echo "📝 Tailing backend logs..."
	@tail -f $(BACKEND_LOG) 2>/dev/null || echo "No backend log file found. Run 'make start-backend' first."

logs-frontend:
	@echo "📝 Tailing frontend logs..."
	@tail -f $(FRONTEND_LOG) 2>/dev/null || echo "No frontend log file found. Run 'make start-frontend' first."

logs-docs:
	@echo "📝 Tailing docs logs..."
	@tail -f $(DOCS_LOG) 2>/dev/null || echo "No docs log file found. Run 'make start-docs' first."

# ============================================
# 清理日志和 PID 文件
# ============================================

clean-logs:
	@echo "🧹 Cleaning log files..."
	@rm -rf $(LOG_DIR)
	@mkdir -p $(LOG_DIR)
	@echo "✅ Logs cleaned"

clean-pids:
	@echo "🧹 Cleaning PID files..."
	@rm -rf $(PID_DIR)
	@echo "✅ PID files cleaned"

# 增强的 clean 命令，同时清理日志和 PID
clean: clean-pids clean-logs
	@echo "🧹 Cleaning build artifacts..."
	@cd frontend && rm -rf node_modules .next dist build 2>/dev/null || true
	@cd backend && rm -rf bin 2>/dev/null || true
	@echo "✅ Clean complete"

# ============================================
# 安装依赖
# ============================================

# Install dependencies
install:
	@echo "📦 Installing dependencies..."
	@cd frontend && pnpm install
	@cd backend && go mod download
	@echo "✅ Dependencies installed"

# Build production versions
build:
	@echo "🔨 Building frontend..."
	@cd frontend && pnpm run build
	@echo "🔨 Building backend..."
	@cd backend && go build -o bin/server ./cmd/server/main.go
	@echo "✅ Build complete"

# Install and build everything
all: install build
	@echo "✅ All tasks completed"

# Run tests
test:
	@echo "🧪 Running tests..."
	@cd frontend && pnpm test || true
	@cd backend && go test ./... || true

# Run both frontend and backend concurrently (foreground)
dev:
	@echo "🚀 Starting both frontend and backend (foreground)..."
	@if command -v concurrently >/dev/null 2>&1; then \
		concurrently "make frontend" "make backend"; \
	else \
		echo "⚠️  'concurrently' not found. Installing..."; \
		pnpm install -g concurrently && concurrently "make frontend" "make backend"; \
	fi

# ============================================
# Docker Compose 管理命令
# ============================================

# Docker 相关变量
DOCKER_COMPOSE := docker-compose
ENV_FILE := .env

# 启动所有服务
docker-up:
	@echo "🐳 Starting all services..."
	@$(DOCKER_COMPOSE) --env-file $(ENV_FILE) up -d
	@$(MAKE) docker-status

# 启动服务并查看日志
docker-up-logs:
	@echo "🐳 Starting all services with logs..."
	@$(DOCKER_COMPOSE) --env-file $(ENV_FILE) up

# 停止所有服务
docker-down:
	@echo "🛑 Stopping all services..."
	@$(DOCKER_COMPOSE) --env-file $(ENV_FILE) down

# 停止并删除 volumes（会丢失数据）
docker-down-volumes:
	@echo "⚠️  Stopping and removing volumes (data will be lost)..."
	@$(DOCKER_COMPOSE) --env-file $(ENV_FILE) down -v

# 重启服务
docker-restart:
	@echo "🔄 Restarting all services..."
	@$(DOCKER_COMPOSE) --env-file $(ENV_FILE) restart
	@$(MAKE) docker-status

# 查看服务状态
docker-status:
	@echo "📊 Docker services status:"
	@$(DOCKER_COMPOSE) --env-file $(ENV_FILE) ps

# 查看日志
docker-logs:
	@echo "📝 Tailing all logs..."
	@$(DOCKER_COMPOSE) --env-file $(ENV_FILE) logs -f

# 查看特定服务日志
docker-logs-backend:
	@$(DOCKER_COMPOSE) --env-file $(ENV_FILE) logs -f backend

docker-logs-frontend:
	@$(DOCKER_COMPOSE) --env-file $(ENV_FILE) logs -f frontend

docker-logs-postgres:
	@$(DOCKER_COMPOSE) --env-file $(ENV_FILE) logs -f postgres

# 进入容器
docker-shell-backend:
	@$(DOCKER_COMPOSE) --env-file $(ENV_FILE) exec backend /bin/sh

docker-shell-frontend:
	@$(DOCKER_COMPOSE) --env-file $(ENV_FILE) exec frontend /bin/sh

docker-shell-postgres:
	@$(DOCKER_COMPOSE) --env-file $(ENV_FILE) exec postgres psql -U ${DB_USER:-postgres} -d ${DB_NAME:-tiny_forum}

# 重新构建服务
docker-build:
	@echo "🔨 Building Docker images..."
	@$(DOCKER_COMPOSE) --env-file $(ENV_FILE) build --no-cache

docker-build-backend:
	@$(DOCKER_COMPOSE) --env-file $(ENV_FILE) build --no-cache backend

docker-build-frontend:
	@$(DOCKER_COMPOSE) --env-file $(ENV_FILE) build --no-cache frontend

# 清理 Docker 资源
docker-clean:
	@echo "🧹 Cleaning Docker resources..."
	@$(DOCKER_COMPOSE) --env-file $(ENV_FILE) down -v --rmi local
	@docker system prune -f

# 备份数据库
docker-backup-db:
	@echo "💾 Backing up database..."
	@$(DOCKER_COMPOSE) --env-file $(ENV_FILE) exec postgres pg_dump -U ${DB_USER:-postgres} ${DB_NAME:-tiny_forum} > backup_$$(date +%Y%m%d_%H%M%S).sql

# 启动开发环境（带热重载）
docker-dev:
	@echo "🛠️  Starting development environment..."
	@$(DOCKER_COMPOSE) --env-file $(ENV_FILE) -f docker-compose.yml -f docker-compose.override.yml up

# 启动工具服务（adminer等）
docker-tools:
	@echo "🔧 Starting tools services..."
	@$(DOCKER_COMPOSE) --env-file $(ENV_FILE) --profile tools up -d