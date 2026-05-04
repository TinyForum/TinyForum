# ========== 防止重复包含 ==========
ifndef _HELP_MK_INCLUDED
_HELP_MK_INCLUDED := 1
_HELP_TITLE_COLOR  ?= \033[1;32m
_HELP_CMD_COLOR    ?= \033[1;36m
_HELP_DESC_COLOR   ?= \033[0m
_HELP_RESET        ?= \033[0m

.PHONY: _print_help
_print_help:
	@awk -v title="$(_HELP_TITLE_COLOR)" -v cmd="$(_HELP_CMD_COLOR)" -v desc="$(_HELP_DESC_COLOR)" -v rst="$(_HELP_RESET)" -v file="$(SELF_MAKEFILE_PATH)" \
	'BEGIN { printf("\n%sAvailable targets (from %s):%s\n\n", title, file, rst); } \
	/^[[:space:]]*# @/ { comment = $$0; sub(/^[[:space:]]*# @[[:space:]]*/, "", comment); last = comment; next; } \
	/^[[:space:]]*$$/ { next; } \
	/^[[:space:]]*[a-zA-Z_][a-zA-Z0-9_.%-]*[[:space:]]*:/ { \
		c = $$1; sub(/:/, "", c); sub(/^[[:space:]]+/, "", c); \
		if (last != "") printf("  %s%-20s%s %s%s%s\n", cmd, c, rst, desc, last, rst); \
		last = ""; next; \
	} \
	{ last = ""; }' "$(SELF_MAKEFILE_PATH)"

# help 目标，需要调用者传递 SELF_MAKEFILE_PATH
.PHONY: help
self-help:
	@echo "call by $(SELF_MAKEFILE_PATH)"
	@$(MAKE) -s _print_help SELF_MAKEFILE_PATH='$(SELF_MAKEFILE_PATH)'
endif   # 结束 guard