package model

import "time"

// Attachment 附件模型
type Attachment struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	FileID       string    `gorm:"column:file_id;type:varchar(64);uniqueIndex;not null" json:"file_id"` // 唯一标识
	UserID       int64     `gorm:"column:user_id;index;not null" json:"user_id"`
	PostID       int64     `gorm:"column:post_id;index;default:0" json:"post_id"`
	ReplyID      int64     `gorm:"column:reply_id;index;default:0" json:"reply_id"` // 关联回复ID
	OriginalName string    `gorm:"column:original_name;type:varchar(255);not null" json:"original_name"`
	StoredName   string    `gorm:"column:stored_name;type:varchar(255);not null" json:"stored_name"`
	StoredPath   string    `gorm:"column:stored_path;type:varchar(500);not null" json:"stored_path"`
	Size         int64     `gorm:"column:size;not null" json:"size"`
	MimeType     string    `gorm:"column:mime_type;type:varchar(100)" json:"mime_type"`
	FileType     string    `gorm:"column:file_type;type:varchar(50);index" json:"file_type"` // avatar, post_image, attachment
	Ext          string    `gorm:"column:ext;type:varchar(20)" json:"ext"`
	Width        int       `gorm:"column:width;default:0" json:"width"`   // 图片宽度
	Height       int       `gorm:"column:height;default:0" json:"height"` // 图片高度
	Status       int       `gorm:"column:status;default:1" json:"status"` // 0:临时 1:正常 2:已删除
	UploadIP     string    `gorm:"column:upload_ip;type:varchar(45)" json:"upload_ip"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Attachment) TableName() string {
	return "attachments"
}
