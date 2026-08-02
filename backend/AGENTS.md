# AGENTS Rules — TinyForum Backend

Go 语言实现的论坛后端服务（Gin + GORM + PostgreSQL + Redis）。

---

## 0. AI 修复协议（最高优先级：根本原因修复）

> 本协议优先级高于所有编码规范。AI 在修复任何后端 Bug 时必须强制执行以下五步，**严禁跳过步骤直接输出 diff**。

### 第一步：5-Why 根因分析

AI 必须回答：

1. **现象**：报错/用户反馈的具体表现？
2. **表层直接原因**：哪一行代码抛出了异常或返回了错误数据？
3. **深层根因（追问 5 次）**：为什么会出现这个数据/状态？——是并发写入未加锁？数据库字段溢出？前端传参与后端 `request` 结构体不匹配？状态机流转遗漏了事件？
4. **影响范围**：孤立模块还是影响全局？

### 第二步：复现测试先行（红灯准入）

- AI **必须先编写单元测试**捕获 Bug，再修改生产代码。
- 在对应模块生成 `_test.go`，模拟边缘条件，运行 `go test` 必须**失败（Red）**。
- **门禁**：未提供可复现的失败测试，生产代码改动不予通过。

### 第三步：强制分层定位（禁止跨层绕过）

AI 必须明确修复发生的**层级**，禁止用上层 `if` 掩盖下层缺陷：

| 层级           | 允许的“根本修复”                                    | 严厉禁止的“补丁行为”                          |
| :------------- | :-------------------------------------------------- | :-------------------------------------------- |
| **Repository** | 修正 SQL 条件、添加索引、处理 NULL/零值、使用乐观锁 | 在 Service 层加 `if` 过滤空数据               |
| **Service**    | 修正状态机流转、引入幂等/重试、添加事务边界         | 在 Handler 用 `recover()` 捕获 panic 返回 200 |
| **Handler**    | 修正参数校验规则 `binding`、修正 VO 字段映射        | 在响应中硬编码 `"data": null` 或 `"code": 0`  |

### 第四步：数据迁移解耦（若涉及 Schema）

- **严禁**在 Go 代码中做类型转换来“适配”错误的数据库字段。
- **必须**提供 `backend/migrations/` 下的 `.up.sql` 和 `.down.sql` 迁移文件。
- **兼容性铁律**：删除列或重命名必须**分两步提交**（先增后弃），确保新旧代码同时兼容。

### 第五步：提交信息自审

- **坏的描述（补丁）**：`fix [auth]: add nil check for token`
- **好的描述（根因）**：`fix [auth]: ensure JWT parser returns ErrInvalidToken on malformed signature instead of panic`

---

## 1. 常用命令

```bash
make run          # 启动服务（go run ./cmd/server）
make build        # 编译到 bin/server
make tidy         # go mod tidy
make wire         # 重新生成依赖注入代码（internal/wire/...）
make docs-api     # 重新生成 Swagger 文档（docs/）
go build ./...    # 编译检查
go test ./...     # 运行测试
golangci-lint run # 静态检查（启用 govet；gofmt 作为 formatter）
```

运行参数（`cmd/server/flags.go`）：`--config-dir`、`--port`、`--env`、`--verbosity`、`--audit`、`--version`。

---

## 2. 项目结构

```
cmd/server            # 入口：解析 flag → 加载配置 → wire 装配 → 启动 HTTP 服务
internal/
  handler/            # HTTP 层：gin handler + 路由注册（RegisterRoutes），按模块分目录
  service/            # 业务逻辑层：面向接口（interface.go + NewXxxService 构造器）
  repository/         # 数据访问层：GORM 操作，面向接口
  model/              # 数据模型：
    do/               #   GORM 实体（database object）
    request/          #   请求入参
    vo/               #   视图对象（响应给前端）
    dto/              #   数据传输对象（层间传递）
    bo/               #   业务对象
    common/           #   通用结构（BasicResponse、PageResult 等）
  routes/             # 页面/静态文件路由
  middleware/         # 鉴权、Casbin RBAC、限流、内容安全等中间件
  infra/              # 基础设施：config、casbinx、lua、ratelimit、sensitive、validator、email
  wire/               # 依赖注入装配（手工 wire，非代码生成）
  storage/            # 文件存储驱动（local）
  strategy/           # 策略注册表
  startup/            # 启动横幅与启动信息打印
  botapi/             # Bot 访问论坛的 API 封装
  job/                # 后台定时任务（如清理临时文件）
pkg/                  # 可复用通用包（errors、response、logger、jwt、mask 等）
config/               # YAML 配置文件（basic/postgres/redis/private/ai/risk_control）
store/                # 插件静态资源
plugins/              # 插件（如 demo）
```

---

## 3. 技术栈与关键依赖

- **Web 框架**：Gin
- **ORM**：GORM（PostgreSQL 驱动）
- **依赖注入**：Wire（手工装配，非代码生成）
- **日志**：Zap
- **鉴权**：JWT + Casbin（RBAC）
- **缓存/消息**：Redis
- **配置**：YAML（支持 SIGHUP 热更新）
- **WebSocket**：`gorilla/websocket`
- **API 文档**：Swagger（`make docs-api`）

---

## 4. 核心架构规范

### 4.1 分层依赖（铁律）

调用链严格单向：`handler → service → repository`，**禁止反向依赖**。

- **handler**：只做参数绑定、调用 service、用 `pkg/response` 写响应。每个模块一个目录，含 `RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet)`。
- **service**：定义 `interface` + `struct` + `NewXxxService(...)` 构造器。依赖以接口/构造器参数注入，结构体字段持有依赖。
- **repository**：只做数据访问，定义接口 + GORM 实现。事务相关走 `repository/transaction`。

### 4.2 模型隔离（严禁混用）

- `do`：GORM 表实体
- `request`：API 入参（绑定与校验）
- `vo`：API 出参（响应给前端）
- `dto`：层间数据传输
- `bo`：业务对象

**禁止**将 `do` 直接作为 API 响应返回，必须转换为 `vo`。

### 4.3 响应与错误处理

- 一律用 `pkg/response`：`response.Success`、`response.SuccessPage`、`response.Created`，以及语义化错误 `response.BadRequest/Unauthorized/Forbidden/NotFound/InternalError/TooManyRequests`。
- handler 捕获错误后统一调用 `response.HandleError(c, err)`。
- 业务错误优先复用 `pkg/errors` 预定义 `AppError`（如 `apperrors.ErrPostNotFound`），附加信息用链式方法 `ErrXxx.WithDetail(...).WithCause(err)`，**不要直接修改全局错误实例**。
- 统一响应格式：`{"code": 0, "message": "success", "data": ..., "timestamp": ..., "request_id": ..., "trace_id": ...}`。

### 4.4 路由与中间件顺序（`internal/wire/routes.go`）

叠加顺序（由外到内）：

1. `mw.Auth()` / `mw.OptionalAuth()` —— 解析 JWT，注入 `user_role`
2. `mw.CasbinAuth()` —— 路由级 RBAC（读 `user_role` 决策）
3. `mw.RateLimit(...)` —— 限流
4. `mw.ContentCheck(...)` —— 内容安全
5. `mw.ModeratorRequired(...)` —— 版主细粒度权限（查库）

- 公开路由用 `OptionalAuth`（无 token 注入 guest）
- 私有路由统一 `Auth` + `CasbinAuth`
- API 基路径为 `/api/v1`

### 4.5 配置管理

- 静态配置在 `config/*.yml`，启动时加载。
- 运行期通过 `internal/infra/config` 的 `DynamicConfig` 支持热更新（修改 yml 后 `kill -SIGHUP <pid>` 或回调自动重建组件）。
- 命令行 flag 会覆盖为环境变量（见 `cmd/server/flags.go`）。
- 敏感配置（JWT 密钥、邮箱密码等）放 `config/private.yml`，**禁止提交真实密钥**。

---

## 5. 代码质量与安全检查清单

### 5.1 质量标准（零容忍）

AI 生成的代码若触及以下红线，**视为无效交付，必须重写**：

1. **类型逃逸**：使用 `interface{}` 且无类型断言说明。
2. **硬编码**：JWT 密钥、DB DSN 等配置硬编码在代码中。
3. **裸日志**：捕获错误后仅 `fmt.Println`，未接入 Zap。
4. **魔数/魔字符串**：业务状态码未定义为常量。
5. **跳过 Migration**：修改了 `do` 但未生成 SQL 迁移文件，或注释了 `AutoMigrate`。

### 5.2 安全检查清单（每次提交前自检）

- [ ] SQL 查询是否使用了 GORM 参数化占位符？是否拼接了用户输入？
- [ ] 涉及权限判断的接口，是否在 Handler 层使用了 `CasbinAuth` 中间件？
- [ ] 用户上传文件/图片的接口，是否校验了 MIME 类型和文件大小限制？
- [ ] 日志中是否明文打印了密码、Token 或手机号？

---

## 6. 后端最佳实践（逐条约束）

### 6.1 Context 与超时

- Service/Repository 方法首参必须为 `context.Context`；GORM 用 `db.WithContext(ctx)`。
- 外部 IO 调用必须设 `context.WithTimeout`（数据库/Redis/AI 接口），防止积压雪崩。

### 6.2 事务

- 事务由 Service 层闭包开启（`repo.Transaction`），Handler 不得开启事务。
- Repository 只执行传入的 `*gorm.DB`，不自建 `Begin/Commit`。

### 6.3 查询防御

- 列表查询禁用 `OFFSET` 深分页，统一改用游标（`WHERE id < ? ORDER BY id DESC`）。
- `Preload` 必须显式指定关联字段，禁止 `clause.Associations`。

### 6.4 日志（Zap）

- 错误日志必带 `zap.Error(err)`（err 需 wrap stack）与 `user_id`、`request_id`。
- 数据库查询 >100ms 记录 Warn；循环内禁止 Debug 日志。

### 6.5 依赖注入（Wire）

- 构造器返回接口类型（便于 Mock）；出现循环依赖时需拆分解耦。
- 新增 handler/service/repository 后，在 `internal/wire/` 对应文件中手工装配（`NewServices`/`NewHandlers`/`NewRepositories` 等）。

### 6.6 启动快速失败

- 关键配置（JWT 密钥、DB DSN）缺失或非法时，启动阶段必须 `panic`，不带病运行。

### 6.7 测试分层

- 单元测试（Service/Utils）用 `mockgen` 生成 Mock，不依赖真实 DB。
- 集成测试加 `//go:build integration` 标签，与单元测试分开运行。

### 6.8 优雅关闭（HTTP）

- 监听信号后调用 `srv.Shutdown(ctx)`，设定最大等待时长（如 10s），配合 WebSocket 一起释放资源。

### 6.9 WebSocket 最佳实践

- **库选型**：`github.com/gorilla/websocket`，升级逻辑封装在 `handler/websocket.go`。
- **架构（Hub + Client）**：全局 Hub 管理连接池与房间订阅；Client 封装 Conn、用户 ID、发送通道，独立读/写协程。
- **认证（握手时校验）**：从 Gin Context 获取 JWT 鉴权后的 `userID`，校验失败返回 401，禁止通过消息体延迟认证。
- **并发安全**：Map 操作用 `sync.RWMutex` 保护；写协程监听 `send` 通道，读协程循环读取；断开连接时 `defer` 清理资源并关闭通道。
- **心跳**：服务端 30s 主动发送 Ping 帧，客户端 Pong 响应；超时（60s）自动断开并清理。
- **消息协议**：统一 JSON 格式 `{"type":"event","room":"...", "data":{...}}`，读协程按 `type` 路由；业务数据必须通过 Service 层处理。
- **安全防护**：单条消息限制 4KB，`send` 通道缓冲 256，消息频率超限（>10条/秒）主动断开连接。
- **水平扩展（预备）**：多 Pod 部署时引入 Redis Pub/Sub 转发广播消息，Hub 内集成 Redis 订阅者。
- **优雅关闭**：监听 `SIGINT/SIGTERM`，停止接受新连接 → 发送关闭帧 → 等待 5s → 释放所有资源后退出。
- **日志**：连接事件必须记录 Zap 结构化日志，关联 `user_id`，便于追踪。
- **Wire 装配**：Hub 和 Redis 订阅者实例在 `internal/wire/` 中完成依赖注入。

---

## 7. 数据库变更规范（Migration）

- **命名格式**：`backend/migrations/{timestamp}_{描述}.up.sql` 与 `down.sql`。
- **不可逆操作**：删除列 (`DROP COLUMN`) 或修改列类型导致数据丢失的，必须先在 `up.sql` 中备份数据到临时表，并在 `down.sql` 提供恢复脚本。
- **测试环境验证**：AI 提供的迁移脚本必须在 `docker-compose` 启动的本地 PostgreSQL 验证通过。

---

## 8. 代码风格

- 模块名为 `tiny-forum`，内层 import 用 `tiny-forum/internal/...`。
- 公开类型、接口方法、handler 需要有注释（代码内普遍用中文注释）。
- 不新增无意义注释；重构时保留既有中文注释习惯。
- 所有 handler 方法加 Swagger 注解（`@Summary` / `@Tags` / `@Router` 等）。
- 修改后运行 `make wire` 生成装配代码，并 `golangci-lint run` 检查。

---

## 9. Git 提交规范

### 9.1 提交信息格式

`<type> [<module>]: <description>`，其中：

- **type**：`feat` / `fix` / `docs` / `update` / `refactor` / `perf` / `test` / `chore`
- **[module]**：影响模块（如 `post`、`auth`），多模块逗号分隔
- **description**：英文，现在时态，不超过 50 字符

示例：

```
fix [post]: fix post like
feat [auth]: add remember me
refactor [auth, api]: extract common validator
```

### 9.2 破坏性变更（BREAKING CHANGE）

必须在 footer 包含 `BREAKING CHANGE:` 说明，或在 `type` 后紧跟 `!`：

```bash
git commit -m "feat! [api]: change response structure" -m "详见迁移指南。"
```

### 9.3 分支策略

- `main` 为生产保护分支，**任何人（包括 AI）严禁直接 push**。
- 修复任务切 `fix-[issue_id]`；新功能切 `feat-[name]`；重构切 `refactor-[name]`。
- 所有变更必须通过 `dev` → PR/MR 合并。

### 9.4 项目规范

- 任何对代码的修改都应该创建新的分支，修改完成后应该 commit
- 每次提交应该包含一个 commit，并且每个 commit 只做一件事情
- 每个提交信息应该简洁明了，并且应该包含一个简短的描述
- 每个提交信息应该使用英文，并且应该使用现在时态
- 提交的代码应该经过代码审查，并且应该通过所有的测试

---

## 10. AI 工作量约束

- **文件修改限制**：单次对话中，AI 不得一次性修改超过 **10 个文件**，除非用户明确要求“全量重构”。
- **依赖引入**：AI 不得主动引入新的第三方依赖，必须先注释说明理由并等待用户批准。
- **代码删除**：删除任何逻辑前，AI 必须输出“备份提示”并标记为临时注释（`// TODO: confirm before removal`），确保不影响现有调用链。

---

## 11. 验证

- 修改业务代码后至少保证 `go build ./...` 通过。
- 涉及配置加载、参数校验、错误映射的逻辑，检查 `golangci-lint run` 与相关 `go test ./...`。
- 修改路由/中间件时对照 `internal/wire/routes.go` 的中间件顺序约定，避免破坏 RBAC 与限流。
