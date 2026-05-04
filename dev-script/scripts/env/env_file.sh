#!/usr/bin/env bash
# =============================================================================
# lib/env_file.sh — .env 文件读写操作
# =============================================================================
[[ -n "${_LIB_ENV_FILE_LOADED:-}" ]] && return 0
readonly _LIB_ENV_FILE_LOADED=1

source "$(dirname "${BASH_SOURCE[0]}")/core.sh"

# ── 解析 .env 文件到关联数组 ─────────────────────────────────────────────────
# env_file::load <file> <target_assoc_array_name>
env_file::load() {
  local file="$1"
  local -n _target="$2"   # nameref

  [[ ! -f "$file" ]] && { log::warn ".env 不存在: $file"; return 0; }

  local lineno=0
  while IFS= read -r line || [[ -n "$line" ]]; do
    (( lineno++ )) || true
    # 去除行尾注释（不在引号内的 #）
    local stripped
    stripped="$(echo "$line" | sed "s/[[:space:]]*#[^'\"]*$//")"
    stripped="$(str::trim "$stripped")"
    [[ -z "$stripped" ]] && continue

    if [[ "$stripped" =~ ^([A-Za-z_][A-Za-z0-9_]*)=(.*)$ ]]; then
      local key="${BASH_REMATCH[1]}"
      local val="${BASH_REMATCH[2]}"
      # 去除外层引号
      val="${val%\"}" ; val="${val#\"}"
      val="${val%\'}" ; val="${val#\'}"
      _target["$key"]="$val"
      log::debug "  loaded: $key=${val}"
    else
      log::warn ".env:$lineno 格式无效，跳过: $line"
    fi
  done < "$file"
}

# ── 写出 .env（从关联数组）──────────────────────────────────────────────────
# env_file::write <file> <source_assoc_array_name> [key_order_array_name]
env_file::write() {
  local file="$1"
  local -n _src="$2"
  local -n _order="${3:-_SCHEMA_ORDER}"   # 可选：指定输出顺序

  {
    echo "# =============================================================="
    echo "# .env — TinyForum 环境配置"
    echo "# 由 scripts/env.sh 自动生成 — $(date '+%Y-%m-%d %H:%M:%S')"
    echo "# 警告：包含敏感信息，请勿提交至版本控制"
    echo "# =============================================================="
    echo ""

    for key in "${_order[@]}"; do
      [[ -z "${_src[$key]+x}" ]] && continue
      local val="${_src[$key]}"
      # 含空格或特殊字符的值加引号
      if [[ "$val" =~ [[:space:]\$\`\"\'] ]]; then
        echo "${key}=\"${val}\""
      else
        echo "${key}=${val}"
      fi
    done
  } | file::atomic_write "$file"

  log::success "已写入: $file"
}

# ── 写出带注释的 .env（含分组 header 和描述）────────────────────────────────
env_file::write_annotated() {
  local file="$1"
  local -n _src="$2"

  source "$(dirname "${BASH_SOURCE[0]}")/schema.sh"

  local tmp; tmp="$(mktemp)"
  {
    echo "# =============================================================="
    echo "# .env — TinyForum 环境配置（带注释版本）"
    echo "# 由 scripts/env.sh generate 命令生成 — $(date '+%Y-%m-%d %H:%M:%S')"
    echo "# ⚠️  此文件包含敏感信息，已加入 .gitignore，请勿手动提交"
    echo "# =============================================================="
    echo ""

    # 分组输出
    local cur_group=""
    for key in $(schema::keys); do
      local group="${key%%_*}"
      if [[ "$group" != "$cur_group" ]]; then
        echo ""
        echo "# ── ${group} ─────────────────────────────────────────────────────"
        cur_group="$group"
      fi

      local desc; desc="$(schema::desc "$key")"
      local required; schema::required "$key" && required="[必填] " || required=""
      local secret;   schema::secret  "$key" && secret="[敏感] " || secret=""
      echo "# ${required}${secret}${desc}"

      local val="${_src[$key]:-$(schema::default "$key")}"
      if schema::secret "$key" && [[ -n "$val" ]]; then
        # 敏感值：写实际值但在注释里标记
        echo "${key}=${val}"
      else
        echo "${key}=${val}"
      fi
    done
  } > "$tmp"
  mv "$tmp" "$file"
  log::success "已生成注释版 .env: $file"
}

# ── 设置单个键 ────────────────────────────────────────────────────────────────
env_file::set_key() {
  local file="$1" key="$2" val="$3"

  if grep -qE "^${key}=" "$file" 2>/dev/null; then
    # 原地替换
    sed -i.bak "s|^${key}=.*|${key}=${val}|" "$file"
    rm -f "${file}.bak"
    log::info "已更新: ${key}=${val}"
  else
    # 追加
    echo "${key}=${val}" >> "$file"
    log::info "已追加: ${key}=${val}"
  fi
}

# ── 删除单个键 ────────────────────────────────────────────────────────────────
env_file::unset_key() {
  local file="$1" key="$2"
  sed -i.bak "/^${key}=/d" "$file"
  rm -f "${file}.bak"
  log::info "已删除: ${key}"
}

# ── 获取单个键 ────────────────────────────────────────────────────────────────
env_file::get_key() {
  local file="$1" key="$2"
  local val
  val="$(grep -E "^${key}=" "$file" 2>/dev/null | head -1 | cut -d= -f2-)"
  val="${val%\"}" ; val="${val#\"}"
  val="${val%\'}" ; val="${val#\'}"
  echo "$val"
}

# ── diff 与 example 文件 ─────────────────────────────────────────────────────
env_file::diff() {
  local env_file="$1" example="$2"
  declare -A env_keys ex_keys

  env_file::load "$env_file" env_keys
  env_file::load "$example"  ex_keys

  local missing=() extra=()
  for k in "${!ex_keys[@]}"; do
    [[ -z "${env_keys[$k]+x}" ]] && missing+=("$k")
  done
  for k in "${!env_keys[@]}"; do
    [[ -z "${ex_keys[$k]+x}" ]] && extra+=("$k")
  done

  local ok=1
  if (( ${#missing[@]} > 0 )); then
    log::warn "缺少的键（在 example 中存在但 .env 中没有）:"
    printf '  %s  %s\n' "${C_RED}✗${C_RESET}" "${missing[@]}"
    ok=0
  fi
  if (( ${#extra[@]} > 0 )); then
    log::warn "多余的键（在 .env 中存在但 example 中没有）:"
    printf '  %s  %s\n' "${C_YELLOW}?${C_RESET}" "${extra[@]}"
  fi
  [[ $ok == 1 ]] && log::success "diff 通过：无缺失键"
  return $(( 1 - ok ))
}