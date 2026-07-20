package sensitive

// ActionPolicy 定义一种动作的配置
type ActionPolicy struct {
	Score      int      `yaml:"score"`
	Categories []string `yaml:"categories"`
}

// PolicyConfig 对应 policy.yml 根结构
type PolicyConfig struct {
	Version string                  `yaml:"version"`
	Actions map[string]ActionPolicy `yaml:"actions"` // key: "block", "review", ...
}
