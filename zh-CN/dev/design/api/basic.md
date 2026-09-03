# 响应的基本结构

## 简单响应

```json
{
  "code": 0,
  "data": "string",
  "message": "string",
  "request_id": "string",
  "timestamp": 0,
  "trace_id": "string"
}
```

详情项：选择性嵌入（只嵌入当前页面必需的）

列表项：极简，所有关联对象只用 ID 或 count
