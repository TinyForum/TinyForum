package bo

import (
	"mime/multipart"
)

// 	UploadFile(ctx context.Context, userID int64, fileHeader *multipart.FileHeader, req *request.UploadPostFileRequest) (*dto.UploadResponse, error)

// PluginListBO 用于 Service 接收查询参数（不含业务返回字段）
type PluginUpdateBO struct {
	UserID     uint                  `json:"user_Id"`
	FileHeader *multipart.FileHeader `json:"file_header"`
}
