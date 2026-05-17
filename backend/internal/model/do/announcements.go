package do

import (
	"time"
	"tiny-forum/internal/model/common"
)

// Announcement 公告模型
type Announcement struct {
	common.BaseModel
	Title       string              `gorm:"not null;size:200" json:"title"`           // 公告标题
	Content     string              `gorm:"type:text;not null" json:"content"`        // 公告内容
	Summary     string              `gorm:"size:500" json:"summary"`                  // 公告摘要
	Cover       string              `gorm:"size:500" json:"cover"`                    // 封面图
	Type        *AnnouncementType   `gorm:"type:varchar;default:normal" json:"type"`  // 公告类型
	Status      *AnnouncementStatus `gorm:"type:varchar;default:draft" json:"status"` // 公告状态
	IsPinned    bool                `gorm:"default:false;index" json:"is_pinned"`     // 是否置顶
	IsGlobal    bool                `gorm:"default:true;index" json:"is_global"`      // 是否全局公告
	BoardID     *uint               `gorm:"index;default:null" json:"board_id"`       // 关联的板块ID（为空则全局）
	PublishedAt *time.Time          `json:"published_at"`                             // 发布时间
	ExpiredAt   *time.Time          `json:"expired_at"`                               // 过期时间
	ViewCount   int                 `gorm:"default:0" json:"view_count"`              // 浏览次数
	CreatedBy   uint                `gorm:"not null;index" json:"created_by"`         // 创建人ID
	UpdatedBy   uint                `json:"updated_by"`                               // 更新人ID

	// 关联关系
	Board   *Board `gorm:"foreignKey:BoardID" json:"board,omitempty" swaggerignore:"true"`
	Creator *User  `gorm:"foreignKey:CreatedBy" json:"creator,omitempty" swaggerignore:"true"`
}

// TableName 指定表名
func (Announcement) TableName() string {
	return "announcements"
}

// AnnouncementType 公告类型（可存储）
type AnnouncementType string

const (
	AnnouncementTypeNormal    AnnouncementType = "normal"    // 普通公告
	AnnouncementTypeImportant AnnouncementType = "important" // 重要公告
	AnnouncementTypeEmergency AnnouncementType = "emergency" // 紧急公告
	AnnouncementTypeEvent     AnnouncementType = "event"     // 活动公告
)

// IsValid 检查公告类型是否合法
func (t AnnouncementType) IsValid() bool {
	return t >= AnnouncementTypeNormal && t <= AnnouncementTypeEvent
}

// AnnouncementStatus 公告状态（可存储）
type AnnouncementStatus string

const (
	AnnouncementStatusDraft     AnnouncementStatus = "draft"     // 草稿
	AnnouncementStatusPublished AnnouncementStatus = "published" // 已发布
	AnnouncementStatusArchived  AnnouncementStatus = "archived"  // 已归档
)

// IsValid 检查公告状态是否合法
func (vs AnnouncementStatus) IsValid() bool {
	switch vs {
	case AnnouncementStatusDraft, AnnouncementStatusPublished, AnnouncementStatusArchived:
		return true
	}
	return false
}

// AnnouncementStatusFilter 公告状态查询筛选（可包含 All）
type AnnouncementStatusFilter string

const (
	AnnouncementStatusFilterAll       AnnouncementStatusFilter = "all" // 全部（仅查询）
	AnnouncementStatusFilterDraft     AnnouncementStatusFilter = "draft"
	AnnouncementStatusFilterPublished AnnouncementStatusFilter = "published"
	AnnouncementStatusFilterArchived  AnnouncementStatusFilter = "archived"
)

// FromStatus 将存储状态转换为筛选状态
func (f AnnouncementStatusFilter) FromStatus(s AnnouncementStatus) AnnouncementStatusFilter {
	return AnnouncementStatusFilter(s)
}

// IsAll 判断是否为“全部”筛选
func (f AnnouncementStatusFilter) IsAll() bool {
	return f == AnnouncementStatusFilterAll
}
