# 快速开始

```bash
# 先创建数据库。
# 无需创建表， GORM 会自动创建和维护
psql -U postgres -h localhost -d postgres -c "CREATE DATABASE tiny_forum;"
```

## （非必要）初始化数据

> 如果你只是使用，你可以自行创建论坛，在数据可直接创建仅仅为了开发测试

```sql
-- 插入世界版块（提供所有必需字段）
INSERT INTO boards (
    name, 
    slug, 
    description, 
    created_at, 
    updated_at, 
    view_role, 
    post_role, 
    reply_role
) VALUES (
    '世界', 
    'world', 
    '世界板块，所有人都可以看到',
    NOW(), 
    NOW(),
    'user',
    'user',
    'user'
);
```

```sql
-- 插入公告版块
INSERT INTO boards (
    name, 
    slug, 
    description, 
    created_at, 
    updated_at,
    view_role, 
    post_role, 
    reply_role
) VALUES (
    '公告', 
    'announcement', 
    '官方公告和重要通知',
    NOW(), 
    NOW(),
    'admin',
    'admin',
    'admin'
);
```


## 配置说明

修改 `backend/config/config.yaml`：

```yaml
server:
  port: 8080
  mode: debug  # debug | release

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