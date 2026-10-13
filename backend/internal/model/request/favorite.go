package request

// AddFavoriteRequest 收藏请求
type AddFavoriteRequest struct {
	TargetID   uint   `json:"target_id" binding:"required"`    // 收藏对象ID
	TargetType string `json:"target_type"`                     // 收藏对象类型，默认 creation
	GroupID    int64  `json:"group_id"`                        // 收藏夹ID，默认 0（默认收藏夹）
}

// ListFavoritesRequest 我的收藏查询参数
type ListFavoritesRequest struct {
	PageRequest
	TargetType string `form:"target_type"` // 收藏对象类型，默认 creation
}
