package sensitive

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// OllamaConfig Ollama 客户端配置
type OllamaConfig struct {
	// BaseURL Ollama 服务地址，默认 http://localhost:11434
	BaseURL string
	// Model 使用的模型名称，如 "qwen3:0.6b"
	Model string
	// Timeout 单次请求超时，默认 60s
	Timeout time.Duration
}

func (c *OllamaConfig) baseURL() string {
	if c.BaseURL == "" {
		return "http://localhost:11434"
	}
	return strings.TrimRight(c.BaseURL, "/")
}

func (c *OllamaConfig) timeout() time.Duration {
	if c.Timeout == 0 {
		return 60 * time.Second
	}
	return c.Timeout
}

// ollamaRequest Ollama /api/generate 请求体
type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// ollamaResponse Ollama /api/generate 响应体（stream=false）
type ollamaResponse struct {
	Response string `json:"response"`
}

// llmJudgment LLM 返回的结构化判断
type llmJudgment struct {
	// Level: "safe" | "review" | "block"
	Level  string `json:"level"`
	Reason string `json:"reason"`
}

// buildReviewPrompt 构建给 LLM 的复判 Prompt。
// 要求模型只输出 JSON，避免额外废话影响解析。
func buildReviewPrompt(text string, hitWords []string) string {
	wordsJSON, _ := json.Marshal(hitWords)
	prompt := fmt.Sprintf(`你是一个严格的内容安全审核助手。
以下文本在关键词系统中命中了词：%s。
请判断该文本在中文社区论坛场景下是否真正存在风险。

文本内容：
"""
%s
"""

判断标准（优先级从高到低）：
1. **block（明确违规）**：文本明确包含以下任何一类内容，且没有歧义：
   - 分裂国家、侮辱英烈、煽动颠覆政权、宣扬恐怖主义/极端主义
   - 明确的人身威胁、严重侮辱（如“我要杀了你”）
   - 色情内容的具体描写或传播渠道
   - 枪支弹药、毒品等违禁品的交易或制作方法
   - 诈骗、赌博等违法广告
   - 其他明显违反中国法律法规的内容
   注意：只有当你**非常确定**文本直接违反上述某一条时，才输出 block。

2. **safe（安全）**：
   - 文本完全无害，只是因包含某些词而被误报（如“垃圾”用于评价物品）
   - 正常讨论、常识分享、温和建议
   - 没有攻击意图或违法意图

3. **review（需人工审核）**：
   - 文本内容模糊，可能擦边，例如暗讽、谐音、轻度攻击但不清
   - 你无法确定是否属于 block 或 safe
   - 或者文本虽然命中词，但整体语境需要人类进一步判断

输出要求：
- 只输出一个 JSON 对象，不要有任何额外内容。
- 格式：{"level":"safe|review|block","reason":"<不超过30字的中文原因>"}
- 如果你无法判断，输出 {"level":"review","reason":"需要人工复核"}

示例（仅供格式参考，不要照抄内容）：
文本：“这个政策真愚蠢” 命中词：["愚蠢"] → {"level":"review","reason":"批评政策但不违规"}
文本：“我要炸了政府大楼” 命中词：["炸","政府"] → {"level":"block","reason":"明确暴力恐怖威胁"}
文本：“你真是个傻瓜” 命中词：["傻瓜"] → {"level":"review","reason":"轻度人身攻击"}
文本：“这部电影很垃圾” 命中词：["垃圾"] → {"level":"safe","reason":"评价物品，无攻击对象"}

现在请只输出你的判断 JSON。`, string(wordsJSON), text)
	log.Printf("[sensitive] LLM 复判 Prompt: %s", prompt)
	return prompt
}

// parseJudgment 从 LLM 返回的原始文本中提取 JSON 并映射到 Level。
// 容错：尝试在响应中定位第一个 JSON 对象。
func parseJudgment(raw string) (Level, error) {
	// 找到第一个 '{' 到最后一个 '}'
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start == -1 || end == -1 || end <= start {
		return LevelReview, fmt.Errorf("no JSON found in LLM response: %q", raw)
	}
	jsonStr := raw[start : end+1]

	var j llmJudgment
	if err := json.Unmarshal([]byte(jsonStr), &j); err != nil {
		return LevelReview, fmt.Errorf("parse LLM JSON %q: %w", jsonStr, err)
	}

	log.Printf("[sensitive] LLM 复判结果: level=%s reason=%s", j.Level, j.Reason)

	switch strings.ToLower(j.Level) {
	case "safe":
		return LevelClean, nil
	case "block":
		return LevelBlock, nil
	default: // "review" 或未知值 → 保守保持 Review
		return LevelReview, nil
	}
}

// ---- Ollama 复判 ----

// llmReview 调用 Ollama 对模糊命中文本进行二次判断。
// 返回最终 Level（LevelClean / LevelReview / LevelBlock）。
func (f *filter) llmReview(text string, hitWords []string) (Level, error) {
	prompt := buildReviewPrompt(text, hitWords)

	reqBody, err := json.Marshal(ollamaRequest{
		Model:  f.ollama.Model,
		Prompt: prompt,
		Stream: false,
	})
	if err != nil {
		return LevelReview, fmt.Errorf("marshal request: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), f.ollama.timeout())
	defer cancel()

	url := f.ollama.baseURL() + "/api/generate"
	log.Printf("[sensitive] LLM 复判请求: %s", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return LevelReview, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return LevelReview, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LevelReview, fmt.Errorf("read body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return LevelReview, fmt.Errorf("ollama status %d: %s", resp.StatusCode, body)
	}

	var ollamaResp ollamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return LevelReview, fmt.Errorf("unmarshal ollama response: %w", err)
	}

	return parseJudgment(ollamaResp.Response)
}
