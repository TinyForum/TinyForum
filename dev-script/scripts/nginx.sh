#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/nginx/setup.sh"

USE_SITES=false
RELOAD=true
DRY_RUN=false

while [[ $# -gt 0 ]]; do
    case "$1" in
        --use-sites) USE_SITES=true; shift ;;
        --no-reload) RELOAD=false; shift ;;
        --dry-run) DRY_RUN=true; shift ;;
        --help|-h)
            echo "用法: $0 [选项]"
            echo "  --use-sites    使用 sites-available 模式"
            echo "  --no-reload    不重载 Nginx"
            echo "  --dry-run      预览配置"
            exit 0 ;;
        *) echo -e "${RED}未知参数: $1${NC}"; exit 1 ;;
    esac
done

main_setup "$USE_SITES" "$RELOAD" "$DRY_RUN"