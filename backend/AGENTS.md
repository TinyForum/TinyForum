# AGENTS.md — TinyForum Backend

Go 语言实现的论坛后端服务（Gin + GORM + PostgreSQL + Redis）。

## 常用命令

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

## 项目结构

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

## 架构分层约定

调用链严格单向：`handler → service → repository`，禁止反向依赖。

- **handler**：只做参数绑定、调用 service、用 `pkg/response` 写响应。每个模块一个目录，含 `RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet)`。
- **service**：定义 `interface` + `struct` + `NewXxxService(...)` 构造器。依赖以接口/构造器参数注入，结构体字段持有依赖。
- **repository**：只做数据访问，定义接口 + GORM 实现。事务相关走 `repository/transaction`。
- **wire**：新增 handler/service/repository 后，在 `internal/wire/` 对应文件中手工装配（`NewServices`/`NewHandlers`/`NewRepositories` 等）。

## 最佳实践

### Context 与超时

- Service/Repository 方法首参必须为 `context.Context`；GORM 用 `db.WithContext(ctx)`。
- 外部 IO 调用必须设 `context.WithTimeout`（如数据库/Redis/AI 接口），防止请求积压雪崩。

### 事务

- 事务由 Service 层闭包开启（`repo.Transaction`），Handler 不得开启事务。
- Repository 只执行传入的 `*gorm.DB`，不自建 `Begin/Commit`。

### 查询防御

- 列表查询禁用 `OFFSET` 深分页，统一改用游标（`WHERE id < ? ORDER BY id DESC`）。
- `Preload` 必须显式指定关联字段，禁止 `clause.Associations`。

### 日志（Zap）

- 错误日志必带 `zap.Error(err)`（err 需 wrap stack）与 `user_id`、`request_id`。
- 数据库查询 >100ms 记录 Warn；循环内禁止 Debug 日志。

### 依赖注入（Wire）

- 构造器返回接口类型（便于 Mock）；出现循环依赖时需拆分解耦。

### 启动快速失败

- 关键配置（JWT 密钥、DB DSN）缺失或非法时，启动阶段必须 `panic`，不带病运行。

### 测试分层

- 单元测试（Service/Utils）用 `mockgen` 生成 Mock，不依赖真实 DB。
- 集成测试加 `//go:build integration` 标签，与单元测试分开运行。

### 优雅关闭（HTTP）

- 监听信号后调用 `srv.Shutdown(ctx)`，设定最大等待时长（如 10s），配合 WebSocket 一起释放资源。

## 代码风格

- 模块名为 `tiny-forum`，内层 import 用 `tiny-forum/internal/...`。
- 公开类型、接口方法、handler 需要有注释（代码内普遍用中文注释）。
- 不新增无意义注释；重构时保留既有中文注释习惯。
- 所有 handler 方法加 Swagger 注解（`@Summary` / `@Tags` / `@Router` 等）。
- 修改后运行 `make wire` 生成装配代码，并 `golangci-lint run` 检查。

## 响应与错误处理

- 一律用 `pkg/response`：`response.Success`、`response.SuccessPage`、`response.Created`，以及语义化错误 `response.BadRequest/Unauthorized/Forbidden/NotFound/InternalError/TooManyRequests`。
- handler 捕获错误后统一调用 `response.HandleError(c, err)`。
- 业务错误优先复用 `pkg/errors` 里预定义的 `AppError` 实例（如 `apperrors.ErrPostNotFound`），附加信息用链式方法 `ErrXxx.WithDetail(...).WithCause(err)`，不要直接改全局变量。
- 响应统一格式：`{"code": 0, "message": "success", "data": ..., "timestamp": ..., "request_id": ..., "trace_id": ...}`。

## 路由与中间件约定（internal/wire/routes.go）

中间件叠加顺序：

1. `mw.Auth()` / `mw.OptionalAuth()` —— 解析 JWT，注入 `user_role`
2. `mw.CasbinAuth()` —— 路由级 RBAC（读 `user_role` 决策）
3. `mw.RateLimit(...)` —— 限流
4. `mw.ContentCheck(...)` —— 内容安全
5. `mw.ModeratorRequired(...)` —— 版主细粒度权限（查库）

公开路由用 `OptionalAuth`（无 token 注入 guest），私有路由统一 `Auth` + `CasbinAuth`。API 基路径为 `/api/v1`。

## 即时通信（WebSocket）最佳实践

- **库选型**：使用 `github.com/gorilla/websocket`，升级逻辑封装在 `handler/websocket.go`，由 Gin 路由承载。
- **架构（Hub + Client）**：全局 Hub 管理连接池与房间订阅；Client 封装 Conn、用户 ID、发送通道，独立读/写协程。
- **认证（握手时校验）**：从 Gin Context 中获取 JWT 鉴权后的 `userID`，校验失败返回 401，禁止通过消息体延迟认证。
- **并发安全**：Map 操作用 `sync.RWMutex` 保护；写协程监听 `send` 通道，读协程循环读取；断开连接时 `defer` 清理资源并关闭通道。
- **心跳**：服务端 30s 主动发送 Ping 帧，客户端 Pong 响应；超时（60s）自动断开并清理。
- **消息协议**：统一 JSON 格式 `{"type":"event","room":"...", "data":{...}}`，读协程按 `type` 路由；业务数据必须通过 Service 层处理。
- **安全防护**：单条消息限制 4KB，`send` 通道缓冲 256，消息频率超限（>10条/秒）主动断开连接。
- **水平扩展（预备）**：多 Pod 部署时引入 Redis Pub/Sub 转发广播消息，Hub 内集成 Redis 订阅者。
- **优雅关闭**：监听 `SIGINT/SIGTERM`，停止接受新连接 → 发送关闭帧 → 等待 5s → 释放所有资源后退出。
- **日志**：连接事件必须记录 Zap 结构化日志，关联 `user_id`，便于追踪。
- **Wire 装配**：Hub 和 Redis 订阅者实例在 `internal/wire/` 中完成依赖注入。

## 配置

- 静态配置在 `config/*.yml`，启动时加载；运行期通过 `internal/infra/config` 的 `DynamicConfig` 支持热更新（修改 yml 后 `kill -SIGHUP <pid>` 或回调自动重建组件）。
- 命令行 flag 会覆盖为环境变量（见 `cmd/server/flags.go`）。
- 敏感配置（JWT 密钥、邮箱密码等）放 `config/private.yml`，提交时不得包含真实密钥。

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

## 验证

- 修改业务代码后至少保证 `go build ./...` 通过。
- 涉及配置加载、参数校验、错误映射的逻辑，检查 `golangci-lint run` 与相关 `go test ./...`。
- 修改路由/中间件时对照 `internal/wire/routes.go` 的中间件顺序约定，避免破坏 RBAC 与限流。
