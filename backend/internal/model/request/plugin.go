package request

type ListPluginsRequest struct {
	PageRequest
	Keyword          string `form:"keyword"`           // 关键字
	Status           string `form:"status"`            // 插件状态
	ModerationStatus string `form:"moderation_status"` // 风控状态审核结果 (normal, pending, rejected)
	Tag              string `form:"tag"`               // 标签名称
}
