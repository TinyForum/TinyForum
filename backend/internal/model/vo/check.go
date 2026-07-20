package vo

import "tiny-forum/internal/infra/sensitive"

// CheckResult 内容检测结果
type CheckResult struct {
	Original  string                   `json:"original"`  // 原始内容
	Masked    string                   `json:"masked"`    // 处理后的内容
	Sensitive bool                     `json:"sensitive"` // 是否包含敏感词
	Level     sensitive.Level          `json:"level"`     // 敏感词级别
	Score     int                      `json:"score"`     // 匹配分数
	Action    sensitive.Action         `json:"action"`    // 处理动作
	Matches   []*sensitive.MatchResult `json:"matches"`   // 匹配结果
	HitWords  []string                 `json:"hit_words"` // 命中的敏感词
	Replaced  string
	Passed    bool
}
