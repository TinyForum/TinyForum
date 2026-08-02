package nocode

// ─── 触发器 ──────────────────────────────────────────────────────────

var BuiltinTriggers = []NodeMeta{
	{Type: string(TriggerOnSchedule), Label: "定时触发", Icon: "clock",
		Description: "按 Cron 表达式定时执行",
		Params: []ParamMeta{
			{Key: "cron", Label: "Cron 表达式", Type: "cron", Required: true, Placeholder: "0 9 * * 1"},
		},
	},
	{Type: string(TriggerOnNewPost), Label: "新帖触发", Icon: "file-text",
		Description: "有新帖子发布时触发",
		Params: []ParamMeta{
			{Key: "board_ids", Label: "板块（空=全部）", Type: "tags"},
		},
	},
	{Type: string(TriggerOnNewComment), Label: "新评论触发", Icon: "message-circle",
		Description: "有新评论时触发",
	},
	{Type: string(TriggerOnUserRegister), Label: "新用户注册", Icon: "user-plus",
		Description: "新用户完成注册时触发",
	},
	{Type: string(TriggerOnKeyword), Label: "关键词触发", Icon: "search",
		Description: "帖子或评论包含关键词时触发",
		Params: []ParamMeta{
			{Key: "keywords", Label: "关键词", Type: "tags", Required: true},
			{Key: "scope", Label: "范围", Type: "select", Required: true, Default: "both",
				Options: []OptionMeta{
					{Label: "帖子", Value: "post"}, {Label: "评论", Value: "comment"}, {Label: "全部", Value: "both"},
				}},
		},
	},
	{Type: string(TriggerOnManual), Label: "手动触发", Icon: "play",
		Description: "仅通过 API 手动触发",
	},
}

// ─── 控制 ────────────────────────────────────────────────────────────

var BuiltinControl = []NodeMeta{
	{Type: "if", Label: "条件分支 (if/else)", Icon: "git-branch", Category: "flow",
		Description: "条件为真执行 true 路径，否则走 false 路径",
		Params: []ParamMeta{
			{Key: "condition_field", Label: "检查字段", Type: "select", Required: true,
				Options: fieldOptions(),
			},
			{Key: "condition_op", Label: "运算符", Type: "select", Required: true, Default: "equals",
				Options: []OptionMeta{
					{Label: "等于", Value: "equals"},
					{Label: "不等于", Value: "not_equals"},
					{Label: "包含", Value: "contains"},
					{Label: "大于", Value: "greater_than"},
					{Label: "小于", Value: "less_than"},
					{Label: "为空", Value: "is_empty"},
					{Label: "非空", Value: "is_not_empty"},
				},
			},
			{Key: "condition_value", Label: "比较值", Type: "text", Placeholder: "支持模板 {{var}}"},
		},
	},
	{Type: "while", Label: "循环 (while)", Icon: "refresh-cw", Category: "flow",
		Description: "条件为真时反复执行 body 内的步骤",
		Params: []ParamMeta{
			{Key: "condition_field", Label: "检查字段", Type: "select", Required: true,
				Options: fieldOptions(),
			},
			{Key: "condition_op", Label: "运算符", Type: "select", Required: true, Default: "less_than",
				Options: []OptionMeta{
					{Label: "小于", Value: "less_than"},
					{Label: "大于", Value: "greater_than"},
					{Label: "不等于", Value: "not_equals"},
				},
			},
			{Key: "condition_value", Label: "比较值", Type: "text", Required: true, Placeholder: "支持 {{var}}"},
			{Key: "max_iter", Label: "最大迭代", Type: "number", Default: 100},
		},
	},
	{Type: "wait", Label: "等待", Icon: "pause", Category: "control",
		Description: "暂停执行指定秒数",
		Params: []ParamMeta{
			{Key: "seconds", Label: "秒数", Type: "number", Required: true, Default: 1},
		},
	},
	{Type: "stop", Label: "终止流程", Icon: "stop-circle", Category: "control",
		Description: "提前结束整个流程",
	},
}

// ─── 变量 ────────────────────────────────────────────────────────────

var BuiltinVariables = []NodeMeta{
	{Type: "set_variable", Label: "设置变量", Icon: "edit-3", Category: "variable",
		Description: "设置一个变量，后续节点可通过 {{name}} 引用",
		Params: []ParamMeta{
			{Key: "name", Label: "变量名", Type: "text", Required: true},
			{Key: "value", Label: "值", Type: "text", Required: true, Placeholder: "支持模板 {{var}}"},
		},
		Outputs: []VarOutput{{Name: "__custom__", Type: "string", Desc: "用户自定义变量名"}},
	},

	// ── 数据获取 ──
	{Type: "get_post_info", Label: "获取帖子信息", Icon: "file-text", Category: "data",
		Description: "从当前事件中提取帖子字段并存入变量",
		Outputs: []VarOutput{
			{Name: "post_id", Type: "number", Desc: "帖子ID"},
			{Name: "post_title", Type: "string", Desc: "帖子标题"},
			{Name: "post_content", Type: "string", Desc: "帖子正文"},
			{Name: "post_author_id", Type: "number", Desc: "作者ID"},
			{Name: "post_author_name", Type: "string", Desc: "作者名"},
			{Name: "board_id", Type: "number", Desc: "板块ID"},
		},
	},
	{Type: "get_user_info", Label: "获取用户信息", Icon: "user", Category: "data",
		Description: "提取当前事件的用户字段到变量",
		Outputs: []VarOutput{
			{Name: "user_id", Type: "number", Desc: "用户ID"},
			{Name: "username", Type: "string", Desc: "用户名"},
			{Name: "user_role", Type: "string", Desc: "用户角色"},
			{Name: "user_post_count", Type: "number", Desc: "发帖数"},
		},
	},
	{Type: "get_comment_info", Label: "获取评论信息", Icon: "message-circle", Category: "data",
		Description: "提取当前事件的评论字段到变量",
		Outputs: []VarOutput{
			{Name: "comment_id", Type: "number", Desc: "评论ID"},
			{Name: "comment_content", Type: "string", Desc: "评论内容"},
			{Name: "comment_author_id", Type: "number", Desc: "评论者ID"},
		},
	},

	// ── 算术 ──
	{Type: "add", Label: "加法", Icon: "plus", Category: "math",
		Description: "a + b，结果存入目标变量",
		Params: []ParamMeta{
			{Key: "target", Label: "结果变量", Type: "text", Required: true},
			{Key: "a", Label: "a", Type: "text", Required: true, Placeholder: "数字或 {{var}}"},
			{Key: "b", Label: "b", Type: "text", Required: true, Placeholder: "数字或 {{var}}"},
		},
		Outputs: []VarOutput{{Name: "__custom__", Type: "number", Desc: "运算结果"}},
	},
	{Type: "subtract", Label: "减法", Icon: "minus", Category: "math",
		Params: []ParamMeta{
			{Key: "target", Label: "结果变量", Type: "text", Required: true},
			{Key: "a", Label: "a", Type: "text", Required: true},
			{Key: "b", Label: "b", Type: "text", Required: true},
		},
		Outputs: []VarOutput{{Name: "__custom__", Type: "number", Desc: "运算结果"}},
	},
	{Type: "multiply", Label: "乘法", Icon: "x", Category: "math",
		Params: []ParamMeta{
			{Key: "target", Label: "结果变量", Type: "text", Required: true},
			{Key: "a", Label: "a", Type: "text", Required: true},
			{Key: "b", Label: "b", Type: "text", Required: true},
		},
		Outputs: []VarOutput{{Name: "__custom__", Type: "number", Desc: "运算结果"}},
	},
	{Type: "divide", Label: "除法", Icon: "divide", Category: "math",
		Params: []ParamMeta{
			{Key: "target", Label: "结果变量", Type: "text", Required: true},
			{Key: "a", Label: "a", Type: "text", Required: true},
			{Key: "b", Label: "b", Type: "text", Required: true},
		},
		Outputs: []VarOutput{{Name: "__custom__", Type: "number", Desc: "运算结果"}},
	},
	{Type: "modulo", Label: "取余", Icon: "percent", Category: "math",
		Params: []ParamMeta{
			{Key: "target", Label: "结果变量", Type: "text", Required: true},
			{Key: "a", Label: "a", Type: "text", Required: true},
			{Key: "b", Label: "b", Type: "text", Required: true},
		},
		Outputs: []VarOutput{{Name: "__custom__", Type: "number", Desc: "运算结果"}},
	},
	{Type: "concat", Label: "字符串拼接", Icon: "align-left", Category: "string",
		Params: []ParamMeta{
			{Key: "target", Label: "结果变量", Type: "text", Required: true},
			{Key: "a", Label: "字符串 a", Type: "text", Required: true, Placeholder: "支持 {{var}}"},
			{Key: "b", Label: "字符串 b", Type: "text", Required: true, Placeholder: "支持 {{var}}"},
		},
		Outputs: []VarOutput{{Name: "__custom__", Type: "string", Desc: "拼接结果"}},
	},
	{Type: "length", Label: "取长度", Icon: "ruler", Category: "string",
		Params: []ParamMeta{
			{Key: "target", Label: "结果变量", Type: "text", Required: true},
			{Key: "source", Label: "源字符串", Type: "text", Required: true, Placeholder: "支持 {{var}}"},
		},
		Outputs: []VarOutput{{Name: "__custom__", Type: "number", Desc: "字符串长度"}},
	},
}

// ─── 动作 ────────────────────────────────────────────────────────────

var BuiltinActions = []NodeMeta{
	// 帖子
	{Type: "reply_post", Label: "回复帖子", Icon: "reply", Category: "post",
		Params: []ParamMeta{
			{Key: "content", Label: "回复内容", Type: "textarea", Required: true, Placeholder: "支持模板 {{username}}"},
		},
	},
	{Type: "delete_post", Label: "删除帖子", Icon: "trash-2", Category: "post"},
	{Type: "hide_post", Label: "隐藏帖子", Icon: "eye-off", Category: "post"},
	{Type: "pin_post", Label: "置顶帖子", Icon: "pin", Category: "post"},
	{Type: "lock_post", Label: "锁定帖子", Icon: "lock", Category: "post"},
	{Type: "create_post", Label: "发布帖子", Icon: "plus-square", Category: "post",
		Params: []ParamMeta{
			{Key: "board_id", Label: "板块ID", Type: "number", Required: true},
			{Key: "title", Label: "标题", Type: "text", Required: true},
			{Key: "content", Label: "正文", Type: "textarea", Required: true},
		},
	},
	// 评论
	{Type: "delete_comment", Label: "删除评论", Icon: "x-circle", Category: "comment"},
	// 用户
	{Type: "ban_user", Label: "封禁用户", Icon: "user-x", Category: "user",
		Params: []ParamMeta{
			{Key: "reason", Label: "原因", Type: "text", Required: true},
			{Key: "duration_sec", Label: "时长（秒）", Type: "number", Default: 86400},
		},
	},
	{Type: "send_message", Label: "发送私信", Icon: "mail", Category: "user",
		Params: []ParamMeta{
			{Key: "to_user_id", Label: "接收者ID", Type: "number", Placeholder: "空=触发者"},
			{Key: "content", Label: "消息内容", Type: "textarea", Required: true},
		},
	},
	// 集成
	{Type: "webhook", Label: "调用 Webhook", Icon: "link", Category: "integration",
		Params: []ParamMeta{
			{Key: "url", Label: "URL", Type: "text", Required: true},
			{Key: "method", Label: "方法", Type: "select", Default: "POST",
				Options: []OptionMeta{{Label: "POST", Value: "POST"}, {Label: "GET", Value: "GET"}}},
			{Key: "body", Label: "请求体", Type: "textarea"},
		},
	},
	{Type: "notify_admin", Label: "通知管理员", Icon: "bell", Category: "integration",
		Params: []ParamMeta{
			{Key: "message", Label: "通知内容", Type: "textarea", Required: true},
		},
	},
}

func fieldOptions() []OptionMeta {
	return []OptionMeta{
		{Label: "帖子ID (post_id)", Value: "post_id"},
		{Label: "帖子标题 (post_title)", Value: "post_title"},
		{Label: "帖子正文 (post_content)", Value: "post_content"},
		{Label: "作者ID (author_id)", Value: "author_id"},
		{Label: "用户名 (username)", Value: "username"},
		{Label: "板块ID (board_id)", Value: "board_id"},
		{Label: "用户角色 (user_role)", Value: "user_role"},
		{Label: "发帖数 (user_post_count)", Value: "user_post_count"},
		{Label: "评论ID (comment_id)", Value: "comment_id"},
		{Label: "评论内容 (comment_content)", Value: "comment_content"},
		{Label: "自定义变量", Value: "__custom__"},
	}
}
