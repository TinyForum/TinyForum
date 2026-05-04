

# 概述

## 1. 说明

该数据库设计用于支持一个**社区论坛/内容管理系统**，具备以下核心功能：

- **用户管理**：注册、登录、角色权限（普通用户/管理员/版主等）
- **版块管理**：支持多级版块（Board），可嵌套父子关系
- **内容发布**：支持多种帖子类型（文章、随笔、问答），含标签系统
- **评论互动**：支持嵌套评论、点赞、投票（赞同/反对）
- **问答功能**：问题可采纳答案，支持积分奖励
- **关注系统**：用户间关注、话题关注
- **通知系统**：各类事件通知
- **举报系统**：用户举报违规内容
- **签到系统**：连续签到奖励


## 2. 通用基础模型

所有实体均继承 `BaseModel`，提供统一的时间戳和软删除能力。

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键，自增 |
| CreatedAt | time.Time | 创建时间 |
| UpdatedAt | time.Time | 更新时间 |
| DeletedAt | gorm.DeletedAt | 软删除时间（带索引） |

---

## 3. 表结构详情

### 3.1 `users` - 用户表

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| username | varchar(50) | NOT NULL, UNIQUE | 用户名 |
| email | varchar(100) | NOT NULL, UNIQUE | 邮箱 |
| password | varchar | NOT NULL | 加密密码（JSON 不可见） |
| avatar | varchar(500) | | 头像 URL |
| bio | varchar(500) | | 个人简介 |
| role | varchar(20) | DEFAULT 'user' | 角色：user / admin / moderator |
| score | int | DEFAULT 0 | 积分 |
| is_active | bool | DEFAULT true | 是否激活 |
| is_blocked | bool | DEFAULT false | 是否被封禁 |
| last_login | timestamp | | 最后登录时间 |

**关系**：
- 一对多：`posts`（作者）、`comments`（作者）
- 多对多：`followers` / `following`（关注关系）

---

### 3.2 `boards` - 版块表（支持多级嵌套）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| name | varchar(50) | NOT NULL, UNIQUE | 版块名称 |
| slug | varchar(50) | NOT NULL, UNIQUE | URL 友好标识 |
| description | varchar(500) | | 描述 |
| icon | varchar(100) | | 图标标识 |
| cover | varchar(500) | | 封面图 |
| parent_id | uint | INDEX, NULL | 父级版块 ID（支持嵌套） |
| sort_order | int | DEFAULT 0 | 排序顺序 |
| view_role | varchar(20) | DEFAULT 'user' | 浏览权限角色 |
| post_role | varchar(20) | DEFAULT 'user' | 发帖权限角色 |
| reply_role | varchar(20) | DEFAULT 'user' | 回复权限角色 |
| post_count | int | DEFAULT 0 | 帖子总数 |
| thread_count | int | DEFAULT 0 | 主题数 |
| today_count | int | DEFAULT 0 | 今日新增 |

**关系**：
- 自关联：`parent` / `children`（树形结构）
- 一对多：`moderators`（版主）
- 一对多：`posts`

---

### 3.3 `posts` - 帖子表

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| title | varchar(200) | NOT NULL | 标题 |
| content | text | NOT NULL | 内容 |
| summary | varchar(500) | | 摘要 |
| cover | varchar(500) | | 封面 |
| type | varchar(20) | DEFAULT 'post' | 类型：post / essay / question / article |
| status | varchar(20) | DEFAULT 'published' | 状态：published / draft / deleted |
| author_id | uint | NOT NULL, INDEX | 作者 ID |
| view_count | int | DEFAULT 0 | 浏览量 |
| like_count | int | DEFAULT 0 | 点赞数 |
| pin_top | bool | DEFAULT false | 全局置顶 |
| board_id | uint | INDEX | 所属版块 ID |
| pin_in_board | bool | DEFAULT false | 版块内置顶 |

**关系**：
- 多对多：`tags`（通过 `post_tags` 中间表）
- 一对多：`comments`
- 一对多：`likes`
- 一对一：`question`（仅当 type='question' 时）

---

### 3.4 `comments` - 评论表

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| content | text | NOT NULL | 评论内容 |
| post_id | uint | NOT NULL, INDEX | 所属帖子 ID |
| author_id | uint | NOT NULL, INDEX | 作者 ID |
| parent_id | uint | INDEX, NULL | 父评论 ID（支持嵌套） |
| like_count | int | DEFAULT 0 | 点赞数 |
| is_answer | bool | DEFAULT false | 是否为答案（问答帖专用） |
| is_accepted | bool | DEFAULT false | 是否为被采纳答案 |
| vote_count | int | DEFAULT 0 | 投票净得分（赞同-反对） |

**关系**：
- 自关联：`parent` / `replies`（嵌套评论树）
- 一对多：`likes`
- 多对一：`post`、`author`

---

### 3.5 `questions` - 问答扩展表

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| post_id | uint | UNIQUE, NOT NULL | 关联的帖子 ID |
| accepted_answer_id | uint | NULL | 被采纳的评论 ID |
| reward_score | int | DEFAULT 0 | 悬赏积分 |
| answer_count | int | DEFAULT 0 | 回答数量 |
| view_count | int | DEFAULT 0 | 浏览量 |

**关系**：
- 一对一：`post`
- 一对一：`accepted_answer`（关联到 comment）

---

### 3.6 `tags` - 标签表

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| name | varchar(50) | UNIQUE, NOT NULL | 标签名称 |
| description | varchar(200) | | 描述 |
| color | varchar(20) | DEFAULT '#6366f1' | 标签颜色（HEX） |
| post_count | int | DEFAULT 0 | 使用该标签的帖子数 |

**关系**：
- 多对多：`posts`（通过 `post_tags`）

---

### 3.7 `likes` - 点赞表

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| user_id | uint | NOT NULL, INDEX | 点赞用户 ID |
| post_id | uint | INDEX, NULL | 帖子 ID（与 comment_id 二选一） |
| comment_id | uint | INDEX, NULL | 评论 ID（与 post_id 二选一） |

**约束**：`post_id` 和 `comment_id` 不能同时为 NULL，也不能同时非 NULL。

---

### 3.8 `votes` - 投票表（问答赞同/反对）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| user_id | uint | INDEX, UNIQUE(idx_user_comment) | 用户 ID |
| comment_id | uint | INDEX, UNIQUE(idx_user_comment) | 评论 ID |
| value | tinyint | | 1：赞同，-1：反对 |

---

### 3.9 `follows` - 用户关注表

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| follower_id | uint | NOT NULL, INDEX | 关注者 ID |
| following_id | uint | NOT NULL, INDEX | 被关注者 ID |

---

### 3.10 `moderators` - 版主表

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| user_id | uint | NOT NULL, UNIQUE(idx_user_board) | 用户 ID |
| board_id | uint | NOT NULL, UNIQUE(idx_user_board) | 版块 ID |
| permissions | json | | 扩展权限（JSON） |
| can_delete_post | bool | DEFAULT false | 可删除帖子 |
| can_pin_post | bool | DEFAULT false | 可置顶帖子 |
| can_edit_any_post | bool | DEFAULT false | 可编辑任何帖子 |
| can_manage_moderator | bool | DEFAULT false | 可管理其他版主 |
| can_ban_user | bool | DEFAULT false | 可封禁用户 |

---

### 3.11 `notifications` - 通知表

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| user_id | uint | NOT NULL, INDEX | 接收用户 ID |
| sender_id | uint | INDEX, NULL | 触发用户 ID |
| type | varchar(30) | | 通知类型（点赞/评论/关注等） |
| content | varchar(500) | | 通知内容 |
| target_id | uint | NULL | 关联目标 ID |
| target_type | varchar(50) | | 目标类型（post/comment/user） |
| is_read | bool | DEFAULT false | 是否已读 |

---

### 3.12 `reports` - 举报表

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| reporter_id | uint | NOT NULL, INDEX | 举报人 ID |
| target_id | uint | NOT NULL | 被举报目标 ID |
| target_type | varchar(50) | NOT NULL | 目标类型（post/comment/user） |
| reason | varchar(500) | NOT NULL | 举报原因 |
| status | varchar(20) | DEFAULT 'pending' | 状态：pending / resolved / rejected |
| handler_id | uint | NULL | 处理人 ID |
| handle_note | varchar(500) | | 处理备注 |

---

### 3.13 `topics` / `topic_posts` / `topic_follows` - 话题系统

#### topics 话题表

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| title | varchar(100) | NOT NULL | 话题标题 |
| description | varchar(500) | | 描述 |
| cover | varchar(500) | | 封面 |
| creator_id | uint | NOT NULL, INDEX | 创建者 ID |
| is_public | bool | DEFAULT true | 是否公开 |
| post_count | int | DEFAULT 0 | 包含帖子数 |
| follower_count | int | DEFAULT 0 | 关注者数量 |

#### topic_posts 话题-帖子关联表（多对多）

| 字段 | 说明 |
|------|------|
| topic_id | 话题 ID |
| post_id | 帖子 ID |

#### topic_follows 话题关注表（多对多）

| 字段 | 说明 |
|------|------|
| topic_id | 话题 ID |
| user_id | 用户 ID |

---

### 3.14 `sign_ins` - 签到记录表

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| user_id | uint | NOT NULL, INDEX | 用户 ID |
| sign_date | timestamp | NOT NULL | 签到日期 |
| score | int | DEFAULT 5 | 获得积分 |
| continued | int | DEFAULT 1 | 连续签到天数 |

---

## 4. 枚举类型说明

| 枚举名 | 可选值 | 说明 |
|--------|--------|------|
| UserRole | user, admin, moderator | 用户角色 |
| PostType | post, essay, question, article | 帖子类型 |
| PostStatus | published, draft, deleted | 帖子状态 |
| NotificationType | 由业务定义 | 通知类型 |
| ReportStatus | pending, resolved, rejected | 举报状态 |

---

## 5. 索引建议

| 表 | 索引字段 | 类型 | 说明 |
|----|----------|------|------|
| users | username, email | UNIQUE | 快速登录查询 |
| posts | author_id, board_id, created_at | INDEX | 列表查询 |
| comments | post_id, author_id, parent_id | INDEX | 评论树查询 |
| notifications | user_id, is_read, created_at | INDEX | 通知列表 |
| reports | status, target_type | INDEX | 待处理举报 |
| sign_ins | user_id, sign_date | INDEX | 每日签到检查 |

---

## 6. 数据一致性约束

1. **点赞唯一性**：同一用户对同一帖子/评论只能点赞一次
2. **投票唯一性**：同一用户对同一评论只能投票一次（赞同或反对）
3. **关注唯一性**：同一对关注关系不可重复
4. **版主唯一性**：同一用户在一个版块只能有一个版主记录
5. **问答采纳**：一个问答帖只能有一个被采纳的答案
6. **帖子与问答**：只有 type='question' 的帖子才关联 question 记录