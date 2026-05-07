# response 包使用说明

`response` 包为 Gin 框架提供了一套统一、规范的 HTTP 响应处理工具。它封装了成功响应、错误响应、分页数据、字段校验错误以及全局错误处理中间件，帮助你快速构建 RESTful API。

## 功能特性

- ✅ **统一响应结构**：所有响应均遵循 `vo.BasicResponse` 格式
- ✅ **语义化方法**：`Success`, `BadRequest`, `Unauthorized`, `NotFound` 等
- ✅ **分页数据集**：`SuccessPage` 自动计算 `has_more`
- ✅ **字段级校验错误**：自动将 `validator` 错误转为友好提示
- ✅ **链路追踪**：支持 `RequestID` / `TraceID` 自动提取与传递
- ✅ **错误集中处理**：`HandleError` 与中间件兜底未处理的错误
- ✅ **Panic 恢复**：`RecoveryMiddleware` 防止服务崩溃

## 安装

```go
import "tiny-forum/internal/pkg/response"
```

## 快速开始

### 1. 注册中间件

```go
r := gin.Default()

// 恢复中间件（建议放在最前面）
r.Use(response.RecoveryMiddleware())

// 尾部错误处理中间件（建议放在最后，兜底未处理的错误）
r.Use(response.ErrorHandlerMiddleware())
```

### 2. 在 Handler 中使用

```go
func GetUserHandler(c *gin.Context) {
    user, err := service.GetUserByID(123)
    if err != nil {
        // 使用统一错误处理
        response.HandleError(c, err)
        return
    }
    // 成功响应
    response.Success(c, user)
}
```

## API 参考

### 成功响应

| 方法                                          | 说明                   | HTTP 状态码                 |
| --------------------------------------------- | ---------------------- | --------------------------- |
| `Success(c, data, opts...)`                   | 通用成功响应           | 200                         |
| `SuccessWithMessage(c, msg, data)`            | 带自定义消息的成功响应 | 200                         |
| `SuccessPage(c, list, total, page, pageSize)` | 分页数据响应           | 200                         |
| `Created(c, data, location)`                  | 资源创建成功           | 201，同时设置 `Location` 头 |
| `NoContent(c)`                                | 无内容响应             | 204                         |

**示例：**

```go
// 普通成功
response.Success(c, map[string]string{"name": "test"})

// 分页数据
response.SuccessPage(c, users, 100, 1, 10)

// 创建资源
response.Created(c, newUser, "/users/100")
```

### 错误响应（语义化）

| 方法                        | 说明         | HTTP 状态码 | 业务码              |
| --------------------------- | ------------ | ----------- | ------------------- |
| `BadRequest(c, msg)`        | 请求参数错误 | 400         | CodeInvalidRequest  |
| `Unauthorized(c, msg)`      | 未授权       | 401         | CodeUnauthorized    |
| `Forbidden(c, msg)`         | 权限不足     | 403         | CodeForbidden       |
| `NotFound(c, msg)`          | 资源不存在   | 404         | CodeNotFound        |
| `Conflict(c, msg)`          | 资源冲突     | 409         | CodeInvalidRequest  |
| `TooManyRequests(c, msg)`   | 请求过于频繁 | 429         | CodeTooManyRequests |
| `InternalError(c, msg)`     | 系统内部错误 | 500         | CodeInternalError   |
| `ValidationFailed(c, errs)` | 参数校验失败 | 400         | CodeValidation      |

**示例：**

```go
response.BadRequest(c, "用户ID不能为空")
response.NotFound(c, "文章不存在")
response.ValidationFailed(c, validationErrors)
```

### 统一错误处理 `HandleError(c, err)`

该函数会根据错误类型自动选择正确的响应方式：

- **\*apperrors.AppError**：直接映射其 HTTP 状态码和业务码
- **validator.ValidationErrors**：展开所有字段错误，返回 `ValidationFailed` 响应（字段级错误列表）
- **context.DeadlineExceeded**：返回 `504 Gateway Timeout`
- **context.Canceled**：静默结束（客户端主动取消，不记录日志）
- **其他错误**：记录日志，返回 `500 Internal Server Error`

**推荐使用方式**：在 Handler 层调用，不再手动判断错误类型。

```go
if err := service.DoSomething(); err != nil {
    response.HandleError(c, err)
    return
}
```

### 自定义选项

可通过 `Option` 函数为响应附加额外字段：

| 选项                   | 作用               |
| ---------------------- | ------------------ |
| `WithTraceID(traceID)` | 设置响应的 TraceID |
| `WithMessage(msg)`     | 覆盖默认的 Message |

**示例：**

```go
response.Success(c, data, response.WithTraceID("abc-123"))
response.Success(c, data, response.WithMessage("自定义成功消息"))
```

## 错误处理中间件

### RecoveryMiddleware

恢复 `panic`，防止服务崩溃，记录日志后返回 500。

### ErrorHandlerMiddleware

放在路由组最末尾，用于兜底处理那些没有直接调用 `response.*` 方法，而是通过 `c.Error(err)` 记录的错误。

**使用方式：**

```go
r := gin.New()
r.Use(response.RecoveryMiddleware())    // 最前
r.Use(gin.Logger())
r.Use(gin.Recovery())                   // 如果使用了 gin.Recovery，可以去掉
// ... 路由注册 ...
r.Use(response.ErrorHandlerMiddleware()) // 最后
```

## 响应结构

所有响应均使用 `vo.BasicResponse`：

```go
type BasicResponse struct {
    Code      int         `json:"code"`
    Message   string      `json:"message"`
    Data      interface{} `json:"data,omitempty"`
    Timestamp int64       `json:"timestamp"`
    RequestID string      `json:"request_id,omitempty"`
    TraceID   string      `json:"trace_id,omitempty"`
}
```

成功时 `code` 为 `0`，错误时根据 `apperrors` 定义返回非零业务码。

## 注意事项

1. **请确保 `vo.BasicResponse` 定义与包内一致**（包含上述字段）。
2. **`HandleError` 会调用 `c.Abort()`**，因此调用后应立即 `return`。
3. 使用 `Created` 时，`location` 参数可以留空（不设置 Location 头）。
4. 校验错误消息目前仅支持中文，可根据需求扩展 `validationMessage`。
5. 日志使用自定义的 `logger` 包（基于 `zap`），请确保已在项目中正确初始化。

## 完整示例

```go
package handler

import (
    "tiny-forum/internal/pkg/response"
    "github.com/gin-gonic/gin"
)

type UserHandler struct{}

func (h *UserHandler) Create(c *gin.Context) {
    var req CreateUserReq
    if err := c.ShouldBindJSON(&req); err != nil {
        response.HandleError(c, err)
        return
    }

    user, err := h.service.Create(req)
    if err != nil {
        response.HandleError(c, err)
        return
    }

    response.Created(c, user, "/users/"+user.ID)
}

func (h *UserHandler) List(c *gin.Context) {
    users, total, err := h.service.List(c.Query("page"), c.Query("size"))
    if err != nil {
        response.HandleError(c, err)
        return
    }
    response.SuccessPage(c, users, total, page, pageSize)
}
```

---

该包极大简化了 Gin 项目的响应处理，推荐在业务代码中统一使用。如有更多定制需求，可以扩展 `Option` 和错误类型映射。
