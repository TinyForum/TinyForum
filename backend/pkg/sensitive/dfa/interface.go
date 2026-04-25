package dfa

import (
	"sync"
)

type DFAFilter interface {
	AddWord(words ...string) error
	DelWord(words ...string) error
	LoadDictFile(path string) error
	LoadDictContent(content string) error
	IsSensitive(text string) bool
	FindOne(text string) string
	FindAllCount(text string) map[string]int
	FindAll(text string) []string
	Replace(text string, replaceChar rune) string
	Remove(text string) string
}

// DFAFilter DFA 敏感词过滤器
type dFAFilter struct {
	root *trieNode
	mu   sync.RWMutex
}

// NewDFAFilter 创建一个新的 DFA 过滤器
func NewDFAFilter() DFAFilter {
	return &dFAFilter{root: newTrieNode()}
}
