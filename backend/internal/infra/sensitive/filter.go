package sensitive

import (
	"net/http"
	"sync"
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
	dfa         DFAFilter
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
