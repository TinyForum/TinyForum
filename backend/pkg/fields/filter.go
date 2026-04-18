package fields

import "strings"

/**
 * @Description: 根据白名单过滤请求的字段，若过滤后为空则返回默认字段
 * @param requested 请求的字段
 * @param allowed 白名单
 * @param defaults 默认字段
 * @return []string
 */
func Filter(requested string, allowed []string, defaults []string) []string {
	if requested == "" {
		return defaults
	}

	// 构建白名单 map 用于快速查找
	allowMap := make(map[string]bool, len(allowed))
	for _, f := range allowed {
		allowMap[f] = true
	}

	parts := strings.Split(requested, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if allowMap[p] {
			result = append(result, p)
		}
	}

	if len(result) == 0 {
		return defaults
	}
	return result
}
