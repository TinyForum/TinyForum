package sensitive

import (
	"strings"
)

type Normalizer struct {
	mappings map[string]string
}

func NewNormalizer(mappings map[string]string) *Normalizer {
	return &Normalizer{mappings: mappings}
}

// Normalize 执行归一化：全角转半角、大小写、映射替换等
func (n *Normalizer) Normalize(text string) string {
	// 全角转半角
	text = fullToHalf(text)
	// 大小写统一（英文）
	text = strings.ToLower(text)
	// 自定义映射替换（如 "微 信" -> "微信"）
	for from, to := range n.mappings {
		text = strings.ReplaceAll(text, from, to)
	}
	// 可增加更多：去除空格、繁简转换等（可按需扩展）
	return text
}

// fullToHalf 全角字符转半角
func fullToHalf(s string) string {
	var buf strings.Builder
	for _, r := range s {
		if r >= 0xFF01 && r <= 0xFF5E { // 全角ASCII
			buf.WriteRune(r - 0xFEE0)
		} else if r == 0x3000 { // 全角空格
			buf.WriteRune(' ')
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}
