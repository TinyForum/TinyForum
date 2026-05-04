#!/usr/bin/env bash
# =============================================================================
# lib/yaml_parser.sh — 从多个 YAML 配置文件提取值并映射到 ENV 键
# =============================================================================
[[ -n "${_LIB_YAML_PARSER_LOADED:-}" ]] && return 0
readonly _LIB_YAML_PARSER_LOADED=1

source "$(dirname "${BASH_SOURCE[0]}")/core.sh"
source "$(dirname "${BASH_SOURCE[0]}")/schema.sh"

# ── Python 驱动的 YAML → flat dict 提取器 ────────────────────────────────────
yaml_parser::flatten() {
  # 将 YAML 文件转为 KEY=value 行，支持嵌套和简单数组
  local file="$1"
  local prefix="${2:-}"   # 可选前缀，如 "BASIC"
  core::require_cmd python3
  python3 - "$file" "$prefix" <<'PYEOF'
import sys, re

def flatten(obj, prefix="", sep="_"):
    items = {}
    if isinstance(obj, dict):
        for k, v in obj.items():
            new_key = f"{prefix}{sep}{k}" if prefix else k
            items.update(flatten(v, new_key, sep))
    elif isinstance(obj, list):
        # join list of scalars as comma-separated; recurse into dicts
        scalars = []
        for item in obj:
            if isinstance(item, dict):
                items.update(flatten(item, prefix, sep))
            elif item is not None:
                scalars.append(str(item))
        if scalars:
            items[prefix] = ",".join(scalars)
    elif isinstance(obj, str) and obj.startswith('{') and obj.endswith('}'):
        # inline dict like {limit: 20, window: 1h}
        for pair in obj[1:-1].split(','):
            pair = pair.strip()
            if ':' in pair:
                pk, pv = pair.split(':', 1)
                new_key = f"{prefix}{sep}{pk.strip()}" if prefix else pk.strip()
                items[new_key] = pv.strip()
    else:
        if obj is not None:
            items[prefix] = str(obj)
    return items

def parse_yaml_simple(path):
    """Pure-Python minimal YAML parser supporting scalars, maps, sequences"""
    import re as _re
    root = {}
    stack = []  # [(indent, dict_ref)]
    
    with open(path, encoding="utf-8") as f:
        lines = f.readlines()

    i = 0
    while i < len(lines):
        raw = lines[i].rstrip('\n')
        stripped = raw.lstrip()
        if not stripped or stripped.startswith('#'):
            i += 1
            continue
        indent = len(raw) - len(stripped)

        # inline comment removal (naive)
        stripped = _re.sub(r'\s+#[^"\']*$', '', stripped)

        if stripped.startswith('- '):
            # sequence item
            val = stripped[2:].strip()
            # handle {k: v, k2: v2} inline maps
            if val.startswith('{') and val.endswith('}'):
                inner = {}
                for pair in val[1:-1].split(','):
                    if ':' in pair:
                        pk, pv = pair.split(':', 1)
                        inner[pk.strip()] = pv.strip()
                if stack:
                    parent = stack[-1][1]
                    if isinstance(parent, list):
                        parent.append(inner)
            else:
                val = val.strip('"\'')
                if stack:
                    parent = stack[-1][1]
                    if isinstance(parent, list):
                        parent.append(val)
            i += 1
            continue

        if ':' in stripped:
            key, _, rest = stripped.partition(':')
            key = key.strip().strip('"\'')
            rest = rest.strip().strip('"\'')

            # pop stack to current level
            while stack and stack[-1][0] >= indent:
                stack.pop()

            if stack:
                parent = stack[-1][1]
            else:
                parent = root

            if rest:
                # scalar value
                parent[key] = rest
            else:
                # peek ahead for sequence
                next_indent = None
                for j in range(i+1, min(i+5, len(lines))):
                    nxt = lines[j].lstrip()
                    if nxt and not nxt.startswith('#'):
                        next_indent = len(lines[j]) - len(lines[j].lstrip())
                        next_stripped = nxt.rstrip()
                        if next_stripped.startswith('- '):
                            parent[key] = []
                            stack.append((indent, parent[key]))
                        else:
                            parent[key] = {}
                            stack.append((indent, parent[key]))
                        break
                else:
                    parent[key] = {}
                    stack.append((indent, parent[key]))
        i += 1
    return root

try:
    data = parse_yaml_simple(sys.argv[1])
    prefix = sys.argv[2].upper() if len(sys.argv) > 2 else ""
    flat = flatten(data, prefix)
    for k, v in sorted(flat.items()):
        key = k.upper().replace("-", "_").replace(".", "_")
        print(f"{key}={v}")
except Exception as e:
    print(f"# ERROR: {e}", file=sys.stderr)
    sys.exit(1)
PYEOF
}

# ── 从多个 YAML 文件加载到关联数组（映射到 schema 键名）────────────────────
yaml_parser::load_all() {
  local -n _out="$1"   # target assoc array
  local config_dir="$2"

  declare -A raw_flat=()

  # ── basic.yml → SERVER_, FRONTEND_, LOG_, OLLAMA_, UPLOAD_, 等 ──
  local basic="${config_dir}/basic.yml"
  if [[ -f "$basic" ]]; then
    while IFS='=' read -r k v; do
      [[ "$k" =~ ^# ]] && continue
      [[ -z "$k" ]] && continue
      raw_flat["$k"]="$v"
    done < <(yaml_parser::flatten "$basic" "BASIC")
    log::debug "已解析 basic.yml ($(wc -l < "$basic") 行)"
  fi

  # ── private.yml → DB_, REDIS_, JWT_, EMAIL_, ADMIN_ ──
  local private="${config_dir}/private.yml"
  if [[ -f "$private" ]]; then
    while IFS='=' read -r k v; do
      [[ "$k" =~ ^# ]] && continue
      [[ -z "$k" ]] && continue
      raw_flat["$k"]="$v"
    done < <(yaml_parser::flatten "$private" "PRIVATE")
    log::debug "已解析 private.yml"
  fi

  # ── postgres.yml ──
  local pg="${config_dir}/postgres.yml"
  if [[ -f "$pg" ]]; then
    while IFS='=' read -r k v; do
      [[ "$k" =~ ^# ]] && continue
      [[ -z "$k" ]] && continue
      raw_flat["PG_${k}"]="$v"
    done < <(yaml_parser::flatten "$pg" "")
    log::debug "已解析 postgres.yml"
  fi

  # ── redis.yml ──
  local redis="${config_dir}/redis.yml"
  if [[ -f "$redis" ]]; then
    while IFS='=' read -r k v; do
      [[ "$k" =~ ^# ]] && continue
      [[ -z "$k" ]] && continue
      raw_flat["REDIS_YML_${k}"]="$v"
    done < <(yaml_parser::flatten "$redis" "")
    log::debug "已解析 redis.yml"
  fi

  # ── config.yml (frontend) ──
  local cfg="${config_dir}/config.yml"
  if [[ -f "$cfg" ]]; then
    while IFS='=' read -r k v; do
      [[ "$k" =~ ^# ]] && continue
      [[ -z "$k" ]] && continue
      raw_flat["CFG_${k}"]="$v"
    done < <(yaml_parser::flatten "$cfg" "")
    log::debug "已解析 config.yml"
  fi

  # ── 映射规则：raw_flat key → schema key ──────────────────────────────────
  # server
  _out[SERVER_HOST]="${raw_flat[BASIC_SERVER_HOST]:-}"
  _out[SERVER_PORT]="${raw_flat[BASIC_SERVER_PORT]:-}"
  _out[SERVER_MODE]="${raw_flat[BASIC_SERVER_MODE]:-}"
  _out[SERVER_READ_TIMEOUT]="${raw_flat[BASIC_SERVER_READ_TIMEOUT]:-}"
  _out[SERVER_WRITE_TIMEOUT]="${raw_flat[BASIC_SERVER_WRITE_TIMEOUT]:-}"
  # frontend
  _out[FRONTEND_PROTOCOL]="${raw_flat[BASIC_FRONTEND_PROTOCOL]:-}"
  _out[FRONTEND_HOST]="${raw_flat[BASIC_FRONTEND_HOST]:-}"
  _out[FRONTEND_PORT]="${raw_flat[BASIC_FRONTEND_PORT]:-}"
  # database (prefer private.yml)
  _out[DB_HOST]="${raw_flat[PRIVATE_DATABASE_HOST]:-${raw_flat[PG_HOST]:-}}"
  _out[DB_PORT]="${raw_flat[PRIVATE_DATABASE_PORT]:-${raw_flat[PG_PORT]:-}}"
  _out[DB_USER]="${raw_flat[PRIVATE_DATABASE_USER]:-${raw_flat[PG_USER]:-}}"
  _out[DB_PASSWORD]="${raw_flat[PRIVATE_DATABASE_PASSWORD]:-${raw_flat[PG_PASSWORD]:-}}"
  _out[DB_NAME]="${raw_flat[PRIVATE_DATABASE_DBNAME]:-${raw_flat[PG_DBNAME]:-}}"
  _out[DB_SSLMODE]="${raw_flat[PRIVATE_DATABASE_SSLMODE]:-${raw_flat[PG_SSLMODE]:-}}"
  _out[DB_TIMEZONE]="${raw_flat[PRIVATE_DATABASE_TIMEZONE]:-}"
  _out[DB_MAX_IDLE_CONNS]="${raw_flat[PG_MAX_IDLE_CONNS]:-}"
  _out[DB_MAX_OPEN_CONNS]="${raw_flat[PG_MAX_OPEN_CONNS]:-}"
  # redis (prefer private.yml)
  _out[REDIS_HOST]="${raw_flat[PRIVATE_REDIS_HOST]:-${raw_flat[REDIS_YML_HOST]:-}}"
  _out[REDIS_PORT]="${raw_flat[PRIVATE_REDIS_PORT]:-${raw_flat[REDIS_YML_PORT]:-}}"
  _out[REDIS_USER]="${raw_flat[PRIVATE_REDIS_USER]:-${raw_flat[REDIS_YML_USER]:-}}"
  _out[REDIS_PASSWORD]="${raw_flat[PRIVATE_REDIS_PASSWORD]:-${raw_flat[REDIS_YML_PASSWORD]:-}}"
  _out[REDIS_DB]="${raw_flat[PRIVATE_REDIS_DB]:-${raw_flat[REDIS_YML_DB]:-}}"
  _out[REDIS_POOL_SIZE]="${raw_flat[REDIS_YML_POOL_SIZE]:-}"
  _out[REDIS_MIN_IDLE_CONNS]="${raw_flat[REDIS_YML_MIN_IDLE_CONNS]:-}"
  _out[REDIS_DIAL_TIMEOUT]="${raw_flat[REDIS_YML_DIAL_TIMEOUT]:-}"
  _out[REDIS_READ_TIMEOUT]="${raw_flat[REDIS_YML_READ_TIMEOUT]:-}"
  _out[REDIS_WRITE_TIMEOUT]="${raw_flat[REDIS_YML_WRITE_TIMEOUT]:-}"
  # jwt
  _out[JWT_SECRET]="${raw_flat[PRIVATE_JWT_SECRET]:-}"
  _out[JWT_EXPIRE]="${raw_flat[PRIVATE_JWT_EXPIRE]:-}"
  _out[JWT_REFRESH_EXPIRE]="${raw_flat[PRIVATE_JWT_REFRESH_EXPIRE]:-}"
  _out[JWT_ISSUER]="${raw_flat[PRIVATE_JWT_ISSUER]:-}"
  # email
  _out[EMAIL_HOST]="${raw_flat[PRIVATE_EMAIL_HOST]:-}"
  _out[EMAIL_PORT]="${raw_flat[PRIVATE_EMAIL_PORT]:-}"
  _out[EMAIL_USERNAME]="${raw_flat[PRIVATE_EMAIL_USERNAME]:-}"
  _out[EMAIL_PASSWORD]="${raw_flat[PRIVATE_EMAIL_PASSWORD]:-}"
  _out[EMAIL_FROM]="${raw_flat[PRIVATE_EMAIL_FROM]:-}"
  _out[EMAIL_FROM_NAME]="${raw_flat[PRIVATE_EMAIL_FROM_NAME]:-}"
  _out[EMAIL_SSL]="${raw_flat[PRIVATE_EMAIL_SSL]:-}"
  _out[EMAIL_TLS]="${raw_flat[PRIVATE_EMAIL_TLS]:-}"
  # log
  _out[LOG_LEVEL]="${raw_flat[BASIC_LOG_LEVEL]:-}"
  _out[LOG_FILENAME]="${raw_flat[BASIC_LOG_FILENAME]:-}"
  _out[LOG_MAX_SIZE]="${raw_flat[BASIC_LOG_MAX_SIZE]:-}"
  _out[LOG_MAX_BACKUPS]="${raw_flat[BASIC_LOG_MAX_BACKUPS]:-}"
  _out[LOG_MAX_AGE]="${raw_flat[BASIC_LOG_MAX_AGE]:-}"
  _out[LOG_COMPRESS]="${raw_flat[BASIC_LOG_COMPRESS]:-}"
  _out[LOG_CONSOLE]="${raw_flat[BASIC_LOG_CONSOLE]:-}"
  _out[LOG_JSON_FORMAT]="${raw_flat[BASIC_LOG_JSON_FORMAT]:-}"
  # ollama / llamacpp
  _out[OLLAMA_BASE_URL]="${raw_flat[BASIC_OLLAMA_BASE_URL]:-}"
  _out[OLLAMA_MODEL]="${raw_flat[BASIC_OLLAMA_MODEL]:-}"
  _out[OLLAMA_NUM_PREDICT]="${raw_flat[BASIC_OLLAMA_NUM_PREDICT]:-}"
  _out[OLLAMA_TEMPERATURE]="${raw_flat[BASIC_OLLAMA_TEMPERATURE]:-}"
  _out[LLAMACPP_BASE_URL]="${raw_flat[BASIC_LLAMACPP_BASE_URL]:-}"
  _out[LLAMACPP_MODEL]="${raw_flat[BASIC_LLAMACPP_MODEL]:-}"
  # upload
  _out[UPLOAD_DIR]="${raw_flat[BASIC_UPLOAD_UPLOAD_DIR]:-}"
  _out[UPLOAD_URL_PREFIX]="${raw_flat[BASIC_UPLOAD_URL_PREFIX]:-}"
  _out[UPLOAD_MAX_SIZE]="${raw_flat[BASIC_UPLOAD_MAX_SIZE]:-}"
  # admin
  _out[ADMIN_EMAIL]="${raw_flat[PRIVATE_ADMIN_EMAIL]:-}"
  _out[ADMIN_USERNAME]="${raw_flat[PRIVATE_ADMIN_USERNAME]:-}"
  _out[ADMIN_PASSWORD]="${raw_flat[PRIVATE_ADMIN_PASSWORD]:-}"
  _out[ADMIN_ROLE]="${raw_flat[PRIVATE_ADMIN_ROLE]:-}"
  # proxy / frontend config
  _out[PROXY_ENABLED]="${raw_flat[CFG_PROXY_ENABLED]:-}"
  _out[PROXY_BACKEND_URL]="${raw_flat[CFG_PROXY_BACKEND_URL]:-}"

  log::debug "yaml_parser::load_all 完成，共映射 ${#_out[@]} 个键"
}