# 快速开始

## 一、检查并配置开发环境

```bash
./start-dev.sh
```

> 如果有需要按照的程序会提示你。
>
> 前端默认使用 pnpm 启动，可以自行修改

启动成功显示如下

```bash
Done in 220ms using pnpm v10.32.1
  Frontend dependencies installed.

==================================
✅ Setup complete!

To start the backend:
  cd backend && go run ./cmd/server/main.go

To start the frontend:
  cd frontend && pnpm run dev

Database connection info:
  Host: localhost:5432
  User: caoyang
  Password: (empty - using trust authentication)
  Database: tiny_forum
==================================

Testing database connection...
✅ Database connection successful

All services are ready! 🎉
```
根据 Database 显示的信息修改 `backend/config/config.yaml` 文件

> 在 mac 上，一般是 mac 的用户名为数据库用户名，psql 允许无密码连接
 
 您也可以创建指定的用户

```sql

-- 创建 tiny_admin 用户
CREATE USER tiny_admin WITH PASSWORD 'tiny_admin_password';

-- 创建数据库（如果不存在）
CREATE DATABASE tiny_forum OWNER tiny_admin;

-- 授予所有权限
GRANT ALL PRIVILEGES ON DATABASE tiny_forum TO tiny_admin;

-- 允许用户创建数据库（可选）
ALTER USER tiny_admin CREATEDB;
```

# 使用当前系统用户执行 SQL

```bash
psql -d postgres -f /tmp/create_user.sql
```

## 你可以使用 Makefile 启动前后端

```bash
# 在这个项目的根目录
make frontend
make backend
```

- 前端 UI: http://localhost:3000/
- 后端 API: http://localhost:8080/api/v1/
- swagger: http://localhost:8080/api/v1/swagger/index.html

可以通过 docsify 生成项目文档，但是需要先安装 docsify

```bash
make docs # 生成文档()
```

