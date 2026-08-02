# AGENTS Rules — TinyForum（全局 AI 协作总控）

> 本文件是 TinyForum 项目的**最高优先级 AI 协作宪章**，覆盖 **Go 后端 + Next.js 前端 + Vite 插件** 三大子项目。  
> **所有 AI 生成的代码、修复建议、重构方案，必须无条件遵守本文件约定。**  
> **核心理念**：**拒绝“症状缓解”（打补丁），强制“根因治疗”（修本质）。**

---

## 0. AI 修复协议（最高优先级：五步根因法）

> 本协议优先级高于所有编码规范。AI 在修复任何 Bug 时必须强制执行以下五步，**严禁跳过步骤直接输出 diff**。

### 第一步：5-Why 根因分析

AI 必须回答：

1. **现象**：报错/用户反馈的具体表现？
2. **表层直接原因**：哪一行代码抛出了异常或返回了错误数据？
3. **深层根因（追问 5 次）**：为什么会出现这个数据/状态？——是并发写入未加锁？数据库字段溢出？前端传参与后端结构体不匹配？状态机流转遗漏了事件？
4. **影响范围**：孤立模块还是影响全局？

### 第二步：复现测试先行（红灯准入）

- AI **必须先编写单元测试**捕获 Bug，再修改生产代码。
- 后端：生成 `_test.go`，运行 `go test` 必须**失败（Red）**。
- 前端：用 Vitest 模拟异常 API 响应或状态。
- **门禁**：未提供可复现的失败测试，生产代码改动不予通过。

### 第三步：强制分层定位（禁止跨层绕过）

AI 必须明确修复发生的**层级**，禁止用上层 `if` 掩盖下层缺陷：

| 层级                   | 允许的“根本修复”                                    | 严厉禁止的“补丁行为”                                    |
| :--------------------- | :-------------------------------------------------- | :------------------------------------------------------ |
| **Repository（后端）** | 修正 SQL 条件、添加索引、处理 NULL/零值、使用乐观锁 | 在 Service 层加 `if` 过滤空数据                         |
| **Service（后端）**    | 修正状态机流转、引入幂等/重试、添加事务边界         | 在 Handler 用 `recover()` 捕获 panic 返回 200           |
| **Handler（后端）**    | 修正参数校验规则 `binding`、修正 VO 字段映射        | 在响应中硬编码 `"data": null` 或 `"code": 0`            |
| **View/Hook（前端）**  | 修正 Zustand 状态更新、修正 TanStack Query 依赖键   | 加 `setTimeout` 延迟等待、用 `try...catch` 静默吞掉错误 |

### 第四步：数据迁移解耦（若涉及 Schema）

- **严禁**在 Go 代码中做类型转换来“适配”错误的数据库字段。
- **必须**提供 `backend/migrations/` 下的 `.up.sql` 和 `.down.sql` 迁移文件。
- **兼容性铁律**：删除列或重命名必须**分两步提交**（先增后弃），确保新旧代码同时兼容。

### 第五步：提交信息自审

- **坏的描述（补丁）**：`fix [auth]: add nil check for token`
- **好的描述（根因）**：`fix [auth]: ensure JWT parser returns ErrInvalidToken on malformed signature instead of panic`

> **子文档落地**：后端 AI 协议细则见 `backend/AGENTS.md` 第 0 章；前端 AI 协议细则见 `frontend/AGENTS.md` 第 0 章。

---

## 1. 项目概览

Tiny Forum 是一个开源论坛项目：**Go（Gin + GORM）后端** × **Next.js（App Router）前端** × **Vite 插件体系** × PostgreSQL 持久化。

| 子项目    | 位置          | 技术栈                                                                                                        |
| --------- | ------------- | ------------------------------------------------------------------------------------------------------------- |
| 后端      | `backend/`    | Go 1.26, Gin, GORM, Wire（手工注入）, JWT, Zap                                                                |
| 前端      | `frontend/`   | Next.js 16（App Router）+ React 19 + TypeScript, Tailwind/DaisyUI, TanStack Query, Zustand, Tiptap, next-intl |
| 插件      | `plugin/`     | Vite + React（独立构建，产物 `dist/main.js` + `manifest.json`）                                               |
| Make 脚本 | `dev-script/` | 用于启动、构建、测试、打包项目                                                                                |
| 文档      | `docs/`       | 项目 docsify 文档，涵盖用户手册、开发手册、架构说明等                                                         |

---

## 2. 跨项目代码规范（总纲）

### 2.1 分层架构

- **后端依赖方向严格单向**：`handler → service → repository`，禁止反向依赖。
- **前端依赖方向**：`View → State/Query → Client`；数据流向 `User Action → View → Query/State → Client → 后端`。
- **层内职责单一**：后端 handler 只绑定参数与写响应，service 只放业务规则，repository 只做数据访问；前端 View 只渲染 UI，不直接调 API。

### 2.2 注释与命名

- **注释**：代码内统一用中文；公开的 Go 类型、接口方法、handler 必须带注释。不新增无意义注释；重构时保留既有中文注释。
- **后端命名**：包名小写单数（`handler/announcement`）；接口 `XxxService` / `XxxRepository` + 实现 `xxxService` / `xxxRepository` + `NewXxxService(...)` 构造器。
- **前端命名**：文件 `PascalCase.tsx`、`camelCase.ts`、hook 以 `use` 开头；领域类型以 `XxxDO` / `XxxVO` / `XxxListParams` 命名。

### 2.3 质量门禁

- **后端**：改动后保证 `go build ./...` 通过；新增/修改装配逻辑后运行 `make wire`；运行 `golangci-lint`（含 govet、gofmt）。
- **前端**：改动后保证 `type-check` 与 `lint` 通过，`strict` 模式下不得用 `any`（仅存量兼容）；提交前执行 Prettier（单引号、无分号）。

### 2.4 SOLID 原则

- **单一职责**：每个类、函数、组件只做一件事。
- **开闭原则**：对扩展开放，对修改封闭。
- **里氏替换**：子类对象必须能替换父类对象。
- **接口隔离**：接口应小而专，避免“胖接口”。
- **依赖倒置**：高层模块不应依赖低层模块，两者都应依赖抽象。

---

## 3. Git 协作底线（跨项目统一）

### 3.1 分支策略

- `main` 为生产保护分支，**任何人（包括 AI）严禁直接 push**。
- 所有变更必须通过 `dev` → PR/MR 合并。
- **前缀强制**：修复任务 `fix-[issue_id]`；新功能 `feat-[name]`；重构 `refactor-[name]`；文档 `docs-[name]`；插件 `update-[name]`。

### 3.2 提交信息格式（全局统一）

`<type> [<module>]: <description>`

- **type**：`feat` / `fix` / `docs` / `update` / `refactor` / `perf` / `test` / `chore`
- **[module]**：影响模块（如 `post`、`auth`），多模块逗号分隔
- **description**：英文，现在时态，不超过 50 字符

示例：

```

fix [post]: fix post like
feat [auth]: add remember me
refactor [auth, api]: extract common validator

```

### 3.3 破坏性变更（BREAKING CHANGE）

必须在 footer 包含 `BREAKING CHANGE:` 说明，或在 `type` 后紧跟 `!`：

```bash
git commit -m "feat! [api]: change response structure" -m "详见迁移指南。"
```

### 3.4 标签（Tag）规范

正式发布使用附注标签：

```bash
git tag -a v1.0.0 -m "v1.0.0 发布（支持多因素认证）" -m "- 新增：OAuth2 登录\n- 修复：高并发下的 session 冲突"
git push origin v1.0.0
```

---

## 4. 跨项目协作约束

- **文件修改限制**：单次对话中，AI 不得一次性修改超过 **10 个文件**，除非用户明确要求“全量重构”。
- **依赖引入**：AI 不得主动引入新的第三方依赖（Go module / npm package），必须先注释说明理由并等待用户批准。
- **代码删除**：删除任何逻辑前，AI 必须输出“备份提示”并标记为临时注释（`// TODO: confirm before removal`），确保不影响现有调用链。
- **环境启动**：根目录提供 `docker-compose.yml` 启动 Postgres；启动前复制 `backend/config/private.yml.example` 并填入真实值。根目录运行 `make dev`（后端）和 `pnpm dev`（前端）分别启动。

---

## 5. 文档编写规范

- **文档存放**：所有项目文档统一放在 `docs/` 目录，使用 docsify 构建。
- **分类结构**：
  - **用户手册**：面向最终用户，介绍论坛功能、操作流程、常见问题等，位于 `docs/zh-CN/user/` 和 `docs/en-US/user/`。
  - **开发手册**：面向开发者，包含项目架构设计、模块说明、API 接口文档、插件开发指南、环境搭建等，位于 `docs/zh-CN/dev/` 和 `docs/en-US/dev/`。
  - **架构设计文档**：描述系统整体架构、技术选型、数据流、部署方案等，位于 `docs/zh-CN/architecture/` 和 `docs/en-US/architecture/`。
- **语言与格式**：文档使用 Markdown 编写，中文优先（同时提供英文版本）。代码示例、配置片段必须与实际代码保持一致。
- **更新时机**：任何代码变更（新增功能、修复 Bug、重构）若影响用户行为或开发流程，必须同步更新相关文档。AI 在提交代码前必须检查是否需要更新文档，并在 Commit Message 中注明 `docs` 类型或关联文档更新。
- **文档审查**：文档变更也需要经过 PR 审查，确保准确性和可读性。
- **索引**：根目录 `docs/README.md` 提供文档导航，所有分类文档均需在此索引中链接。

> 子项目内置注释（如 Go 的 godoc）不属于 `docs/` 文档体系，但 API 的 Swagger 注解应与开发手册中的接口说明保持同步。

---

## 6. 子项目专属规范（引用入口）

> 所有技术实现细节（具体命令、代码示例、测试框架配置）均在子项目 AGENTS.md 中维护，全局文档仅定义跨项目原则。

- **后端**：`backend/AGENTS.md`（含响应/错误处理、路由/中间件、WebSocket、Migration、Wire 装配细则）
- **前端**：`frontend/AGENTS.md`（含 Client 层、Query/Zustand 状态边界、国际化、文件上传、实时通信细则）
- **插件**：`plugin/`（见 `manifest.json` 与 `docs/zh-CN/dev/plugin.md`）

---

## 7. 质量红线（全局零容忍）

AI 生成的代码若触及以下红线，**视为无效交付，必须重写**：

1. **类型逃逸**：后端 `interface{}` 无类型断言，前端 `any` 逃逸 TypeScript（存量兼容除外）。
2. **硬编码**：前端文案未通过 `next-intl`；后端配置（JWT 密钥、DB DSN）硬编码在代码中。
3. **裸日志**：捕获错误后仅 `fmt.Println` / `console.log`，未接入 Zap / 统一错误上报。
4. **魔数/魔字符串**：业务状态码未定义为常量或枚举。
5. **跳过 Migration**：修改 GORM 模型 `do` 但未生成 SQL 迁移文件，或注释了 `AutoMigrate`。

---

## 8. Make 脚本与 Shell 规范

项目使用 `dev-script/` 目录集中管理所有构建、运行、测试、部署相关的 Makefile 和 Shell 脚本。

### 8.1 目录结构

```
dev-script/
├── backend/               # 后端相关配置（如 config/ 下的环境模板）
├── scripts/               # 可执行脚本（按功能分层）
│   ├── dev/               # 开发环境脚本（启动、检查、环境检测等）
│   ├── env/               # 环境变量解析与验证
│   ├── db/                # 数据库操作（mock 数据、清理等）
│   └── nginx/             # Nginx 配置辅助脚本
├── help.mk                # 帮助信息（make help）
├── Makefile.bench         # 基准测试
├── Makefile.cfg           # 配置生成/检查
├── Makefile.check         # 代码检查（lint、format）
├── Makefile.clean         # 清理产物
├── Makefile.code          # 代码生成（wire、swagger）
├── Makefile.common        # 公共变量与通用规则
├── Makefile.dev           # 开发环境启动
├── Makefile.docker        # Docker 镜像构建
├── Makefile.env           # 环境变量加载
├── Makefile.log           # 日志操作
├── Makefile.main          # 主入口（包含 .DEFAULT_GOAL 或 core targets）
├── Makefile.nginx         # Nginx 相关
├── Makefile.podman        # Podman 替代 Docker 的命令
├── Makefile.prod          # 生产环境构建/部署
└── scripts/               # （已列）
```

### 8.2 命名与职责

- **Makefile.\***：每个文件聚焦单一职责（如 `Makefile.dev` 只含开发环境目标，`Makefile.prod` 只含生产目标）。根目录的 `Makefile`（通常位于项目根）通过 `include dev-script/*.mk` 引入这些片段。
- **脚本文件**：位于 `scripts/` 下的子目录，使用 `.sh` 后缀，功能专一，避免大杂烩。
  - `dev/`：开发辅助（`backend.sh`、`frontend.sh`、`postgres.sh`、`redis.sh` 等），由 Makefile.dev 调用。
  - `env/`：环境变量加载、验证、解析（`core.sh`、`validator.sh`、`yaml_parser.sh` 等）。
  - `db/`：数据库数据准备（`mock_data.sh`、`clean_data.sh`）。
  - `nginx/`：Nginx 配置或代理设置。

### 8.3 编码规范

#### Shell 脚本（Bash）

- **解释器**：统一使用 `#!/usr/bin/env bash`，并启用严格模式：
  ```bash
  set -euo pipefail
  IFS=$'\n\t'
  ```
- **函数命名**：使用 `snake_case`，并添加注释说明功能及参数。
- **变量**：全部大写常量，局部变量使用 `local` 修饰，引用时加双引号（`"$var"`）防止分词。
- **错误处理**：捕获错误时输出明确信息到 stderr，非零退出码。
- **可移植性**：尽量使用 POSIX 兼容语法，避免 GNU 扩展（如 `bash` 特有的数组用法仅在确认可用时使用）。
- **依赖检查**：脚本开头应检查所需命令（如 `docker`、`go`、`pnpm`）是否存在，缺失时给出友好提示并退出。
- **日志**：关键操作输出带时间戳的日志，便于调试。

#### Makefile

- **目标命名**：使用 `snake_case` 或 `kebab-case`（统一即可），推荐 `kebab-case`（如 `run-dev`、`clean-all`）。
- **伪目标**：所有非文件目标须声明为 `.PHONY`。
- **变量**：统一在 `Makefile.common` 中定义全局变量（如 `GOBIN`、`NODE_VERSION`），其他文件引用。
- **依赖关系**：明确声明目标间的依赖，避免冗余执行。
- **安全性**：使用 `$(MAKE)` 递归调用子 Makefile，避免硬编码 `make` 路径。
- **帮助信息**：在 `help.mk` 中定义 `help` 目标，自动提取各 Makefile 中带 `##` 注释的目标说明，提供统一帮助入口。

### 8.4 使用原则

- **入口**：用户在根目录执行 `make <target>`，实际解析由根目录 `Makefile` 通过 `include` 组合完成。
- **环境隔离**：所有脚本应能独立运行（即不依赖当前工作目录），通常以 `dev-script/` 为基准路径。
- **跨平台**：优先考虑 Linux/macOS 兼容性；若需 Windows 支持，使用 Git Bash 或 WSL。

### 8.5 AI 协作约束

- 修改任何 `.mk` 或 `.sh` 文件时，AI 必须遵循第 0 章的“五步根因法”，不得以“先能跑起来”为由引入临时补丁。
- 新增脚本或目标须先在 `help.mk` 中添加说明，并更新 `docs/` 中对应的开发手册。
- 任何脚本变量或路径变更，必须同步更新 `Makefile.common` 和 `scripts/env/` 下的环境加载逻辑。

---

**生效声明**：本文件覆盖项目内所有 `*.md` 文档。当其他文档与本文冲突时，以本文的“AI 修复协议”和“质量红线”为准。
