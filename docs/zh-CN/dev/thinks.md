# 设计思想

```bash
      HTTP Request
           │
           ▼
    ┌─────────────┐
    │   Handler   │  ← 解析请求参数
    └─────────────┘
           │
           │ 调用 Service 方法
           ▼
    ┌─────────────┐
    │   Service   │  ← 业务逻辑、事务
    └─────────────┘
           │
           │ 调用 Repository 方法
           ▼
    ┌─────────────┐
    │ Repository  │  ← 数据库操作
    └─────────────┘
           │
           ▼
       Database
```

Handler：收请求、发响应（网络请求解析响应）
Service：执行业务规则（业务）
Repository：存取数据（数据库）

三者形成单向依赖链，上层依赖下层接口，下层不感知上层。

依赖方向：Handler → Service → Repository
数据流向：Request → Handler → Service → Repository → DB

遵循依赖倒置原则。核心关系是：Handler 依赖 Service，Service 依赖 Repository。

## repository

1. 内部定义接口，依赖倒置（正在修改）。
2. 定义 Query 作为查询方法，给 Service 验证类型




## service
1. 定义 Service 接口，依赖倒置（正在修改）。
2. 定义 DTO 作为操作项目，给 Handler 验证类型。

## handler
1. 定义 Handler 接口，依赖倒置（正在修改）。
2. 定义 Request 和 Response 作为请求和响应，验证 Http 请求和响应类型。


## model
1. 定义 Model 实体
2. 定义白名

# pkg
1. fields，过滤字段，独立纯函数，可被多个 Service 复用。