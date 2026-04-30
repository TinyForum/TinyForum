include dev-script/Makefile.common
include dev-script/Makefile.dev
include dev-script/Makefile.podman
include dev-script/Makefile.docker
include dev-script/Makefile.main
include dev-script/Makefile.bench
include dev-script/Makefile.code
include dev-script/Makefile.env
help:
	@echo "=============="
	@echo "TinyForum help"
	@echo "=============="
	@echo "  make code-help     - 代码检查帮助"
	@echo "  make docker-help   - docker 帮助"
	@echo "  make podman-help   - podman 帮助"
	@echo "  make dev-help      - 开发帮助 "
	@echo "  make bench-help    - 性能测试帮助 "
	@echo "  make dev-help      - 环境变量帮助 "