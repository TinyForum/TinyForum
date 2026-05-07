# response 方法返回数据

在标准响应中，如果是成功，返回 `response.Success(c,data,options)`



## 标准数据

方法：`response.Success(c, data, options)`

- `c`：`gin` 上下文
- `data`：`vo` 结构体
- `options`：可选项（来自 `BasicResponse`）



## 分页数据

方法：`response.SuccessPage(c, list, total, page, pageSize, options)`

- `c`：`gin` 上下文
- `list`：`vo` 类型的数组
- `total`：总页数
- `page`：当前页数
- `pageSize`：页面数据数量
- `options`：可选项（来自 `BasicResponse`）



## BasicResponse 基本响应数据

1. **`Code`：业务状态码（默认返回）**
2. `Message`：业务返回信息
3. **`Data`：业务数据（默认返回）**
4. `Timestamp`：时间戳
5. `RequestID`：请求 ID
6. `TraceID`：追踪 ID





