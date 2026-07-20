# TinyForum Sensitive Dictionary

TinyForum 内容安全敏感词库。

基于 **规则 + 标签 + 风险等级 + 正则表达式** 的审核规则系统，用于论坛、社区、评论、私信等用户生成内容（UGC）的风险检测。

支持：

- YAML 配置化规则
- 正则匹配
- 分类管理
- 风险评分
- 自动拦截
- 人工审核
- 敏感词替换
- 规则热更新

---

# 目录结构

```bash
sensitive-dicts/
├── normalize.yml
├── whitelist.yml
├── blacklist.yml
├── README.md
├── adult/
│   ├── abuse.yml
│   ├── pornography.yml
│   └── sexual.yml
├── finance/
│   ├── business.yml
│   ├── crypto.yml
│   ├── fraud.yml
│   ├── gambling.yml
│   ├── investment.yml
│   └── stock.yml
├── health/
│   ├── addiction.yml
│   ├── drugs.yml
│   ├── efficacy.yml
│   └── medical.yml
├── marketing/
│   ├── advertisement.yml
│   ├── authority.yml
│   └── promotion.yml
├── religion/
│   ├── cultivation.yml
│   ├── fengshui.yml
│   ├── fortune.yml
│   ├── occult.yml
│   ├── religion.yml
│   └── superstition.yml
├── security/
│   ├── dangerous.yml
│   ├── self_harm.yml
│   ├── violence.yml
│   └── weapons.yml
├── social/
│   ├── funeral.yml
│   └── minor.yml         
├── society/
│   ├── contact.yml
│   ├── illegal.yml
│   ├── politics.yml
│   └── teenager.yml
├── spam/
│   ├── spam.yml
│   └── urls.yml
└── misc/                 
```

---

# 设计原则

## 1. 分类与策略分离

规则文件只负责：

- 内容分类
- 风险识别
- 标签定义

策略文件负责：

- 是否拦截
- 是否审核
- 是否替换

例如：

```text
pornography.yml

发现色情内容


↓

policy/block.yml

决定直接拦截
```

---

# 规则格式

所有规则统一采用 YAML：

```yaml
id:
name:
category:
priority:
action:
score:
regex:
tags:
replace:
description:
examples:
```

示例：

```yaml
version: 1

category: pornography

rules:
  - id: PORN_001

    name: 色情内容

    priority: 900

    action: block

    score: 100

    regex:
      - "色情"

      - "裸聊"

      - "成人视频"

    tags:
      - 色情

      - 成人内容

    replace: true

    description: 色情低俗内容

    examples:
      - 裸聊网站

      - 成人视频
```

---

# Rule 数据结构

Go 定义：

```go
package sensitive


type Rule struct {

	ID string `yaml:"id"`

	Name string `yaml:"name"`

	Category string `yaml:"category"`

	Priority int `yaml:"priority"`


	Action Action `yaml:"action"`

	Score int `yaml:"score"`


	Regex []string `yaml:"regex"`


	Tags []string `yaml:"tags"`


	Replace bool `yaml:"replace"`


	Description string `yaml:"description"`


	Examples []string `yaml:"examples"`
}
```

---

# MatchResult

规则命中结果：

```go
type MatchResult struct {

	RuleID string


	Category string


	Label string


	Word string


	Start int


	End int


	Action Action


	Score int


	Tags []string
}
```

字段说明：

| 字段     | 说明     |
| -------- | -------- |
| RuleID   | 规则编号 |
| Category | 风险分类 |
| Word     | 命中内容 |
| Start    | 开始位置 |
| End      | 结束位置 |
| Score    | 风险分数 |
| Action   | 处理动作 |
| Tags     | 风险标签 |

---

# CheckResult

最终审核结果：

```go
type CheckResult struct {


	Original string


	Masked string


	Sensitive bool


	Level Level


	Score int


	Action Action


	Matches []*MatchResult
}
```

示例：

输入：

```text
添加微信领取福利
```

输出：

```json
{
  "sensitive": true,

  "masked": "添加****领取福利",

  "score": 40,

  "action": "review",

  "matches": [
    {
      "category": "contact",
      "word": "微信",
      "score": 40
    }
  ]
}
```

---

# Action 类型

审核动作：

```text
pass

↓

review

↓

block
```

说明：

| Action | 说明         |
| ------ | ------------ |
| pass   | 正常通过     |
| review | 进入人工审核 |
| block  | 直接拒绝     |

---

# 风险等级

根据 Score 计算：

| Score  | 等级     | 处理     |
| ------ | -------- | -------- |
| 0-20   | Normal   | 通过     |
| 21-50  | Low      | 记录     |
| 51-80  | Medium   | 审核     |
| 81-100 | High     | 拦截     |
| 100+   | Critical | 强制拦截 |

---

# 审核流程

```
                    Text
                      |
                      v

              normalize.go

                      |
                      v

              regexp matcher

                      |
                      v


             +----------------+

             | MatchResult[]  |

             +----------------+

                      |
                      v

             risk calculator

                      |
                      v


        +-------------+-------------+

        |                           |

      block                      review

        |                           |

    reject API              admin queue

```

---

# 文本标准化

normalize.yml 用于处理：

- 大小写转换
- 繁简转换
- 全角半角转换
- 特殊字符清理
- Unicode 归一化

例如：

输入：

```
微 信
微-信
VX
ＶＸ
```

统一：

```
微信
```

---

# 白名单机制

whitelist.yml

用于降低误伤：

例如：

```yaml
rules:
  - word: 微信支付

    action: pass

  - word: 赌博小说

    action: pass
```

---

# 黑名单机制

blacklist.yml

高风险强制规则：

```yaml
rules:
  - word: xxx

    action: block

    score: 100
```

---

# 性能优化

推荐：

## 规则加载

启动时：

```
YAML

↓

Rule Loader

↓

Regex Cache

↓

Matcher
```

## 匹配优化

规则数量较大时：

推荐：

- Trie 前缀树
- Aho-Corasick 自动机
- Regex 编译缓存

流程：

```
Rule

 |

 +-- exact words

 |

 +-- regex

 |

 +-- keywords index

```

---

# 版本管理

建议：

```
sensitive-dicts

v1.0

v1.1

v1.2
```

每次修改：

记录：

- 修改人
- 修改时间
- 新增规则
- 删除规则
- 命中变化

---

# 使用场景

适用于：

- 社区论坛
- 评论系统
- 私信系统
- 博客平台
- 内容管理系统
- AI 对话系统

---


# 测试
```bash
curl -X POST http://localhost:8080/check -H "Content-Type: application/json" -d '{"text":"自从Go 1.18正式引入泛型后，我们团队一直在探索如何在实际项目中合理使用这一特性。本文分享了我们在重构数据结构和算法库时的经验，包括类型约束的设计、性能影响分析，以及常见陷阱的规避。通过三个 真实案例，展示了泛型如何提升代码复用性和类型安全性。"}'
```

# License

TinyForum Sensitive Dictionary

内部内容安全规则库。
