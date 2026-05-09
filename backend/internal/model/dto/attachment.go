package dto

import "time"

// // FileInfo 文件信息（用于数据库存储）
// type FileInfo struct {

// 	// ID           string              `json:"id"`
// 	UserID       int64               `json:"user_id"`
// 	PostID       int64               `json:"post_id"`
// 	OriginalName string              `json:"original_name"`
// 	StoredName   string              `json:"stored_name"`
// 	StoredPath   string              `json:"stored_path"`
// 	Size         int64               `json:"size"`
// 	MimeType     string              `json:"mime_type"`
// 	FileType     do.FileType         `json:"file_type"`
// 	Ext          string              `json:"ext"`
// 	Status       do.AttachmentStatus `json:"status"`
// 	// CreatedAt    string              `json:"created_at"`
// }

type FileInfo struct {
	FileID       string    `json:"file_id"`
	OriginalName string    `json:"original_name"`
	Size         int64     `json:"size"`
	FileType     string    `json:"file_type"`
	MimeType     string    `json:"mime_type"`
	URL          string    `json:"url"`
	CreatedAt    time.Time `json:"created_at"`
}
