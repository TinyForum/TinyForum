# Tiny Forum API
一个基于 Gin 的论坛系统 API

## Version: 1.0

### Security

**ApiKeyAuth**  

| apiKey | *API Key* |
| ------ | --------- |
| Description | Type "Bearer" followed by a space and the JWT token. |
| Name | Authorization |
| In | header |

---
### /admin/posts

#### GET
##### Summary

管理员获取帖子列表

##### Description

管理员分页获取所有帖子列表，支持关键词搜索

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| page | query | 页码 | No | integer |
| page_size | query | 每页数量 | No | integer |
| keyword | query | 搜索关键词 | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [response.PageData](#responsepagedata-model) & { **"list"**: [ [model.Post](#modelpost-model) ] } } |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /admin/posts/{id}/pin

#### PUT
##### Summary

切换帖子置顶状态

##### Description

管理员切换指定帖子的置顶状态

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 帖子ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 操作成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的帖子ID | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /admin/users

#### GET
##### Summary

管理员获取用户列表

##### Description

管理员分页获取所有用户列表，支持关键词搜索

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| page | query | 页码 | No | integer |
| page_size | query | 每页数量 | No | integer |
| keyword | query | 搜索关键词（用户名、邮箱） | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [response.PageData](#responsepagedata-model) & { **"list"**: [ [model.User](#modeluser-model) ] } } |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /admin/users/{id}/active

#### PUT
##### Summary

管理员设置用户状态

##### Description

管理员启用或禁用指定用户账号

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 用户ID | Yes | integer |
| body | body | 状态信息 | Yes | [handler.SetUserActiveRequest](#handlersetuseractiverequest-model) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 操作成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的用户ID或请求参数错误 | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |
| 404 | 用户不存在 | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /admin/users/{id}/role

#### PUT
##### Summary

管理员设置用户角色

##### Description

管理员设置指定用户的角色（user/moderator/admin）

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 用户ID | Yes | integer |
| body | body | 角色信息 | Yes | [handler.SetUserRoleRequest](#handlersetuserrolerequest-model) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 操作成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的用户ID或角色 | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |
| 404 | 用户不存在 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

---
### /auth/login

#### POST
##### Summary

用户登录

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| body | body | 登录信息 | Yes | [service.LoginInput](#servicelogininput-model) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [response.Response](#responseresponse-model) |

### /auth/me

#### GET
##### Summary

获取当前用户信息

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /auth/register

#### POST
##### Summary

用户注册

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| body | body | 注册信息 | Yes | [service.RegisterInput](#serviceregisterinput-model) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [response.Response](#responseresponse-model) |

---
### /boards

#### GET
##### Summary

获取板块列表

##### Description

分页获取板块列表

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| page | query | 页码 | No | integer |
| page_size | query | 每页数量 | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [response.PageData](#responsepagedata-model) & { **"list"**: [ [model.Board](#modelboard-model) ] } } |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

#### POST
##### Summary

创建板块（仅管理员）

##### Description

创建一个新的板块，需要管理员权限

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| body | body | 板块信息 | Yes | [service.CreateBoardInput](#servicecreateboardinput-model) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 创建成功 | [response.Response](#responseresponse-model) & { **"data"**: [model.Board](#modelboard-model) } |
| 400 | 请求参数错误 | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /boards/slug/{slug}

#### GET
##### Summary

根据Slug获取板块

##### Description

根据板块标识符（slug）获取板块信息

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| slug | path | 板块标识符 | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [model.Board](#modelboard-model) } |
| 404 | 板块不存在 | [response.Response](#responseresponse-model) |

### /boards/tree

#### GET
##### Summary

获取板块树形结构

##### Description

获取所有板块的树形层级结构

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [ [model.BoardTree](#modelboardtree-model) ] } |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

### /boards/{id}

#### GET
##### Summary

获取板块详情

##### Description

根据ID获取板块详细信息

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 板块ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [model.Board](#modelboard-model) } |
| 400 | 无效的板块ID | [response.Response](#responseresponse-model) |
| 404 | 板块不存在 | [response.Response](#responseresponse-model) |

#### PUT
##### Summary

更新板块（仅管理员）

##### Description

更新指定板块的信息，需要管理员权限

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 板块ID | Yes | integer |
| body | body | 板块信息 | Yes | [service.CreateBoardInput](#servicecreateboardinput-model) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 更新成功 | [response.Response](#responseresponse-model) & { **"data"**: [model.Board](#modelboard-model) } |
| 400 | 请求参数错误 | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |
| 404 | 板块不存在 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

#### DELETE
##### Summary

删除板块（仅管理员）

##### Description

删除指定板块，需要管理员权限

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 板块ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 删除成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的板块ID | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |
| 404 | 板块不存在 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /boards/{id}/posts

#### GET
##### Summary

获取板块下的帖子列表

##### Description

分页获取指定板块下的所有帖子

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 板块ID | Yes | integer |
| page | query | 页码 | No | integer |
| page_size | query | 每页数量 | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [response.PageData](#responsepagedata-model) & { **"list"**: [ [model.Post](#modelpost-model) ] } } |
| 400 | 无效的板块ID | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

---
### /boards/{id}/bans

#### POST
##### Summary

禁言用户

##### Description

在指定板块禁言用户，需要版主权限

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 板块ID | Yes | integer |
| body | body | 禁言信息 | Yes | [handler.BanUserRequest](#handlerbanuserrequest-model) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 禁言成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 请求参数错误 | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /boards/{id}/bans/{user_id}

#### DELETE
##### Summary

解除禁言

##### Description

解除用户在指定板块的禁言，需要版主权限

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 板块ID | Yes | integer |
| user_id | path | 用户ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 解除成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的ID | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

---
### /boards/{id}/moderators

#### GET
##### Summary

获取板块版主列表

##### Description

获取指定板块的所有版主信息

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 板块ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [ [model.Moderator](#modelmoderator-model) ] } |
| 400 | 无效的板块ID | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

#### POST
##### Summary

添加版主

##### Description

为指定板块添加版主，需要版主管理权限

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 板块ID | Yes | integer |
| body | body | 版主信息 | Yes | [handler.AddModeratorRequest](#handleraddmoderatorrequest-model) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 添加成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 请求参数错误 | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /boards/{id}/moderators/{user_id}

#### DELETE
##### Summary

移除版主

##### Description

移除指定板块的版主，需要版主管理权限

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 板块ID | Yes | integer |
| user_id | path | 用户ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 移除成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的ID | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

---
### /comments

#### POST
##### Summary

创建评论

##### Description

创建一条新的评论

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| body | body | 评论信息 | Yes | [service.CreateCommentInput](#servicecreatecommentinput-model) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 创建成功 | [response.Response](#responseresponse-model) & { **"data"**: [model.Comment](#modelcomment-model) } |
| 400 | 请求参数错误 | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /comments/post/{post_id}

#### GET
##### Summary

获取帖子的评论列表

##### Description

分页获取指定帖子的所有评论

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| post_id | path | 帖子ID | Yes | integer |
| page | query | 页码 | No | integer |
| page_size | query | 每页数量 | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [response.PageData](#responsepagedata-model) & { **"list"**: [ [model.Comment](#modelcomment-model) ] } } |
| 400 | 无效的帖子ID | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

### /comments/{id}

#### DELETE
##### Summary

删除评论

##### Description

删除指定的评论（用户可以删除自己的评论，管理员可以删除任何评论）

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 评论ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 删除成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的评论ID或删除失败 | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |
| 404 | 评论不存在 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

---
### /comments/post/{post_id}/answers

#### GET
##### Summary

获取帖子的答案列表

##### Description

获取指定帖子的所有答案（仅限问答帖）

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| post_id | path | 帖子ID | Yes | integer |
| page | query | 页码 | No | integer |
| page_size | query | 每页数量 | No | integer |
| sort | query | 排序方式 | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [response.PageData](#responsepagedata-model) & { **"list"**: [ [model.Comment](#modelcomment-model) ] } } |
| 400 | 无效的帖子ID | [response.Response](#responseresponse-model) |

### /comments/{id}/accept

#### POST
##### Summary

采纳答案

##### Description

采纳某个回答作为最佳答案（仅帖子作者可操作）

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 评论ID | Yes | integer |
| post_id | query | 帖子ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 采纳成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 请求参数错误 | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /comments/{id}/answer

#### PUT
##### Summary

标记为答案

##### Description

将评论标记为问题的答案（帖子作者或版主可操作）

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 评论ID | Yes | integer |
| body | body | 是否标记为答案 | Yes | object |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 操作成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 请求参数错误 | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /comments/{id}/vote

#### GET
##### Summary

获取答案投票状态

##### Description

获取当前用户对指定答案的投票状态

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 评论ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的评论ID | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

#### POST
##### Summary

对答案投票

##### Description

对问答帖的答案进行投票（赞成up/反对down），重复投票会取消

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 评论ID | Yes | integer |
| body | body | 投票类型 | Yes | object |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 投票成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 请求参数错误 | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 不能给自己的答案投票 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /questions/answer/{comment_id}/vote

#### POST
##### Summary

投票回答

##### Description

对问题的回答进行投票（赞同或反对）

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| comment_id | path | 回答评论ID | Yes | integer |
| body | body | 投票信息 | Yes | [handler.VoteAnswerRequest](#handlervoteanswerrequest-model) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 投票成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的评论ID或投票类型 | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 404 | 回答不存在 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /questions/{post_id}/answer/{comment_id}/accept

#### POST
##### Summary

采纳答案

##### Description

采纳某个回答作为问题的正确答案（仅问题作者可操作）

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| post_id | path | 问题帖子ID | Yes | integer |
| comment_id | path | 回答评论ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 采纳成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的ID或操作失败 | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限（非问题作者） | [response.Response](#responseresponse-model) |
| 404 | 问题或回答不存在 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /questions/{post_id}/answers

#### GET
##### Summary

获取问题的回答列表

##### Description

分页获取指定问题的所有回答，已采纳的回答会排在前面

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| post_id | path | 问题帖子ID | Yes | integer |
| page | query | 页码 | No | integer |
| page_size | query | 每页数量 | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的帖子ID | [response.Response](#responseresponse-model) |
| 404 | 问题不存在 | [response.Response](#responseresponse-model) |

---
### /notifications

#### GET
##### Summary

获取通知列表

##### Description

分页获取当前用户的通知列表

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| page | query | 页码 | No | integer |
| page_size | query | 每页数量 | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [response.PageData](#responsepagedata-model) & { **"list"**: [ [model.Notification](#modelnotification-model) ] } } |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /notifications/read-all

#### POST
##### Summary

标记所有通知为已读

##### Description

将当前用户的所有未读通知标记为已读

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 标记成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /notifications/unread-count

#### GET
##### Summary

获取未读通知数量

##### Description

获取当前用户未读通知的总数

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

---
### /posts

#### GET
##### Summary

获取帖子列表

##### Description

分页获取帖子列表，支持多种筛选和排序

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| page | query | 页码 | No | integer |
| page_size | query | 每页数量 | No | integer |
| keyword | query | 搜索关键词 | No | string |
| sort_by | query | 排序方式 | No | string |
| type | query | 帖子类型 | No | string |
| author_id | query | 作者ID | No | integer |
| tag_id | query | 标签ID | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [response.PageData](#responsepagedata-model) & { **"list"**: [ [model.Post](#modelpost-model) ] } } |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

#### POST
##### Summary

创建帖子

##### Description

创建新的帖子（支持普通帖和问答帖）

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| body | body | 帖子信息 | Yes | [service.CreatePostInput](#servicecreatepostinput-model) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 创建成功 | [response.Response](#responseresponse-model) & { **"data"**: [model.Post](#modelpost-model) } |
| 400 | 请求参数错误 | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /posts/{id}

#### GET
##### Summary

获取帖子详情

##### Description

根据ID获取帖子的详细信息，包括点赞状态

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 帖子ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的帖子ID | [response.Response](#responseresponse-model) |
| 404 | 帖子不存在 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

#### PUT
##### Summary

更新帖子

##### Description

更新自己的帖子（管理员可以更新任何帖子）

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 帖子ID | Yes | integer |
| body | body | 帖子信息 | Yes | [service.UpdatePostInput](#serviceupdatepostinput-model) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 更新成功 | [response.Response](#responseresponse-model) & { **"data"**: [model.Post](#modelpost-model) } |
| 400 | 请求参数错误 | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |
| 404 | 帖子不存在 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

#### DELETE
##### Summary

删除帖子

##### Description

删除自己的帖子（管理员可以删除任何帖子）

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 帖子ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 删除成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的帖子ID | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |
| 404 | 帖子不存在 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /posts/{id}/like

#### POST
##### Summary

点赞帖子

##### Description

为指定帖子点赞

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 帖子ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 点赞成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的帖子ID | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

#### DELETE
##### Summary

取消点赞帖子

##### Description

取消对指定帖子的点赞

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 帖子ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 取消点赞成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的帖子ID | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

---
### /tags

#### GET
##### Summary

获取所有标签

##### Description

获取系统中所有的标签列表

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [ [model.Tag](#modeltag-model) ] } |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

#### POST
##### Summary

创建标签（仅管理员）

##### Description

创建一个新的标签，需要管理员权限

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| body | body | 标签信息 | Yes | [service.CreateTagInput](#servicecreatetaginput-model) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 创建成功 | [response.Response](#responseresponse-model) & { **"data"**: [model.Tag](#modeltag-model) } |
| 400 | 请求参数错误 | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /tags/{id}

#### PUT
##### Summary

更新标签（仅管理员）

##### Description

更新指定标签的信息，需要管理员权限

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 标签ID | Yes | integer |
| body | body | 标签信息 | Yes | [service.CreateTagInput](#servicecreatetaginput-model) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 更新成功 | [response.Response](#responseresponse-model) & { **"data"**: [model.Tag](#modeltag-model) } |
| 400 | 请求参数错误或无效的标签ID | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |
| 404 | 标签不存在 | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

#### DELETE
##### Summary

删除标签（仅管理员）

##### Description

删除指定标签，需要管理员权限

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 标签ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 删除成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的标签ID | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |
| 404 | 标签不存在 | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

---
### /timeline/following

#### GET
##### Summary

获取关注时间线

##### Description

获取当前用户关注的人的内容时间线

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| page | query | 页码 | No | integer |
| page_size | query | 每页数量 | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [response.PageData](#responsepagedata-model) & { **"list"**: [ [model.TimelineEvent](#modeltimelineevent-model) ] } } |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /timeline/home

#### GET
##### Summary

获取首页时间线

##### Description

获取当前用户首页的时间线（推荐内容）

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| page | query | 页码 | No | integer |
| page_size | query | 每页数量 | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [response.PageData](#responsepagedata-model) & { **"list"**: [ [model.TimelineEvent](#modeltimelineevent-model) ] } } |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /timeline/subscribe/{user_id}

#### POST
##### Summary

关注用户

##### Description

关注指定用户，接收其动态更新

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| user_id | path | 要关注的用户ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 关注成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的用户ID或不能关注自己 | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 409 | 已关注该用户 | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

#### DELETE
##### Summary

取消关注用户

##### Description

取消关注指定用户

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| user_id | path | 要取消关注的用户ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 取消关注成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的用户ID | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 404 | 未关注该用户 | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /timeline/subscribe/{user_id}/status

#### GET
##### Summary

检查是否已关注

##### Description

检查当前用户是否已关注指定用户

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| user_id | path | 要检查的用户ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的用户ID | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /timeline/subscriptions

#### GET
##### Summary

获取关注列表

##### Description

获取当前用户关注的所有用户列表

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [ [model.User](#modeluser-model) ] } |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

---
### /topics

#### GET
##### Summary

获取专题列表

##### Description

分页获取所有专题列表

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| page | query | 页码 | No | integer |
| page_size | query | 每页数量 | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [response.PageData](#responsepagedata-model) & { **"list"**: [ [model.Topic](#modeltopic-model) ] } } |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

#### POST
##### Summary

创建专题

##### Description

创建一个新的专题（需要管理员权限）

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| body | body | 专题信息 | Yes | [service.CreateTopicInput](#servicecreatetopicinput-model) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 创建成功 | [response.Response](#responseresponse-model) & { **"data"**: [model.Topic](#modeltopic-model) } |
| 400 | 请求参数错误 | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /topics/creator/{creator_id}

#### GET
##### Summary

获取用户创建的专题

##### Description

获取指定用户创建的所有专题

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| creator_id | path | 创建者用户ID | Yes | integer |
| page | query | 页码 | No | integer |
| page_size | query | 每页数量 | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [response.PageData](#responsepagedata-model) & { **"list"**: [ [model.Topic](#modeltopic-model) ] } } |
| 400 | 无效的用户ID | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

### /topics/{id}

#### GET
##### Summary

获取专题详情

##### Description

根据ID获取专题详细信息

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 专题ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [model.Topic](#modeltopic-model) } |
| 400 | 无效的专题ID | [response.Response](#responseresponse-model) |
| 404 | 专题不存在 | [response.Response](#responseresponse-model) |

#### PUT
##### Summary

更新专题

##### Description

更新指定专题的信息（需要管理员权限）

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 专题ID | Yes | integer |
| body | body | 专题信息 | Yes | [service.CreateTopicInput](#servicecreatetopicinput-model) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 更新成功 | [response.Response](#responseresponse-model) & { **"data"**: [model.Topic](#modeltopic-model) } |
| 400 | 请求参数错误或无效的专题ID | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |
| 404 | 专题不存在 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

#### DELETE
##### Summary

删除专题

##### Description

删除指定专题（需要管理员权限）

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 专题ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 删除成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的专题ID | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |
| 404 | 专题不存在 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /topics/{id}/follow

#### POST
##### Summary

关注专题

##### Description

关注指定专题，接收专题更新通知

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 专题ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 关注成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的专题ID | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 404 | 专题不存在 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

#### DELETE
##### Summary

取消关注专题

##### Description

取消关注指定专题

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 专题ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 取消关注成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的专题ID | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 404 | 专题不存在或未关注 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /topics/{id}/follow/status

#### GET
##### Summary

检查是否关注专题

##### Description

检查当前用户是否已关注指定专题

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 专题ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的专题ID | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /topics/{id}/followers

#### GET
##### Summary

获取专题关注者列表

##### Description

分页获取关注指定专题的用户列表

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 专题ID | Yes | integer |
| page | query | 页码 | No | integer |
| page_size | query | 每页数量 | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [response.PageData](#responsepagedata-model) & { **"list"**: [ [model.User](#modeluser-model) ] } } |
| 400 | 无效的专题ID | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

### /topics/{id}/posts

#### GET
##### Summary

获取专题帖子列表

##### Description

分页获取指定专题下的所有帖子

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 专题ID | Yes | integer |
| page | query | 页码 | No | integer |
| page_size | query | 每页数量 | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [response.PageData](#responsepagedata-model) & { **"list"**: [ [model.Post](#modelpost-model) ] } } |
| 400 | 无效的专题ID | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

#### POST
##### Summary

添加帖子到专题

##### Description

将指定帖子添加到专题中（需要管理员权限）

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 专题ID | Yes | integer |
| body | body | 帖子信息 | Yes | [handler.AddPostToTopicRequest](#handleraddposttotopicrequest-model) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 添加成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 请求参数错误 | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |
| 404 | 专题或帖子不存在 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /topics/{id}/posts/{post_id}

#### DELETE
##### Summary

从专题移除帖子

##### Description

将指定帖子从专题中移除（需要管理员权限）

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 专题ID | Yes | integer |
| post_id | path | 帖子ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 移除成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的ID | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 403 | 无权限 | [response.Response](#responseresponse-model) |
| 404 | 专题或帖子不存在 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

---
### /users/leaderboard

#### GET
##### Summary

获取用户排行榜

##### Description

根据用户积分或活跃度获取排行榜

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| limit | query | 返回数量 | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [ [model.User](#modeluser-model) ] } |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

### /users/profile

#### PUT
##### Summary

更新用户资料

##### Description

更新当前登录用户的个人资料

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| body | body | 用户资料信息 | Yes | [service.UpdateProfileInput](#serviceupdateprofileinput-model) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 更新成功 | [response.Response](#responseresponse-model) & { **"data"**: [service.UserProfileResponse](#serviceuserprofileresponse-model) } |
| 400 | 请求参数错误 | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /users/{id}

#### GET
##### Summary

获取用户资料

##### Description

根据用户ID获取用户详细资料

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 用户ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [service.UserProfileResponse](#serviceuserprofileresponse-model) } |
| 400 | 无效的用户ID | [response.Response](#responseresponse-model) |
| 404 | 用户不存在 | [response.Response](#responseresponse-model) |

### /users/{id}/follow

#### POST
##### Summary

关注用户

##### Description

关注指定用户

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 要关注的用户ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 关注成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的用户ID或不能关注自己 | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 409 | 已关注该用户 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

#### DELETE
##### Summary

取消关注用户

##### Description

取消关注指定用户

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 要取消关注的用户ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 取消关注成功 | [response.Response](#responseresponse-model) & { **"data"**: object } |
| 400 | 无效的用户ID | [response.Response](#responseresponse-model) |
| 401 | 未授权 | [response.Response](#responseresponse-model) |
| 404 | 未关注该用户 | [response.Response](#responseresponse-model) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| ApiKeyAuth |  |

### /users/{id}/followers

#### GET
##### Summary

获取粉丝列表

##### Description

获取指定用户的粉丝列表

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 用户ID | Yes | integer |
| page | query | 页码 | No | integer |
| page_size | query | 每页数量 | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [response.PageData](#responsepagedata-model) & { **"list"**: [ [model.User](#modeluser-model) ] } } |
| 400 | 无效的用户ID | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

### /users/{id}/following

#### GET
##### Summary

获取关注列表

##### Description

获取指定用户关注的用户列表

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | 用户ID | Yes | integer |
| page | query | 页码 | No | integer |
| page_size | query | 每页数量 | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | 获取成功 | [response.Response](#responseresponse-model) & { **"data"**: [response.PageData](#responsepagedata-model) & { **"list"**: [ [model.User](#modeluser-model) ] } } |
| 400 | 无效的用户ID | [response.Response](#responseresponse-model) |
| 500 | 服务器内部错误 | [response.Response](#responseresponse-model) |

---
### Models

#### handler.AddModeratorRequest Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| can_ban_user | boolean | *Example:* `true` | No |
| can_delete_post | boolean | *Example:* `true` | No |
| can_edit_any_post | boolean | *Example:* `false` | No |
| can_manage_moderator | boolean | *Example:* `false` | No |
| can_pin_post | boolean | *Example:* `true` | No |
| user_id | integer | *Example:* `1` | Yes |

#### handler.AddPostToTopicRequest Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| post_id | integer | 帖子ID<br>*Example:* `123` | Yes |
| sort_order | integer | 排序顺序<br>*Example:* `0` | No |

#### handler.BanUserRequest Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| expires_at | string | *Example:* `"2024-12-31T23:59:59Z"` | No |
| reason | string | *Example:* `"发布违规内容"` | Yes |
| user_id | integer | *Example:* `1` | Yes |

#### handler.SetUserActiveRequest Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| active | boolean | 用户状态：true-启用，false-禁用<br>*Example:* `true` | No |

#### handler.SetUserRoleRequest Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| role | string | 用户角色<br>*Enum:* `"user"`, `"moderator"`, `"admin"`<br>*Example:* `"moderator"` | Yes |

#### handler.VoteAnswerRequest Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| vote_type | string | 投票类型：up-赞同，down-反对<br>*Enum:* `"up"`, `"down"`<br>*Example:* `"up"` | Yes |

#### model.ActionType Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| model.ActionType | string |  |  |

#### model.Board Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| cover | string |  | No |
| created_at | string |  | No |
| description | string |  | No |
| icon | string |  | No |
| id | integer |  | No |
| name | string |  | No |
| parent_id | integer |  | No |
| post_count | integer |  | No |
| post_role | [model.UserRole](#modeluserrole-model) |  | No |
| reply_role | [model.UserRole](#modeluserrole-model) |  | No |
| slug | string |  | No |
| sort_order | integer |  | No |
| thread_count | integer |  | No |
| today_count | integer |  | No |
| updated_at | string |  | No |
| view_role | [model.UserRole](#modeluserrole-model) |  | No |

#### model.BoardTree Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| children | [ [model.BoardTree](#modelboardtree-model) ] |  | No |
| cover | string |  | No |
| description | string |  | No |
| icon | string |  | No |
| id | integer |  | No |
| name | string |  | No |
| parent_id | integer |  | No |
| post_count | integer |  | No |
| post_role | [model.UserRole](#modeluserrole-model) |  | No |
| reply_role | [model.UserRole](#modeluserrole-model) |  | No |
| slug | string |  | No |
| sort_order | integer |  | No |
| thread_count | integer |  | No |
| today_count | integer |  | No |
| view_role | [model.UserRole](#modeluserrole-model) |  | No |

#### model.Comment Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| author | [model.User](#modeluser-model) |  | No |
| author_id | integer |  | No |
| content | string |  | No |
| created_at | string |  | No |
| id | integer |  | No |
| is_accepted | boolean | 新增 | No |
| is_answer | boolean | 新增 | No |
| like_count | integer |  | No |
| parent | [model.Comment](#modelcomment-model) |  | No |
| parent_id | integer |  | No |
| post_id | integer |  | No |
| replies | [ [model.Comment](#modelcomment-model) ] |  | No |
| updated_at | string |  | No |
| vote_count | integer | 新增 | No |

#### model.Moderator Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| board | [model.Board](#modelboard-model) |  | No |
| board_id | integer |  | No |
| can_ban_user | boolean |  | No |
| can_delete_post | boolean |  | No |
| can_edit_any_post | boolean |  | No |
| can_manage_moderator | boolean |  | No |
| can_pin_post | boolean |  | No |
| created_at | string |  | No |
| id | integer |  | No |
| permissions | string |  | No |
| updated_at | string |  | No |
| user | [model.User](#modeluser-model) |  | No |
| user_id | integer |  | No |

#### model.Notification Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| content | string |  | No |
| created_at | string |  | No |
| id | integer |  | No |
| is_read | boolean |  | No |
| sender | [model.User](#modeluser-model) |  | No |
| sender_id | integer |  | No |
| target_id | integer |  | No |
| target_type | string |  | No |
| type | [model.NotificationType](#modelnotificationtype-model) |  | No |
| updated_at | string |  | No |
| user_id | integer |  | No |

#### model.NotificationType Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| model.NotificationType | string |  |  |

#### model.Post Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| author | [model.User](#modeluser-model) |  | No |
| author_id | integer |  | No |
| board | [model.Board](#modelboard-model) | 新增 | No |
| board_id | integer | 新增 | No |
| content | string |  | No |
| cover | string |  | No |
| created_at | string |  | No |
| id | integer |  | No |
| is_question | boolean | 新增：标记是否为问答帖 | No |
| like_count | integer |  | No |
| pin_in_board | boolean | 新增 | No |
| pin_top | boolean |  | No |
| question | [model.Question](#modelquestion-model) | 新增 | No |
| status | [model.PostStatus](#modelpoststatus-model) |  | No |
| summary | string |  | No |
| tags | [ [model.Tag](#modeltag-model) ] |  | No |
| title | string |  | No |
| type | [model.PostType](#modelposttype-model) |  | No |
| updated_at | string |  | No |
| view_count | integer |  | No |

#### model.PostStatus Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| model.PostStatus | string |  |  |

#### model.PostType Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| model.PostType | string |  |  |

#### model.Question Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| accepted_answer | [model.Comment](#modelcomment-model) |  | No |
| accepted_answer_id | integer |  | No |
| answer_count | integer |  | No |
| created_at | string |  | No |
| id | integer |  | No |
| post | [model.Post](#modelpost-model) |  | No |
| post_id | integer |  | No |
| reward_score | integer |  | No |
| updated_at | string |  | No |

#### model.Tag Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| color | string |  | No |
| created_at | string |  | No |
| description | string |  | No |
| id | integer |  | No |
| name | string |  | No |
| post_count | integer |  | No |
| updated_at | string |  | No |

#### model.TimelineEvent Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| action | [model.ActionType](#modelactiontype-model) |  | No |
| actor | [model.User](#modeluser-model) |  | No |
| actor_id | integer |  | No |
| created_at | string |  | No |
| id | integer |  | No |
| payload | string |  | No |
| score | integer |  | No |
| target_id | integer |  | No |
| target_type | string |  | No |
| updated_at | string |  | No |
| user | [model.User](#modeluser-model) |  | No |
| user_id | integer |  | No |

#### model.Topic Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| cover | string |  | No |
| created_at | string |  | No |
| creator | [model.User](#modeluser-model) |  | No |
| creator_id | integer |  | No |
| description | string |  | No |
| follower_count | integer |  | No |
| id | integer |  | No |
| is_public | boolean |  | No |
| post_count | integer |  | No |
| title | string |  | No |
| updated_at | string |  | No |

#### model.User Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| avatar | string |  | No |
| bio | string |  | No |
| created_at | string |  | No |
| email | string |  | No |
| id | integer |  | No |
| is_active | boolean | 激活，是否可以登录 | No |
| is_blocked | boolean | 封禁，是否可以发言 | No |
| last_login | string |  | No |
| role | [model.UserRole](#modeluserrole-model) |  | No |
| score | integer |  | No |
| updated_at | string |  | No |
| username | string |  | No |

#### model.UserRole Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| model.UserRole | string |  |  |

#### response.PageData Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| list |  |  | No |
| page | integer |  | No |
| page_size | integer |  | No |
| total | integer |  | No |

#### response.Response Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| code | integer |  | No |
| data |  |  | No |
| message | string |  | No |

#### service.CreateBoardInput Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| cover | string |  | No |
| description | string |  | No |
| icon | string |  | No |
| name | string |  | Yes |
| parent_id | integer |  | No |
| post_role | string |  | No |
| reply_role | string |  | No |
| slug | string |  | Yes |
| sort_order | integer |  | No |
| view_role | string |  | No |

#### service.CreateCommentInput Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| content | string |  | Yes |
| parent_id | integer |  | No |
| post_id | integer |  | Yes |

#### service.CreatePostInput Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| content | string |  | Yes |
| cover | string |  | No |
| summary | string |  | No |
| tag_ids | [ integer ] |  | No |
| title | string |  | Yes |
| type | string |  | No |

#### service.CreateTagInput Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| color | string |  | No |
| description | string |  | No |
| name | string |  | Yes |

#### service.CreateTopicInput Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| cover | string |  | No |
| description | string |  | No |
| is_public | boolean |  | No |
| title | string |  | Yes |

#### service.LoginInput Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| email | string |  | Yes |
| password | string |  | Yes |

#### service.RegisterInput Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| email | string |  | Yes |
| password | string |  | Yes |
| username | string |  | Yes |

#### service.UpdatePostInput Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| content | string |  | No |
| cover | string |  | No |
| summary | string |  | No |
| tag_ids | [ integer ] |  | No |
| title | string |  | No |

#### service.UpdateProfileInput Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| avatar | string |  | No |
| bio | string |  | No |

#### service.UserProfileResponse Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| avatar | string |  | No |
| bio | string |  | No |
| created_at | string |  | No |
| email | string |  | No |
| follower_count | integer |  | No |
| following_count | integer |  | No |
| id | integer |  | No |
| is_active | boolean | 激活，是否可以登录 | No |
| is_blocked | boolean | 封禁，是否可以发言 | No |
| is_following | boolean |  | No |
| last_login | string |  | No |
| role | [model.UserRole](#modeluserrole-model) |  | No |
| score | integer |  | No |
| updated_at | string |  | No |
| username | string |  | No |
