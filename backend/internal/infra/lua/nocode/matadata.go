package nocode

// NocodeMetadata 包含所有内置节点的元数据，通过 API 下发给前端。
type NocodeMetadata struct {
	Triggers  []NodeMeta `json:"triggers"`
	Control   []NodeMeta `json:"control"`
	Variables []NodeMeta `json:"variables"`
	Actions   []NodeMeta `json:"actions"`
}

type NodeMeta struct {
	Type        string      `json:"type"`
	Label       string      `json:"label"`
	Description string      `json:"description,omitempty"`
	Icon        string      `json:"icon,omitempty"`
	Category    string      `json:"category,omitempty"`
	Params      []ParamMeta `json:"params,omitempty"`
	Outputs     []VarOutput `json:"outputs,omitempty"` // 该节点产生的变量
}

type ParamMeta struct {
	Key         string       `json:"key"`
	Label       string       `json:"label"`
	Type        string       `json:"type"` // text|number|boolean|select|textarea|tags|cron
	Required    bool         `json:"required"`
	Default     interface{}  `json:"default,omitempty"`
	Placeholder string       `json:"placeholder,omitempty"`
	Options     []OptionMeta `json:"options,omitempty"`
}

type OptionMeta struct {
	Label string `json:"label"`
	Value any    `json:"value"`
}
