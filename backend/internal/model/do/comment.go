package do

import "tiny-forum/internal/model/common"

// type CommentStatus string

// const (
// 	CommentStatusVisible CommentStatus = "visible" // 正常可见
// 	CommentStatusPending CommentStatus = "pending" // 待审核
// 	CommentStatusHidden  CommentStatus = "hidden"  // 已隐藏
// )

//	type Comment struct {
//		common.BaseModel
//		WorksID  uint   `gorm:"index" json:"works_id"`                     // 作品ID，Article、Posts
//		AuthorID uint   `gorm:"index" json:"author_id"`                    // 评论作者ID
//		ReplyID  uint   `gorm:"index" json:"reply_id"`                     // 回复ID
//		Reply    *Reply `gorm:"foreignKey:ReplyID" json:"reply,omitempty"` // 回复
//	}
//
// Comment 作品（文章/帖子）评论场景下的回复扩展（依附于 Reply）

// ==================== 评论（作品/回答下的评论，支持嵌套） ====================

// Comment 评论扩展表（依附于 Reply，TargetType = "post" 或 "answer"）
type Comment struct {
	common.BaseModel
	WorksID uint `gorm:"not null;index" json:"works_id"`       // 作品ID，Article、Posts
	ReplyID uint `gorm:"not null;uniqueIndex" json:"reply_id"` // 关联的回复ID（一对一）

	// 关联回复内容
	Reply *Reply `gorm:"foreignKey:ReplyID;references:ID" json:"reply,omitempty"`
	// ---------- 社区互动字段 ----------
	IsPinned     bool `gorm:"default:false;index" json:"is_pinned"`    // 是否置顶
	IsAnonymous  bool `gorm:"default:false;index" json:"is_anonymous"` // 是否匿名（隐藏作者）
	DislikeCount int  `gorm:"default:0" json:"dislike_count"`          // 点踩/反对数
	ReportCount  int  `gorm:"default:0" json:"report_count"`           // 举报计数（用于风控）
	// ---------- 展示与排序辅助 ----------
	SortWeight int    `gorm:"default:0;index" json:"sort_weight"`   // 排序权重（可组合：置顶*1000 + 点赞 - 反对 + 时间衰减）
	IpLocation string `gorm:"size:50" json:"ip_location,omitempty"` // IP属地（如“中国 北京”）
}
