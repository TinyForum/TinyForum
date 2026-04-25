# tiny-forum 容器化部署指南

## 目录结构

```
project/
├── backend/
│   ├── Dockerfile
│   └── configs/
│       ├── basic.yaml
│       ├── private.yaml        ← 不要提交到 Git
│       └── risk_control.yaml
├── frontend/
│   └── Dockerfile
├── deploy/
│   └── nginx/
│       ├── nginx.conf
│       └── conf.d/
│           └── default.conf
├── docker-compose.yml
├── docker-compose.override.yml  ← 仅开发环境使用
├── .env                         ← 从 .env.example 复制，不要提交到 Git
└── .env.example
```

## 快速开始

### 1. 准备环境变量

```bash
cp .env.example .env
# 编辑 .env，至少修改 DB_PASSWORD 和 JWT_SECRET
```

### 2. 创建 Nginx 目录

```bash
mkdir -p deploy/nginx/conf.d
cp deploy/nginx/nginx.conf deploy/nginx/nginx.conf
cp deploy/nginx/conf.d/default.conf deploy/nginx/conf.d/default.conf
```

### 3. 启动服务

**Docker：**
```bash
make dc-up
# 或直接
docker compose up -d
```

**Podman：**
```bash
make dc-up
# 或直接
podman-compose up -d
# 或（Podman v4+ 内置 compose）
podman compose up -d
```

## 常用命令

| 命令 | 说明 |
|------|------|
| `make dc-up` | 启动所有服务（后台） |
| `make dc-down` | 停止所有服务 |
| `make dc-status` | 查看服务状态 |
| `make dc-logs` | 跟踪全部日志 |
| `make dc-logs-backend` | 只看后端日志 |
| `make dc-build` | 重新构建镜像 |
| `make dc-backup-db` | 备份数据库 |
| `make dc-dev` | 开发模式（含 override） |
| `make dc-clean` | 销毁全部资源 |

## Podman 兼容说明

本配置完全兼容 Podman，需注意：

- **rootless Podman**：默认不能监听 80/443 端口，可改用 `HTTP_PORT=8080` 或启用 `net.ipv4.ip_unprivileged_port_start`。
- **podman-compose**：`pip install podman-compose` 安装，用法与 docker-compose 一致。
- **Podman v4+ 内置 compose**：`podman compose` 命令，无需额外安装。
- **Volume 权限**：rootless 模式下容器内 UID 与主机不同，如遇权限问题，在 volume 挂载后执行 `podman unshare chown -R 1001:1001 ./data`。

## 服务端口

| 服务 | 容器内端口 | 宿主机默认端口 |
|------|-----------|--------------|
| Nginx | 80 | 80 (可通过 HTTP_PORT 改) |
| Backend | 8080 | 8080 (可通过 BACKEND_PORT 改) |
| PostgreSQL | 5432 | 5432 (可通过 DB_PORT 改) |
| Redis | 6379 | 6379 (可通过 REDIS_PORT 改) |

## 注意事项

- `private.yaml` 和 `.env` 含敏感信息，务必加入 `.gitignore`。
- 生产环境请修改 `JWT_SECRET` 为 32 字符以上的随机字符串。
- PostgreSQL 数据存储在 Docker/Podman volume `tiny-forum_postgres_data` 中，`dc-down` 不会删除数据；`dc-down-v` 会永久删除。