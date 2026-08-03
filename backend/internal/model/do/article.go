package do

//
import (
	"tiny-forum/internal/model/common"
)

type Article struct {
	common.BaseModel
	CreationID uint     `gorm:"uniqueIndex;not null" json:"creation_id"` // 外键，唯一索引保证一对一
	Creation   Creation `gorm:"foreignKey:CreationID;references:ID" json:"creation,omitempty"`
	// 如果有 Article 特有字段，加在这里；若无，可保留空结构体
}
type CreationType string

const (
	CreationTypeImageText CreationType = "image_text" // 图文
	CreationTypeShortVideo CreationType = "short_video" // 短视频
	CreationTypeLongVideo  CreationType = "long_video"  // 长视频
	CreationTypeImage     CreationType = "image"       // 图片
	CreationTypeArticle   CreationType = "article"     // 文章
	CreationTypeQuestion  CreationType = "question"    // 问答
	CreationTypeTopic     CreationType = "topic"       // 话题
	CreationTypePost     CreationType = "post"        // 帖子
)

// 合法的帖子类型集合
var validPostTypes = map[CreationType]bool{
	CreationTypeImageText:  true,
	CreationTypeShortVideo: true,
	CreationTypeLongVideo:  true,
	CreationTypeImage:     true,
	CreationTypeArticle:   true,
	CreationTypeQuestion:  true,
	CreationTypeTopic:     true,
	CreationTypePost:      true,
}

var validCreationStatuses = map[CreationStatus]bool{
	CreationStatusDraft:     true,
	CreationStatusPending:   true,
	CreationStatusPublished: true,
	CreationStatusHidden:    true,
}

type CreationStatus string

// const (
// 	PostTypePost    PostType = "post"
// 	PostTypeArticle PostType = "article"
// 	PostTypeTopic   PostType = "topic"
// )

// 用户主动控制的状态（用户能感知、能操作）

const (
	CreationStatusDraft     CreationStatus = "draft"     // 草稿（用户保存未发布）
	CreationStatusPending   CreationStatus = "pending"   // 待用户确认/提交（如编辑后重新提交）
	CreationStatusPublished CreationStatus = "published" // 已发布（用户主动发布）
	CreationStatusHidden    CreationStatus = "hidden"    // 用户隐藏（如自己删除/隐藏，或管理员操作但以用户视角展示）
)

// enum [draft pending published hidden]

// 系统风控状态（由内容安全模块自动判定或管理员审核结果）

// IsValid 检查帖子类型是否合法
func (pt CreationType) IsValid() bool {
	return validPostTypes[pt]
}

// 可选：从字符串安全转换
func ParsePostType(s string) CreationType {
	pt := CreationType(s)
	if pt.IsValid() {
		return pt
	}
	return CreationTypeArticle // 默认值
}

func ParseCreationStatus(s string) CreationStatus {
	ps := CreationStatus(s)
	if ps.IsValid() {
		return ps
	}
	return CreationStatusPublished // 默认值
}

func (ps CreationStatus) IsValid() bool {
	return validCreationStatuses[ps]
}
