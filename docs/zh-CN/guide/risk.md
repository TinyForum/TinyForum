# 风控系统集成指南

## 1. 新增依赖（go.mod）

```bash
go get github.com/redis/go-redis/v9
```

> 敏感词过滤使用内置 DFA 实现，无需额外依赖。
> 词库 > 10万条时可替换为 `github.com/BobuSumisu/aho-corasick`。

---

## 2. 数据库迁移

在 `AutoMigrate` 调用中增加新模型：

```go
// internal/wire/wire.go 或 main.go 的 DB 初始化处
db.AutoMigrate(
    // ... 现有模型 ...
    &model.ContentAuditTask{},
    &model.AuditLog{},
    &model.UserRiskRecord{},
    // Report 已存在，GORM 会自动补充新字段 (type, handle_at)
    // Comment 已存在，GORM 会自动补充 status 字段
    // Post 已存在，pending 只是常量，无需迁移
)
```

---

## 3. Redis 初始化

```go
// config/config.go 中增加 Redis 配置项
type RedisConfig struct {
    Addr     string `yaml:"addr"`      // e.g. "localhost:6379"
    Password string `yaml:"password"`
    DB       int    `yaml:"db"`
}

// pkg/redis/client.go（新建）
package redis

import (
    "github.com/redis/go-redis/v9"
)

func NewClient(cfg config.RedisConfig) *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:     cfg.Addr,
        Password: cfg.Password,
        DB:       cfg.DB,
    })
}
```

在 `basic.yaml` 中添加：
```yaml
redis:
  addr: "localhost:6379"
  password: ""
  db: 0
```

---

## 4. 依赖注入（wire.go）

```go
// internal/wire/wire.go

// 初始化顺序
rdb         := redis.NewClient(cfg.Redis)
limiter     := ratelimit.NewLimiter(rdb)
sensitiveFilter := sensitive.NewFilter()
riskRepo    := riskrepo.NewRiskRepository(db)
riskSvc     := riskservice.NewRiskService(riskRepo, limiter)
checkSvc    := riskservice.NewContentCheckService(riskRepo, sensitiveFilter)
riskHandler := riskhandler.NewRiskHandler(checkSvc, riskSvc)
```

---

## 5. 路由接入

### 5a. 限流中间件（在现有路由上叠加）

```go
// cmd/server/main.go 或 router 初始化处

authorized := r.Group("/api")
authorized.Use(middleware.RequireAuth())
{
    // 发帖：加限流
    authorized.POST("/posts",
        middleware.RateLimitMiddleware(db, riskSvc, ratelimit.ActionCreatePost),
        middleware.ContentCheckMiddleware(checkSvc, []string{"title", "content"}),
        postHandler.CreatePost,
    )

    // 评论：加限流
    authorized.POST("/posts/:id/comments",
        middleware.RateLimitMiddleware(db, riskSvc, ratelimit.ActionCreateComment),
        middleware.ContentCheckMiddleware(checkSvc, []string{"content"}),
        commentHandler.CreateComment,
    )

    // 举报：加限流（防止刷举报）
    authorized.POST("/reports",
        middleware.RateLimitMiddleware(db, riskSvc, ratelimit.ActionSendReport),
        reportHandler.CreateReport,
    )
}

// 管理端审核队列
adminGroup := r.Group("/api/admin")
adminGroup.Use(middleware.RequireAuth(), middleware.RequirePermission(model.PermAdmin))
{
    riskHandler.RegisterRoutes(adminGroup)
}
```

### 5b. Handler 内读取审核标记

```go
// internal/handler/post/post_crud.go - CreatePost handler 中

func (h *PostHandler) CreatePost(c *gin.Context) {
    var input CreatePostInput
    if err := c.ShouldBindJSON(&input); err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    // 读取内容审核中间件注入的标记
    reviewRequired, hitWords := middleware.IsReviewRequired(c)

    status := model.PostStatusPublished
    if reviewRequired {
        status = model.PostStatusPending  // 命中 review 级词汇，进入待审核
    }

    post := &model.Post{
        Title:    input.Title,
        Content:  input.Content,
        AuthorID: c.GetUint(model.ContextUserID),
        Status:   status,
    }

    if err := h.postSvc.Create(post); err != nil {
        response.InternalError(c, err.Error())
        return
    }

    // 异步创建审核任务（不阻塞响应）
    if reviewRequired {
        go checkSvc.CreateAuditTaskForPost(post.ID, "sensitive_word", hitWords)
    }

    response.Success(c, post)
}
```

---

## 6. 举报创建时触发聚合检测

```go
// internal/service/report 或 handler/report - 举报创建后

func (s *ReportService) Create(report *model.Report) error {
    if err := s.repo.Create(report); err != nil {
        return err
    }

    // 异步检查是否触发聚合审核
    go func() {
        targetType := model.AuditTargetType(report.TargetType)
        triggered, _ := checkSvc.HandleReportAggregate(targetType, report.TargetID)
        if triggered {
            // 可选：推送通知给版主
            _ = notificationSvc.NotifyModerators("有内容触发举报聚合审核")
        }
    }()

    return nil
}
```

---

## 7. 管理员封禁用户时写入审计日志

```go
// internal/service/user/user_admin_status.go

func (s *UserService) BlockUser(operatorID, targetID uint, reason, ip string) error {
    var before, after string
    // ... 获取 before 状态，执行封禁 ...

    // 写审计日志
    _ = riskSvc.WriteAuditLog(
        operatorID,
        model.AuditActionBlockUser,
        "user", targetID,
        before, after,
        reason, ip,
    )
    return nil
}
```

---

## 8. 动态扩展词库（运营人员使用）

```go
// 通过管理 API 或初始化时从数据库加载额外词库
filter.AddWords(sensitive.LevelBlock, []string{"新增违禁词1", "新增违禁词2"})
filter.AddWords(sensitive.LevelReview, []string{"新增观察词1"})
```

---

## 架构总览

```
HTTP 请求
    │
    ├─ RequireAuth()                    # 现有：认证
    ├─ RateLimitMiddleware()            # 新增：频率限制（Redis 滑动窗口）
    ├─ ContentCheckMiddleware()         # 新增：敏感词 Pre-check（同步）
    │       │
    │       ├─ LevelBlock  → 400 直接拒绝
    │       └─ LevelReview → context 注入标记，放行
    │
    └─ Handler
            │
            ├─ 读取 IsReviewRequired() → status = pending
            ├─ 写库
            └─ go CreateAuditTask()    # 异步：创建审核任务


举报提交
    └─ go HandleReportAggregate()      # 异步：聚合检测 → 创建审核任务


管理员操作
    ├─ WriteAuditLog()                 # 每次操作记录
    └─ GET /admin/risk/audit/tasks     # 审核队列
```