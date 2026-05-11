package mask

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// MaskFunc 脱敏函数签名
type MaskFunc func(string) string

// 内置策略实现
var builtinStrategies = map[string]MaskFunc{
	"name":     maskName,
	"mobile":   maskMobile,
	"email":    maskEmail,
	"idcard":   maskIDCard,
	"bankcard": maskBankCard,
	"address":  maskAddressWithKeep,
	"full":     maskFull,
	"regex":    nil, // 特殊处理，需要额外参数
}

// maskName 中文姓名脱敏：张*，欧阳* -> 欧**
func maskName(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	l := len(runes)
	if l == 1 {
		return "*"
	}
	if l == 2 {
		return string(runes[0]) + "*"
	}
	// 三个字及以上：欧阳疯 -> 欧**
	return string(runes[0]) + strings.Repeat("*", l-1)
}

// maskMobile 手机号：保留前3后4
func maskMobile(s string) string {
	if len(s) < 7 {
		return maskFull(s)
	}
	return s[:3] + "****" + s[len(s)-4:]
}

// maskEmail 邮箱：保留第一个字符和@域名
func maskEmail(s string) string {
	at := strings.Index(s, "@")
	if at <= 0 {
		return maskFull(s)
	}
	local := s[:at]
	if len(local) == 0 {
		return maskFull(s)
	}
	first := local[:1]
	return first + "***" + s[at:]
}

// maskIDCard 身份证：保留前6后4，中间*
func maskIDCard(s string) string {
	if len(s) < 10 {
		return maskFull(s)
	}
	return s[:6] + strings.Repeat("*", len(s)-10) + s[len(s)-4:]
}

// maskBankCard 银行卡：保留前6后4，中间*
func maskBankCard(s string) string {
	if len(s) < 10 {
		return maskFull(s)
	}
	return s[:6] + strings.Repeat("*", len(s)-10) + s[len(s)-4:]
}

// maskFull 完全隐藏：全部替换为 *
func maskFull(s string) string {
	if s == "" {
		return ""
	}
	return strings.Repeat("*", utf8.RuneCountInString(s))
}

// maskAddressWithKeep 地址脱敏，支持 keep 参数：保留前N个字符
// 需要在处理时解析 tag 属性，策略函数本身只负责基础逻辑
func maskAddressWithKeep(s string) string {
	// 默认保留前6个字符
	return maskAddressKeepN(s, 6)
}

func maskAddressKeepN(s string, keep int) string {
	runes := []rune(s)
	if len(runes) <= keep {
		return s
	}
	return string(runes[:keep]) + strings.Repeat("*", len(runes)-keep)
}

// applyRegex 应用正则表达式脱敏
func applyRegex(s, pattern, replace string) string {
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(s, replace)
}
