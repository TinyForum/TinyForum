# API 设计规范最佳实践（RESTful）

## 1. URL 设计

- **使用名词复数**：`/posts`, `/users`, `/comments`。
- **资源嵌套**：例如 `/posts/{post_id}/comments` 表示某帖子的评论列表。
- **避免动词**：用 HTTP 方法表达操作，如 `POST /posts` 创建，`DELETE /posts/{id}` 删除。若必须使用动词，放在 URL 末尾：`/posts/{id}/like`。
- **使用小写字母和短横线**：`/user-profiles`（但常用单/复数名词即可）。
- **版本号**：放在路径开头 `/api/v1/...`。

## 2. HTTP 方法与语义

| 方法   | 语义           | 示例                              | 是否幂等 |
|--------|----------------|-----------------------------------|----------|
| GET    | 获取资源       | `GET /posts`<br>`GET /posts/1`    | 是       |
| POST   | 创建资源       | `POST /posts`                     | 否       |
| PUT    | 全量更新       | `PUT /posts/1`                    | 是       |
| PATCH  | 部分更新       | `PATCH /posts/1`                  | 否       |
| DELETE | 删除资源       | `DELETE /posts/1`                 | 是       |

- **POST vs PUT**：POST 用于“创建”，服务端决定新资源 ID；PUT 用于“替换”或“创建在已知 ID”。
- **具体到你的项目**：点赞用 `POST /posts/{id}/like`，取消点赞用 `DELETE /posts/{id}/like` —— 虽非严格 RESTful，但被广泛接受。

## 3. 状态码规范

| 状态码 | 含义                     | 使用场景                                         |
|--------|--------------------------|--------------------------------------------------|
| 200    | OK                       | GET、PUT、PATCH、DELETE 成功                     |
| 201    | Created                  | POST 成功创建资源，响应头应含 `Location`         |
| 204    | No Content               | 删除成功或更新成功但无需返回 body                |
| 400    | Bad Request              | 请求参数缺失、格式错误、业务校验失败             |
| 401    | Unauthorized             | 未认证（未提供 token 或 token 无效）             |
| 403    | Forbidden                | 已认证但无权限（如普通用户调用管理员接口）       |
| 404    | Not Found                | 资源不存在                                       |
| 409    | Conflict                 | 资源冲突（如重复创建、状态冲突）                 |
| 422    | Unprocessable Entity     | 语义错误（比如参数校验失败，更细粒度）           |
| 429    | Too Many Requests        | 触发限流                                         |
| 500    | Internal Server Error    | 服务端内部错误                                   |

- 你的项目中已有 `response` 包，可以统一封装成功/失败响应。

## 4. 请求与响应格式

- **统一使用 JSON**：`Content-Type: application/json`，`Accept: application/json`。
- **请求体**：字段使用 `camelCase` 或 `snake_case`，但推荐 `snake_case` 与数据库字段一致（配合 JSON tag）。
- **成功响应结构**（建议统一）：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": { ... }
  }
  ```
  或直接返回数据（简单场景）。你的项目已有 `response.Success(c, data)`，保持全局统一即可。

- **错误响应结构**：
  ```json
  {
    "code": 40001,
    "message": "参数错误：title不能为空",
    "details": { ... }
  }
  ```

- **分页响应**：
  ```json
  {
    "data": [...],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5
    }
  }
  ```

## 5. 参数传递

- **路径参数**：资源 ID，如 `/posts/{id}`。
- **查询参数**：过滤、排序、分页。如 `/posts?page=1&page_size=20&sort=-created_at`。
- **请求体**：复杂对象（创建、更新）。
- **分页参数**：建议 `page` 和 `page_size`，最大限制 100 或可配置。

## 6. 认证与授权

- 使用 **JWT**（你的项目已有 `jwt` 包），放在 `Authorization: Bearer <token>` 头。
- 中间件分类：`AuthMW`（强制认证）、`OptionalAuthMW`（可选，用户未认证时游客模式）。
- 角色权限：通过中间件 `AdminRequiredMW`、`ModeratorRequiredMW` 处理。

## 7. 限流与幂等

- **限流**：使用 `RateLimitMW` 如你的项目，按用户 + 接口维度。
- **幂等性**：非幂等方法（POST）可引入 `Idempotency-Key` 头防止重复提交。

## 8. 版本管理

- URL 版本：`/api/v1/...`，当 API 发生不兼容变更时递增版本号（如 `/api/v2/...`）。
- 也可以使用 `Accept` 头 `application/vnd.your-api.v1+json`，但 URL 版本最简单。

## 9. 文档

- 使用 **Swagger/OpenAPI**（你的项目已集成），注解写在 Handler 方法上，自动生成文档。
- 示例（gin-swagger）：
  ```go
  // ListPosts 获取帖子列表
  // @Summary 获取帖子列表
  // @Tags posts
  // @Produce json
  // @Param page query int false "页码"
  // @Param page_size query int false "每页数量"
  // @Success 200 {object} response.Response{data=[]dto.Post}
  // @Router /posts [get]
  ```

## 10. 一致性原则

- **资源路径**：始终使用复数名词。
- **字段命名**：统一风格（如 `user_id` 而非 `userId`）。
- **错误码**：定义业务错误码表，便于前端处理。
- **时间格式**：使用 ISO8601，如 `2025-04-25T10:30:00Z`。

