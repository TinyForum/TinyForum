# AI 辅助

---

## AI 在 TinyForum 中的应用

目前 AI 主要用于**内容安全的敏感词复核**环节。

### 为什么需要 AI 复核？

传统 DFA（确定性有限自动机）敏感词匹配存在以下问题：

- **误判率高**：无法理解上下文语义，容易将正常讨论误标为违规
- **变体绕过**：无法识别谐音、拆字、拼音等变体形式
- **方言/俚语**：无法理解特定语境下的真实含义

引入 LLM 进行语义层面的二次判断，可以显著降低误判率，提升用户体验。

---

## 配置 AI 服务

### 方式一：Ollama（本地部署，推荐）

```bash
# 安装 Ollama
curl -fsSL https://ollama.com/install.sh | sh

# 下载模型（推荐中文模型）
ollama pull qwen3:0.6b

# 配置为监听所有网络接口（Docker 需要）
sudo systemctl edit ollama.service
```

添加：

```ini
[Service]
Environment="OLLAMA_HOST=0.0.0.0:11434"
```

```bash
sudo systemctl daemon-reload
sudo systemctl restart ollama
```

### 方式二：配置 OpenAI 兼容 API

编辑 `backend/config/ai.yml`：

```yaml
provider: openai
openai:
  api_key: "sk-xxxxxxxx"
  base_url: "https://api.openai.com/v1"
  model: "gpt-4o-mini"
  timeout: 10s
```

---

## AI 配置文件

```yaml
# backend/config/ai.yml
enabled: true
provider: ollama           # ollama / openai / custom
ollama:
  base_url: "http://localhost:11434"
  model: "qwen3:0.6b"
  timeout: 10s
  max_retries: 2

# 提示词配置
prompts:
  content_check: |
    请判断以下文本是否包含违规内容（色情、暴力、违法、广告）：
    
    文本：{{.Content}}
    
    请仅回复：safe（安全）、unsafe（违规）、uncertain（不确定）
```

---

## AI 复核流程

```
用户发布内容
    ↓
DFA 敏感词扫描
    ↓
匹配到敏感词？── 否 → 放行
    ↓ 是
提交 LLM 复核
    ↓
LLM 判断结果
    ├── safe → 放行（DFA 误判）
    ├── unsafe → 拦截（确认违规）
    └── uncertain → 进入人工审核队列
```

---

## 性能考量

- **超时设置**：建议 AI 请求超时 5-10s，避免阻塞用户请求
- **降级策略**：AI 服务不可用时，自动降级为纯 DFA 模式
- **缓存**：相同内容的 AI 判断结果可缓存一定时间
- **异步处理**：非紧急场景可异步调用 AI，先返回提交成功

---

## 扩展方向

未来 AI 可扩展到以下场景：

- **智能推荐**：基于用户兴趣推荐帖子和主题
- **内容摘要**：自动生成长文章的摘要
- **自动标签**：根据帖子内容自动推荐标签
- **情感分析**：分析评论的情感倾向，辅助社区氛围管理
- **智能客服**：Bot 机器人接入 LLM 实现智能问答
