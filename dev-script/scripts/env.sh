#!/usr/bin/env bash
# =============================================================================
# scripts/env.sh — TinyForum 环境变量管理主入口
#
# 命令:
#   init      从 YAML 配置文件自动生成 .env（可用 --force 覆盖）
#   generate  同 init，但输出带注释的版本
#   edit      用 $EDITOR 打开 .env 编辑
#   show      显示当前 .env 内容（敏感字段脱敏）
#   check     必填项校验（快速版）
#   validate  全量 schema 校验 + 报告
#   list-keys 列出所有 schema 键及默认值
#   get       获取单个键的值
#   set       设置单个键的值
#   unset     删除单个键
#   diff      对比 .env 与 .env.example
#   export    打印 export 语句（供 eval 使用）
#   sync-yaml 从当前 YAML 文件重新同步 .env（保留手动覆盖值）
#
# 用法:
#   ./scripts/env.sh <command> [options]
#   make env-init          # 通过 Makefile 调用
# =============================================================================
set -euo pipefail

# ── 路径解析 ─────────────────────────────────────────────────────────────────
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
LIB_DIR="${SCRIPT_DIR}/env"
CONFIG_DIR="${PROJECT_ROOT}/config"   # YAML 文件目录，可通过 --config-dir 覆盖

# ── 加载库 ───────────────────────────────────────────────────────────────────
source "${LIB_DIR}/core.sh"
source "${LIB_DIR}/schema.sh"
source "${LIB_DIR}/env_file.sh"
source "${LIB_DIR}/validator.sh"
source "${LIB_DIR}/yaml_parser.sh"

# ── 全局默认参数 ─────────────────────────────────────────────────────────────
ENV_FILE="${PROJECT_ROOT}/.env"
ENV_EXAMPLE="${PROJECT_ROOT}/.env.example"
FORCE=0
ANNOTATED=0
VERBOSE=0

# ── 参数解析（getopt 风格，兼容无 GNU getopt 环境）─────────────────────────
_parse_args() {
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --env-file)    ENV_FILE="$2";    shift 2 ;;
      --example)     ENV_EXAMPLE="$2"; shift 2 ;;
      --config-dir)  CONFIG_DIR="$2";  shift 2 ;;
      --key|-k)      _KEY="$2";        shift 2 ;;
      --value|-v)    _VALUE="$2";      shift 2 ;;
      --required)    _REQUIRED="$2";   shift 2 ;;
      --force|-f)    FORCE=1;          shift   ;;
      --annotated|-a) ANNOTATED=1;     shift   ;;
      --verbose)     VERBOSE=1; LOG_VERBOSE=1; shift ;;
      --help|-h)     _cmd_help;        exit 0  ;;
      *)             shift ;;  # 忽略未知参数
    esac
  done
}

# =============================================================================
# 命令实现
# =============================================================================

# ── init / generate ──────────────────────────────────────────────────────────
_cmd_init() {
  log::section "初始化 .env"

  if [[ -f "$ENV_FILE" && "$FORCE" == "0" ]]; then
    log::warn ".env 已存在: $ENV_FILE"
    log::warn "使用 --force 覆盖，或使用 sync-yaml 做增量同步"
    return 1
  fi

  [[ -f "$ENV_FILE" ]] && file::backup "$ENV_FILE"

  # 从 YAML 提取
  declare -A extracted=()
  if [[ -d "$CONFIG_DIR" ]]; then
    log::step "从 YAML 文件提取配置: $CONFIG_DIR"
    yaml_parser::load_all extracted "$CONFIG_DIR"
  else
    log::warn "config 目录不存在: ${CONFIG_DIR}，将使用 schema 默认值"
  fi

  # 填充缺失键用 schema 默认值
  for key in $(schema::keys); do
    [[ -z "${extracted[$key]+x}" || -z "${extracted[$key]}" ]] && \
      extracted["$key"]="$(schema::default "$key")"
  done

  log::step "写入 .env"
  if [[ "$ANNOTATED" == "1" ]]; then
    env_file::write_annotated "$ENV_FILE" extracted
  else
    env_file::write "$ENV_FILE" extracted _SCHEMA_ORDER
  fi

  # 自动生成 .env.example（敏感字段置空）
  _cmd_gen_example

  log::success "初始化完成"
}

# ── 生成 .env.example ────────────────────────────────────────────────────────
_cmd_gen_example() {
  declare -A example=()
  for key in $(schema::keys); do
    if schema::secret "$key"; then
      example["$key"]=""
    else
      example["$key"]="$(schema::default "$key")"
    fi
  done

  {
    echo "# .env.example — 提交到版本控制的模板，敏感字段留空"
    echo "# 由 env.sh gen-example 生成 — $(date '+%Y-%m-%d %H:%M:%S')"
    echo ""
    for key in $(schema::keys); do
      local desc; desc="$(schema::desc "$key")"
      local req; schema::required "$key" && req="[必填] " || req=""
      echo "# ${req}${desc}"
      echo "${key}=${example[$key]}"
    done
  } | file::atomic_write "$ENV_EXAMPLE"
  log::info "已生成: $ENV_EXAMPLE"
}

# ── sync-yaml ────────────────────────────────────────────────────────────────
_cmd_sync_yaml() {
  log::section "从 YAML 同步 .env（增量更新）"

  declare -A current=() yaml_vals=()

  # 1. 加载现有 .env
  [[ -f "$ENV_FILE" ]] && env_file::load "$ENV_FILE" current

  # 2. 从 YAML 提取最新值
  [[ -d "$CONFIG_DIR" ]] || { log::error "config 目录不存在: $CONFIG_DIR"; return 1; }
  yaml_parser::load_all yaml_vals "$CONFIG_DIR"

  # 3. 合并策略：YAML 的值优先，但手动设置的 secret 不覆盖（除非 --force）
  local updated=0
  for key in $(schema::keys); do
    local yaml_val="${yaml_vals[$key]:-}"
    local cur_val="${current[$key]:-}"

    if [[ -z "$yaml_val" ]]; then continue; fi   # YAML 中无值，跳过

    if schema::secret "$key" && [[ -n "$cur_val" ]] && [[ "$FORCE" == "0" ]]; then
      log::debug "保留敏感键手动值: $key"
      continue
    fi

    if [[ "$yaml_val" != "$cur_val" ]]; then
      log::step "更新: ${key}  ${cur_val:-<empty>} → ${yaml_val}"
      current["$key"]="$yaml_val"
      (( updated++ )) || true
    fi
  done

  if (( updated > 0 )); then
    file::backup "$ENV_FILE"
    env_file::write "$ENV_FILE" current _SCHEMA_ORDER
    log::success "同步完成，更新了 $updated 个键"
  else
    log::info "无变更"
  fi
}

# ── show ─────────────────────────────────────────────────────────────────────
_cmd_show() {
  core::require_file "$ENV_FILE"
  declare -A vals=()
  env_file::load "$ENV_FILE" vals

  log::section "当前 .env 配置"
  local cur_group=""
  for key in $(schema::keys); do
    local group="${key%%_*}"
    [[ "$group" != "$cur_group" ]] && { echo ""; printf "${C_CYAN}[%s]${C_RESET}\n" "$group"; cur_group="$group"; }

    local val="${vals[$key]:-}"
    local display="$val"
    [[ -z "$val" ]] && display="${C_GRAY}(未设置)${C_RESET}"
    schema::secret "$key" && [[ -n "$val" ]] && display="${C_YELLOW}***${C_RESET}"
    printf "  %-35s = %b\n" "$key" "$display"
  done
  echo ""
}

# ── validate ─────────────────────────────────────────────────────────────────
_cmd_validate() {
  core::require_file "$ENV_FILE"
  declare -A vals=()
  env_file::load "$ENV_FILE" vals

  log::section "Schema 校验"
  validator::report vals
  validator::check_all vals
}

# ── check（轻量必填项检查）──────────────────────────────────────────────────
_cmd_check() {
  core::require_file "$ENV_FILE"
  declare -A vals=()
  env_file::load "$ENV_FILE" vals

  local -a required_keys
  IFS=' ' read -ra required_keys <<< "${_REQUIRED:-}"

  # 若未传 --required，从 schema 自动获取
  if [[ ${#required_keys[@]} -eq 0 ]]; then
    for k in $(schema::keys); do
      schema::required "$k" && required_keys+=("$k")
    done
  fi

  log::section "必填项检查"
  local failed=0
  for key in "${required_keys[@]}"; do
    local val="${vals[$key]:-}"
    if [[ -z "$val" ]]; then
      log::error "  ✗ ${key} 未设置"
      (( failed++ )) || true
    else
      log::success "  ✓ ${key}"
    fi
  done
  (( failed == 0 )) || { log::error "共 $failed 个必填项未设置"; return 1; }
}

# ── get / set / unset ────────────────────────────────────────────────────────
_cmd_get() {
  core::require_file "$ENV_FILE"
  [[ -z "${_KEY:-}" ]] && { log::error "用法: env.sh get --key KEY"; return 1; }
  env_file::get_key "$ENV_FILE" "$_KEY"
}

_cmd_set() {
  [[ -z "${_KEY:-}" || -z "${_VALUE:-}" ]] && { log::error "用法: env.sh set --key KEY --value VALUE"; return 1; }
  [[ ! -f "$ENV_FILE" ]] && touch "$ENV_FILE"
  env_file::set_key "$ENV_FILE" "$_KEY" "$_VALUE"
}

_cmd_unset() {
  core::require_file "$ENV_FILE"
  [[ -z "${_KEY:-}" ]] && { log::error "用法: env.sh unset --key KEY"; return 1; }
  env_file::unset_key "$ENV_FILE" "$_KEY"
}

# ── list-keys ─────────────────────────────────────────────────────────────────
_cmd_list_keys() {
  printf "%-40s %-30s %s\n" "KEY" "DEFAULT" "描述"
  printf '%0.s─' {1..90}; echo
  for key in $(schema::keys); do
    local def; def="$(schema::default "$key")"
    schema::secret "$key" && def="***"
    [[ ${#def} -gt 28 ]] && def="${def:0:25}..."
    printf "%-40s %-30s %s\n" "$key" "$def" "$(schema::desc "$key")"
  done
}

# ── diff ─────────────────────────────────────────────────────────────────────
_cmd_diff() {
  core::require_file "$ENV_FILE"
  [[ ! -f "$ENV_EXAMPLE" ]] && { log::warn ".env.example 不存在，先运行 init"; return 1; }
  env_file::diff "$ENV_FILE" "$ENV_EXAMPLE"
}

# ── export ───────────────────────────────────────────────────────────────────
_cmd_export() {
  core::require_file "$ENV_FILE"
  declare -A vals=()
  env_file::load "$ENV_FILE" vals
  for key in $(schema::keys); do
    local val="${vals[$key]:-$(schema::default "$key")}"
    printf 'export %s="%s"\n' "$key" "$val"
  done
}

# ── edit ─────────────────────────────────────────────────────────────────────
_cmd_edit() {
  core::require_file "$ENV_FILE"
  local editor="${EDITOR:-${VISUAL:-vi}}"
  "$editor" "$ENV_FILE"
}

# ── help ─────────────────────────────────────────────────────────────────────
_cmd_help() {
  
   echo -e "$(cat <<HELP
${C_BOLD}TinyForum env.sh — 环境变量管理工具${C_RESET}

${C_CYAN}用法:${C_RESET}
  $(basename "$0") <command> [options]

${C_CYAN}命令:${C_RESET}
  ${C_GREEN}init${C_RESET}          从 YAML 配置文件生成 .env（存在时报错，--force 覆盖）
  ${C_GREEN}generate${C_RESET}      同 init，输出带注释的 .env（-a/--annotated）
  ${C_GREEN}sync-yaml${C_RESET}     从 YAML 增量同步 .env（保留手动覆盖的敏感值）
  ${C_GREEN}show${C_RESET}          显示当前 .env（敏感字段脱敏）
  ${C_GREEN}check${C_RESET}         必填项快速检查
  ${C_GREEN}validate${C_RESET}      全量 schema 校验 + 格式检查
  ${C_GREEN}list-keys${C_RESET}     列出所有 schema 键、默认值、说明
  ${C_GREEN}get${C_RESET}           获取单个键值  --key KEY
  ${C_GREEN}set${C_RESET}           设置单个键值  --key KEY --value VALUE
  ${C_GREEN}unset${C_RESET}         删除单个键    --key KEY
  ${C_GREEN}diff${C_RESET}          对比 .env 与 .env.example
  ${C_GREEN}export${C_RESET}        输出 export 语句（eval 使用）
  ${C_GREEN}edit${C_RESET}          用 \$EDITOR 打开 .env

${C_CYAN}通用选项:${C_RESET}
  --env-file <path>     指定 .env 路径（默认: .env）
  --example <path>      指定 .env.example 路径
  --config-dir <path>   YAML 配置文件目录（默认: ./config）
  --force / -f          强制覆盖（init 和 sync-yaml）
  --annotated / -a      生成带注释的 .env
  --verbose             开启 debug 日志
  --help / -h           显示此帮助

${C_CYAN}示例:${C_RESET}
  $(basename "$0") init --config-dir ./config
  $(basename "$0") sync-yaml --force
  $(basename "$0") set --key DB_PASSWORD --value "my-secret"
  $(basename "$0") validate
  eval "\$($(basename "$0") export)"
HELP
)"
}

# =============================================================================
# 主入口
# =============================================================================
main() {
  core::bash_version_check

  local cmd="${1:-help}"
  shift || true

  _parse_args "$@"

  case "$cmd" in
    init)          _cmd_init ;;
    generate)      ANNOTATED=1 _cmd_init ;;
    sync-yaml)     _cmd_sync_yaml ;;
    show)          _cmd_show ;;
    check)         _cmd_check ;;
    validate)      _cmd_validate ;;
    list-keys)     _cmd_list_keys ;;
    get)           _cmd_get ;;
    set)           _cmd_set ;;
    unset)         _cmd_unset ;;
    diff)          _cmd_diff ;;
    export)        _cmd_export ;;
    edit)          _cmd_edit ;;
    gen-example)   _cmd_gen_example ;;
    help|--help|-h) _cmd_help ;;
    *)
      log::error "未知命令: $cmd"
      _cmd_help
      exit 1 ;;
  esac
}

main "$@"