# 开发手册

欢迎阅读 TinyForum 开发手册。本文档面向开发者，涵盖项目架构设计、编码规范、API 文档、插件开发和部署运维。

---

## 目录导航

### 基础环境
- [快速开始](/zh-CN/dev/first) — 本地开发环境搭建
- [架构设计](/zh-CN/dev/architecture) — 前后端分层架构详解
- [数据库设计](/zh-CN/dev/database/intro) — PostgreSQL 配置与 GORM 使用规范

### API 开发
- [RESTful API 设计规范](/zh-CN/dev/restful_api) — URL 设计、HTTP 方法、状态码
- [统一后端响应](/zh-CN/dev/uni_response) — response 包使用说明
- [服务层级调用方法](/zh-CN/dev/service_call) — Handler → Service → Repository 调用链
- [Swagger 测试](/zh-CN/dev/swagger) — API 接口测试方法

### 数据层
- [数据模型设计规范](/zh-CN/dev/object_model) — Request/VO/BO/DO/DTO 定义
- [数据转换规范](/zh-CN/dev/data_transform) — 各层数据边界与 Converter 模式
- [命名规范](/zh-CN/dev/named) — Go 代码命名最佳实践

### 高级功能
- [机器人系统](/zh-CN/dev/robot) — Lua 脚本 + 零代码 Flow 引擎
- [插件系统](/zh-CN/dev/plugin) — 前端插件加载与开发
- [AI 辅助功能](/zh-CN/dev/ai) — 敏感词 LLM 复核

### 前端开发
- [TanStack Query 指南](/zh-CN/dev/hooks) — 前端数据获取最佳实践

### 运维
- [Redis 操作指南](/zh-CN/dev/redis/guide) — Redis 限流键管理
- [PostgreSQL 小贴士](/zh-CN/dev/tips) — 常用数据库操作

---

## 相关资源

- [项目 GitHub](https://github.com/caoyang2002/TinyForum)
- [后端 AGENTS.md](https://github.com/caoyang2002/TinyForum/blob/main/backend/AGENTS.md)
- [前端 AGENTS.md](https://github.com/caoyang2002/TinyForum/blob/main/frontend/AGENTS.md)
- [全局 AGENTS.md](https://github.com/caoyang2002/TinyForum/blob/main/AGENTS.md)
