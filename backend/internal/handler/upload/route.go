package upload

import (
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
func (h *UploadHandler) RegisterRoutes(upload *gin.RouterGroup) {
	g := upload.Group("/:file_id")
	{
		g.GET("/:file_id", h.GetFile)
		g.DELETE("/:file_id", h.DeleteFile)
	}
	upload.POST("", h.Upload)

	upload.GET("/user/files", h.GetUserFiles)
}

// 公开文件访问（无认证）
// r.GET("/files/:file_id", handlers.UploadHandler.ServeFile)

// 静态文件服务（开发环境）
// r.Static("/uploads", "./uploads")
