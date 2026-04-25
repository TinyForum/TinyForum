package sensitive

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"unicode"
)

// trieNode DFA 树节点
type trieNode struct {
	children map[rune]*trieNode
	isEnd    bool
}

func newTrieNode() *trieNode {
	return &trieNode{children: make(map[rune]*trieNode)}
}

// DFAFilter DFA 敏感词过滤器
type DFAFilter struct {
	root *trieNode
	mu   sync.RWMutex
}

// NewDFAFilter 创建一个新的 DFA 过滤器
func NewDFAFilter() *DFAFilter {
	return &DFAFilter{root: newTrieNode()}
}

// normalize 规范化字符：去除不可见字符、统一大小写、全角转半角
func normalize(r rune) (rune, bool) {
	// 跳过不可见 / 零宽 / 控制字符
	if r == '\u200b' || r == '\u200c' || r == '\u200d' ||
		r == '\ufeff' || r == '\u00ad' || unicode.Is(unicode.Cf, r) {
		return 0, false
	}
	// 全角字母数字 → 半角
	if r >= 0xFF01 && r <= 0xFF5E {
		r = r - 0xFEE0
	}
	// 统一小写
	r = unicode.ToLower(r)
	return r, true
}

// normalizeRunes 将文本规范化为 rune 切片，同时保留原始位置映射
// 返回：规范化后的 rune 列表，以及每个规范化 rune 对应原始 rune 的索引
func normalizeRunes(text string) ([]rune, []int) {
	runes := []rune(text)
	var norm []rune
	var idx []int
	for i, r := range runes {
		nr, ok := normalize(r)
		if !ok {
			continue
		}
		norm = append(norm, nr)
		idx = append(idx, i)
	}
	return norm, idx
}

// AddWord 动态添加敏感词（支持多个）
func (f *DFAFilter) AddWord(words ...string) error {
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
func (f *DFAFilter) DelWord(words ...string) error {
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

func (f *DFAFilter) delWord(node *trieNode, runes []rune, depth int) bool {
	if depth == len(runes) {
		if !node.isEnd {
			return false
		}
		node.isEnd = false
		return len(node.children) == 0
	}
	r := runes[depth]
	nr, ok := normalize(r)
	if !ok {
		return false
	}
	child, exists := node.children[nr]
	if !exists {
		return false
	}
	shouldDelete := f.delWord(child, runes, depth+1)
	if shouldDelete {
		delete(node.children, nr)
		return !node.isEnd && len(node.children) == 0
	}
	return false
}

// LoadDictFile 从文件加载词库（每行一个词）
func (f *DFAFilter) LoadDictFile(path string) error {
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
func (f *DFAFilter) LoadDictContent(content string) error {
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

// match 从规范化 rune 列表的 start 位置开始尝试最长匹配
// 返回：匹配到的原始文本（通过原始 rune 列表重建），以及匹配结束的规范化索引（不含），-1 表示无匹配
func (f *DFAFilter) match(norm []rune, origRunes []rune, origIdx []int, start int) (string, int) {
	node := f.root
	lastEnd := -1
	lastEndOrig := -1

	for i := start; i < len(norm); i++ {
		child, ok := node.children[norm[i]]
		if !ok {
			break
		}
		node = child
		if node.isEnd {
			lastEnd = i + 1
			if i < len(origIdx) {
				lastEndOrig = origIdx[i] + 1
			}
		}
	}

	if lastEnd == -1 {
		return "", -1
	}

	// 通过原始 rune 重建命中词
	origStart := origIdx[start]
	word := string(origRunes[origStart:lastEndOrig])
	return word, lastEnd
}

// IsSensitive 判断文本中是否存在敏感词
func (f *DFAFilter) IsSensitive(text string) bool {
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
func (f *DFAFilter) FindOne(text string) string {
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
func (f *DFAFilter) FindAll(text string) []string {
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
func (f *DFAFilter) FindAllCount(text string) map[string]int {
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
func (f *DFAFilter) Replace(text string, replaceChar rune) string {
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
func (f *DFAFilter) Remove(text string) string {
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
