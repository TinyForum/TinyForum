# 开发指南

本指南为 TinyForum 开发者提供从环境搭建到编码规范的完整参考。

---

## 快速导航

| 分类 | 文档 | 说明 |
|------|------|------|
| **架构** | [架构设计](/zh-CN/dev/architecture) | 前后端分层架构详解 |
| **环境** | [快速开始](/zh-CN/dev/first) | 本地开发环境搭建 |
| **API** | [RESTful API 设计规范](/zh-CN/dev/restful_api) | URL 设计、HTTP 方法、状态码 |
| **服务** | [服务层级调用方法](/zh-CN/dev/service_call) | Handler → Service → Repository 调用链规范 |
| **响应** | [统一后端响应](/zh-CN/dev/uni_response) | response 包使用说明 |
| **数据** | [数据模型设计规范](/zh-CN/dev/object_model) | Request/VO/BO/DO/DTO 定义与转换 |
| **转换** | [数据转换规范](/zh-CN/dev/data_transform) | 各层数据边界与 Converter 模式 |
| **数据库** | [数据库设计](/zh-CN/dev/database/intro) | PostgreSQL 配置、GORM 使用、迁移规范 |
| **机器人** | [机器人系统](/zh-CN/dev/robot) | Lua 脚本 / 零代码机器人开发 |
| **插件** | [插件系统](/zh-CN/dev/plugin) | 前端插件加载与开发指南 |
| **AI** | [AI 辅助功能](/zh-CN/dev/ai) | 敏感词 LLM 复核 |
| **前端** | [TanStack Query 指南](/zh-CN/dev/hooks) | 前端数据获取最佳实践 |
| **测试** | [Swagger 测试](/zh-CN/dev/swagger) | API 接口测试方法 |
| **规范** | [命名规范](/zh-CN/dev/named) | Go 代码命名最佳实践 |
| **运维** | [Redis 操作](/zh-CN/dev/redis/guide) | Redis 限流键管理 |

---

## 项目技术栈

| 层级 | 技术 | 用途 |
|------|------|------|
| 后端语言 | Go 1.24+ | 服务端核心逻辑 |
| Web 框架 | Gin | HTTP 路由与中间件 |
| ORM | GORM | 数据库访问 |
| 数据库 | PostgreSQL | 关系型数据持久化 |
| 缓存 | Redis | 会话、限流、排行榜 |
| 鉴权 | JWT + Casbin | 认证与 RBAC 权限 |
| 日志 | Zap | 结构化日志 |
| API 文档 | Swagger | 自动生成 API 文档 |
| 依赖注入 | Wire（手工装配） | 模块解耦 |
| 前端框架 | Next.js 16 (App Router) | SSR/SSG 前端 |
| 前端语言 | TypeScript + React 19 | 类型安全组件开发 |
| 状态管理 | Zustand + TanStack Query | 客户端状态与服务端缓存 |

---

## 开发流程

### 1. 环境准备

```bash
# 确保 Go >= 1.24, Node.js >= 20, PostgreSQL >= 16, Redis >= 7
make init-dev
```

### 2. 启动开发服务

```bash
# 后端 (http://localhost:8080)
cd backend && make run

# 前端 (http://localhost:3000)  
cd frontend && pnpm dev
```

### 3. 代码变更工作流

```
创建分支 → 编写代码 → 自测 → 提交 PR → Code Review → 合并
```

分支命名遵循 `feat-[name]` / `fix-[issue_id]` / `refactor-[name]` 规范。

### 4. 提交前检查清单

- [ ] `go build ./...` 编译通过
- [ ] `go test ./...` 测试通过
- [ ] `golangci-lint run` 零告警
- [ ] 修改 `do` 模型时，已添加对应的 Migration SQL
- [ ] API 变更已更新 Swagger 注解（`make docs-api`）

---

## 项目结构速查

```
backend/
├── cmd/server/          # 程序入口
├── internal/
│   ├── handler/         # HTTP 层：参数绑定、路由注册
│   ├── service/         # 业务层：逻辑编排、事务管理
│   ├── repository/      # 数据层：GORM 数据库操作
│   ├── model/
│   │   ├── do/          # 数据库实体 (GORM)
│   │   ├── request/     # API 入参
│   │   ├── vo/          # API 出参
│   │   ├── dto/         # 层间传输对象
│   │   └── bo/          # 业务对象
│   ├── middleware/      # 中间件 (Auth, Casbin, RateLimit, ContentCheck)
│   ├── infra/           # 基础设施 (配置, 邮件, 敏感词)
│   └── wire/            # 依赖注入装配
├── config/              # YAML 配置文件
├── pkg/                 # 可复用工具包
└── migrations/          # SQL 迁移脚本
```
