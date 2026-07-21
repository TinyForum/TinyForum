# 机器人

目前 Tiny Forum 已支持机器人进行自动化功能。

机器人是一种身份。

## 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                     HTTP API 层                              │
│  POST /bots          创建（Lua脚本 / 零代码Flow）             │
│  POST /bots/:id/run  手动触发执行                            │
│  GET  /bots/nocode/metadata  零代码节点元数据（前端用）        │
│  POST /bots/nocode/validate  校验Flow配置                    │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│                  Service 层 (service.go)                     │
│  ┌──────────────────┐    ┌────────────────────────────────┐ │
│  │  Lua 脚本路径     │    │      零代码流程路径              │ │
│  │                  │    │                                │ │
│  │  bot.script_code │    │  bot.config_values["flow"]     │ │
│  │        ↓         │    │           ↓                    │ │
│  │   LuaSandbox     │    │      FlowEngine                │ │
│  │        ↓         │    │           ↓                    │ │
│  │    BotSDK        │    │  evalCondition → execAction    │ │
│  └────────┬─────────┘    └────────────┬───────────────────┘ │
│           │                           │                     │
│           └──────────┬────────────────┘                     │
│                      ↓                                      │
│               ForumAPIImpl                                  │
│        (PostRepo / UserRepo / MsgRepo ...)                  │
└─────────────────────────────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                  调度器 (StartScheduler)                      │
│    cron.Cron  ← TriggerSchedule (CronExpr)                  │
│    EventBus   ← TriggerEvent    (EventFilter)               │
│    手动触发   ← TriggerManual   (RunNow API)                 │
└─────────────────────────────────────────────────────────────┘
```

## 机器人两种模式对比

| 特性     | Lua 脚本机器人       | 零代码机器人                  |
| -------- | -------------------- | ----------------------------- |
| 创建方式 | 编写 Lua 代码        | 拖拽节点填写参数              |
| 灵活性   | 极高（完整编程能力） | 中等（内置节点组合）          |
| 上手门槛 | 需要编程基础         | 零门槛                        |
| 存储位置 | `bot.script_code`    | `bot.config_values["flow"]`   |
| 执行引擎 | LuaSandbox           | FlowEngine                    |
| 模板变量 | Lua 变量             | `{{event.field}}` Go template |

## Lua SDK 完整 API

### forum.*

```lua
-- 帖子
local post, err = forum.getPost(post_id)
local post, err = forum.createPost(title, content, section_id)
local comment, err = forum.replyPost(post_id, content)
local ok, err = forum.deletePost(post_id)
local ok, err = forum.moderatePost(post_id, action, reason)
-- action: "hide" | "pin" | "lock" | "delete"

-- 评论
local comment, err = forum.getComment(comment_id)
local ok, err = forum.deleteComment(comment_id)

-- 用户
local user, err = forum.getUser(user_id)
local ok, err = forum.banUser(user_id, reason, duration_sec)
local ok, err = forum.sendMessage(to_user_id, content)

-- 统计
local stats, err = forum.getStats()
-- stats.post_count / user_count / comment_count / active_today
```

### http.*

```lua
local body, status, err = http.get(url, headers_table?)
local body, status, err = http.post(url, body_str, headers_table?)
```

### util.*

```lua
local ts = util.now()                        -- Unix 时间戳
local s  = util.format_time(ts, layout)      -- 格式化时间
util.sleep(ms)                               -- 最多 5000ms
local b  = util.contains(str, sub)          -- 字符串包含
local t  = util.split(str, sep)             -- 字符串分割
local s  = util.trim(str)
local s  = util.lower(str)
local s  = util.upper(str)
```

### json.*

```lua
local str, err = json.encode(table)
local tbl, err = json.decode(str)
```

### 全局变量

```lua
config  -- bot.config_values（机器人配置，只读）
event   -- 触发事件数据（只读）
log(msg)          -- 记录日志
logf(fmt, ...)    -- 格式化日志
```

## 权限系统

机器人创建时声明 `permissions` 列表，SDK 调用时强制校验：

| 权限             | 允许的 SDK 调用                                              |
| ---------------- | ------------------------------------------------------------ |
| `read:posts`     | forum.getPost                                                |
| `write:posts`    | forum.createPost, forum.replyPost                            |
| `read:comments`  | forum.getComment                                             |
| `write:comments` | forum.deleteComment                                          |
| `read:user`      | forum.getUser                                                |
| `manage:content` | forum.deletePost, forum.moderatePost, forum.banUser, forum.deleteComment |
| `send:message`   | forum.sendMessage                                            |
| `read:stats`     | forum.getStats                                               |

## 零代码 Flow JSON 结构

```json
{
  "version": "1",
  "trigger": {
    "type": "on_keyword",
    "params": { "keywords": ["广告"], "scope": "both" }
  },
  "conditions": [
    {
      "type": "user_role_is",
      "negate": true,
      "params": { "role": "admin" }
    }
  ],
  "actions": [
    { "type": "hide_post", "params": {} },
    { "type": "send_message", "params": { "content": "警告：{{matched_kw}}" } },
    { "type": "ban_user", "params": { "reason": "发布广告", "duration_sec": 86400 } }
  ]
}
```

## 内置触发器

| 类型               | 说明       | 关键 event 字段                                       |
| ------------------ | ---------- | ----------------------------------------------------- |
| `on_schedule`      | cron 定时  | 无                                                    |
| `on_new_post`      | 新帖发布   | post_id, user_id, section_id, post_title              |
| `on_new_comment`   | 新评论     | comment_id, post_id, user_id, content                 |
| `on_user_register` | 新用户注册 | user_id, username, email                              |
| `on_keyword`       | 关键词命中 | matched_kw, content_type, post_id/comment_id, user_id |
| `on_manual`        | 手动触发   | 自定义                                                |

## 内置条件

| 类型                    | 说明             | 关键参数               |
| ----------------------- | ---------------- | ---------------------- |
| `post_title_contains`   | 标题含关键词     | keywords[]             |
| `post_content_contains` | 正文含关键词     | keywords[]             |
| `user_role_is`          | 用户角色匹配     | role                   |
| `user_post_count_gte`   | 用户发帖数 ≥ N   | count                  |
| `section_id_in`         | 所在板块         | ids[]                  |
| `time_range`            | 当前时间在区间内 | start, end, tz         |
| `custom_expr`           | 简单表达式       | expr (如 event.x > 10) |

## 内置动作

| 类型             | 分类        | 说明             |
| ---------------- | ----------- | ---------------- |
| `reply_post`     | post        | 回复帖子         |
| `delete_post`    | post        | 删除帖子         |
| `hide_post`      | post        | 隐藏帖子         |
| `pin_post`       | post        | 置顶帖子         |
| `lock_post`      | post        | 锁定帖子         |
| `create_post`    | post        | 发布新帖         |
| `delete_comment` | comment     | 删除评论         |
| `ban_user`       | user        | 封禁用户         |
| `warn_user`      | user        | 警告（私信）     |
| `send_message`   | user        | 发送私信         |
| `webhook`        | integration | 调用外部 Webhook |
| `notify_admin`   | integration | 通知管理员       |
| `wait`           | control     | 等待 N 秒        |
| `set_variable`   | control     | 设置流程变量     |
| `stop_if`        | control     | 条件提前结束     |

## 接入步骤

### 1. 实现 ForumAPI

```go
// 将真实的 repo 注入到 ForumAPIImpl
api := botapi.NewForumAPI(bot, systemBotUserID,
    postRepo, commentRepo, userRepo, msgRepo, statsRepo)
```

### 2. 初始化 Service

```go
botSvc := botservice.NewService(botRepo, api)
```

### 3. 在事件发生处发布事件

```go
// 新用户注册后
botSvc.PublishEvent("on_user_register", map[string]any{
    "user_id":  user.ID,
    "username": user.Username,
    "email":    user.Email,
})

// 新帖子发布后
botSvc.PublishEvent("on_new_post", map[string]any{
    "post_id":    post.ID,
    "user_id":    post.AuthorID,
    "section_id": post.SectionID,
    "post_title": post.Title,
})
```

### 4. 启动调度器

```go
botSvc.StartScheduler()
defer botSvc.StopScheduler()
```