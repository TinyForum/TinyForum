# ============================================
# TinyForum 主 Makefile
# ============================================

# 获取当前 Makefile 自身的绝对路径（必须在任何 include 之前）
SELF_MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))

# dev-script 目录绝对路径
DEVS_SCRIPT_DIR := $(abspath dev-script)

# 定义各个模块 Makefile 的绝对路径
MAKEFILE_COMMON_PATH := $(DEVS_SCRIPT_DIR)/Makefile.common
MAKEFILE_CLEAN_PATH  := $(DEVS_SCRIPT_DIR)/Makefile.clean
MAKEFILE_CHECK_PATH  := $(DEVS_SCRIPT_DIR)/Makefile.check
MAKEFILE_DEV_PATH    := $(DEVS_SCRIPT_DIR)/Makefile.dev
MAKEFILE_PODMAN_PATH := $(DEVS_SCRIPT_DIR)/Makefile.podman
MAKEFILE_DOCKER_PATH := $(DEVS_SCRIPT_DIR)/Makefile.docker
MAKEFILE_MAIN_PATH   := $(DEVS_SCRIPT_DIR)/Makefile.main
MAKEFILE_BENCH_PATH  := $(DEVS_SCRIPT_DIR)/Makefile.bench
MAKEFILE_CODE_PATH   := $(DEVS_SCRIPT_DIR)/Makefile.code
MAKEFILE_ENV_PATH    := $(DEVS_SCRIPT_DIR)/Makefile.env
MAKEFILE_CFG_PATH    := $(DEVS_SCRIPT_DIR)/Makefile.cfg
MAKEFILE_LOG_PATH    := $(DEVS_SCRIPT_DIR)/Makefile.log

# logo
BANNER_PATH			 := $(DEVS_SCRIPT_DIR)/scripts/dev/banner.txt

# dev shell
SHELL_DEV_PATH		 := $(DEVS_SCRIPT_DIR)/scripts/dev.sh

# 包含 help.mk（提供 _print_help 目标）
include dev-script/help.mk

# 现在 include 所有模块（通过变量路径，方便调试）
include $(MAKEFILE_COMMON_PATH)
include $(MAKEFILE_CLEAN_PATH)
include $(MAKEFILE_CHECK_PATH)
include $(MAKEFILE_DEV_PATH)
include $(MAKEFILE_PODMAN_PATH)
include $(MAKEFILE_DOCKER_PATH)
include $(MAKEFILE_MAIN_PATH)
include $(MAKEFILE_BENCH_PATH)
include $(MAKEFILE_CODE_PATH)
include $(MAKEFILE_ENV_PATH)
include $(MAKEFILE_CFG_PATH)
include $(MAKEFILE_LOG_PATH)




banner:
	@cat $(BANNER_PATH)

help: banner
	@echo ""
	@echo "$(GREEN)────────────────────────────────────────────────────────$(NC)"
	@echo "$(GREEN) TinyForum 可用命令组$(NC)"
	@echo "$(GREEN)────────────────────────────────────────────────────────$(NC)"
	@echo "  $(GREEN)main-help$(NC)        主要帮助信息"
	@echo "  $(GREEN)dev-help$(NC)         开发帮助信息"
	@echo "  $(GREEN)cfg-help$(NC)         配置帮助信息"
	@echo "  $(GREEN)check-help$(NC)       检查帮助信息"
	@echo "  $(GREEN)clean-help$(NC)       清理帮助信息"
	@echo "  $(GREEN)code-help$(NC)        代码帮助信息"
	@echo "  $(GREEN)log-help$(NC)         日志帮助信息"
	@echo "  $(GREEN)bench-help$(NC)       性能帮助信息"
	@echo "  $(GREEN)env-help$(NC)         环境帮助信息"
	@echo "  $(GREEN)docker-help$(NC)      Docker 帮助信息"
	@echo "  $(GREEN)podman-help$(NC)      Podman 帮助信息"
	@echo "$(GREEN)────────────────────────────────────────────────────────$(NC)"
	@echo ""