// package nocode

// // NocodeMetadata 是前端零代码编辑器所需的全部节点元数据
// type NocodeMetadata struct {
// 	Triggers   []NodeMeta `json:"triggers"`
// 	Conditions []NodeMeta `json:"conditions"`
// 	Actions    []NodeMeta `json:"actions"`
// }

// // BuiltinConditionMetas 内置条件元数据
// var BuiltinConditionMetas = []NodeMeta{
// 	{
// 		Type: string(CondPostTitleContains), Label: "标题包含关键词", Icon: "type",
// 		Description: "帖子标题包含任意一个指定关键词",
// 		Params: []ParamMeta{
// 			{Key: "keywords", Label: "关键词列表", Type: "tags", Required: true},
// 		},
// 	},
// 	{
// 		Type: string(CondPostContentContains), Label: "正文包含关键词", Icon: "file-text",
// 		Description: "帖子或评论正文包含任意一个指定关键词",
// 		Params: []ParamMeta{
// 			{Key: "keywords", Label: "关键词列表", Type: "tags", Required: true},
// 		},
// 	},
// 	{
// 		Type: string(CondUserRoleIs), Label: "用户角色是", Icon: "shield",
// 		Params: []ParamMeta{
// 			{Key: "role", Label: "角色", Type: "select", Required: true,
// 				Options: []struct {
// 					Label string `json:"label"`
// 					Value any    `json:"value"`
// 				}{
// 					{Label: "普通用户", Value: "user"},
// 					{Label: "版主", Value: "moderator"},
// 					{Label: "管理员", Value: "admin"},
// 					{Label: "已封禁", Value: "banned"},
// 				},
// 			},
// 		},
// 	},
// 	{
// 		Type: string(CondUserPostCountGte), Label: "用户发帖数 ≥", Icon: "bar-chart-2",
// 		Params: []ParamMeta{
// 			{Key: "count", Label: "发帖数阈值", Type: "number", Required: true, Default: 10},
// 		},
// 	},
// 	{
// 		Type: string(CondSectionIDIn), Label: "所在板块是", Icon: "folder",
// 		Params: []ParamMeta{
// 			{Key: "ids", Label: "板块 ID 列表", Type: "tags", Required: true},
// 		},
// 	},
// 	{
// 		Type: string(CondTimeRange), Label: "当前时间在区间内", Icon: "clock",
// 		Params: []ParamMeta{
// 			{Key: "start", Label: "开始时间 (HH:mm)", Type: "text", Required: true, Placeholder: "09:00"},
// 			{Key: "end", Label: "结束时间 (HH:mm)", Type: "text", Required: true, Placeholder: "18:00"},
// 			{Key: "tz", Label: "时区", Type: "text", Default: "Asia/Shanghai"},
// 		},
// 	},
// 	{
// 		Type: string(CondCustomExpr), Label: "自定义表达式", Icon: "code",
// 		Description: `简单比较表达式，如：event.score > 80`,
// 		Params: []ParamMeta{
// 			{Key: "expr", Label: "表达式", Type: "text", Required: true,
// 				Placeholder: "event.score > 80"},
// 		},
// 	},
// }

package nocode
