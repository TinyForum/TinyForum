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

// 辅助函数：隐藏邮箱中间部分
func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	local := parts[0]
	if len(local) <= 2 {
		return "***@" + parts[1]
	}

	masked := local[:2] + strings.Repeat("*", len(local)-2) + "@" + parts[1]
	return masked
}

// 辅助函数
func getUnifiedMessage(locale string) string {
	if locale == "zh-CN" || locale == "zh" {
		return "如果您的邮箱已注册，您将收到密码重置链接"
	}
	return "If your email is registered, you will receive a password reset link"
}
