package do

import "tiny-forum/internal/model/common"

// ==================== 通用回复内容表（多态） ====================

// ReplyStatus 回复状态
type ReplyStatus string

const (
	ReplyStatusVisible ReplyStatus = "visible" // 正常可见
	ReplyStatusPending ReplyStatus = "pending" // 待审核
	ReplyStatusHidden  ReplyStatus = "hidden"  // 已隐藏
)

// TargetType 目标类型常量（避免魔法字符串）
const (
	TargetReplyTypeQuestion = "question" // 对应问题（即回答）
	TargetReplyTypePost     = "post"     // 对应作品/帖子（即根评论）
	TargetReplyTypeAnswer   = "answer"   // 对应某个回答（即回答的子评论）
)

// Reply 通用回复实体（统一管理内容和层级）
type Reply struct {
	common.BaseModel
	Content    string      `gorm:"not null;type:text" json:"content"`                      // 回复内容
	AuthorID   uint        `gorm:"not null;index" json:"author_id"`                        // 回复者ID
	TargetType string      `gorm:"type:varchar(20);not null;index" json:"target_type"`     // 目标类型: question / post / answer
	TargetID   uint        `gorm:"not null;index" json:"target_id"`                        // 目标实体ID
	ParentID   *uint       `gorm:"index" json:"parent_id"`                                 // 父回复ID（仅评论支持嵌套，回答必须为 NULL）
	LikeCount  int         `gorm:"default:0" json:"like_count"`                            // 点赞数
	Status     ReplyStatus `gorm:"type:varchar(20);default:'visible';index" json:"status"` // 状态

	// 关联（只做数据加载，不参与业务外键约束逻辑）
	Author  User    `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Parent  *Reply  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Replies []Reply `gorm:"foreignKey:ParentID" json:"replies,omitempty"` // 子回复列表（仅评论使用）
}
