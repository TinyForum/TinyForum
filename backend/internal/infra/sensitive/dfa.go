package sensitive

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
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
type dfaFilter struct {
	root *trieNode
	mu   sync.RWMutex
}

// NewDFAFilter 创建一个新的 DFA 过滤器
func NewDFAFilter() DFAFilter {
	return &dfaFilter{root: newTrieNode()}
}

// AddWord 动态添加敏感词（支持多个）
func (f *dfaFilter) AddWord(words ...string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, word := range words {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}
		node := f.root
		for _, r := range []rune(word) {
			nr, ok := normalize(r)
			if !ok {
				continue
			}
			if _, exists := node.children[nr]; !exists {
				node.children[nr] = newTrieNode()
			}
			node = node.children[nr]
		}
		node.isEnd = true
	}
	return nil
}

// DelWord 动态删除敏感词（支持多个）
func (f *dfaFilter) DelWord(words ...string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, word := range words {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}
		f.delWord(f.root, []rune(word), 0)
	}
	return nil
}

// LoadDictFile 从文件加载词库（每行一个词）
func (f *dfaFilter) LoadDictFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
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
	if err := scanner.Err(); err != nil {
		return err
	}
	return f.AddWord(words...)
}

// LoadDictContent 从字符串内容加载词库
func (f *dfaFilter) LoadDictContent(content string) error {
	var words []string
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		words = append(words, line)
	}
	return f.AddWord(words...)
}

// IsSensitive 判断文本中是否存在敏感词
func (f *dfaFilter) IsSensitive(text string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	origRunes := []rune(text)
	norm, origIdx := normalizeRunes(text)

	for i := 0; i < len(norm); i++ {
		if _, end := f.match(norm, origRunes, origIdx, i); end != -1 {
			return true
		}
	}
	return false
}

// FindOne 查找文本中第一个敏感词
func (f *dfaFilter) FindOne(text string) string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	origRunes := []rune(text)
	norm, origIdx := normalizeRunes(text)

	for i := 0; i < len(norm); i++ {
		if word, end := f.match(norm, origRunes, origIdx, i); end != -1 {
			_ = end
			return word
		}
	}
	return ""
}

// FindAll 查找文本中所有敏感词（去重），跳过已匹配区域
func (f *dfaFilter) FindAll(text string) []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	origRunes := []rune(text)
	norm, origIdx := normalizeRunes(text)

	seen := make(map[string]bool)
	var result []string

	i := 0
	for i < len(norm) {
		if word, end := f.match(norm, origRunes, origIdx, i); end != -1 {
			if !seen[word] {
				seen[word] = true
				result = append(result, word)
			}
			i = end
		} else {
			i++
		}
	}
	return result
}

// FindAllCount 查找所有敏感词及其出现次数
func (f *dfaFilter) FindAllCount(text string) map[string]int {
	f.mu.RLock()
	defer f.mu.RUnlock()

	origRunes := []rune(text)
	norm, origIdx := normalizeRunes(text)

	result := make(map[string]int)

	i := 0
	for i < len(norm) {
		if word, end := f.match(norm, origRunes, origIdx, i); end != -1 {
			result[word]++
			i = end
		} else {
			i++
		}
	}
	return result
}

// Replace 替换所有敏感词为指定字符（按敏感词长度替换对应数量的字符）
func (f *dfaFilter) Replace(text string, replaceChar rune) string {
	log.Printf(text)

	f.mu.RLock()
	defer f.mu.RUnlock()

	origRunes := []rune(text)
	norm, origIdx := normalizeRunes(text)

	// 标记原始 rune 中哪些位置需要被替换
	mask := make([]bool, len(origRunes))

	i := 0
	for i < len(norm) {
		if word, end := f.match(norm, origRunes, origIdx, i); end != -1 {
			// 标记从 origIdx[i] 到 origIdx[end-1]+1 的原始 rune
			fmt.Println("word:", word, "end:", end)
			origStart := origIdx[i]
			origEnd := origIdx[end-1] + 1
			for k := origStart; k < origEnd; k++ {
				mask[k] = true
			}
			i = end
		} else {
			i++
		}
	}

	// 构建结果
	var sb strings.Builder
	for i, r := range origRunes {
		if mask[i] {
			sb.WriteRune(replaceChar)
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// Remove 从文本中删除所有敏感词
func (f *dfaFilter) Remove(text string) string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	origRunes := []rune(text)
	norm, origIdx := normalizeRunes(text)

	mask := make([]bool, len(origRunes))

	i := 0
	for i < len(norm) {
		if _, end := f.match(norm, origRunes, origIdx, i); end != -1 {
			origStart := origIdx[i]
			origEnd := origIdx[end-1] + 1
			for k := origStart; k < origEnd; k++ {
				mask[k] = true
			}
			i = end
		} else {
			i++
		}
	}

	var sb strings.Builder
	for i, r := range origRunes {
		if !mask[i] {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
