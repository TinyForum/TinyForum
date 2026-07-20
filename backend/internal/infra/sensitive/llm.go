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

// AIConfig LLM 配置
type AIConfig struct {
	Enable      bool          // 是否启用
	Provider    string        // 固定 "openai"
	APIKey      string        // API Key
	BaseURL     string        // 默认 https://api.openai.com/v1
	Model       string        // 默认 gpt-3.5-turbo
	Timeout     time.Duration // 请求超时（秒）
	MaxRetries  int
	Temperature float64
	MaxTokens   int
}

func (c *AIConfig) baseURL() string {
	if c.BaseURL == "" {
		return "https://api.openai.com/v1"
	}
	return strings.TrimRight(c.BaseURL, "/")
}

func (c *AIConfig) effectiveTimeout() time.Duration {
	if c.Timeout <= 0 {
		return 60 * time.Second
	}
	// 假设配置中 timeout 为毫秒（与 config.yml 一致）
	return time.Duration(c.Timeout) * time.Millisecond
}

// LLMClient 封装 API 调用
type LLMClient struct {
	config AIConfig
	client *http.Client
}

func NewLLMClient(cfg AIConfig) *LLMClient {
	return &LLMClient{
		config: cfg,
		client: &http.Client{
			Timeout: cfg.effectiveTimeout(),
		},
	}
}

// ChatRequest 兼容 OpenAI 格式
type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Call 发送 Prompt 并返回模型回复
func (c *LLMClient) Call(ctx context.Context, prompt string) (string, error) {
	url := c.config.baseURL() + "/chat/completions"
	reqBody := chatRequest{
		Model: c.config.Model,
		Messages: []message{
			{Role: "user", Content: prompt},
		},
		Temperature: c.config.Temperature,
		MaxTokens:   c.config.MaxTokens,
	}
	jsonData, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	var lastErr error
	for i := 0; i <= c.config.MaxRetries; i++ {
		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("LLM API error: %d %s", resp.StatusCode, string(body))
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		var chatResp chatResponse
		if err := json.Unmarshal(body, &chatResp); err != nil {
			lastErr = err
			continue
		}
		if chatResp.Error != nil {
			lastErr = fmt.Errorf("LLM error: %s", chatResp.Error.Message)
			continue
		}
		if len(chatResp.Choices) == 0 {
			lastErr = fmt.Errorf("LLM returned empty response")
			continue
		}
		return chatResp.Choices[0].Message.Content, nil
	}
	return "", lastErr
}

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
