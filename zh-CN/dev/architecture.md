# 架构说明

## 后端

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

---

---

---

# 前端

```bash
      HTTP Response
           ↑
           │
    ┌─────────────┐
    │   View      │  ← 渲染 UI，处理用户交互
    └─────────────┘
           ↑
           │ 订阅数据变化
    ┌─────────────┐
    │   State     │  ← 管理组件状态（本地 + 全局）
    └─────────────┘
           ↑
           │ 调用 Query/Mutation
    ┌─────────────┐
    │   Query     │  ← 数据获取、缓存、同步（React Query/SWR）
    └─────────────┘
           ↑
           │ 调用 API 客户端
    ┌─────────────┐
    │   Client    │  ← HTTP 请求封装（拦截器、错误处理）
    └─────────────┘
           ↑
           │
      HTTP Request
```

**核心分层职责：**

| 层级       | 职责                               | 不做什么                       |
| ---------- | ---------------------------------- | ------------------------------ |
| **View**   | 渲染 UI、绑定事件、组合组件        | 不包含业务逻辑、不直接调用 API |
| **State**  | 管理全局/共享状态（用户、UI 状态） | 不管理服务端数据               |
| **Query**  | 服务端数据获取、缓存、后台同步     | 不处理 UI 渲染                 |
| **Client** | HTTP 请求封装、类型安全            | 不包含业务逻辑                 |

**依赖方向：** View → State/Query → Client

**数据流向：** User Action → View → Query/State → Client → 后端

## 目录结构映射

```bash
src/
├── app/                    # 路由层（Next.js App Router）
│   └── [locale]/           # 国际化路由
│       ├── admin/          # 后台管理路由
│       ├── dashboard/      # 用户面板路由
│       └── forum/          # 论坛路由
│
├── features/               # 功能模块（按领域拆分）
│   ├── admin/              # 后台管理功能
│   │   ├── components/     # View 层：模块专用 UI 组件
│   │   ├── hooks/          # View + Query 层：组合 Query 和 State
│   │   ├── services/       # Service 层：复杂业务逻辑
│   │   └── types/          # Model 层：模块类型定义
│   ├── auth/               # 认证功能
│   ├── forum/              # 论坛功能
│   └── moderation/         # 审核功能
│
├── shared/                 # 共享资源
│   ├── ui/                 # View 层：通用 UI 组件
│   ├── lib/                # 工具函数（纯函数）
│   ├── api/                # Client 层：HTTP 客户端
│   ├── query/              # Query 层：数据获取（React Query）
│   ├── store/              # State 层：全局状态（Zustand）
│   └── types/              # Model 层：全局类型
│
├── layouts/                # View 层：布局组件
└── config/                 # 配置
```

**对应关系：**

| 后端       | 前端   | 职责                           |
| ---------- | ------ | ------------------------------ |
| Handler    | View   | 接收请求/用户操作，返回响应/UI |
| Service    | State  | 业务逻辑、状态管理             |
| Repository | Client | 数据存取（HTTP 请求）          |

**依赖方向：** View → State → Client

**数据流向：** User Action → View → State → Client → 后端

---

## View

1. 渲染 UI，处理用户交互
2. 薄层，不做业务逻辑
3. 页面即 View

## State

1. 管理数据状态（useState + Zustand + React Query）
2. 对应后端的 Service，处理业务逻辑
3. Hook 封装状态逻辑

## Client

1. 封装 HTTP 请求
2. 对应后端的 Repository，负责数据存取
3. 统一错误处理、类型安全

## Model

1. 定义类型（对应后端的 Model）
2. 前后端共享类型定义
