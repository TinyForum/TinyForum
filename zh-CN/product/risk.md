# 行为风控

---

## 风控体系

TinyForum 通过多维度行为分析，识别和防范恶意用户行为。

### 风控维度

| 维度 | 指标 | 用途 |
|------|------|------|
| **用户频率** | 单位时间操作次数 | 检测刷帖、暴力破解 |
| **内容质量** | 帖子重复度、广告特征 | 检测垃圾内容 |
| **IP 信誉** | IP 关联账号数、历史违规记录 | 检测批量注册 |
| **行为序列** | 注册→发帖→私信的时序模式 | 检测营销号 |
| **社交图谱** | 关注关系、互动模式 | 检测刷粉、水军 |

---

## 限流控制

基于 Redis 的滑动窗口限流，按 **用户 + 接口** 维度计数。

### 限流配置

```yaml
# backend/config/risk_control.yml
rate_limit:
  enabled: true
  default:
    window: 60    # 窗口时长（秒）
    max: 100      # 最大请求数
  actions:
    register:
      window: 3600
      max: 3      # 每小时最多注册 3 次
    login:
      window: 300
      max: 10     # 每5分钟最多登录 10 次
    create_post:
      window: 60
      max: 5      # 每分钟最多发 5 帖
    create_comment:
      window: 60
      max: 20     # 每分钟最多发 20 条评论
```

### 限流中间件

```go
// 对登录接口限流
api.POST("/auth/login", mw.RateLimit("login"), h.Login)

// 对发帖接口限流
api.POST("/posts", mw.Auth(), mw.RateLimit("create_post"), h.Create)
```

### 限流响应

触发限流时返回：

```json
{
  "code": 429,
  "message": "请求过于频繁，请稍后再试",
  "retry_after": 45
}
```

---

## IP 风险评分

### 评分规则

| 行为 | 加分 | 说明 |
|------|------|------|
| 新 IP 注册 | +10 | 初始分 |
| 多账号同 IP | +20/账号 | 超出阈值后线性累加 |
| 被举报确认 | +30 | 举报属实 |
| 发布违规内容 | +50 | 被审核拒绝 |
| 短时间内大量操作 | +15 | 频率异常 |

### 风险等级

| 分数范围 | 等级 | 处理策略 |
|----------|------|----------|
| 0-20 | 低风险 | 正常放行 |
| 21-40 | 中风险 | 内容进入审核队列 |
| 41-60 | 高风险 | 限制发帖频率 |
| 61+ | 极高风险 | 临时封禁 IP |

---

## 用户风险评分

类似 IP 评分，但追踪到具体用户账号：

| 行为 | 加分 | 说明 |
|------|------|------|
| 新注册 | +5 | 初始分 |
| 违规内容 | +30 | 审核确认违规 |
| 被举报 | +10 | 每次被独立用户举报 |
| 频繁操作 | +20 | 触达频率阈值 |
| 申诉成功 | -20 | 降低评分 |

---

## 风控事件日志

所有风控事件记录在以下表中：

- `ip_risk_records`：IP 维度的风险记录
- `user_risk_records`：用户维度的风险记录
- `audit_logs`：审核操作的审计日志
- `blocked_ips`：被封禁的 IP 列表

### 日志查询

管理员可通过后台查看风控日志：

```
GET /api/v1/admin/risk/records?type=ip&risk_level=high
GET /api/v1/admin/risk/records?type=user&user_id=123
```

---

## 风控与内容安全联动

```
用户操作
  ├── 内容检查（敏感词 DFA + LLM）
  ├── 频率检查（Redis 限流计数器）
  ├── IP 信誉检查（IP 风险评分）
  └── 用户信誉检查（用户风险评分）
       └── 综合评分 → 放行 / 审核 / 拦截 / 封禁
```

---

## 配置说明

```yaml
# backend/config/risk_control.yml
risk_control:
  enabled: true
  ip_risk:
    max_accounts_per_ip: 5    # 同 IP 最大账号数
    score_short_url: 30       # 触发短期封禁的分数
    score_long_ban: 60        # 触发长期封禁的分数
  user_risk:
    score_audit: 20           # 触发内容审核的分数
    score_limit: 40           # 触发频率限制的分数
    score_ban: 60             # 触发封号的分数
```
