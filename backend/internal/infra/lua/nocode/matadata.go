package nocode

// NocodeMetadata 包含所有内置节点的元数据，通过 API 下发给前端。
type NocodeMetadata struct {
	Triggers   []NodeMeta `json:"triggers"`
	Conditions []NodeMeta `json:"conditions"`
	Actions    []NodeMeta `json:"actions"`
}

// NodeMeta 描述一种节点类型的展示信息和参数定义
type NodeMeta struct {
	Type        string      `json:"type"`
	Label       string      `json:"label"`
	Description string      `json:"description,omitempty"`
	Icon        string      `json:"icon,omitempty"`
	Category    string      `json:"category,omitempty"`
	Params      []ParamMeta `json:"params,omitempty"`
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
