#!/usr/bin/env bash
# =============================================================================
# lib/core.sh — 核心工具库：日志、颜色、错误处理、工具函数
# =============================================================================
set -euo pipefail

# ── Guard（防止重复 source）─────────────────────────────────────────────────
[[ -n "${_LIB_CORE_LOADED:-}" ]] && return 0
readonly _LIB_CORE_LOADED=1

# ── 颜色（tty 时启用）────────────────────────────────────────────────────────
if [[ -t 2 ]]; then
  C_RESET='\033[0m'; C_BOLD='\033[1m'
  C_RED='\033[31m';  C_GREEN='\033[32m'
  C_YELLOW='\033[33m'; C_BLUE='\033[34m'
  C_CYAN='\033[36m'; C_GRAY='\033[90m'
else
  C_RESET=''; C_BOLD=''; C_RED=''; C_GREEN=''
  C_YELLOW=''; C_BLUE=''; C_CYAN=''; C_GRAY=''
fi
readonly C_RESET C_BOLD C_RED C_GREEN C_YELLOW C_BLUE C_CYAN C_GRAY

# ── 日志级别 ─────────────────────────────────────────────────────────────────
declare -g LOG_VERBOSE=${LOG_VERBOSE:-0}   # 1 = 开启 debug 输出

log::info()    { printf "${C_GREEN}[INFO]${C_RESET}  %s\n"  "$*" >&2; }
log::warn()    { printf "${C_YELLOW}[WARN]${C_RESET}  %s\n" "$*" >&2; }
log::error()   { printf "${C_RED}[ERROR]${C_RESET} %s\n"    "$*" >&2; }
log::debug()   { [[ "$LOG_VERBOSE" == "1" ]] && printf "${C_GRAY}[DEBUG]${C_RESET} %s\n" "$*" >&2 || true; }
log::success() { printf "${C_GREEN}${C_BOLD}[OK]${C_RESET}    %s\n" "$*" >&2; }
log::section() { printf "\n${C_CYAN}${C_BOLD}══ %s ══${C_RESET}\n" "$*" >&2; }
log::step()    { printf "${C_BLUE}  →${C_RESET} %s\n" "$*" >&2; }

# ── 断言工具 ─────────────────────────────────────────────────────────────────
core::require_cmd() {
  local cmd="$1"
  if ! command -v "$cmd" &>/dev/null; then
    log::error "必须安装命令: $cmd"
    return 1
  fi
}

core::require_file() {
  local file="$1"
  if [[ ! -f "$file" ]]; then
    log::error "文件不存在: $file"
    return 1
  fi
}

core::require_var() {
  local name="$1" val="${!1:-}"
  if [[ -z "$val" ]]; then
    log::error "必填变量未设置: $name"
    return 1
  fi
}

# ── 字符串工具 ────────────────────────────────────────────────────────────────
str::trim()    { local s="$1"; s="${s#"${s%%[![:space:]]*}"}"; echo "${s%"${s##*[![:space:]]}"}"; }
str::upper()   { echo "${1^^}"; }
str::lower()   { echo "${1,,}"; }
str::slug()    { echo "$1" | tr '[:upper:]' '[:lower:]' | tr -cs 'a-z0-9' '_' | sed 's/^_//;s/_$//'; }

# ── 文件工具 ─────────────────────────────────────────────────────────────────
file::backup() {
  local src="$1"
  [[ ! -f "$src" ]] && return 0
  local bak="${src}.bak.$(date +%Y%m%d_%H%M%S)"
  cp "$src" "$bak"
  log::debug "备份: $src → $bak"
  echo "$bak"
}

file::atomic_write() {
  # 原子写入：先写临时文件再 mv，避免写一半被读到
  local dest="$1"; shift
  local tmp
  tmp="$(mktemp "${dest}.XXXXXX")"
  # 把调用方传入的内容写入 tmp
  cat > "$tmp"
  mv "$tmp" "$dest"
}

# ── YAML 轻量解析（纯 bash，无依赖）─────────────────────────────────────────
# yaml::get <file> <key.path>  → 输出标量值（不支持复杂嵌套列表）
yaml::get() {
  local file="$1" keypath="$2"
  core::require_cmd python3
  python3 - "$file" "$keypath" <<'PYEOF'
import sys, re

def load_simple_yaml(path):
    """极简 YAML 加载器：支持标量、嵌套 mapping、简单 sequence"""
    result = {}
    stack = [(-1, result)]
    with open(path) as f:
        for line in f:
            line = line.rstrip('\n')
            stripped = line.lstrip()
            if not stripped or stripped.startswith('#'):
                continue
            indent = len(line) - len(stripped)
            # sequence item
            if stripped.startswith('- '):
                val = stripped[2:].strip().strip('"\'')
                _, parent = stack[-1]
                if isinstance(parent, list):
                    parent.append(val)
                continue
            if ':' in stripped:
                key, _, val = stripped.partition(':')
                key = key.strip()
                val = val.strip().strip('"\'')
                # pop stack to correct indent level
                while len(stack) > 1 and stack[-1][0] >= indent:
                    stack.pop()
                _, parent = stack[-1]
                if isinstance(parent, dict):
                    if val:
                        parent[key] = val
                    else:
                        parent[key] = {}
                        stack.append((indent, parent[key]))
    return result

def get_nested(d, path):
    keys = path.split('.')
    cur = d
    for k in keys:
        if isinstance(cur, dict) and k in cur:
            cur = cur[k]
        else:
            return None
    return cur

data = load_simple_yaml(sys.argv[1])
val = get_nested(data, sys.argv[2])
if val is not None:
    print(val)
PYEOF
}

# ── 版本检查 ─────────────────────────────────────────────────────────────────
core::bash_version_check() {
  local major="${BASH_VERSINFO[0]}"
  if (( major < 4 )); then
    log::error "需要 Bash 4.0+，当前版本: $BASH_VERSION"
    return 1
  fi
}