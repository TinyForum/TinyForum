# 数据转换规范

## 一、核心原则

- **职责分离**：每一层只处理自己的数据形态，通过转换器（Converter）进行边界隔离。
- **统一响应**：所有接口返回一致的 `Result<T>` 结构，错误与成功均包装其中。
- **分页泛型**：分页数据统一使用 `PageResult<T>`，所有层共用此结构。
- **无歧义命名**：避免 `Response` 与 `Request` 混淆，视图对象不使用 `*Response` 后缀。




## 二、各层详细规范

### 1. Handler 层
- **职责**：

  1. 解析请求，进行参数校验：由于绑定函数 `SouldBindxxx()` 已经进行了参数验证，因而之后不在验证参数类型，仅检查值是否合规。

     > 注意：`Request` 结构体中存在类型绑定限制，为避免绑定失败，所有兜底策略都使用 `common.NewDefault(&reqestName)` 方法，保证所有字段都有合规值。

  2. 调用 Service：每个 `Handler` 都有自己的 `Service` ，因而直接调用，将业务交由 `Service` 处理

  3. 包装 BO：Service 的传入参数类型必须是 BO，Handler 需要转换 Request 为 BO，转换函数定义在 `converter` 包中。

  4. 包装 VO：为了数据安全，返回前端的数据必须是 VO（直接结构体） 或者 common（分页结构体），严格禁止 DO 返回

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

#### 情况一：分页查询

```go
ListPosts(ctx context.Context, ListPostsBO *bo.PageQuery[bo.ListPosts]) ([]do.Post, int64, error) 
```

保持上下文传递，并将参数聚合为一个 `BO`，分页数据使用 `bo.PageQuery` 包装。

#### 情况二：基本查询





### 3. Repo 层
- **职责**：数据库访问、缓存操作，仅做数据存取。
- **输入/输出**：`XxxDO` 或 `PageResult<XxxDO>`。
- **禁止**：
  1. 接收 BO（应该接受 DO）
  2. 返回 BO（应该返回 DO）
  3. 不包含任何业务逻辑。
- **跨 Repo 调用**：不允许直接调用，应由 Service 层编排。

#### 参数情况一：分页查询

```go
List(ctx context.Context, query *dto.PluginQueryDTO, pageParam common.PageParam) (*common.PageResult[do.PluginMeta], error)
```

保持上下文传递，并将参数聚合为  `DO` 作为入库标准。



```bash
AdminList(ctx context.Context, ListPostsDO *common.PageQuery[do.Post]) ([]do.Post, int64, error)
```

接收：使用 common 包装分页请求，CRUD 数据合规为 `DO`

返回：DO，绝对信任 Service 层处理。





### 4. DTO
- **使用场景**：
  - 调用外部系统的接口，其数据模型与内部 BO 差异较大。
  - 当存在自定义查询项目的时候使用 DTO，避免污染 DO。
- **转换位置**：在 Service 边界处（如调用外部客户端的方法内）完成 BO ↔ DTO 转换。
- **命名**：`XxxDTO`，不嵌入 VO/PO 的分页结构，应使用独立的 `PageResult<DTO>`。

## 三、转换器（Converter）模式

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

## 四、异常与错误处理

- 业务异常（如用户不存在、余额不足）使用自定义业务错误码，通过 `Result` 返回给前端。
- 系统异常（如数据库连接失败）由全局异常处理器捕获，统一返回 `Result` 格式（Code=500）。
- Service 层不直接返回 `Result`，而是返回 `(BO, error)` 或 `(PageResult[BO], error)`，由 Handler 组装 `Result`。

## 五、最佳实践检查清单

- [ ] Handler 只负责参数绑定和调用 Service，不写业务逻辑。
- [ ] Service 公开方法只使用 BO 作为参数和返回值（除非有明确的外部契约需要 DTO）。
- [ ] Repo 方法只接收和返回 DO。
- [ ] 所有跨层转换都通过 Converter 完成，禁止在 Service/Handler 里手动逐字段赋值。
- [ ] 分页数据一律使用 `PageResult<T>`，不返回裸 `[]T`。
- [ ] 统一响应格式为 `Result<T>`，错误也通过 `Result` 返回，不自定义结构。
- [ ] 包结构按功能模块（如 `user`）组织，内部再按分层（`controller`、`service`、`repository`）划分。

