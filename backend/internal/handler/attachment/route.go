package attachment

import (
	"tiny-forum/internal/middleware"
	"tiny-forum/internal/service/attachment"

	"github.com/gin-gonic/gin"
)

type AttachmentHandler struct {
	svc attachment.AttachmentService
}

func NewAttachmentHandler(svc attachment.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{svc: svc}
}

// RegisterRoutes 注册附件相关路由
func (h *AttachmentHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	// 需要认证的附件操作
	attachments := api.Group("/attachments")
	attachments.Use(mw.Auth())
	{
		// 通用文件上传（通过 type 字段区分业务类型，如 post_image, comment_image, plugin 等）
		attachments.POST("", h.UploadFile) // POST /api/v1/attachments

		// 获取当前用户的所有文件列表（支持 type 过滤）
		attachments.GET("/user/me", h.ListMyFiles) // GET /api/v1/attachments/user/me 获取当前用户的所有文件列表

		// 单个文件操作
		attachments.GET("/:file_id", h.GetFile)       // GET /api/v1/attachments/:file_id 获取文件信息
		attachments.DELETE("/:file_id", h.DeleteFile) // DELETE /api/v1/attachments/:file_id 删除文件
	}

	// 公开文件访问（无认证，注意权限控制：只能访问公有文件或带签名 URL）
	public := api.Group("/files")
	{
		public.GET("/:file_id", h.ServeFile) // GET /api/v1/files/:file_id
	}
}

// package attachment
