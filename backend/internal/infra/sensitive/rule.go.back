package sensitive

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// StringSlice 支持字符串或字符串列表的 YAML 解析
type StringSlice []string

func (s *StringSlice) UnmarshalYAML(value *yaml.Node) error {
	// 尝试解码为字符串
	var str string
	if err := value.Decode(&str); err == nil {
		*s = []string{str}
		return nil
	}
	// 尝试解码为字符串切片
	var slice []string
	if err := value.Decode(&slice); err == nil {
		*s = slice
		return nil
	}
	return fmt.Errorf("期望字符串或字符串列表")
}

type Action string

// 处理动作
const (
	ActionPass    Action = "pass"    // 直接通过
	ActionReview  Action = "review"  // 人工审核
	ActionBlock   Action = "block"   // 直接拦截
	ActionShadow  Action = "shadow"  // 屏蔽
	ActionReplace Action = "replace" // 替换处理

)

type Level string

// 敏感词级别
const (
	LevelNormal   Level = "Normal"
	LevelLow      Level = "Low"
	LevelMedium   Level = "Medium"
	LevelHigh     Level = "High"
	LevelCritical Level = "Critical"
)

// 单条规则
type Rule struct {
	ID          string      `yaml:"id"`             // 唯一标识
	Name        string      `yaml:"name"`           // 规则名称
	Category    string      `yaml:"category"`       // 规则分类
	Priority    int         `yaml:"priority"`       // 优先级
	Action      Action      `yaml:"action"`         // 处理动作
	Score       int         `yaml:"score"`          // 匹配分数
	Regex       StringSlice `yaml:"regex"`          // 支持字符串或数组
	Tags        StringSlice `yaml:"tags"`           // 支持字符串或数组
	Replace     bool        `yaml:"replace"`        // 是否替换
	Description string      `yaml:"description"`    // 规则描述
	Examples    StringSlice `yaml:"examples"`       // 支持字符串或数组
	Word        string      `yaml:"word,omitempty"` // 白/黑名单专用
}

// 规则容器
type RuleContainer struct {
	Version     int    `yaml:"version"`     // 版本号
	Name        string `yaml:"name"`        // 规则集名称
	Description string `yaml:"description"` // 规则集描述
	Rules       []Rule `yaml:"rules"`       // 规则列表
}

// 匹配结果
type MatchResult struct {
	RuleID   string   `json:"rule_id"`  // 规则ID
	Category string   `json:"category"` // 规则分类
	Word     string   `json:"word"`     // 匹配到的词
	Start    int      `json:"start"`    // 开始位置
	End      int      `json:"end"`      // 结束位置
	Action   Action   `json:"action"`   // 处理动作
	Score    int      `json:"score"`    // 匹配分数
	Tags     []string `json:"tags"`     // 使用时将 StringSlice 转为 []string
}

// 检查结果
type CheckResult struct {
	Original  string         `json:"original"`  // 原始内容
	Masked    string         `json:"masked"`    // 处理后的内容
	Sensitive bool           `json:"sensitive"` // 是否包含敏感词
	Level     Level          `json:"level"`     // 敏感词级别
	Score     int            `json:"score"`     // 匹配分数
	Action    Action         `json:"action"`    // 处理动作
	Matches   []*MatchResult `json:"matches"`   // 匹配结果
}
