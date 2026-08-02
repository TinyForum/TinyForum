# AGENTS.md — TinyForum

Tiny Forum 是一个开源论坛项目：**Go（Gin + GORM）后端** × **Next.js（App Router）前端** × **Vite 插件体系** × PostgreSQL 持久化。

本文件约定**跨子项目的代码编写规范与最佳实践**。运行命令、数据库相关请勿在这里维护，分别见各子项目文档：

- 后端：`backend/AGENTS.md`
- 前端：`frontend/AGENTS.md`
- 插件：`plugin/`（见 `manifest.json` 与 `docs/zh-CN/dev/plugin.md`）

## 仓库组成

| 子项目    | 位置          | 技术栈                                                                                                        |
| --------- | ------------- | ------------------------------------------------------------------------------------------------------------- |
| 后端      | `backend/`    | Go 1.26, Gin, GORM, Wire（手工注入）, JWT, Zap                                                                |
| 前端      | `frontend/`   | Next.js 16（App Router）+ React 19 + TypeScript, Tailwind/DaisyUI, TanStack Query, Zustand, Tiptap, next-intl |
| 插件      | `plugin/`     | Vite + React（独立构建，产物 `dist/main.js` + `manifest.json`）                                               |
| Make 脚本 | `dev-script/` | 用于启动、构建、测试、打包项目                                                                                |
| 文档      | `docs/`       | 项目 docsify 文档，涵盖用户手册、开发手册、架构说明等                                                         |

## 全局代码规范

### 分层架构

- 后端依赖方向严格单向：`handler → service → repository`，禁止反向依赖。
- 前端依赖方向：`View → State/Query → Client`；数据流向 `User Action → View → Query/State → Client → 后端`。
- 层内职责单一：后端 handler 只绑定参数与写响应，service 只放业务规则，repository 只做数据访问；前端 View 只渲染 UI，不直接调 API。

### 注释

- 代码内注释统一用中文；公开的 Go 类型、接口方法、handler 必须带注释。
- 不新增无意义注释；重构时保留既有中文注释，不随意改写他人注释风格。

### 命名

- 后端：包名小写单数（`handler/announcement`、`service/announcement`、`repository/announcement`）；接口 `XxxService` / `XxxRepository` + 实现 `xxxService` / `xxxRepository` + `NewXxxService(...)` 构造器。
- 前端：文件 `PascalCase.tsx`、`camelCase.ts`、hook 以 `use` 开头；领域类型以 `XxxDO` / `XxxVO` / `XxxListParams` 命名。

### 质量门禁

- 后端：改动后保证 `go build ./...` 通过；新增/修改装配逻辑后运行 `make wire` 重新生成依赖注入代码；运行 `golangci-lint`（含 govet、gofmt）。
- 前端：改动后保证 `type-check` 与 `lint` 通过，`strict` 模式下不得用 `any` 逃逸类型检查（仅存量兼容）。
- 前端提交前执行格式化（Prettier：单引号、无分号）。

## 后端最佳实践（详见 backend/AGENTS.md）

- **响应统一**：一律用 `pkg/response`（`Success` / `SuccessPage` / `Created` / `HandleError` 等），格式为 `{"code": 0, "message": "success", "data": ...}`。
- **错误处理**：业务错误复用 `pkg/errors` 预定义 `AppError`（如 `apperrors.ErrPostNotFound`），附加信息用链式方法 `ErrXxx.WithDetail(...).WithCause(err)`，**不要直接修改全局错误实例**；handler 层统一 `response.HandleError(c, err)` 收口。
- **依赖注入**：新增 handler/service/repository 后，在 `internal/wire/` 手工装配（`NewServices`/`NewHandlers`/`NewRepositories`）。
- **接口文档**：所有 handler 方法加 Swagger 注解（`@Summary` / `@Tags` / `@Router` 等），与前端响应类型保持一致。
- **模型分层**：`do`（GORM 实体）/ `request`（入参）/ `vo`（出参）/ `dto`（层间传递）/ `bo`（业务对象）各归其位，不混放。
- **安全**：JWT 密钥等敏感配置放 `config/private.yml`，禁止提交真实密钥；私有路由统一 `Auth` + `CasbinAuth`。

## 前端最佳实践（详见 frontend/AGENTS.md）

- **数据请求**：统一经 `shared/api/` Client 层（axios 实例，baseURL `/api/v1`，响应 `{ code, message, data }`）；API 方法返回 `apiClient.get<ApiResponse<T>>(...)`，hook 解出 `response.data.data` 消费。
- **状态分层**：跨组件共享状态才进 Zustand `src/store/`（如 auth）；局部状态留在组件内；服务端数据用 TanStack Query 的 Query/Mutation 管理。
- **Query key**：用集中定义的 `xxxKeys` 工厂（如 `postKeys.list(params)`）；mutation 成功后 `invalidateQueries` 刷新列表/详情。
- **新增模块流程**：`shared/api/types/<module>.model.ts` 定义类型 → `shared/api/modules/<module>.ts` 定义 api 对象 → `features/<module>/hooks` 封装 Query/Mutation。
- **国际化**：文案经 next-intl 消息文件（`messages/zh-CN`、`messages/en-US`），不硬编码到组件。
- **组件样式**：遵循 Tailwind + DaisyUI 现有模式，不用内联样式；import 别名统一 `@/`。

## 插件开发

- 插件独立构建（Vite + React），产物为 `dist/main.js`，由 `manifest.json` 声明 `slug`、`slots`、`configSchema`、`permissions`。
- 插件通过 PluginSlot 槽位（如 `after-header`、`before-footer`）注入 UI，`slotProps` 传入页面数据（如 `postId`）。
- 插件改动不依赖前后端重新构建；扩展后端能力遵循 `docs/zh-CN/dev/plugin.md` 约定。

## 代码更新规范

- 所有的 fix 任务必须在 "fix-[name]" 中进行，避免污染 main 分支。
- 所有的功能更新必须在 "feat-[name]" 中进行，避免污染 main 分支。
- 所有的文档更新必须在 "docs-[name]" 中进行，避免污染 main 分支。
- 所有的插件更新必须在 "update-[name]" 中进行，避免污染 main 分支。
- 所有的重构任务必须在 "refactor-[name]" 中进行，避免污染 main 分支。
- dev 分支用于开发，main 分支用于发布，不允许直接向 main 分支提交代码。
- 禁止向 main 分支提交代码，必须先向 dev 分支提交代码，然后合并到 main 分支。

## Git 提交规范

### 提交信息格式

提交信息采用 `<type> [<module>]: <description>` 格式，其中：

- **`type`**：表示提交类别，必须为以下之一：
  - `feat`：新增功能（对应语义化版本 MINOR 递增）
  - `fix`：修复 Bug（对应语义化版本 PATCH 递增）
  - `docs`：仅文档变更（README、注释等）
  - `update`：更新依赖库、配置文件、环境变量等非业务逻辑调整
  - `refactor`：代码重构（不改变外部行为，不新增功能，也不修复 Bug）
  - `perf`：性能优化
  - `test`：增加或修改测试代码
  - `chore`：构建工具、CI/CD、辅助脚本等杂务

- **`[module]`**：方括号内填写影响模块，如 `post`、`comment`、`auth`、`plugin` 等。若同时影响多个模块，使用逗号分隔，例如 `[auth, post]`。模块名请统一使用英文小写。

- **`description`**：使用英文，简明扼要描述本次提交（建议不超过 50 个字符）。

#### 示例

```
fix [post]: fix post like
feat [auth]: add remember me
update [comment]: update nested reply
refactor [auth, api]: extract common validator
perf [image]: lazy load thumbnails
```

### 提交命令

使用 `-m` 参数提供摘要与详细描述（摘要与正文之间 Git 会自动添加空行）：

```bash
git commit -m "<type> [<module>]: <description>" -m "<details changes>"
```

如果详细描述内容较多（超过 3 行），建议直接运行 `git commit` 进入编辑器编写，并保持每行不超过 **72 个字符**，以便在终端下排版整洁。

### 破坏性变更（BREAKING CHANGE）

当提交包含不兼容的 API 改动时（需要升级主版本号 MAJOR），必须在提交信息的**正文或页脚**中以 `BREAKING CHANGE:` 开头说明，或在 `type` 后紧跟 `!` 标记（两种方式择一即可）：

```bash
# 方式一：页脚注明
git commit -m "refactor [auth]: rewrite permission interceptor" -m "BREAKING CHANGE: 移除旧版 checkRole 方法，请改用 verifyScope。"

# 方式二：感叹号标记（推荐）
git commit -m "feat! [api]: change response structure" -m "详见迁移指南。"
```

### 标签（Tag）规范

发布正式版本时，使用**附注标签（annotated tag）**并附带结构化的发布说明：

```bash
# 创建标签（标题 + 详细变更列表）
git tag -a v1.0.0 -m "v1.0.0 发布（支持多因素认证）" -m "- 新增：OAuth2 登录\n- 修复：高并发下的 session 冲突\n- 移除：旧版短信网关接口"

# 推送标签到远程仓库
git push origin v1.0.0
```

建议同时推送代码和标签时使用 `git push --follow-tags`。

- 提交粒度：每次提交应只解决一个逻辑单元，避免混合多个不相关的改动。
- 分支命名：除您已有的前缀（fix-、feat- 等），建议加上 issue 编号或简短描述，如 fix-answer-vote-bug。
- 合并前检查：在合并到 dev 或 main 前，必须通过 CI 测试（构建、单元测试、代码规范检查）。
- 代码审查：所有合并请求（PR/MR）至少一名同事审阅，尤其是对核心模块的修改。
- 回滚预案：确保每次发布都可快速回滚（如保留 tag 或 release 版本号）。
- 代码修复：紧急修复必须先在 fix-[main] 分支进行，然后合并到 dev 分支，并立即 push 到远程仓库。
- 提交与推送：使用 `git commit -m "<type> [<module>]: <描述>"` 提交代码，然后使用 `git push` 推送到远程仓库。
- 提交描述：描述应该使用英文

## 代码删除

- 删除代码，前应备份，标记为临时删除，并在测试系统无影响后，进行彻底删除。
- 如果需要删除代码，必须先进行代码审查，确保删除的代码不会影响系统的正常运行。
- 如果删除的代码是必要的，必须提供充分的理由和计划，确保删除后的系统仍然可以正常运行。

## 测试与质量红线

- 后端 Service 层改动必须附带对应 `_test.go` 单元测试（覆盖率增量不得低于 70%），运行 `go test ./... -cover` 验证。
- 前端核心逻辑（API Client、Utils）需用 Vitest 编写单测；UI 交互变更需自测通过。
- 禁止使用 `any` 类型（后端 interface{} / 前端 any），存量代码逐步替换，新增代码零容忍。

## 数据库变更（Migration）

- 禁用 `AutoMigrate` 操作生产库；表结构变更需在 `backend/migrations/` 生成 SQL 迁移文件，命名格式 `YYYYMMDDHHMMSS_描述.up.sql` / `down.sql`。
- 删除列或重命名列必须分两步走（先代码兼容，后迁移清理）。

## 环境与启动

- 根目录提供 `docker-compose.yml` 启动 Postgres 依赖；启动前必须复制 `backend/config/private.yml.example` 并填入真实值。
- 根目录运行 `make dev`（后端）和 `pnpm dev`（前端）分别启动。

## AI 协作约束

- 任何改动必须先给出实施计划，不得一次修改超过 10 个文件（除非用户明确要求）。
- 新增依赖需先注释说明理由，待用户批准后再引入。
