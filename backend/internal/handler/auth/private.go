package auth

import "strings"

func parseAcceptLanguage(header string) string {
	if header == "" {
		return ""
	}
	// 解析 "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7"
	parts := strings.Split(header, ",")
	if len(parts) == 0 {
		return ""
	}
	// 取第一个，并去掉权重部分
	first := strings.Split(parts[0], ";")[0]
	// 转换为短代码（可选）
	return strings.ToLower(first[:2]) // "en" or "zh"
}
