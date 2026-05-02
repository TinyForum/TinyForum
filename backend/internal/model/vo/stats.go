package vo

type UserStatsInfo struct {
	TotalPost    int64 `json:"total_post"`    // 总帖子数
	TotalComment int64 `json:"total_comment"` // 总评论数
	// TotalBoard   int64 `json:"total_board"`   // 总板块数
	// TotalTag     int64 `json:"total_tag"`     // 总标签数
	TotalFavorite  int64 `json:"total_favorites"` // 总收藏数
	TotalLike      int64 `json:"total_like"`      // 总点赞数
	TotalFollower  int64 `json:"total_follower"`  // 总关注数
	TotalFollowing int64 `json:"total_following"` // 总粉丝数
	TotalReport    int64 `json:"total_report"`    // 总举报数
	TotalViolation int64 `json:"total_violation"` // 总违规数
	TotalQuestion  int64 `json:"total_question"`  // 总问题数
	TotalAnswer    int64 `json:"total_answer"`    // 总回答数
	TotalUpload    int64 `json:"total_upload"`    // 总上传数
	TotalScore     int   `json:"total_score"`     // 总积分
}
