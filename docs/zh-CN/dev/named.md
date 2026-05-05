# 数据转换规范（修订版）

## 一、核心原则

- **职责分离**：每一层只处理自己的数据形态，通过转换器（Converter）进行边界隔离。
- **统一响应**：所有接口返回一致的 `Result<T>` 结构，错误与成功均包装其中。
- **分页泛型**：分页数据统一使用 `PageResult<T>`，所有层共用此结构。
- **无歧义命名**：避免 `Response` 与 `Request` 混淆，视图对象不使用 `*Response` 后缀。

## 二、数据对象定义

| 对象类型    | 命名建议     | 所属层                    | 用途                                             | 转换发生位置                                  |
| ----------- | ------------ | ------------------------- | ------------------------------------------------ | --------------------------------------------- |
| **Request** | `XxxRequest` | Handler                   | 接收原始请求参数，仅做绑定与基础校验             | Handler → Service（转 BO）                    |
| **VO**      | `XxxVO`      | Handler（返回值）         | 视图对象，直接序列化返回给前端                   | Service → Handler（由 Converter 将 BO 转 VO） |
| **BO**      | `XxxBO`      | Service 核心              | 业务对象，Service 内部及横向调用的标准模型       | Service 入口/出口处转换                       |
| **DO**      | `XxxDO`      | Repo（数据库）            | 数据对象，与数据库表一一对应                     | Service ↔ Repo（由 Converter 转 BO ↔ DO）     |
| **DTO**     | `XxxDTO`     | 可选（跨服务/跨模块边界） | 数据传输对象，用于解耦不同微服务或模块之间的模型 | Service 边界处转换（BO ↔ DTO）                |

> **说明**：`DTO` 仅在以下场景使用：  
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

## 五、各层详细规范

### 1. Handler 层
- **职责**：解析请求、参数校验、调用 Service、包装响应。
- **输入**：`XxxRequest`（仅用于绑定，不附加任何业务逻辑）。
- **输出**：`Result<XxxVO>`，其中 `Result` 是全局响应包装器。
- **禁止**：在 Handler 内进行 BO ↔ DO 转换，或直接返回 DO。

**响应包装器定义**（与 VO 命名不冲突）：
```go
type Result[T any] struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    T      `json:"data"`   // 可为 XxxVO 或 PageResult<XxxVO>
}
```

### 2. Service 层
- **职责**：核心业务逻辑，事务边界，调用 Repo 及其他 Service。
- **输入/输出**：`XxxBO`（公开方法统一使用 BO，非特殊情况不使用 DTO）。
- **横向调用**：同一项目内 Service 之间直接传递 BO，无需额外 DTO。
- **转换**：在入口处将 Request 转为 BO，在出口处将 BO 转为 VO；调用 Repo 前将 BO 转为 DO，接收 DO 后转回 BO。
- **建议**：所有转换逻辑抽取到独立的 `Converter` 组件（如 `XxxConverter`），避免 Service 代码膨胀。

### 3. Repo 层
- **职责**：数据库访问、缓存操作，仅做数据存取。
- **输入/输出**：`XxxDO` 或 `PageResult<XxxDO>`。
- **禁止**：接收 BO 或 返回 BO，不包含任何业务逻辑。
- **跨 Repo 调用**：不允许直接调用，应由 Service 层编排。

### 4. DTO
- **使用场景**：
  - 调用外部系统的接口，其数据模型与内部 BO 差异较大。
  - 当存在自定义查询项目的时候使用 DTO，避免污染 DO。
- **转换位置**：在 Service 边界处（如调用外部客户端的方法内）完成 BO ↔ DTO 转换。
- **命名**：`XxxDTO`，不嵌入 VO/PO 的分页结构，应使用独立的 `PageResult<DTO>`。

## 六、转换器（Converter）模式

推荐为每个聚合根或模块创建 Converter，集中管理转换逻辑：

```go
type XxxConverter interface {
    RequestToBO(req *XxxRequest) *XxxBO
    BOToVO(bo *XxxBO) *XxxVO
    BOToDO(bo *XxxBO) *XxxDO
    DOToBO(do *XxxDO) *XxxBO
    // 可选：BOPageToVOPage(page *PageResult[XxxBO]) *PageResult[XxxVO]
}
```

- 使用依赖注入或工厂模式提供 Converter 实例。
- 避免在循环中重复创建 Converter，保持无状态。

## 七、异常与错误处理

- 业务异常（如用户不存在、余额不足）使用自定义业务错误码，通过 `Result` 返回给前端。
- 系统异常（如数据库连接失败）由全局异常处理器捕获，统一返回 `Result` 格式（Code=500）。
- Service 层不直接返回 `Result`，而是返回 `(BO, error)` 或 `(PageResult[BO], error)`，由 Handler 组装 `Result`。

## 八、命名约定总览

| 类型       | 命名示例            | 说明                             |
| ---------- | ------------------- | -------------------------------- |
| 请求对象   | `CreateUserRequest` | 位于 `model/request` 包          |
| 视图对象   | `UserVO`            | 位于 `model/vo` 包               |
| 业务对象   | `UserBO`            | 位于 `model/bo` 包               |
| 数据对象   | `UserDO`            | 位于 `model/do` 包               |
| 传输对象   | `UserDTO`           | 位于 `model/dto` 包              |
| 分页结果   | `PageResult[T]`     | 通用泛型，放在 `model/common` 包 |
| 响应包装器 | `Result[T]`         | 通用泛型，放在 `model/common` 包 |

## 九、最佳实践检查清单

- [ ] Handler 只负责参数绑定和调用 Service，不写业务逻辑。
- [ ] Service 公开方法只使用 BO 作为参数和返回值（除非有明确的外部契约需要 DTO）。
- [ ] Repo 方法只接收和返回 DO。
- [ ] 所有跨层转换都通过 Converter 完成，禁止在 Service/Handler 里手动逐字段赋值。
- [ ] 分页数据一律使用 `PageResult<T>`，不返回裸 `[]T`。
- [ ] 统一响应格式为 `Result<T>`，错误也通过 `Result` 返回，不自定义结构。
- [ ] 包结构按功能模块（如 `user`）组织，内部再按分层（`controller`、`service`、`repository`）划分。

