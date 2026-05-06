package do

import (
	"database/sql/driver"
	"fmt"
	"tiny-forum/internal/model/common"
)

// ========== 枚举类型定义 ==========

// FileType 附件业务类型
type FileType string

const (
	FileTypeAvatar       FileType = "avatar"        // 用户头像
	FileTypePostImage    FileType = "post_image"    // 帖子图片
	FileTypePostFile     FileType = "post_file"     // 帖子普通附件
	FileTypeCommentImage FileType = "comment_image" // 评论图片
	FileTypePluginAsset  FileType = "plugin_asset"  // 插件资源文件
)

func (f FileType) IsValid() bool {
	switch f {
	case FileTypeAvatar, FileTypePostImage, FileTypePostFile, FileTypePluginAsset:
		return true
	}
	return false
}

// MimeTypeMajor 媒体类型大类（做简单约束）
type MimeTypeMajor string

const (
	MimeImage    MimeTypeMajor = "image"
	MimeVideo    MimeTypeMajor = "video"
	MimeAudio    MimeTypeMajor = "audio"
	MimeDocument MimeTypeMajor = "document"
	MimeArchive  MimeTypeMajor = "archive"
	MimeOther    MimeTypeMajor = "other"
)

// AttachmentStatus 附件状态
type AttachmentStatus int

const (
	StatusTemp    AttachmentStatus = 0 // 临时（未关联业务对象）
	StatusNormal  AttachmentStatus = 1 // 正常已使用
	StatusDeleted AttachmentStatus = 2 // 已删除（软删除标记）
)

func (s AttachmentStatus) String() string {
	switch s {
	case StatusTemp:
		return "temporary"
	case StatusNormal:
		return "normal"
	case StatusDeleted:
		return "deleted"
	}
	return "unknown"
}

// ========== 主模型 ==========

// Attachment 附件模型（包含插件上传支持）
type Attachment struct {
	common.BaseModel

	// ---------- 唯一标识与归属 ----------
	FileID   string `gorm:"column:file_id;type:varchar(64);uniqueIndex;not null;comment:附件唯一标识(通常为UUID)" json:"file_id"`
	UserID   int64  `gorm:"column:user_id;index;not null;comment:上传用户ID" json:"user_id"`
	PluginID string `gorm:"column:plugin_id;type:varchar(64);index;comment:插件ID（若由插件上传）" json:"plugin_id,omitempty"` // 支持插件上传

	// ---------- 业务关联 ----------
	PostID  int64 `gorm:"column:post_id;index;default:0;comment:关联的帖子ID" json:"post_id"`
	ReplyID int64 `gorm:"column:reply_id;index;default:0;comment:关联的回复ID" json:"reply_id"`

	// ---------- 文件元信息 ----------
	OriginalName string `gorm:"column:original_name;type:varchar(255);not null;comment:用户上传时的原始文件名" json:"original_name"`
	StoredName   string `gorm:"column:stored_name;type:varchar(255);not null;comment:系统存储的唯一文件名" json:"stored_name"`
	StoredPath   string `gorm:"column:stored_path;type:varchar(500);not null;comment:存储相对路径" json:"stored_path"`
	Size         int64  `gorm:"column:size;not null;comment:文件大小(字节)" json:"size"`

	// 使用明确的枚举类型替代模糊 string
	FileType  FileType      `gorm:"column:file_type;type:varchar(50);index;comment:附件业务类型" json:"file_type"`
	MimeType  string        `gorm:"column:mime_type;type:varchar(100);comment:完整MIME类型，例如 image/png" json:"mime_type"`
	MimeMajor MimeTypeMajor `gorm:"column:mime_major;type:varchar(20);comment:媒体大类(image/video/audio等);index" json:"mime_major"` // 便于索引查询
	Ext       string        `gorm:"column:ext;type:varchar(20);comment:文件扩展名(不含点)" json:"ext"`

	// 图片专用字段
	Width  int `gorm:"column:width;default:0;comment:图片宽度(px)" json:"width"`
	Height int `gorm:"column:height;default:0;comment:图片高度(px)" json:"height"`

	// 状态与来源
	Status   AttachmentStatus `gorm:"column:status;default:1;index;comment:0临时 1正常 2已删除" json:"status"`
	UploadIP string           `gorm:"column:upload_ip;type:varchar(45);comment:上传时的客户端IP" json:"upload_ip"`

	// 可选：插件自定义元数据 (JSON)
	PluginMeta map[string]any `gorm:"column:plugin_meta;type:json;comment:插件附加数据, 例如缩略图处理参数等" json:"plugin_meta,omitempty" serializer:json"`
}

// TableName 指定表名
func (Attachment) TableName() string {
	return "attachments"
}

// ===== 实现自定义类型扫描/值转换（可选，保证枚举存储为字符串） =====

// Value 实现 driver.Valuer，将 FileType 存储为字符串
func (f FileType) Value() (driver.Value, error) {
	return string(f), nil
}

// Scan 实现 sql.Scanner，从数据库字符串解析为 FileType
func (f *FileType) Scan(value interface{}) error {
	if value == nil {
		*f = ""
		return nil
	}
	switch v := value.(type) {
	case string:
		*f = FileType(v)
	case []byte:
		*f = FileType(string(v))
	default:
		return fmt.Errorf("unsupported type for FileType: %T", v)
	}
	if !f.IsValid() {
		return fmt.Errorf("invalid FileType value: %s", *f)
	}
	return nil
}

// 同样为 AttachmentStatus 实现 Value/Scan（存储为 int）
func (s AttachmentStatus) Value() (driver.Value, error) {
	return int(s), nil
}

func (s *AttachmentStatus) Scan(value interface{}) error {
	if value == nil {
		*s = 0
		return nil
	}
	switch v := value.(type) {
	case int64:
		*s = AttachmentStatus(v)
	case int:
		*s = AttachmentStatus(v)
	default:
		return fmt.Errorf("unsupported type for AttachmentStatus: %T", v)
	}
	return nil
}
