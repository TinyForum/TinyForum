# 用户认证与授权

---

## 认证体系

TinyForum 使用 **JWT (JSON Web Token)** 进行无状态用户认证。

### 认证流程

```
注册 → 邮箱验证 → 登录 → 获得 JWT → 后续请求携带 Token
```

### 注册

```
POST /api/v1/auth/register
```

- 必填：用户名、邮箱、密码
- 可选：昵称、头像
- 密码要求：6-32 位字符
- 注册后自动发送欢迎邮件（如已配置邮件服务）

### 登录

```
POST /api/v1/auth/login
```

- 支持用户名或邮箱登录
- 登录成功返回 JWT Token，存入 HttpOnly Cookie
- Token 有效期可配置（默认 7 天）
- 支持 Refresh Token 续期

### 登出

```
POST /api/v1/auth/logout
```

- 清除客户端 Cookie
- 将当前 Token 加入黑名单（Redis）

### 密码管理

| 功能 | 接口 | 说明 |
|------|------|------|
| 修改密码 | `PUT /api/v1/auth/password` | 需提供旧密码 |
| 忘记密码 | `POST /api/v1/auth/forgot-password` | 发送重置邮件 |
| 重置密码 | `POST /api/v1/auth/reset-password` | 通过邮件中的 Token 重置 |
| 注销账号 | `DELETE /api/v1/auth/account` | 永久删除账号及关联数据 |

---

## 授权体系

### 角色层级

```
super_admin > admin > moderator > reviewer > member > user > guest
```

| 角色 | 说明 | 权限范围 |
|------|------|----------|
| `super_admin` | 超级管理员 | 全局所有权限，唯一最高权限者 |
| `admin` | 管理员 | 全局管理（用户、帖子、板块、积分） |
| `moderator` | 版主 | 所在板块的内容管理、用户封禁 |
| `reviewer` | 审核员 | 内容审核权限 |
| `member` | 会员 | 基础功能 + 高积分权重 |
| `user` | 注册用户 | 发帖、评论、点赞、关注 |
| `guest` | 游客 | 仅浏览公开内容 |

### RBAC 权限控制

基于 Casbin 的细粒度权限控制：

- **路由级权限**：通过 `CasbinAuth` 中间件，根据用户角色 + 请求路径 + HTTP 方法决策
- **资源级权限**：版主仅能管理自己板块的内容
- **操作级权限**：Bot 机器人声明所需权限列表，执行时校验

### 中间件链

```
Auth → CasbinAuth → RateLimit → ContentCheck → ModeratorRequired
```

- `Auth`：解析 JWT，注入 `user_id` 和 `user_role`
- `OptionalAuth`：可选认证（游客模式，未登录注入 guest）
- `CasbinAuth`：路由级 RBAC 权限检查
- `RateLimit`：按用户/接口维度限流
- `ContentCheck`：敏感词过滤
- `ModeratorRequired`：版主角色的细粒度权限

### 管理员路由

管理员接口集中注册在 `/api/v1/admin/*` 下，需要 `admin` 或更高角色。

---

## 安全机制

- **JWT 密钥**：通过配置文件管理，启动时校验长度 >= 32 位
- **HttpOnly Cookie**：防止 XSS 窃取 Token
- **SameSite=Lax**：防止 CSRF 攻击
- **生产环境 Secure**：HTTPS 下自动启用 Secure 标志
- **密码加密**：bcrypt 哈希存储
- **操作日志**：管理员操作全量记录到 `audit_logs` 表
