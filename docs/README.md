# Tiny Forum — 全栈技术交流社区

> Go (Gin + GORM) 后端 × Next.js 14 (App Router) 前端 × PostgreSQL

---

## 技术栈

| 层     | 技术                                                                           |
| ------ | ------------------------------------------------------------------------------ |
| 后端   | Go 1.21, Gin, GORM, Wire (手动注入), JWT, Zap                                  |
| 前端   | Next.js 16, TypeScript, Tailwind CSS, DaisyUI, TanStack Query, Zustand, Tiptap |
| 数据库 | PostgreSQL 16                                                                  |
| 部署   | Docker + Docker Compose                                                        |

## 功能列表

- ✅ 用户注册 / 登录 / JWT 鉴权
- ✅ 发帖（帖子 / 文章 / 话题）、富文本编辑器
- ✅ 评论 & 嵌套回复
- ✅ 点赞 / 取消点赞
- ✅ 标签系统
- ✅ 关注 / 取消关注
- ✅ 积分系统 & 排行榜
- ✅ 站内消息通知
- ✅ 个人主页 / 编辑资料
- ✅ 管理后台（用户管理、封禁、置顶）
- ✅ 全文搜索（标题 & 内容）
- ✅ 深色 / 浅色主题切换（DaisyUI）

---

## 快速启动

### 方式一：Docker Compose（推荐，一键启动）

```bash
# 克隆后直接运行
docker compose up -d

# 访问
# 前端：http://localhost:3000
# 后端 API：http://localhost:8080/api/v1
```

### 方式二：本地开发

#### 前置要求

- Go 1.21+
- Node.js 20+
- PostgreSQL 16（本地或 Docker）

#### 初始化数据库

```bash
# 先创建数据库，然后执行 SQL 文件
psql -U postgres -h localhost -d postgres -c "CREATE DATABASE tiny_forum;"
```

#### 启动 PostgreSQL（Docker 单独启动）

```bash
docker run -d \
  --name bbs_postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=tiny_forum \
  -p 5432:5432 \
  postgres:16-alpine
```

#### 启动后端

```bash
cd backend

# 安装依赖
go mod tidy

# 修改配置（如有需要）
# vim config/config.yaml

# 运行
go run ./cmd/server
```

> 首次启动会自动 AutoMigrate 建表，无需手动执行 SQL。

#### 启动前端

```bash
cd frontend

# 安装依赖
npm install

# 配置 API 地址（默认 localhost:8080）
# vim .env.local

# 开发模式
npm run dev
```

访问 http://localhost:3000

---

## 项目结构

```
tiny-forum/
├── backend/
│   ├── cmd/server/main.go          # 入口
│   ├── config/
│   │   ├── config.go               # 配置结构
│   │   └── config.yaml             # 配置文件
│   ├── internal/
│   │   ├── handler/                # HTTP 处理层
│   │   │   ├── auth.go
│   │   │   ├── post.go
│   │   │   ├── comment.go
│   │   │   ├── user.go
│   │   │   ├── tag.go
│   │   │   └── notification.go
│   │   ├── service/                # 业务逻辑层
│   │   │   ├── user.go
│   │   │   ├── post.go
│   │   │   ├── comment.go
│   │   │   ├── tag.go
│   │   │   └── notification.go
│   │   ├── repository/             # 数据访问层
│   │   │   ├── user.go
│   │   │   ├── post.go
│   │   │   ├── comment.go
│   │   │   ├── tag.go
│   │   │   └── notification.go
│   │   ├── model/                  # GORM 数据模型
│   │   │   └── model.go
│   │   ├── middleware/             # Gin 中间件
│   │   │   └── auth.go
│   │   └── wire/                   # 依赖注入 & 路由
│   │       └── wire.go
│   ├── pkg/
│   │   ├── jwt/                    # JWT 工具
│   │   ├── logger/                 # Zap 日志
│   │   └── response/               # 统一响应
│   ├── Dockerfile
│   ├── Makefile
│   └── go.mod
│
├── frontend/
│   └── src/
│       ├── app/                    # Next.js App Router
│       │   ├── page.tsx            # 首页
│       │   ├── auth/login/         # 登录
│       │   ├── auth/register/      # 注册
│       │   ├── posts/              # 帖子列表
│       │   ├── posts/new/          # 发帖
│       │   ├── posts/[id]/         # 帖子详情
│       │   ├── posts/[id]/edit/    # 编辑帖子
│       │   ├── users/[id]/         # 用户主页
│       │   ├── notifications/      # 通知
│       │   ├── leaderboard/        # 排行榜
│       │   ├── settings/           # 个人设置
│       │   └── admin/              # 管理后台
│       ├── components/
│       │   ├── layout/             # Navbar, Providers
│       │   └── post/               # PostCard, CommentSection, RichEditor
│       ├── lib/
│       │   ├── api-client.ts       # Axios 实例
│       │   ├── api.ts              # API 函数
│       │   └── utils.ts            # 工具函数
│       ├── store/
│       │   └── auth.ts             # Zustand auth store
│       └── types/
│           └── index.ts            # TypeScript 类型
│
├── docker-compose.yml
└── README.md
```

## API 文档

### 认证

| 方法 | 路径                  | 说明         |
| ---- | --------------------- | ------------ |
| POST | /api/v1/auth/register | 注册         |
| POST | /api/v1/auth/login    | 登录         |
| GET  | /api/v1/auth/me       | 当前用户信息 |

### 帖子

| 方法   | 路径                   | 说明                                   |
| ------ | ---------------------- | -------------------------------------- |
| GET    | /api/v1/posts          | 列表（支持分页、搜索、排序、标签过滤） |
| GET    | /api/v1/posts/:id      | 详情                                   |
| POST   | /api/v1/posts          | 发布                                   |
| PUT    | /api/v1/posts/:id      | 编辑                                   |
| DELETE | /api/v1/posts/:id      | 删除                                   |
| POST   | /api/v1/posts/:id/like | 点赞                                   |
| DELETE | /api/v1/posts/:id/like | 取消点赞                               |

### 评论

| 方法   | 路径                           | 说明          |
| ------ | ------------------------------ | ------------- |
| GET    | /api/v1/comments/post/:post_id | 帖子评论列表  |
| POST   | /api/v1/comments               | 发表评论/回复 |
| DELETE | /api/v1/comments/:id           | 删除评论      |

### 用户

| 方法   | 路径                      | 说明     |
| ------ | ------------------------- | -------- |
| GET    | /api/v1/users/:id         | 用户主页 |
| PUT    | /api/v1/users/profile     | 更新资料 |
| POST   | /api/v1/users/:id/follow  | 关注     |
| DELETE | /api/v1/users/:id/follow  | 取消关注 |
| GET    | /api/v1/users/leaderboard | 积分排行 |

### 标签

| 方法 | 路径         | 说明               |
| ---- | ------------ | ------------------ |
| GET  | /api/v1/tags | 所有标签           |
| POST | /api/v1/tags | 创建标签（管理员） |

### 通知

| 方法 | 路径                               | 说明     |
| ---- | ---------------------------------- | -------- |
| GET  | /api/v1/notifications              | 通知列表 |
| GET  | /api/v1/notifications/unread-count | 未读数量 |
| POST | /api/v1/notifications/read-all     | 全部已读 |

## 积分规则

| 行为     | 积分 |
| -------- | ---- |
| 注册     | 0    |
| 发帖     | +10  |
| 发表评论 | +3   |
| 点赞他人 | +2   |

## 配置说明

修改 `backend/config/config.yaml`：

```yaml
server:
  port: 8080
  mode: debug # debug | release

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: tiny_forum

jwt:
  secret: "your-secret-key-at-least-32-chars"
  expire: 72h
```
