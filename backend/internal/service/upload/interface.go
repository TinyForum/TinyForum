// internal/service/upload/interface.go
package upload

import (
	"context"
	"mime/multipart"
	"tiny-forum/internal/model/do"
)

type Engine interface {
	// Upload 仅负责存储文件，返回存储结果，不操作数据库
	Upload(ctx context.Context, req *UploadRequest) (*UploadResult, error)
	DeleteFile(ctx context.Context, storedPath string) error   // 新增
}

type UploadRequest struct {
	UserID    int64
	PluginID  string
	File      *multipart.FileHeader
	FileType  do.FileType
	PostID    int64
	ReplyID   int64
	ClientIP  string
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