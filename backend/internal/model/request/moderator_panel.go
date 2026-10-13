package request

// ListModeratorPostsRequest 版主面板帖子列表查询参数
type ListModeratorPostsRequest struct {
	PageRequest
	Keyword string `form:"keyword"` // 关键词（匹配标题/内容）
	Status  string `form:"status"`  // 帖子状态：pending/published/deleted/hidden，留空为全部
}

// ListModeratorReportsRequest 版主面板举报列表查询参数
type ListModeratorReportsRequest struct {
	PageRequest
	Status string `form:"status"` // 举报状态：pending/resolved/rejected，留空为全部
}

// ListModeratorBansRequest 版主面板封禁列表查询参数
type ListModeratorBansRequest struct {
	PageRequest
}