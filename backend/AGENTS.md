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

## 配置

- 静态配置在 `config/*.yml`，启动时加载；运行期通过 `internal/infra/config` 的 `DynamicConfig` 支持热更新（修改 yml 后 `kill -SIGHUP <pid>` 或回调自动重建组件）。
- 命令行 flag 会覆盖为环境变量（见 `cmd/server/flags.go`）。
- 敏感配置（JWT 密钥、邮箱密码等）放 `config/private.yml`，提交时不得包含真实密钥。

## Git 提交规范

提交信息采用 `<type> [<module>]: <描述>` 形式，例如：

```
fix [answer]: fix answer vote
feat [flags]: add command line flags
docs [plugin]: update plugin dev description
update [name]: update post to creation
```

- type：`feat` / `fix` / `docs` / `update` / `refactor` 等
- module：方括号内填写影响模块（answer、plugin、question、error handle 等）

## 验证

- 修改业务代码后至少保证 `go build ./...` 通过。
- 涉及配置加载、参数校验、错误映射的逻辑，检查 `golangci-lint run` 与相关 `go test ./...`。
- 修改路由/中间件时对照 `internal/wire/routes.go` 的中间件顺序约定，避免破坏 RBAC 与限流。
