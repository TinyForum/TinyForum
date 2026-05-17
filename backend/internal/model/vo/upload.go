package vo

import "tiny-forum/internal/model/do"

type FileInfoVO struct {
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

// UploadResult 上传结果，由调用方负责保存到数据库
type UploadResult struct {
	FileHash     string
	StoredPath   string
	StoredName   string
	MimeType     string
	MimeMajor    do.MimeTypeMajor
	Ext          string
	Size         int64
	OriginalName string
}
