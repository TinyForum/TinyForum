# 快速开始

## 安装方式

TinyForum 支持多种部署方式：

| 方式 | 适用场景 | 难度 |
|------|----------|------|
| Docker Compose | 生产/测试环境 | 低 |
| Podman | macOS/Linux 无 Docker 环境 | 中 |
| 本地开发 | 开发调试 | 中 |
| 脚本一键部署 | Ubuntu 生产环境 | 低 |

---

## 脚本工具

项目使用 `dev-script/` 目录管理构建脚本，通过根目录 `Makefile` 调用。

```bash
# 查看所有可用命令
make help

# 开发环境初始化
make init-dev

# 启动后端开发服务
make run-backend

# 启动前端开发服务
make run-frontend

# 代码检查
make check

# 代码生成（Wire + Swagger）
make code-gen
```

## 环境要求

- Go >= 1.24
- Node.js >= 20 + pnpm
- PostgreSQL >= 16
- Redis >= 7
- Docker / Podman（可选）

## 联系方式

如有问题，请通过 GitHub Issues 反馈：
https://github.com/caoyang2002/TinyForum/issues
