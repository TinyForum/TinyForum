package nocode

// BuiltinTriggers 内置触发器元数据
var BuiltinTriggers = []NodeMeta{
	{Type: string(TriggerOnSchedule), Label: "定时触发", Icon: "clock",
		Description: "按 Cron 表达式定时执行",
		Params: []ParamMeta{
			{Key: "cron", Label: "Cron 表达式", Type: "cron", Required: true,
				Placeholder: "0 9 * * 1（每周一 9:00）"},
		},
	},
	{Type: string(TriggerOnNewPost), Label: "新帖触发", Icon: "file-text",
		Description: "有新帖子发布时触发",
		Params: []ParamMeta{
			{Key: "board_ids", Label: "板块（空=全部）", Type: "tags"},
		},
	},
	{Type: string(TriggerOnNewComment), Label: "新评论触发", Icon: "message-circle",
		Description: "有新评论时触发"},
	{Type: string(TriggerOnUserRegister), Label: "新用户注册", Icon: "user-plus",
		Description: "新用户完成注册时触发"},
	{Type: string(TriggerOnKeyword), Label: "关键词触发", Icon: "search",
		Description: "帖子或评论包含关键词时触发",
		Params: []ParamMeta{
			{Key: "keywords", Label: "关键词列表", Type: "tags", Required: true},
			{Key: "scope", Label: "检测范围", Type: "select", Required: true, Default: "both",
				Options: []OptionMeta{
					{Label: "帖子", Value: "post"},
					{Label: "评论", Value: "comment"},
					{Label: "全部", Value: "both"},
				}},
		},
	},
	{Type: string(TriggerOnManual), Label: "手动触发", Icon: "play",
		Description: "仅通过 API 手动触发"},
}

// BuiltinConditions 内置条件元数据
var BuiltinConditions = []NodeMeta{
	{Type: string(CondPostTitleContains), Label: "标题含关键词", Icon: "type",
		Params: []ParamMeta{{Key: "keywords", Label: "关键词列表", Type: "tags", Required: true}}},
	{Type: string(CondPostContentContains), Label: "正文含关键词", Icon: "file-text",
		Params: []ParamMeta{{Key: "keywords", Label: "关键词列表", Type: "tags", Required: true}}},
	{Type: string(CondUserRoleIs), Label: "用户角色是", Icon: "shield",
		Params: []ParamMeta{
			{Key: "role", Label: "角色", Type: "select", Required: true,
				Options: []OptionMeta{
					{Label: "普通用户", Value: "user"},
					{Label: "版主", Value: "moderator"},
					{Label: "管理员", Value: "admin"},
				}},
		},
	},
	{Type: string(CondUserPostCountGte), Label: "用户发帖数 ≥", Icon: "bar-chart-2",
		Params: []ParamMeta{{Key: "count", Label: "发帖数阈值", Type: "number", Required: true, Default: 10}}},
	{Type: string(CondBoardIDIn), Label: "所在板块是", Icon: "folder",
		Params: []ParamMeta{{Key: "ids", Label: "板块 ID 列表", Type: "tags", Required: true}}},
	{Type: string(CondTimeRange), Label: "时间在区间内", Icon: "clock",
		Params: []ParamMeta{
			{Key: "start", Label: "开始时间 (HH:mm)", Type: "text", Required: true, Placeholder: "09:00"},
			{Key: "end", Label: "结束时间 (HH:mm)", Type: "text", Required: true, Placeholder: "18:00"},
			{Key: "tz", Label: "时区", Type: "text", Default: "Asia/Shanghai"},
		},
	},
	{Type: string(CondCustomExpr), Label: "自定义表达式", Icon: "code",
		Description: "支持简单比较：event.score > 80",
		Params: []ParamMeta{
			{Key: "expr", Label: "表达式", Type: "text", Required: true, Placeholder: "event.score > 80"},
		},
	},
}

// BuiltinActions 内置动作元数据
var BuiltinActions = []NodeMeta{
	{Type: string(ActionReplyPost), Label: "回复帖子", Icon: "reply", Category: "post",
		Params: []ParamMeta{
			{Key: "content", Label: "回复内容", Type: "textarea", Required: true,
				Placeholder: "支持模板变量：{{.username}}"},
		},
	},
	{Type: string(ActionDeletePost), Label: "删除帖子", Icon: "trash-2", Category: "post"},
	{Type: string(ActionHidePost), Label: "隐藏帖子", Icon: "eye-off", Category: "post"},
	{Type: string(ActionPinPost), Label: "置顶帖子", Icon: "pin", Category: "post"},
	{Type: string(ActionLockPost), Label: "锁定帖子", Icon: "lock", Category: "post"},
	{Type: string(ActionCreatePost), Label: "发布新帖", Icon: "plus-square", Category: "post",
		Params: []ParamMeta{
			{Key: "board_id", Label: "板块 ID", Type: "number", Required: true},
			{Key: "title", Label: "标题", Type: "text", Required: true},
			{Key: "content", Label: "正文", Type: "textarea", Required: true},
		},
	},
	{Type: string(ActionDeleteComment), Label: "删除评论", Icon: "x-circle", Category: "comment"},
	{Type: string(ActionBanUser), Label: "封禁用户", Icon: "user-x", Category: "user",
		Params: []ParamMeta{
			{Key: "reason", Label: "封禁原因", Type: "text", Required: true},
			{Key: "duration_sec", Label: "时长（秒）", Type: "number", Default: 86400},
		},
	},
	{Type: string(ActionSendMessage), Label: "发送私信", Icon: "mail", Category: "user",
		Params: []ParamMeta{
			{Key: "to_user_id", Label: "接收用户 ID（空=触发者）", Type: "number"},
			{Key: "content", Label: "消息内容", Type: "textarea", Required: true},
		},
	},
	{Type: string(ActionWebhook), Label: "调用 Webhook", Icon: "link", Category: "integration",
		Params: []ParamMeta{
			{Key: "url", Label: "Webhook URL", Type: "text", Required: true},
			{Key: "method", Label: "HTTP 方法", Type: "select", Default: "POST",
				Options: []OptionMeta{{Label: "POST", Value: "POST"}, {Label: "GET", Value: "GET"}}},
			{Key: "body", Label: "请求体（支持模板）", Type: "textarea"},
			{Key: "headers", Label: "请求头 JSON", Type: "textarea"},
		},
	},
	{Type: string(ActionNotifyAdmin), Label: "通知管理员", Icon: "bell", Category: "integration",
		Params: []ParamMeta{
			{Key: "message", Label: "通知内容", Type: "textarea", Required: true},
		},
	},
	{Type: string(ActionWait), Label: "等待", Icon: "pause", Category: "control",
		Params: []ParamMeta{
			{Key: "seconds", Label: "等待秒数", Type: "number", Required: true, Default: 1},
		},
	},
	{Type: string(ActionSetVariable), Label: "设置变量", Icon: "edit-3", Category: "control",
		Params: []ParamMeta{
			{Key: "name", Label: "变量名", Type: "text", Required: true},
			{Key: "value", Label: "值（支持模板）", Type: "text", Required: true},
		},
	},
	{Type: string(ActionStopIf), Label: "条件停止", Icon: "stop-circle", Category: "control",
		Params: []ParamMeta{
			{Key: "expr", Label: "停止条件", Type: "text", Required: true, Placeholder: "event.score > 80"},
		},
	},
}
