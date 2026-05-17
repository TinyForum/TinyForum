package vo

import "tiny-forum/internal/infra/sensitive"

// CheckResult 内容检测结果
type CheckResult struct {
	Passed   bool            // false 表示直接拦截
	Level    sensitive.Level // 命中等级
	HitWords []string        // 命中词
	Replaced string          // 替换后的内容
}
