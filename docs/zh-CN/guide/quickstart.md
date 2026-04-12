# 快速开始

```bash
# 先创建数据库。
# 无需创建表， GORM 会自动创建和维护
psql -U postgres -h localhost -d postgres -c "CREATE DATABASE tiny_forum;"
```


## 积分规则

| 行为 | 积分 |
|------|------|
| 注册 | 0 |
| 发帖 | +10 |
| 发表评论 | +3 |
| 点赞他人 | +2 |

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