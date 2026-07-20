package sensitive

import (
	"log"
	"regexp"
	"strings"
	"sync"
)

// Matcher 管理所有规则的正则和精确词
type Matcher struct {
	rules      []Rule           // 所有普通规则
	regexRules []*regexp.Regexp // 编译后的正则（对应规则索引）
	wordMap    map[string][]int // 精确词 -> 规则索引（用于白/黑）
	mu         sync.RWMutex
}

// NewMatcher 从规则列表构建匹配器
func NewMatcher(rules []Rule) *Matcher {
	m := &Matcher{
		rules:   rules,
		wordMap: make(map[string][]int),
	}
	m.compile()
	return m
}

// compile 编译正则并构建词索引
func (m *Matcher) compile() {
	m.regexRules = make([]*regexp.Regexp, len(m.rules))
	for i, rule := range m.rules {
		if len(rule.Regex) == 0 {
			continue
		}
		var pattern string
		if len(rule.Regex) == 1 {
			raw := rule.Regex[0]
			// 移除所有空白字符
			pattern = strings.Map(func(r rune) rune {
				if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
					return -1
				}
				return r
			}, raw)
			// 如果模式为空，跳过编译
			if pattern == "" {
				log.Printf("警告：规则 %s (%s) 的正则模式为空，跳过", rule.ID, rule.Name)
				continue
			}
		} else {
			pattern = strings.Join(rule.Regex, "|")
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			log.Printf("编译规则 %s (%s) 失败: %v", rule.ID, rule.Name, err)
			continue
		}
		m.regexRules[i] = re
	}
}

// Match 匹配文本，返回所有命中结果
func (m *Matcher) Match(text string) []*MatchResult {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var results []*MatchResult
	// 1. 精确词匹配（白/黑名单）
	for word, indices := range m.wordMap {
		// 简单查找所有出现位置（可优化为更高效，但数据量通常不大）
		pos := 0
		for {
			idx := strings.Index(text[pos:], word)
			if idx == -1 {
				break
			}
			start := pos + idx
			end := start + len(word)
			for _, ruleIdx := range indices {
				rule := m.rules[ruleIdx]
				results = append(results, &MatchResult{
					RuleID:   rule.ID,
					Category: rule.Category,
					Word:     word,
					Start:    start,
					End:      end,
					Action:   rule.Action,
					Score:    rule.Score,
					Tags:     rule.Tags,
				})
			}
			pos = start + len(word)
		}
	}
	// 2. 正则匹配
	for i, re := range m.regexRules {
		if re == nil {
			continue
		}
		rule := m.rules[i]
		matches := re.FindAllStringIndex(text, -1)
		for _, loc := range matches {
			start, end := loc[0], loc[1]
			if start == end {
				continue // 跳过空匹配
			}
			word := text[start:end]
			results = append(results, &MatchResult{
				RuleID:   rule.ID,
				Category: rule.Category,
				Word:     word,
				Start:    start,
				End:      end,
				Action:   rule.Action,
				Score:    rule.Score,
				Tags:     rule.Tags,
			})
		}
	}

	return results
}
