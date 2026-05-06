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
	attachments := api.Group("/attachments")
	attachments.Use(mw.Auth()) // 需要认证使用
	{

		postFile := attachments.Group("/post")
		{
			postFile.POST("/:post_id", h.UploadPostFile)    // POST /api/v1/attachments/post/:post_id - 上传帖子文件
			postFile.GET("/post/:post_id", h.GetFile)       // GET /api/v1/attachments/post/:post_id - 获取帖子文件信息
			postFile.DELETE("/post/:file_id", h.DeleteFile) // DELETE /api/v1/attachments/post/:file_id - 删除帖子文件
		}

		commentFile := attachments.Group("/comment")
		{
			commentFile.POST("/:comment_id", h.UploadCommentFile) // POST /api/v1/attachments/comment/:comment_id - 上传评论文件
			commentFile.GET("/comment/:file_id", h.GetFile)       // GET /api/v1/attachments/comment/:file_id - 获取评论文件信息
			commentFile.DELETE("/comment/:file_id", h.DeleteFile) // DELETE /api/v1/attachments/comment/:file_id - 删除评论文件
		}

		pluginFile := attachments.Group("/plugin")
		{
			pluginFile.POST("", h.UploadPluginFile)             // POST /api/v1/attachments/plugin - 上传插件文件
			pluginFile.GET("/plugin/:file_id", h.GetFile)       // GET /api/v1/attachments/plugin/:file_id - 获取插件文件信息
			pluginFile.DELETE("/plugin/:file_id", h.DeleteFile) // DELETE /api/v1/attachments/plugin/:file_id - 删除插件文件
		}

		attachments.GET("/users/me/files", h.GetUserFiles) // GET /api/v1/attachments/users/me/files - 获取当前用户的文件列表
	}

	// 公开文件访问（无认证）
	public := api.Group("/files")
	{
		public.GET("/:file_id", h.ServeFile) // GET /api/v1/files/:file_id - 公开访问文件
	}

}
