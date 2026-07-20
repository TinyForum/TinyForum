package sensitive

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadRules 加载目录下所有 .yml 文件（除特殊文件外）为规则列表
func LoadRules(dir string) ([]Rule, error) {
	var allRules []Rule
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".yml") && !strings.HasSuffix(d.Name(), ".yaml") {
			return nil
		}
		base := filepath.Base(path)
		if base == "normalize.yml" || base == "whitelist.yml" || base == "blacklist.yml" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var container RuleContainer
		if err := yaml.Unmarshal(data, &container); err != nil {
			return fmt.Errorf("文件 %s 解析失败: %w", path, err)
		}
		for i := range container.Rules {
			if container.Rules[i].Category == "" {
				// 从文件父目录名继承
				container.Rules[i].Category = filepath.Base(filepath.Dir(path))
			}
			if container.Rules[i].Action == "" {
				container.Rules[i].Action = ActionReview // 默认审核
			}
		}
		allRules = append(allRules, container.Rules...)
		return nil
	})
	return allRules, err
}

// LoadWhitelist 加载白名单（仅包含 word 和 action 的简单规则）
func LoadWhitelist(dir string) ([]Rule, error) {
	return loadSimpleRules(dir, "whitelist.yml")
}

// LoadBlacklist 加载黑名单
func LoadBlacklist(dir string) ([]Rule, error) {
	return loadSimpleRules(dir, "blacklist.yml")
}

func loadSimpleRules(dir, filename string) ([]Rule, error) {
	path := filepath.Join(dir, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var container RuleContainer
	if err := yaml.Unmarshal(data, &container); err != nil {
		return nil, err
	}
	return container.Rules, nil
}

// LoadNormalizeMappings 加载归一化映射（假设格式为 mappings 列表）
func LoadNormalizeMappings(dir string) (map[string]string, error) {
	path := filepath.Join(dir, "normalize.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var raw struct {
		Mappings []struct {
			From string `yaml:"from"`
			To   string `yaml:"to"`
		} `yaml:"mappings"`
	}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	m := make(map[string]string)
	for _, pair := range raw.Mappings {
		m[pair.From] = pair.To
	}
	return m, nil
}

// LoadPolicy 加载 policy.yml 并构建 category -> (action, score) 映射
func LoadPolicy(policyPath string) (map[string]struct {
	Action string
	Score  int
}, error) {
	data, err := os.ReadFile(policyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // 策略文件可选
		}
		return nil, err
	}
	var config PolicyConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	// 构建反向映射：category -> (actionName, score)
	categoryMap := make(map[string]struct {
		Action string
		Score  int
	})
	for actionName, policy := range config.Actions {
		for _, cat := range policy.Categories {
			categoryMap[cat] = struct {
				Action string
				Score  int
			}{Action: actionName, Score: policy.Score}
		}
	}
	return categoryMap, nil
}
