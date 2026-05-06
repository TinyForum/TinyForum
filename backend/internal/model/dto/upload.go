package dto

import "tiny-forum/internal/model/do"

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
	ID           string              `json:"id"`
	UserID       int64               `json:"user_id"`
	PostID       int64               `json:"post_id"`
	OriginalName string              `json:"original_name"`
	StoredName   string              `json:"stored_name"`
	StoredPath   string              `json:"stored_path"`
	Size         int64               `json:"size"`
	MimeType     string              `json:"mime_type"`
	FileType     do.FileType         `json:"file_type"`
	Ext          string              `json:"ext"`
	Status       do.AttachmentStatus `json:"status"`
	CreatedAt    string              `json:"created_at"`
}
