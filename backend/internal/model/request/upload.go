package request

import (
	"mime/multipart"
	"tiny-forum/internal/model/do"
)

// UploadRequest 上传请求
// type UploadPostFileRequest struct {
// 	PostID   int64       `form:"post_id"` // 关联帖子ID（可选）
// 	FileType do.FileType `form:"type" binding:"required,oneof=avatar post_image post_file comment_image plugin"`
// }

// type UploadPluginRequest struct {
// 	File *multipart.FileHeader `form:"file" binding:"required"`
// }

type UploadPostFileRequest struct {
	Type    string `form:"type" binding:"required"` // post_image, avatar, etc.
	PostID  int64  `form:"post_id"`
	ReplyID int64  `form:"reply_id"`
}

type UploadRequest struct {
	UserID   uint
	PluginID string
	File     *multipart.FileHeader
	FileType do.FileType
	PostID   int64
	ReplyID  int64
	ClientIP string
}
