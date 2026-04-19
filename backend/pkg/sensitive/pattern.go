package sensitive

import (
	"regexp"
	"strings"
	"unicode"
)

// ---- Email ----

var emailRe = regexp.MustCompile(
	`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`,
)

// HasEmail 判断字符串中是否存在邮箱地址
func HasEmail(s string) bool {
	return emailRe.MatchString(s)
}

// MaskEmail 将字符串中存在的邮箱地址替换成 "*"
func MaskEmail(s string) string {
	return emailRe.ReplaceAllStringFunc(s, func(m string) string {
		return strings.Repeat("*", len([]rune(m)))
	})
}

// ---- URL ----

var urlRe = regexp.MustCompile(
	`(?i)(https?://|ftp://|www\.)[^\s<>"'，。！？；：、\x00-\x1f]+`,
)

// HasURL 判断字符串中是否存在网址
func HasURL(s string) bool {
	return urlRe.MatchString(s)
}

// MaskURL 将字符串中存在的网址替换成 "*"
func MaskURL(s string) string {
	return urlRe.ReplaceAllStringFunc(s, func(m string) string {
		return strings.Repeat("*", len([]rune(m)))
	})
}

// ---- Digit ----

// countDigits 统计字符串中的数字字符数量
func countDigits(s string) int {
	count := 0
	for _, r := range s {
		if unicode.IsDigit(r) {
			count++
		}
	}
	return count
}

// HasDigit 判断字符串中是否存在大于等于 count 个数字字符
func HasDigit(s string, count int) bool {
	return countDigits(s) >= count
}

// MaskDigit 将字符串中存在的数字替换成 "*"
func MaskDigit(s string) string {
	var sb strings.Builder
	for _, r := range s {
		if unicode.IsDigit(r) {
			sb.WriteRune('*')
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// ---- WechatID ----
// 微信号规则：6-20 位，字母开头，由字母/数字/下划线/短横线组成
// 为避免误伤，前面必须有边界（空格、中文标点、字符串开头等）

var wechatRe = regexp.MustCompile(
	`(?:^|[\s,，。！？；：、【】（）\(\)\[\]])([a-zA-Z][a-zA-Z0-9_\-]{5,19})(?:$|[\s,，。！？；：、【】（）\(\)\[\]])`,
)

// HasWechatID 判断字符串中是否存在微信号
func HasWechatID(s string) bool {
	return wechatRe.MatchString(s)
}

// MaskWechatID 将字符串中存在的微信号替换成 "*"
func MaskWechatID(s string) string {
	return wechatRe.ReplaceAllStringFunc(s, func(m string) string {
		// 保留首尾的边界字符，只替换捕获组内容
		runes := []rune(m)
		n := len(runes)
		if n == 0 {
			return m
		}
		// 确定首尾是否为边界字符（非字母）
		start := 0
		end := n
		prefix := ""
		suffix := ""
		if !isAlpha(runes[0]) {
			prefix = string(runes[0])
			start = 1
		}
		if n > 1 && !isAlpha(runes[n-1]) {
			suffix = string(runes[n-1])
			end = n - 1
		}
		id := runes[start:end]
		return prefix + strings.Repeat("*", len(id)) + suffix
	})
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}