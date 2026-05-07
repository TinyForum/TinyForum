# 命名约定总览

由于 `response` 为包名，在全局应避免使用 `response` 作为变量或者结构体名称



| 类型       | 命名示例        | 变量名              | 位置                             | 备注 |
| ---------- | --------------- | ------------------- | -------------------------------- | ---- |
| 请求对象   | `CreateUser`    | `createUserRequest` | 位于 `model/request` 包          |      |
| 视图对象   | `User`          | `userVO`            | 位于 `model/vo` 包               |      |
| 业务对象   | `User`          | `UserBO`            | 位于 `model/bo` 包               |      |
| 数据对象   | `User`          | `UserDO`            | 位于 `model/do` 包               |      |
| 传输对象   | `User`          | `userDTO`           | 位于 `model/dto` 包              |      |
| 分页结果   | `PageResult[T]` | `pageResult`        | 通用泛型，放在 `model/common` 包 |      |
| 响应包装器 | `Result[T]`     | `result`            | 通用泛型，放在 `model/common` 包 |      |