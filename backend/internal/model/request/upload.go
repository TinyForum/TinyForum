package request

import (
	"mime/multipart"
	"tiny-forum/internal/model/do"
)

// UploadRequest 上传请求
type UploadPostFileRequest struct {
	PostID   int64       `form:"post_id"` // 关联帖子ID（可选）
	FileType do.FileType `form:"type" binding:"required,oneof=avatar post_image post_file comment_image plugin_asset"`
}

type UploadPluginRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}
