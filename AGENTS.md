# AGENTS.md — TinyForum

Tiny Forum 是一个开源论坛项目：**Go（Gin + GORM）后端** × **Next.js（App Router）前端** × **Vite 插件体系** × PostgreSQL 持久化。

本文件约定**跨子项目的代码编写规范与最佳实践**。运行命令、数据库相关请勿在这里维护，分别见各子项目文档：

- 后端：`backend/AGENTS.md`
- 前端：`frontend/AGENTS.md`
- 插件：`plugin/`（见 `manifest.json` 与 `docs/zh-CN/dev/plugin.md`）

## 仓库组成

| 子项目 | 位置 | 技术栈 |
| ------ | ---- | ------ |
| 后端 | `backend/` | Go 1.26, Gin, GORM, Wire（手工注入）, JWT, Zap |
| 前端 | `frontend/` | Next.js 16（App Router）+ React 19 + TypeScript, Tailwind/DaisyUI, TanStack Query, Zustand, Tiptap, next-intl |
| 插件 | `plugin/` | Vite + React（独立构建，产物 `dist/main.js` + `manifest.json`） |

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

## Git 提交规范

提交信息统一 `<type> [<module>]: <描述>`（type：`feat` / `fix` / `docs` / `update` / `refactor` 等；module 为影响模块，如 post、comment、answer、plugin）：

```
fix [answer]: fix answer vote
feat [auth]: add remember me
docs [plugin]: update plugin dev description
```
