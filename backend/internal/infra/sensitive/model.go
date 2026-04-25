package sensitive

// trieNode DFA 树节点
type trieNode struct {
	children map[rune]*trieNode
	isEnd    bool
}

// Level 敏感词命中等级
type Level int

const (
	LevelClean    Level = 0 // [安全] 无敏感词
	LevelReplace  Level = 1 // 【攻击】替换敏感词
	LevelReview   Level = 2 // 【风险】人工审核
	LevelShadowed Level = 3 // 【防对抗】遮蔽词
	LevelBlock    Level = 4 // 【禁止】直接拦截
)

// CheckResult 检测结果
type CheckResult struct {
	Level    Level    // 最高命中等级
	HitWords []string // 命中的词列表
	Text     string   // 替换后的文本（* 号替换）
}

// DictLoadResult 目录加载结果摘要
type DictLoadResult struct {
	BlockFiles   []string
	ReviewFiles  []string
	ShadowFiles  []string
	ReplaceFiles []string
	Errors       []string
	TotalWords   int
}
