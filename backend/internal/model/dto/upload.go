package dto

// UploadRequest 上传请求
type UploadRequest struct {
	PostID   int64  `form:"post_id"` // 关联帖子ID（可选）
	FileType string `form:"type" binding:"required,oneof=avatar post_image comment_attachment"`
}

// UploadResponse 上传响应
type UploadResponse struct {
	FileID       string `json:"file_id"` // 存储标识
	URL          string `json:"url"`     // 访问URL
	OriginalName string `json:"original_name"`
	Size         int64  `json:"size"`
	MimeType     string `json:"mime_type"`
}

// FileInfo 文件信息（用于数据库存储）
type FileInfo struct {
	ID           string `json:"id"`
	UserID       int64  `json:"user_id"`
	PostID       int64  `json:"post_id"`
	OriginalName string `json:"original_name"`
	StoredName   string `json:"stored_name"`
	StoredPath   string `json:"stored_path"`
	Size         int64  `json:"size"`
	MimeType     string `json:"mime_type"`
	FileType     string `json:"file_type"` // avatar, post_image, attachment
	Ext          string `json:"ext"`
	Status       int    `json:"status"` // 0:临时 1:正式 2:已删除
	CreatedAt    string `json:"created_at"`
}
