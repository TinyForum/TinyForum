package do

// Attachment 附件模型
type Attachment struct {
	BaseModel
	FileID       string `gorm:"column:file_id;type:varchar(64);uniqueIndex;not null" json:"file_id"` // 唯一标识
	UserID       int64  `gorm:"column:user_id;index;not null" json:"user_id"`                        // 上传用户ID
	PostID       int64  `gorm:"column:post_id;index;default:0" json:"post_id"`
	ReplyID      int64  `gorm:"column:reply_id;index;default:0" json:"reply_id"`                      // 关联回复ID
	OriginalName string `gorm:"column:original_name;type:varchar(255);not null" json:"original_name"` // 原始名称
	StoredName   string `gorm:"column:stored_name;type:varchar(255);not null" json:"stored_name"`     // 存储名称
	StoredPath   string `gorm:"column:stored_path;type:varchar(500);not null" json:"stored_path"`     // 存储路径
	Size         int64  `gorm:"column:size;not null" json:"size"`                                     // 文件大小
	MimeType     string `gorm:"column:mime_type;type:varchar(100)" json:"mime_type"`                  // mime 类型
	// avatar, post_image, attachment
	FileType string `gorm:"column:file_type;type:varchar(50);index" json:"file_type"` // 文件类型
	Ext      string `gorm:"column:ext;type:varchar(20)" json:"ext"`                   // 文件扩展名
	Width    int    `gorm:"column:width;default:0" json:"width"`                      // 图片宽度
	Height   int    `gorm:"column:height;default:0" json:"height"`                    // 图片高度
	Status   int    `gorm:"column:status;default:1" json:"status"`                    // 0:临时 1:正常 2:已删除
	UploadIP string `gorm:"column:upload_ip;type:varchar(45)" json:"upload_ip"`       // 上传IP
}

func (Attachment) TableName() string {
	return "attachments"
}
