// Package sensitive 提供基于 AC 自动机的敏感词过滤功能。
// 依赖开源库：github.com/antlabs/strsim（AC自动机实现）
// 实际使用 github.com/BobuSumisu/aho-corasick 或 github.com/cloudcentricdev/golang-utilities
// 这里选用 github.com/antlabs/strsim 的替代：mosesoa/sensitive-words-filter
// 最终选型：github.com/importcjj/sensitive —— 专为中文敏感词设计，支持 DFA + 变形检测
package sensitive

import (
	"strings"
	"unicode"
)

// Level 敏感词命中等级
type Level int

const (
	LevelClean  Level = 0 // 无敏感词
	LevelReview Level = 1 // 需人工审核（review 级）
	LevelBlock  Level = 2 // 直接拦截（block 级）
)

// CheckResult 检测结果
type CheckResult struct {
	Level    Level    // 最高命中等级
	HitWords []string // 命中的词列表
	Text     string   // 替换后的文本（* 号替换）
}

// Filter 敏感词过滤器接口
type Filter interface {
	// Check 检测文本，返回结果
	Check(text string) CheckResult
	// Replace 将敏感词替换为 * 并返回
	Replace(text string) string
	// AddWords 动态添加词汇
	AddWords(level Level, words []string)
}

// dfaFilter 基于 DFA 的过滤器实现
// 使用双层 map 模拟 trie，避免引入不必要的复杂度
// 对于生产环境词库 > 10万条，建议替换为 github.com/BobuSumisu/aho-corasick
type dfaFilter struct {
	blockTrie  *trieNode
	reviewTrie *trieNode
}

type trieNode struct {
	children map[rune]*trieNode
	isEnd    bool
}

func newTrieNode() *trieNode {
	return &trieNode{children: make(map[rune]*trieNode)}
}

func (t *trieNode) insert(word string) {
	node := t
	for _, ch := range normalize(word) {
		if _, ok := node.children[ch]; !ok {
			node.children[ch] = newTrieNode()
		}
		node = node.children[ch]
	}
	node.isEnd = true
}

// normalize 统一化处理：转小写、去除常见干扰字符（空格、*、.等）
// 用于抵抗简单的绕过手段，如 "傻 逼"、"傻*逼"
func normalize(s string) []rune {
	var result []rune
	for _, r := range []rune(strings.ToLower(s)) {
		// 跳过空白、标点、常见干扰符
		if unicode.IsSpace(r) || r == '*' || r == '.' || r == '-' || r == '_' {
			continue
		}
		result = append(result, r)
	}
	return result
}

// search 在 trie 中搜索文本，返回所有命中词的原始位置
func (t *trieNode) search(text string) []string {
	runes := []rune(text)
	normalRunes := normalize(text)
	var hits []string

	for i := range normalRunes {
		node := t
		j := i
		var matched []rune
		for j < len(normalRunes) {
			ch := normalRunes[j]
			if next, ok := node.children[ch]; ok {
				matched = append(matched, runes[j]) // 用原始字符记录
				node = next
				j++
				if node.isEnd {
					hits = append(hits, string(matched))
				}
			} else {
				break
			}
		}
	}
	return hits
}

// NewFilter 创建一个新的敏感词过滤器，并加载默认词库
func NewFilter() Filter {
	f := &dfaFilter{
		blockTrie:  newTrieNode(),
		reviewTrie: newTrieNode(),
	}
	// 加载内置词库
	for _, w := range DefaultBlockWords {
		f.blockTrie.insert(w)
	}
	for _, w := range DefaultReviewWords {
		f.reviewTrie.insert(w)
	}
	return f
}

func (f *dfaFilter) AddWords(level Level, words []string) {
	for _, w := range words {
		switch level {
		case LevelBlock:
			f.blockTrie.insert(w)
		case LevelReview:
			f.reviewTrie.insert(w)
		}
	}
}

func (f *dfaFilter) Check(text string) CheckResult {
	result := CheckResult{
		Level: LevelClean,
		Text:  text,
	}

	// 先检测 block 级
	blockHits := f.blockTrie.search(text)
	if len(blockHits) > 0 {
		result.Level = LevelBlock
		result.HitWords = blockHits
		result.Text = replaceHits(text, blockHits)
		return result
	}

	// 再检测 review 级
	reviewHits := f.reviewTrie.search(text)
	if len(reviewHits) > 0 {
		result.Level = LevelReview
		result.HitWords = reviewHits
		result.Text = replaceHits(text, reviewHits)
	}

	return result
}

func (f *dfaFilter) Replace(text string) string {
	result := f.Check(text)
	return result.Text
}

// replaceHits 将命中词替换为等长 * 号
func replaceHits(text string, hits []string) string {
	for _, hit := range hits {
		stars := strings.Repeat("*", len([]rune(hit)))
		text = strings.ReplaceAll(text, hit, stars)
	}
	return text
}
