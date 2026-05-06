package request

import "tiny-forum/internal/model/do"

// UploadRequest 上传请求
type UploadPostFileRequest struct {
	PostID   int64       `form:"post_id"` // 关联帖子ID（可选）
	FileType do.FileType `form:"type" binding:"required,oneof=avatar post_image post_file comment_image plugin_asset"`
}
