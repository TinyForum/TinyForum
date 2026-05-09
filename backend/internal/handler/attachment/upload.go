package attachment

import (
	"tiny-forum/internal/model/request"
	"tiny-forum/pkg/logger"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// UploadFile 上传文件
// @Summary 上传文件
// @Description 上传附件，支持图片、文档等，通过 type 字段区分业务类型（post_image, comment_file, avatar 等）
// @Tags attachment
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Param type formData string true "业务类型：post_image / comment_file / avatar / plugin_asset"
// @Param post_id formData int false "关联的帖子ID（当 type = post_image 时需提供）"
// @Param reply_id formData int false "关联的评论ID（当 type = comment_file 时需提供）"
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Router /attachments [post]
func (h *AttachmentHandler) UploadFile(c *gin.Context) {
	userID := c.GetInt64("user_id")
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required")
		return
	}
	var req request.UploadPostFileRequest
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	resp, err := h.svc.UploadFile(c.Request.Context(), userID, file, &req, c.ClientIP())
	if err != nil {
		response.HandleError(c,  err)
		return
	}
	response.Success(c, resp)
}

// GetFile 获取文件元信息
// @Summary 获取文件元信息
// @Description 根据文件ID获取文件的元数据（不包含文件内容）
// @Tags attachment
// @Accept json
// @Produce json
// @Param file_id path string true "文件ID"
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Failure 404 {object} common.BasicResponse
// @Router /attachments/{file_id} [get]
func (h *AttachmentHandler) GetFile(c *gin.Context) {
	fileID := c.Param("file_id")
	fileInfo, err := h.svc.GetFile(c.Request.Context(), fileID)
	if err != nil {
		logger.Errorf("获取文件源信息错误: ", err)
		response.HandleError(c, err)
		return
	}
	response.Success(c, fileInfo)
}

// DeleteFile 删除文件
// @Summary 删除文件
// @Description 删除用户上传的文件（仅文件所有者可操作）
// @Tags attachment
// @Accept json
// @Produce json
// @Param file_id path string true "文件ID"
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Failure 403 {object} common.BasicResponse
// @Router /attachments/{file_id} [delete]
func (h *AttachmentHandler) DeleteFile(c *gin.Context) {
	userID := c.GetInt64("user_id")
	fileID := c.Param("file_id")
	if err := h.svc.DeleteFile(c.Request.Context(), userID, fileID); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, gin.H{"message": "deleted"})
}

// ServeFile 公开访问文件
// @Summary 公开文件下载
// @Description 根据文件ID直接返回文件内容（适用于图片、附件下载）
// @Tags attachment
// @Accept json
// @Param file_id path string true "文件ID"
// @Success 200 {file} binary "文件内容"
// @Failure 400 {object} common.BasicResponse
// @Failure 404 {object} common.BasicResponse
// @Router /files/{file_id} [get]
func (h *AttachmentHandler) ServeFile(c *gin.Context) {
	fileID := c.Param("file_id")
	// 获取文件路径
	fileInfo, err := h.svc.GetFile(c.Request.Context(), fileID)
	if err != nil {
			response.HandleError(c, err)
		return
	}
	// 注意：这里假设 fileInfo.URL 是相对路径或绝对路径，需要映射到实际存储位置
	// 简化处理：读取物理文件返回。实际应根据 config 中的存储目录拼接。
	// 这里仅示意，你需要实现真正的文件读取逻辑。
	c.File(fileInfo.URL) // 如果 URL 是绝对路径
}