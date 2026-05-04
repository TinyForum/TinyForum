#!/usr/bin/env bash
# =============================================================================
# lib/validator.sh — 运行时校验：对已加载的 env 值执行 schema 规则
# =============================================================================
[[ -n "${_LIB_VALIDATOR_LOADED:-}" ]] && return 0
readonly _LIB_VALIDATOR_LOADED=1

source "$(dirname "${BASH_SOURCE[0]}")/core.sh"
source "$(dirname "${BASH_SOURCE[0]}")/schema.sh"

# ── 校验单个键 ────────────────────────────────────────────────────────────────
# validator::check_key <KEY> <value>  → 0=pass, 1=fail
validator::check_key() {
  local key="$1" val="$2"

  # 必填检查
  if schema::required "$key" && [[ -z "$val" ]]; then
    echo "必填项为空: ${key}"
    return 1
  fi

  # 自定义校验函数
  local fn; fn="$(schema::validator "$key")"
  if [[ -n "$fn" && -n "$val" ]]; then
    if ! "$fn" "$val" 2>/dev/null; then
      local desc; desc="$(schema::desc "$key")"
      echo "${key}='${val}' 不符合规则 ${fn}（${desc}）"
      return 1
    fi
  fi

  return 0
}

# ── 校验全部键（从关联数组）─────────────────────────────────────────────────
# validator::check_all <assoc_array_name> → 汇总并返回 0/1
validator::check_all() {
  local -n _vals="$1"
  local -a errors=()
  local -a warnings=()

  for key in $(schema::keys); do
    local val="${_vals[$key]:-$(schema::default "$key")}"
    local err
    err="$(validator::check_key "$key" "$val" 2>&1)" || {
      errors+=("$err")
    }
  done

  # 未知键警告
  for k in "${!_vals[@]}"; do
    if ! grep -qx "$k" <(schema::keys); then
      warnings+=("未知键（不在 schema 中）: $k")
    fi
  done

  # 特殊安全检查
  local jwt="${_vals[JWT_SECRET]:-$(schema::default "JWT_SECRET")}"
  if [[ "$jwt" == "tiny-forum-secret-change-in-production-32chars" ]]; then
    local mode="${_vals[SERVER_MODE]:-debug}"
    [[ "$mode" == "release" ]] && errors+=("生产模式下 JWT_SECRET 不能使用默认值！") \
                                || warnings+=("JWT_SECRET 仍为默认值，上生产前必须修改")
  fi

  # 输出结果
  local ok=1
  if (( ${#warnings[@]} > 0 )); then
    log::section "校验警告"
    for w in "${warnings[@]}"; do log::warn "  $w"; done
  fi
  if (( ${#errors[@]} > 0 )); then
    log::section "校验错误"
    for e in "${errors[@]}"; do log::error "  $e"; done
    ok=0
  fi

  if [[ $ok == 1 ]]; then
    log::success "所有必填项和格式校验通过（共 ${#_SCHEMA_ORDER[@]} 个键）"
  fi
  return $(( 1 - ok ))
}

# ── 打印校验报告（彩色表格）─────────────────────────────────────────────────
validator::report() {
  local -n _vals="$1"
  printf "\n${C_BOLD}%-40s %-14s %s${C_RESET}\n" "KEY" "STATUS" "VALUE/DESC"
  printf '%0.s─' {1..80}; echo

  for key in $(schema::keys); do
    local val="${_vals[$key]:-$(schema::default "$key")}"
    local err
    local status icon
    err="$(validator::check_key "$key" "$val" 2>&1)"
    if [[ $? -ne 0 ]]; then
      icon="${C_RED}✗ FAIL  ${C_RESET}"
    elif [[ -z "$val" ]] && ! schema::required "$key"; then
      icon="${C_GRAY}○ EMPTY ${C_RESET}"
    else
      icon="${C_GREEN}✓ OK    ${C_RESET}"
    fi

    local display_val="$val"
    schema::secret "$key" && display_val="***(已隐藏)***"
    [[ ${#display_val} -gt 30 ]] && display_val="${display_val:0:27}..."
    printf "%-40s %b%-14s%b %s\n" "$key" "" "$icon" "" "$display_val"
  done
  echo
}