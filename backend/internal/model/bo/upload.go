package bo

import (
	"mime/multipart"
	"tiny-forum/internal/model/do"
)

// 	UploadFile(ctx context.Context, userID int64, fileHeader *multipart.FileHeader, req *request.UploadPostFileRequest) (*dto.UploadResponse, error)

// PluginListBO 用于 Service 接收查询参数（不含业务返回字段）
type PluginUpdateBO struct {
	UserID     uint                  `json:"user_Id"`
	FileHeader *multipart.FileHeader `json:"file_header"`
}

type FileInfoBO struct {
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
