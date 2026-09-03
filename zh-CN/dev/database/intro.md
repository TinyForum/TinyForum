# 数据库设计

---

## 数据库选型

TinyForum 使用 **PostgreSQL** 作为主数据库，通过 **GORM** ORM 框架进行数据访问。

| 组件 | 版本 | 用途 |
|------|------|------|
| PostgreSQL | 16+ | 关系型数据持久化 |
| GORM | v1.31+ | Go ORM 框架 |
| Redis | 7+ | 缓存、限流、排行榜、会话 |

---

## 配置文件

数据库连接配置位于 `backend/config/postgres.yml`：

```yaml
host: localhost
port: 5432
user: postgres
password: your_password
dbname: tinyforum
sslmode: disable
max_open_conns: 80
max_idle_conns: 20
conn_max_lifetime: 5m
```

对应 Go 结构体：`internal/infra/config/model.go` 中的 `ConfigPostgres`。

连接池参数说明：

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `max_open_conns` | 80 | 最大打开连接数（建议不超过 PG 的 `max_connections`） |
| `max_idle_conns` | 20 | 空闲连接池大小 |
| `conn_max_lifetime` | 5m | 连接最大存活时间 |
| `conn_max_idle_time` | 2m | 空闲连接超时时间 |

---

## GORM 使用规范

### 1. Context 传递

所有数据库操作必须传递 `context.Context`，用于超时控制和链路追踪：

```go
// ✅ 正确
func (r *userRepo) FindByID(ctx context.Context, id uint) (*do.User, error) {
    var user do.User
    err := r.db.WithContext(ctx).First(&user, id).Error
    return &user, err
}

// ❌ 错误：未使用 WithContext
func (r *userRepo) FindByID(id uint) (*do.User, error) {
    var user do.User
    err := r.db.First(&user, id).Error
    return &user, err
}
```

### 2. 参数化查询

GORM 默认使用参数化占位符，**严禁拼接 SQL**：

```go
// ✅ 正确：参数化
db.Where("username = ?", username).First(&user)

// ✅ 正确：LIKE 查询
db.Where("name LIKE ?", "%"+keyword+"%").Find(&users)

// ❌ 错误：字符串拼接（SQL 注入风险）
db.Where("username = '" + username + "'").First(&user)
```

### 3. 动态排序字段白名单

当需要动态 `ORDER BY` 时，必须使用白名单验证：

```go
var allowedSortFields = map[string]bool{
    "created_at": true,
    "updated_at": true,
    "like_count": true,
}

func (r *repo) List(sortBy, order string) ([]do.Post, error) {
    if !allowedSortFields[sortBy] {
        sortBy = "created_at"
    }
    if order != "asc" && order != "desc" {
        order = "desc"
    }
    // 安全：已通过白名单验证
    return r.db.Order(sortBy + " " + order).Find(&posts)
}
```

### 4. Preload 预加载

必须显式指定关联字段，禁止使用 `clause.Associations`：

```go
// ✅ 正确：显式指定
db.Preload("Creation.Author").Preload("Creation.Tags").Find(&posts)

// ❌ 错误：全量预加载（性能隐患）
db.Preload(clause.Associations).Find(&posts)
```

### 5. 分页查询

优先使用游标分页，避免深分页 OFFSET 性能问题：

```go
// ✅ 推荐：游标分页
db.Where("id < ?", cursorID).Order("id DESC").Limit(pageSize).Find(&list)

// ⚠️ 浅分页可接受（page < 100）
db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)
```

---

## 数据迁移规范

### 命名格式

```
backend/migrations/{timestamp}_{description}.up.sql
backend/migrations/{timestamp}_{description}.down.sql
```

示例：

```
20250101000001_add_user_avatar.up.sql
20250101000001_add_user_avatar.down.sql
```

### 迁移规则

1. **每个迁移必须提供 `up.sql` 和 `down.sql`**，确保可回滚。
2. **不可逆操作**（如 DROP COLUMN）必须先在 `up.sql` 中备份数据到临时表。
3. **兼容性铁律**：删除列或重命名必须分两步提交（先增后弃），确保新旧代码同时兼容。
4. **本地验证**：迁移脚本必须在 Docker Compose 启动的本地 PostgreSQL 上验证通过。

### 迁移脚本示例

```sql
-- 20250101_add_user_bio.up.sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS bio TEXT DEFAULT '';

-- 20250101_add_user_bio.down.sql  
ALTER TABLE users DROP COLUMN IF EXISTS bio;
```

### 当前状态

项目目前使用 GORM `AutoMigrate` 进行开发阶段的表结构同步。生产环境应逐步迁移到 SQL 迁移文件方式。

---

## 核心表结构

| 表名 | 对应 DO | 说明 |
|------|---------|------|
| `creations` | `Creation` | 内容基表（帖子/文章/问题共用） |
| `users` | `User` | 用户表 |
| `articles` | `Article` | 文章扩展表 |
| `posts` | `Post` | 帖子扩展表 |
| `questions` | `Question` | 问题扩展表 |
| `answers` | `Answer` | 回答表 |
| `comments` | `Comment` | 评论表 |
| `replies` | `Reply` | 回复表 |
| `boards` | `Board` | 板块表 |
| `topics` | `Topic` | 主题表 |
| `tags` | `Tag` | 标签表 |
| `notifications` | `Notification` | 通知表 |
| `follows` | `Follow` | 关注关系表 |
| `likes` | `Like` | 点赞记录表 |
| `attachments` | `Attachment` | 附件表 |
| `moderators` | `Moderator` | 版主表 |
| `bots` | `Bot` | 机器人表 |

完整表结构请查阅 `internal/model/do/` 目录下的 GORM 实体定义。

---

## Redis 数据结构

| Key 模式 | 类型 | 用途 |
|----------|------|------|
| `rl:{user_id}:{action}` | Sorted Set | 用户操作限流计数器 |
| `token:blacklist:{token_hash}` | String | JWT Token 黑名单 |
| `leaderboard:score` | Sorted Set | 用户积分排行榜 |
| `sensitive:dict:{category}` | Set | 敏感词词典缓存 |
