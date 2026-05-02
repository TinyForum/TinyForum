package vo

import (
	"time"
	"tiny-forum/internal/model/do"
)

// LeaderboardItemResponse 排行榜条目响应
//
//	type LeaderboardItemVO struct {
//		ID       uint   `json:"id"`
//		Username string `json:"username"`
//		Avatar   string `json:"avatar"`
//		Score    int    `json:"score"`
//		Rank     int    `json:"rank"`
//	}
//
// SimpleLeaderboardItem 精简版（仅核心字段）
type SimpleLeaderboardItem struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Score    int    `json:"score"`
	Rank     int    `json:"rank"`
}

// type Statistics struct {
// 	TotalPosts          int64 `json:"total_posts"`
// 	TotalComments       int64 `json:"total_comments"`
// 	TotalFavorites      int64 `json:"total_favorites"`
// 	UnreadNotifications int64 `json:"unread_notifications"`
// }

type UserPosts struct {
	ID               uint                `json:"id"`                                          // 帖子ID
	CreatedAt        time.Time           `json:"created_at"`                                  // 创建时间
	UpdatedAt        time.Time           `json:"updated_at"`                                  // 更新时间
	Title            string              `gorm:"not null;size:200" json:"title"`              // 标题
	Summary          string              `gorm:"size:500" json:"summary"`                     // 摘要
	Cover            string              `gorm:"size:500" json:"cover"`                       // 封面
	Type             do.PostType         `gorm:"type:varchar(20);default:'post'" json:"type"` // 帖子类型
	PostStatus       do.PostStatus       `json:"post_status"`                                 // 文章状态：主动状态
	ModerationStatus do.ModerationStatus `json:"moderation_status"`                           // 审核状态：被动状态
	ViewCount        int                 `json:"view_count"`                                  // 浏览数
	LikeCount        int                 `json:"likes_count"`                                 // 点赞数
	CommentCount     int64               `json:"comment_count"`                               // 新增评论数
	PinTop           bool                `json:"pin_top"`                                     // 用户主页置顶
	Tags             []string            `json:"tags"`                                        // 标签列表
	BoardName        string              `gorm:"index" json:"board_name"`                     // 所属板块
	PinInBoard       bool                `gorm:"default:false" json:"pin_in_board"`           // 板块置顶
}
