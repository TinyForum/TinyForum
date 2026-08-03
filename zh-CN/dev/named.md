# Go 方法命名最佳实践

方法命名在 Go 中不仅要符合语言惯例，还要体现接收者类型和方法的职责。以下是从社区和官方指南总结的最佳实践。

## 1. 接收器命名

- **使用类型名的首字母小写**（通常 1~2 个字母）：  
  ```go
  func (u *User) Name() string       // 接收器 u
  func (s *Service) Start() error    // 接收器 s
  ```
- **避免使用 `this`、`self`**：Go 更推荐短接收器名，以增强代码可读性。
- **同一类型的方法接收器名称应保持一致**：不要 `func (u *User)` 混用 `func (usr *User)`。
- **指针 vs 值接收器**：
  - 需要修改接收器状态 → 指针 `*T`
  - 接收器较大（如大型结构体） → 指针避免拷贝
  - 小类型且不需要修改 → 值接收器 `T`

## 2. Getter / Setter 命名

- **Getter 不使用 `Get` 前缀**，直接用字段名（首字母大写为导出）：
  ```go
  func (u *User) Name() string { return u.name }
  ```
- **Setter 使用 `Set` 前缀**：
  ```go
  func (u *User) SetName(name string) { u.name = name }
  ```
- 例外：当需要与字段名区分（如字段是 `Name`，但 Getter 需要额外逻辑）时，仍推荐 `Name()` 而非 `GetName()`。

## 3. 谓词方法（返回 bool）

- 使用 `Is`、`Has`、`Can`、`Should` 等前缀：
  ```go
  func (u *User) IsActive() bool
  func (r *Request) HasBody() bool
  func (p *Policy) CanRead(user *User) bool
  ```

## 4. 方法名应清晰表达行为

- **动词或动宾短语**：`Create`, `FindByID`, `SendEmail`, `MarshalJSON`
- **避免冗长**：`GetUserByID` 直接写成 `UserByID` 或 `Find`（在明确上下文中）。
- **使用特定领域动词**：`Push`/`Pop`、`Enqueue`/`Dequeue`、`Lock`/`Unlock`、`Open`/`Close`

## 5. 工厂方法与构造函数

- 构造函数命名：`New` + 类型名（或 `New` 当包内只有一个主要类型）：
  ```go
  func NewUser(name string) *User
  ```
- 返回接口的构造函数：`New` + 接口名（或 `New` + 具体实现名）：
  ```go
  func NewReader() io.Reader
  ```
- 返回非导出类型的工厂：可自由命名，如 `newHTTPClient`，但包外不可见。

## 6. 转换方法

- 将接收者转换为其他类型：`To` + 目标类型：
  ```go
  func (d Duration) ToMinutes() int
  func (u *User) ToDTO() *UserDTO
  ```
- 从其他类型构造：`From` + 源类型：
  ```go
  func FromDTO(dto *UserDTO) *User
  ```

## 7. 标准接口的方法命名

- `String() string` – 实现 `fmt.Stringer`
- `Error() string` – 实现 `error`
- `Write(p []byte) (n int, err error)` – 实现 `io.Writer`
- `Read(p []byte) (n int, err error)` – 实现 `io.Reader`
- `Close() error` – 实现 `io.Closer`

## 8. 避免的方法名

- **不要与方法所在包名重复**：如 `pkg.User` 包中定义 `User.User()` 易混淆。
- **不要使用与内置函数/类型相同的名称**：如 `len`, `cap`, `make`, `nil`, `true`。
- **不要使用仅大小写区分的名称**：Go 不推荐 `GetName` 与 `GetNAME` 同时存在。

## 9. 示例对比

| 不推荐             | 推荐                | 原因                           |
| ------------------ | ------------------- | ------------------------------ |
| `func (u *User) GetName() string` | `func (u *User) Name() string` | Getter 省略 Get 前缀            |
| `func (s *Service) ExecuteService()` | `func (s *Service) Execute()` | 省略冗余的 Service              |
| `func (c *Calculator) AddNumbers(a,b int)` | `func (c *Calculator) Add(a,b int)` | 参数名已说明含义，方法名应精简 |
| `func (t *Task) RunFunction()` | `func (t *Task) Run()` | 避免重复上下文                  |
| `func (b *Board) Register(router *gin.RouterGroup)` | `func (b *Board) RegisterRoutes(router *gin.RouterGroup)` | 明确注册的是路由                |

## 10. 特殊情况：处理器的路由挂载方法

在 Web 项目中，常见模式：
- **`RegisterRoutes`**：明确表示注册 HTTP 路由。
- **`SetupRoutes`**：同样可接受，但 `RegisterRoutes` 更直接。

```go
func (h *PostHandler) RegisterRoutes(group *gin.RouterGroup, mw *MiddlewareSet) {
    group.GET("", h.List)
    group.POST("", mw.AuthMW(), h.Create)
}
```

## 总结

方法命名应遵循：
1. **短小精悍**：接收器名短，方法名描述准确。
2. **一致性强**：同类操作使用相同的前缀/后缀。
3. **遵循社区惯例**：多阅读标准库（`io`, `http`, `strings`）的设计。
4. **避免冗余**：不要重复包名、接收器类型名。

好的命名让代码像文档一样清晰，减少注释需求。