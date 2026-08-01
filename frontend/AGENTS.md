# AGENTS.md — TinyForum Frontend

Next.js 16（App Router）+ React 19 + TypeScript 实现的论坛前端。技术栈：Tailwind + DaisyUI、TanStack Query、Zustand、Tiptap、next-intl。

## 常用命令

```bash
pnpm run dev         # 开发服务器（:3000，经 next.config.ts rewrites 代理 /api/v1 到后端）
pnpm run build       # 生产构建
pnpm run start       # 启动生产构建
pnpm run lint        # ESLint（husky pre-commit 也会跑）
pnpm run lint:fix    # 自动修复
pnpm run type-check  # tsc --noEmit
pnpm run fmt         # Prettier 格式化 src/**/*.{ts,tsx,css}
```

包管理器为 **pnpm**（存在 `pnpm-lock.yaml`），新增依赖用 `pnpm add`。

## 项目结构

```
src/app/[locale]/    # 路由层（Next.js App Router，页面即 View；国际化前缀 /zh-CN、/en-US）
src/features/        # 功能模块（按领域拆分：post、comment、auth、admin、board、user...）
  <module>/components/   # View：模块专用 UI 组件
  <module>/hooks/        # View+Query：组合 React Query / Zustand 状态
src/layout/          # 布局与跨页组件（Providers、导航、RichEditor、首页/帖子等页面骨架）
src/shared/          # 共享资源
  api/               #   Client 层：client.ts（axios 实例）+ modules/ + types/
  ui/                #   通用 UI 组件（common、editor、upload、nav、rank...）
  lib/               #   工具函数（纯函数：format、helpers、utils、yaml）
  constant/          #   常量（built-in_plugins.ts）
  localdb/           #   IndexedDB（Dexie）
src/store/           # State 层：Zustand 全局状态（auth、token、login、logout...）
src/i18n/            # 国际化（next-intl + next-intl-split，消息在 messages/）
src/middleware.ts    # 中间件：i18n + 路由鉴权 + JWT 校验
next.config.ts       # rewrites 代理 /api/v1、/store、/uploads → 后端；图片 remotePatterns
config.yml           # 前端构建配置（由 make init-dev 生成，勿手工编辑）
```

### 关键依赖

- 数据请求：TanStack Query（`@tanstack/react-query`）+ axios
- 全局状态：Zustand（`src/store/`）
- 表单：react-hook-form + zod
- 富文本：Tiptap（`shared/ui/editor`）
- UI：Tailwind CSS + DaisyUI + Headless UI + Heroicons/Lucide
- 国际化：next-intl（locales：`zh-CN` 默认、`en-US`）
- 本地存储：Dexie（IndexedDB）

## 架构分层约定

依赖方向：`View → State/Query → Client`；数据流向：`User Action → View → Query/State → Client → 后端`。

- **View**（`app/[locale]/` 页面、`features/*/components`、`layout/`）只渲染 UI，不直接调 API、不放业务逻辑。
- **hooks**（`features/*/hooks`）封装 React Query / Zustand 逻辑，对应后端 service。Query key 用集中定义的 `xxxKeys` 工厂（如 `postKeys.list(params)`）；mutation 成功后 `invalidateQueries` 刷新列表/详情。
- **Client 层**（`shared/api/`）：`client.ts` 统一 axios 实例（baseURL `/api/v1`，`withCredentials`，401 自动跳 `/auth/login`）；API 方法按模块放 `shared/api/modules/*.ts`，类型放 `shared/api/types/*.model.ts`（含 `.do.ts` 领域对象）。
- **State 层**（`src/store/`）：Zustand，跨组件共享的全局状态（如 auth）才放这里；局部状态留在组件内。

### API 约定

- 统一请求基路径 `/api/v1`，与后端响应格式对应：`{ code, message, data, ... }`。
- API 方法返回 `apiClient.get<ApiResponse<T>>(...)`，hook 中解出 `response.data.data` 后再消费。
- 新增模块步骤：`shared/api/types/<module>.model.ts` 定义类型 → `shared/api/modules/<module>.ts` 定义 api 对象 → `features/<module>/hooks` 封装 Query/Mutation。

### Hooks 约定

使用 TanStack Query，最佳实践：

- 拥抱声明式：用useQuery和useMutation声明数据依赖，而不是在useEffect中命令式获取。
- 设计好查询键：建立清晰的工厂模式，这是高效缓存的基础。
- 配置好过期时间：为不同数据设置合理的staleTime，减少不必要的请求。
- 完善变更处理：总是使相关查询失效，并对关键操作使用乐观更新。

## 代码风格

- import 别名用 `@/`（`@/shared/...`、`@/features/...`、`@/store/...`）。
- 代码注释习惯用中文；不新增无意义注释，重构时保留既有中文注释。
- 新增文件遵守 Prettier（单引号、无分号），提交前跑 `pnpm run fmt`。
- ESLint 禁用 `react/prop-types`（TS 类型已足够）、`react-hooks/set-state-in-effect`；新增组件按现有模式写。
- TypeScript 保持 `strict`；不要在组件里用 `any` 逃逸类型检查（现有代码中的 `as any` 仅作兼容，新代码避免）。

## Git 提交规范

提交信息采用 `<type> [<module>]: <描述>`（husky 会跑 commitlint conventional commits + `pnpm run lint`）：

```
fix [post]: fix post like
feat [auth]: add remember me
update [comment]: update nested reply
```

type：`feat` / `fix` / `docs` / `update` / `refactor` 等；module 方括号内填写影响模块（post、comment、auth、plugin...）。

## 配置

- 前端运行时配置在 `config.yml`（由 `make init-dev` 生成，勿手工编辑）；`next.config.ts` 从中读取 output、图片 remotePatterns、代理目标。
- 环境变量见 `.env.example`：`NEXT_PUBLIC_API_URL`、`JWT_SECRET`、`NEXT_PUBLIC_AVATAR_BASE_URL`、`BACKEND_URL`。
- 开发期 API/静态资源经 `next.config.ts` rewrites 代理到后端（`/api/v1`、`/store`、`/uploads`），`BACKEND_URL` 未配置时代理关闭。

## 验证

- 前端改动至少保证 `pnpm run type-check` 与 `pnpm run lint` 通过。
- 涉及路由/中间件改动，对照 `src/middleware.ts` 的鉴权与 i18n 逻辑，避免破坏认证跳转。
- 涉及 API 模块改动，确认响应类型与后端 Swagger/文档一致（`{ code, message, data }`）。
