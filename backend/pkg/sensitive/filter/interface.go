package filter

import (
	"net/http"
	"sync"
	"tiny-forum/pkg/sensitive/dfa"
)

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
	dfa         dfa.DFAFilter
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
		dfa:         dfa.NewDFAFilter(),
		blockWords:  make(map[string]bool),
		reviewWords: make(map[string]bool),
		ollama:      ollama,
	}
	if ollama != nil {
		f.httpClient = &http.Client{Timeout: ollama.timeout()}
	}
	return f
}
