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

### Server/Client 组件边界

- `app/` 下页面/布局默认服务端（RSC）。首屏数据用 RSC `fetch` 获取，交互数据用 TanStack Query。
- `"use client"` 必须加在叶子节点（小组件），禁止在根布局或整个页面套用 Client Wrapper。

### 表单处理（react-hook-form + zod）

- 使用 `zodResolver` 解析 schema，TS 类型从 `z.infer<typeof schema>` 推导，杜绝重复定义 interface。
- 错误信息通过 `formState.errors.field` 统一渲染，提交按钮必须绑定 `isSubmitting` 防止重复点击。

### useEffect 红线

- **严格禁止**在 `useEffect` 中调用 API 获取数据（由 TanStack Query / RSC fetch 承担）。
- `useEffect` 仅限：订阅外部事件、DOM 操作、与第三方非 React 库集成。

### 图片与媒体（Core Web Vitals）

- 一律使用 `next/image` 处理头像/封面/上传图片，显式指定 `width`/`height` 或 `fill`，杜绝 `<img>`。
- 外部图片源需提前在 `next.config.ts` 的 `remotePatterns` 注册。

### 国际化（next-intl）

- 页面标题/描述必须通过 `getTranslations` 或 `useTranslations` 生成元数据，禁止硬编码。
- 动态路由的 `generateStaticParams` 必须包含 `locale` 参数。
- 国际化文件在 `src/i18n/messages`，格式为 `json`。

### 组件分层复用

- `shared/ui/`：无业务依赖的纯展示组件（通过 props 通信）。
- `features/*/components`：依赖 Query/Store 的容器组件，可组合 `shared/ui`。
- 禁止跨 `features` 模块直接导入组件，公共逻辑抽至 `shared/lib`。

## 代码风格

- import 别名用 `@/`（`@/shared/...`、`@/features/...`、`@/store/...`）。
- 代码注释习惯用中文；不新增无意义注释，重构时保留既有中文注释。
- 新增文件遵守 Prettier（单引号、无分号），提交前跑 `pnpm run fmt`。
- ESLint 禁用 `react/prop-types`（TS 类型已足够）、`react-hooks/set-state-in-effect`；新增组件按现有模式写。
- TypeScript 保持 `strict`；不要在组件里用 `any` 逃逸类型检查（现有代码中的 `as any` 仅作兼容，新代码避免）。

## 生产级执行规范

### 状态边界（Query vs Zustand）

- TanStack Query：管理所有服务端数据（API 响应），利用缓存和失效机制。
- Zustand：仅管理客户端 UI 状态（弹窗开关、当前 Tab、主题、分页页码）。
- 禁止在 Zustand Store 中存储 API 列表/详情数据。

### 导航与路由

- 内部跳转必须用 `next/link` 的 `<Link>` 或 `useRouter().push()`，禁止使用 `window.location.href` 或 `<a href>`（仅限站外链接）。
- 站外链接必须加 `rel="noopener noreferrer"`。

### 环境变量安全

- 浏览器端（`"use client"` / `shared/api/`）只能引用 `NEXT_PUBLIC_*` 变量。
- 服务端敏感变量（`JWT_SECRET` 等）仅限 `middleware.ts` / RSC / `app/api/` 使用，严禁泄露到客户端。

### 文件上传

- 文件必须封装为 `FormData` 传输；请求拦截器检测到 `FormData` 时自动移除 `Content-Type`，交由浏览器生成 boundary。
- 大文件上传必须绑定 `onUploadProgress` 展示进度条。

### 样式（Tailwind + DaisyUI）

- 条件类名统一使用 `clsx` + `twMerge`（封装为 `cn` 工具函数），禁止字符串拼接。
- 颜色优先使用 DaisyUI 语义变量（`bg-primary`、`text-base-content`），禁止硬编码色值。

## 即时通信（WebSocket/SSE）最佳实践

- **连接收口**：所有实时逻辑必须封装在 `useWebSocket` / `useRealtime` 自定义 Hook 中，组件卸载时强制 `close()`。
- **重连策略**：实现指数退避重连（间隔递增至最大 30s），监听 `navigator.onLine` 在网络恢复时即时重连。
- **心跳保活**：每 30s 发送心跳，超时无响应则断开重连；心跳定时器随组件卸载清除。
- **缓存更新（铁律）**：收到的实时数据**禁止**存入 Zustand，必须通过 `queryClient.setQueryData` 直接更新 TanStack Query 缓存（追加列表或更新详情字段），或 `invalidateQueries` 触发后台刷新。
- **多 Tab 管理**：使用 `BroadcastChannel` API 在标签页间同步消息，避免重复通知；独立连接或仅活跃 Tab 维持长连接。
- **离线队列**：用户离线时发出的消息暂存 Zustand + IndexedDB，连接恢复后按序补发。
- **类型安全**：在 `shared/api/types/realtime.model.ts` 定义严格的事件类型联合（Union Type），解析时使用类型守卫拦截未知事件。

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

## 配置

- 前端运行时配置在 `config.yml`（由 `make init-dev` 生成，勿手工编辑）；`next.config.ts` 从中读取 output、图片 remotePatterns、代理目标。
- 环境变量见 `.env.example`：`NEXT_PUBLIC_API_URL`、`JWT_SECRET`、`NEXT_PUBLIC_AVATAR_BASE_URL`、`BACKEND_URL`。
- 开发期 API/静态资源经 `next.config.ts` rewrites 代理到后端（`/api/v1`、`/store`、`/uploads`），`BACKEND_URL` 未配置时代理关闭。

## 验证

- 前端改动至少保证 `pnpm run type-check` 与 `pnpm run lint` 通过。
- 涉及路由/中间件改动，对照 `src/middleware.ts` 的鉴权与 i18n 逻辑，避免破坏认证跳转。
- 涉及 API 模块改动，确认响应类型与后端 Swagger/文档一致（`{ code, message, data }`）。
