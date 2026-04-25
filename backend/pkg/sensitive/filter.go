package sensitive

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"tiny-forum/pkg/utils"
)

// Level 敏感词命中等级
type Level int

const (
	LevelClean    Level = 0 // [安全] 无敏感词
	LevelReplace  Level = 1 // 【攻击】替换敏感词
	LevelReview   Level = 2 // 【风险】人工审核
	LevelShadowed Level = 3 // 【防对抗】遮蔽词
	LevelBlock    Level = 4 // 【禁止】直接拦截
)

// CheckResult 检测结果
type CheckResult struct {
	Level    Level    // 最高命中等级
	HitWords []string // 命中的词列表
	Text     string   // 替换后的文本（* 号替换）
}

// DictLoadResult 目录加载结果摘要
type DictLoadResult struct {
	BlockFiles   []string
	ReviewFiles  []string
	ShadowFiles  []string
	ReplaceFiles []string
	Errors       []string
	TotalWords   int
}

// OllamaConfig Ollama 客户端配置
type OllamaConfig struct {
	// BaseURL Ollama 服务地址，默认 http://localhost:11434
	BaseURL string
	// Model 使用的模型名称，如 "qwen2.5:7b"
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

// Filter 敏感词过滤器接口
type Filter interface {
	// DFA 基础能力
	IsSensitive(text string) bool
	FindOne(text string) string
	FindAll(text string) []string
	FindAllCount(text string) map[string]int
	Replace(text string, replaceChar rune) string
	Remove(text string) string
	AddWord(words ...string) error
	DelWord(words ...string) error

	// 分级能力
	Check(text string) CheckResult
	AddLevelWords(level Level, words []string)

	// 加载词库
	LoadDictFile(path string, level Level) error
	LoadDictContent(content string, level Level) error
	LoadDictDir(dir string) (DictLoadResult, error)
}

// filter 分级过滤器
type filter struct {
	dfa         *DFAFilter
	blockWords  map[string]bool
	reviewWords map[string]bool
	mu          sync.RWMutex

	// Ollama 相关
	ollama     *OllamaConfig
	httpClient *http.Client
}

// NewFilter 创建过滤器。ollama 为 nil 时禁用 LLM 复判，退化为纯 DFA 模式。
func NewFilter(ollama *OllamaConfig) Filter {
	f := &filter{
		dfa:         NewDFAFilter(),
		blockWords:  make(map[string]bool),
		reviewWords: make(map[string]bool),
		ollama:      ollama,
	}
	if ollama != nil {
		f.httpClient = &http.Client{Timeout: ollama.timeout()}
	}
	return f
}

var (
	globalFilter Filter
	globalOnce   sync.Once
)

// GetGlobalFilter 获取全局单例过滤器（无 LLM）
func GetGlobalFilter() Filter {
	globalOnce.Do(func() {
		globalFilter = NewFilter(nil)
	})
	return globalFilter
}

// ---- DFA 基础能力转发 ----

func (f *filter) IsSensitive(text string) bool            { return f.dfa.IsSensitive(text) }
func (f *filter) FindOne(text string) string              { return f.dfa.FindOne(text) }
func (f *filter) FindAll(text string) []string            { return f.dfa.FindAll(text) }
func (f *filter) FindAllCount(text string) map[string]int { return f.dfa.FindAllCount(text) }
func (f *filter) Replace(text string, r rune) string      { return f.dfa.Replace(text, r) }
func (f *filter) Remove(text string) string               { return f.dfa.Remove(text) }

func (f *filter) AddWord(words ...string) error {
	if err := f.dfa.AddWord(words...); err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, w := range words {
		if !f.reviewWords[w] {
			f.blockWords[w] = true
		}
	}
	return nil
}

func (f *filter) DelWord(words ...string) error {
	if err := f.dfa.DelWord(words...); err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, w := range words {
		delete(f.blockWords, w)
		delete(f.reviewWords, w)
	}
	return nil
}

func (f *filter) AddLevelWords(level Level, words []string) {
	_ = f.dfa.AddWord(words...)
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, w := range words {
		switch level {
		case LevelBlock:
			f.blockWords[w] = true
			delete(f.reviewWords, w)
		case LevelReview:
			f.reviewWords[w] = true
			delete(f.blockWords, w)
		}
	}
}

// ---- 词库加载 ----

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		words = append(words, line)
	}
	return words, scanner.Err()
}

func (f *filter) LoadDictFile(path string, level Level) error {
	words, err := readLines(path)
	if err != nil {
		return fmt.Errorf("sensitive: load %s: %w", path, err)
	}
	f.AddLevelWords(level, words)
	return nil
}

func (f *filter) LoadDictContent(content string, level Level) error {
	var words []string
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		words = append(words, line)
	}
	f.AddLevelWords(level, words)
	return nil
}

func fileLevel(name string) Level {
	lower := strings.ToLower(name)
	if strings.HasPrefix(lower, "review_") {
		return LevelReview
	}
	return LevelBlock
}

func (f *filter) LoadDictDir(dir string) (DictLoadResult, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return DictLoadResult{}, fmt.Errorf("sensitive: open dir %s: %w", dir, err)
	}

	var result DictLoadResult
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".txt") {
			continue
		}

		path := filepath.Join(dir, name)
		level := fileLevel(name)

		words, err := readLines(path)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", name, err))
			continue
		}

		f.AddLevelWords(level, words)
		result.TotalWords += len(words)

		if level == LevelReview {
			result.ReviewFiles = append(result.ReviewFiles, name)
		} else {
			result.BlockFiles = append(result.BlockFiles, name)
		}
	}
	return result, nil
}

// ---- Check（DFA + 可选 LLM 复判）----

// MARK: Check
// Check 检测文本，返回分级结果。
//
// 流程：
//  1. DFA 扫描，无命中 → LevelClean，直接返回
//  2. 有 LevelBlock 命中 → 直接返回 LevelBlock，不调 LLM（节省资源）
//  3. 仅有 LevelReview 命中 → 调 Ollama 复判：
//     - LLM 判定"安全" → 降级为 LevelClean
//     - LLM 判定"风险" → 保持 LevelReview
//     - LLM 判定"违规" → 升级为 LevelBlock
//     - LLM 调用失败  → 保守保持 LevelReview（不降级）
func (f *filter) Check(html string) CheckResult {
	text, err := utils.HTMLToText(html)
	if err != nil {
		return CheckResult{
			Level:    LevelBlock,
			HitWords: []string{},
			Text:     "HTML 解析失败",
		}
	}

	result := CheckResult{
		Level:    LevelClean,
		HitWords: []string{},
		Text:     text,
	}

	// --- Step 1: DFA 扫描 ---
	hitMap := f.dfa.FindAllCount(text)
	if len(hitMap) == 0 {
		return result
	}

	f.mu.RLock()
	var hits []string
	hasBlock := false
	hasReview := false
	for word := range hitMap {
		hits = append(hits, word)
		switch {
		case f.blockWords[word]:
			hasBlock = true
		case f.reviewWords[word]:
			hasReview = true
		default:
			hasBlock = true // 未显式分级 → 默认 block
		}
	}
	f.mu.RUnlock()

	result.HitWords = hits

	// --- Step 2: LevelBlock → 直接拦截，跳过 LLM ---
	if hasBlock {
		result.Level = LevelBlock
		result.Text = f.dfa.Replace(text, '*')
		return result
	}

	// --- Step 3: LevelReview → LLM 复判 ---
	if hasReview {
		result.Level = LevelReview
		result.Text = f.dfa.Replace(text, '*')

		if f.ollama != nil {
			llmLevel, err := f.llmReview(text, hits)
			if err != nil {
				// LLM 失败：保守保持 Review，记录日志
				log.Printf("[sensitive] LLM 复判失败，保守保持 LevelReview: %v", err)
			} else {
				result.Level = llmLevel
				if llmLevel == LevelClean {
					result.Text = text // 安全内容恢复原文
					result.HitWords = []string{}
				}
			}
		}
		return result
	}

	return result
}

// ---- Ollama 复判 ----

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
