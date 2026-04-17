.PHONY: help frontend backend docs install clean build dev all init-config
.PHONY: start stop restart status logs stop-all start-all
.PHONY: start-backend stop-backend start-frontend stop-frontend
.PHONY: start-docs stop-docs
.PHONY: bench bench-wrk bench-oha bench-k6 bench-all bench-ci bench-compare install-tools

# PID 文件存储目录
PID_DIR := .pids
BACKEND_PID_FILE := $(PID_DIR)/backend.pid
FRONTEND_PID_FILE := $(PID_DIR)/frontend.pid
DOCS_PID_FILE := $(PID_DIR)/docs.pid

# 日志目录
LOG_DIR := applogs
BACKEND_LOG := $(LOG_DIR)/backend.log
FRONTEND_LOG := $(LOG_DIR)/frontend.log
DOCS_LOG := $(LOG_DIR)/docs.log

# 压测配置
BASE_URL      := http://localhost:8080
API_PREFIX    := /api/v1
TEST_DURATION := 30s
CONCURRENCY   := 100
REQUESTS      := 10000

# 检测可用的压测工具
WHICH_WRK     := $(shell command -v wrk 2>/dev/null)
WHICH_OHA     := $(shell command -v oha 2>/dev/null)
WHICH_K6      := $(shell command -v k6 2>/dev/null)

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



# 主入口：智能选择工具
bench:
	@if [ -n "$(WHICH_OHA)" ]; then \
		$(MAKE) bench-oha; \
	elif [ -n "$(WHICH_WRK)" ]; then \
		$(MAKE) bench-wrk; \
	else \
		echo "❌ 未找到压测工具，运行 'make install-tools' 安装"; \
		exit 1; \
	fi

# 安装压测工具（macOS/Linux）
install-tools:
	@echo "📦 安装压测工具..."
	@if [ "$$(uname)" = "Darwin" ]; then \
		brew install oha wrk k6 2>/dev/null || true; \
	else \
		echo "Linux 请手动安装: cargo install oha | 编译 wrk | 安装 k6"; \
	fi
	@echo "✅ 安装完成"

# --- oha 测试（推荐，实时可视化）---

bench-oha: bench-oha-light bench-oha-medium bench-oha-heavy

bench-oha-light:
	@echo "🚀 oha 轻量测试 (100并发, 1万请求)"
	oha -n $(REQUESTS) -c 100 $(BASE_URL)$(API_PREFIX)/posts

bench-oha-medium:
	@echo "🚀 oha 中等压力 (500并发, 2万请求)"
	oha -n 20000 -c 500 $(BASE_URL)$(API_PREFIX)/posts

bench-oha-heavy:
	@echo "🚀 oha 高压力 (1000并发, 5万请求)"
	oha -n 50000 -c 1000 $(BASE_URL)$(API_PREFIX)/posts

# --- wrk 测试（极限压力）---

bench-wrk: bench-wrk-light bench-wrk-medium bench-wrk-heavy

bench-wrk-light:
	@echo "🔥 wrk 轻量测试 (12线程, 100连接, 30秒)"
	wrk -t12 -c100 -d$(TEST_DURATION) $(BASE_URL)$(API_PREFIX)/posts

bench-wrk-medium:
	@echo "🔥 wrk 中等压力 (12线程, 400连接, 30秒)"
	wrk -t12 -c400 -d$(TEST_DURATION) $(BASE_URL)$(API_PREFIX)/posts

bench-wrk-heavy:
	@echo "🔥 wrk 极限压力 (12线程, 1000连接, 30秒)"
	wrk -t12 -c1000 -d$(TEST_DURATION) $(BASE_URL)$(API_PREFIX)/posts

# --- k6 测试（场景化，CI友好）---

bench-k6:
	@if [ -z "$(WHICH_K6)" ]; then \
		echo "❌ k6 未安装"; \
		exit 1; \
	fi
	@mkdir -p $(LOG_DIR)
	k6 run \
		-e BASE_URL=$(BASE_URL) \
		-e API_PREFIX=$(API_PREFIX) \
		--out json=$(LOG_DIR)/k6-result.json \
		scripts/bench.js

# --- 全工具对比测试 ---

bench-all: bench-check-tools
	@echo "=========================================="
	@echo "   全工具性能对比测试"
	@echo "=========================================="
	@mkdir -p $(LOG_DIR)
	@echo ""
	@if [ -n "$(WHICH_OHA)" ]; then \
		echo ">>> oha 1000并发测试" | tee $(LOG_DIR)/bench-oha.log; \
		oha -n 50000 -c 1000 $(BASE_URL)$(API_PREFIX)/posts 2>&1 | tee -a $(LOG_DIR)/bench-oha.log; \
	fi
	@echo ""
	@if [ -n "$(WHICH_WRK)" ]; then \
		echo ">>> wrk 1000连接测试" | tee $(LOG_DIR)/bench-wrk.log; \
		wrk -t12 -c1000 -d30s $(BASE_URL)$(API_PREFIX)/posts 2>&1 | tee -a $(LOG_DIR)/bench-wrk.log; \
	fi
	@echo ""
	@echo "✅ 测试结果保存至 $(LOG_DIR)/bench-*.log"

# --- CI/CD 自动化测试（非交互式）---

bench-ci:
	@echo "🤖 CI 模式压测 (快速验证)"
	@mkdir -p $(LOG_DIR)
	@if [ -n "$(WHICH_OHA)" ]; then \
		oha -n 5000 -c 50 --no-tui $(BASE_URL)$(API_PREFIX)/posts > $(LOG_DIR)/bench-ci.log 2>&1; \
		echo "✅ CI 测试完成，结果: $(LOG_DIR)/bench-ci.log"; \
	else \
		echo "⚠️  oha 未安装，跳过 CI 测试"; \
	fi

# --- 对比测试（优化前后）---

bench-compare:
	@echo "📊 优化前后对比测试"
	@mkdir -p $(LOG_DIR)
	@echo "第1轮测试 (当前性能)..." 
	-oha -n 10000 -c 100 $(BASE_URL)$(API_PREFIX)/posts > $(LOG_DIR)/before.log 2>&1
	@echo "请应用优化后按回车继续..."; read dummy
	@echo "第2轮测试 (优化后)..."
	-oha -n 10000 -c 100 $(BASE_URL)$(API_PREFIX)/posts > $(LOG_DIR)/after.log 2>&1
	@echo ""
	@echo "=== 对比结果 ==="
	@grep "Requests/sec" $(LOG_DIR)/before.log $(LOG_DIR)/after.log || true
	@grep "Success rate" $(LOG_DIR)/before.log $(LOG_DIR)/after.log || true

# --- 辅助命令 ---

bench-check-tools:
	@echo "检查压测工具..."
	@[ -n "$(WHICH_OHA)" ] && echo "✅ oha: $(WHICH_OHA)" || echo "❌ oha 未安装"
	@[ -n "$(WHICH_WRK)" ] && echo "✅ wrk: $(WHICH_WRK)" || echo "❌ wrk 未安装"
	@[ -n "$(WHICH_K6)" ] && echo "✅ k6:  $(WHICH_K6)" || echo "❌ k6 未安装"

bench-help:
	@echo "=========================================="
	@echo "   性能测试命令"
	@echo "=========================================="
	@echo ""
	@echo "快速开始:"
	@echo "  make bench              智能选择工具运行压测"
	@echo "  make install-tools      安装 oha, wrk, k6"
	@echo ""
	@echo "分级压测 (oha):"
	@echo "  make bench-oha-light    100并发, 1万请求"
	@echo "  make bench-oha-medium   500并发, 2万请求"  
	@echo "  make bench-oha-heavy    1000并发, 5万请求"
	@echo ""
	@echo "极限压测 (wrk):"
	@echo "  make bench-wrk-light    12线程, 100连接"
	@echo "  make bench-wrk-medium   12线程, 400连接"
	@echo "  make bench-wrk-heavy    12线程, 1000连接"
	@echo ""
	@echo "场景测试:"
	@echo "  make bench-k6           渐进加压场景测试"
	@echo "  make bench-all          全工具对比测试"
	@echo ""
	@echo "CI/CD:"
	@echo "  make bench-ci           非交互式快速测试"
	@echo "  make bench-compare      优化前后对比"
	@echo ""
	@echo "配置变量:"
	@echo "  make bench BASE_URL=http://localhost:3000 API_PREFIX=/api"