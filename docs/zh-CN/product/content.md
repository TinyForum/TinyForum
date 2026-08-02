# 内容安全

---

## 内容安全策略

TinyForum 通过多层防护机制保障社区内容质量与安全合规。

### 防护层级

```
用户输入 → DFA 快速匹配 → LLM 语义复核 → 内容审核队列 → 风险评分 → 决策执行
```

---

## 敏感词过滤

### DFA（确定性有限自动机）

- **原理**：构建敏感词 Trie 树，O(n) 时间复杂度快速扫描
- **词典管理**：
  - 词典文件位于 `backend/sensitive-dicts/`
  - 分类：政治、色情、赌博、医疗、广告等
  - 支持黑白名单、正则匹配
- **配置**：
  ```yaml
  # backend/config/risk_control.yml
  sensitive:
    enabled: true
    dict_dir: "./sensitive-dicts"
    policy: "block"  # block / warn / audit
  ```

### LLM 语义复核

由于 DFA 算法存在上下文误判（如"东风导弹"可能被 DFA 误判为军事敏感词），引入 LLM 进行二次判断：

- **配置**：`backend/config/ai.yml` 中配置 Ollama 或其他 LLM 服务地址
- **模型推荐**：`qwen3:0.6b` 或其他轻量级中文模型
- **工作流程**：
  1. DFA 匹配到疑似敏感词
  2. 将上下文提交给 LLM 判断是否真实违规
  3. LLM 返回 `safe` / `unsafe` / `uncertain`
  4. 根据结果决定放行/拦截/人工审核

```yaml
# backend/config/ai.yml
provider: ollama
ollama:
  base_url: "http://localhost:11434"
  model: "qwen3:0.6b"
  timeout: 10s
```

---

## 内容安全检查中间件

`ContentCheck` 中间件在请求进入业务逻辑前进行拦截：

```go
// 对帖子标题和内容进行敏感词检查
api.POST("/posts", mw.ContentCheck([]string{"title", "content"}), h.Create)
```

### 检查策略

| 策略 | 说明 |
|------|------|
| `block` | 直接拦截，返回 400 |
| `warn` | 允许发布，但记录日志并通知审核 |
| `audit` | 帖子进入待审核队列 |

---

## 文件上传安全

### 限制措施

| 限制项 | 值 | 说明 |
|--------|-----|------|
| 最大文件大小 | 10 MB | 超过限制拒绝 |
| 允许的图片格式 | jpg, jpeg, png, gif, webp, svg | MIME 类型校验 |
| 允许的文档格式 | pdf, txt, md | 白名单机制 |
| 禁止的类型 | exe, sh, php, html | 安全拦截 |

### 存储策略

- 文件存储于 `backend/uploads/` 目录
- 生产环境建议使用 OSS / S3 对象存储
- 文件访问通过 `GET /api/v1/files/:file_id` 公开接口

---

## 行为风控

与内容安全配合，行为风控关注用户操作模式。详见 [行为风控](/zh-CN/product/risk)。

### 联动机制

```
新用户 + 短时间内 >3 条帖子 = 自动进入审核队列
高敏感词评分 + 新注册账号 = 自动拦截并记录风控日志
同一 IP + 多账号注册 = IP 风险标记
```
