# 数据对象定义

| 对象类型    | 命名建议     | 所属层                    | 用途                                             | 转换发生位置                                  |
| ----------- | ------------ | ------------------------- | ------------------------------------------------ | --------------------------------------------- |
| **Request** | `XxxRequest` | Handler                   | 接收原始请求参数，仅做绑定与基础校验             | Handler → Service（转 BO）                    |
| **VO**      | `XxxVO`      | Handler（返回值）         | 视图对象，直接序列化返回给前端                   | Service → Handler（由 Converter 将 BO 转 VO） |
| **BO**      | `XxxBO`      | Service 核心              | 业务对象，Service 内部及横向调用的标准模型       | Service 入口/出口处转换                       |
| **DO**      | `XxxDO`      | Repo（数据库）            | 数据对象，与数据库表一一对应                     | Service ↔ Repo（由 Converter 转 BO ↔ DO）     |
| **DTO**     | `XxxDTO`     | 可选（跨服务/跨模块边界） | 数据传输对象，用于解耦不同微服务或模块之间的模型 | Service 边界处转换（BO ↔ DTO）                |

> **说明**：`DTO` 仅在以下场景使用：  
>
> - 调用的外部服务（第三方 API、消息队列）数据结构与内部 BO 不一致。  
> - 微服务之间需要裁剪字段、避免暴露内部 BO 细节。  
> - 单体项目中若模块间耦合不深，可省略 DTO，直接用 BO。

## 三、分页模型

统一使用泛型 `PageResult<T>`，定义如下：

```go
type PageResult[T any] struct {
    Total    int64 `json:"total"`
    Page     int   `json:"page"`
    PageSize int   `json:"pageSize"`
    List     []T   `json:"list"`
}
```

- 所有分页请求参数统一使用 `PageParam`（包含 `Page`、`PageSize`），不依赖任何特定层。
- 分页响应统一返回 `PageResult<VO>` 给前端，`PageResult<DO>` 仅在 Repo 内部使用。

## 四、分层数据流

```
[前端]
   │ Request (JSON)
   ▼
Handler
   │ 1. 绑定 & 校验 Request
   │ 2. Converter: Request → BO
   ▼
Service
   │ 3. 业务处理：接收 BO，返回 BO
   │    （横向调用其他 Service 时仍使用 BO）
   │ 4. Converter: BO → VO
   ▼
Handler
   │ 5. 组装 Result<VO> 返回
   ▼
[前端]
```

内部数据访问流：

```
Service (BO)
   │ Converter: BO → DO（如需查询/写入）
   ▼
Repo
   │ 接收 DO，返回 DO（或 PageResult<DO>）
   ▼
Database
   │ 回向路径：DO → Converter → BO
```
