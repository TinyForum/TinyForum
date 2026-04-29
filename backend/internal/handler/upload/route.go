package upload

import (
	"tiny-forum/internal/middleware"
	uploadService "tiny-forum/internal/service/upload"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	service uploadService.UploadService
}

func NewUploadHandler(service uploadService.UploadService) *UploadHandler {
	return &UploadHandler{service: service}
}

// 在路由注册函数中
func (h *UploadHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	// 需要认证的上传接口
	upload := api.Group("/upload")
	upload.Use(mw.Auth())
	{
		upload.POST("", h.Upload)                 // POST /api/v1/upload - 上传文件
		upload.GET("/user/files", h.GetUserFiles) // GET /api/v1/upload/user/files - 获取用户文件列表
		upload.GET("/:file_id", h.GetFile)        // GET /api/v1/upload/:file_id - 获取文件信息
		upload.DELETE("/:file_id", h.DeleteFile)  // DELETE /api/v1/upload/:file_id - 删除文件
	}

	// 公开文件访问（无认证）
	public := api.Group("/files")
	{
		public.GET("/:file_id", h.ServeFile) // GET /api/v1/files/:file_id - 公开访问文件
	}
}

// 静态文件服务（开发环境）
// r.Static("/uploads", "./uploads")
