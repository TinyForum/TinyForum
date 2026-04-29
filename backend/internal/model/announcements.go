package model

import (
	"time"
)

// Announcement 公告模型
type Announcement struct {
	BaseModel
	Title       string             `gorm:"not null;size:200" json:"title"`                 // 公告标题
	Content     string             `gorm:"type:text;not null" json:"content"`              // 公告内容
	Summary     string             `gorm:"size:500" json:"summary"`                        // 公告摘要
	Cover       string             `gorm:"size:500" json:"cover"`                          // 封面图
	Type        AnnouncementType   `gorm:"type:varchar(20);default:'normal'" json:"type"`  // 公告类型
	Status      AnnouncementStatus `gorm:"type:varchar(20);default:'draft'" json:"status"` // 状态
	IsPinned    bool               `gorm:"default:false;index" json:"is_pinned"`           // 是否置顶
	IsGlobal    bool               `gorm:"default:true;index" json:"is_global"`            // 是否全局公告
	BoardID     *uint              `gorm:"index;default:null" json:"board_id"`             // 关联的板块ID（为空则全局）
	PublishedAt *time.Time         `json:"published_at"`                                   // 发布时间
	ExpiredAt   *time.Time         `json:"expired_at"`                                     // 过期时间
	ViewCount   int                `gorm:"default:0" json:"view_count"`                    // 浏览次数
	CreatedBy   uint               `gorm:"not null;index" json:"created_by"`               // 创建人ID
	UpdatedBy   uint               `json:"updated_by"`                                     // 更新人ID

	// 关联关系
	Board   *Board `gorm:"foreignKey:BoardID" json:"board,omitempty" swaggerignore:"true"`
	Creator *User  `gorm:"foreignKey:CreatedBy" json:"creator,omitempty" swaggerignore:"true"`
}

// AnnouncementType 公告类型
type AnnouncementType string

const (
	AnnouncementTypeNormal    AnnouncementType = "normal"    // 普通公告
	AnnouncementTypeImportant AnnouncementType = "important" // 重要公告
	AnnouncementTypeEmergency AnnouncementType = "emergency" // 紧急公告
	AnnouncementTypeEvent     AnnouncementType = "event"     // 活动公告
)

// AnnouncementStatus 公告状态
type AnnouncementStatus string

const (
	AnnouncementStatusAll       AnnouncementStatus = "all"       // 所有
	AnnouncementStatusDraft     AnnouncementStatus = "draft"     // 草稿
	AnnouncementStatusPublished AnnouncementStatus = "published" // 已发布
	AnnouncementStatusArchived  AnnouncementStatus = "archived"  // 已归档
)

// TableName 指定表名
func (Announcement) TableName() string {
	return "announcements"
}
