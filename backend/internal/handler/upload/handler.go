package upload

import (
	"strconv"

	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetFile 获取文件信息
// @Summary 获取文件信息
// @Tags 上传管理
// @Produce json
// @Param file_id path string true "文件ID"
// @Success 200 {object} common.BasicResponse
// @Router /upload/{file_id} [get]
func (h *UploadHandler) GetFile(c *gin.Context) {
	fileID := c.Param("file_id")
	if fileID == "" {
		response.BadRequest(c, "缺少文件ID")
		return
	}

	info, err := h.service.GetFile(c.Request.Context(), fileID)
	if err != nil {
		response.NotFound(c, "文件不存在")
		return
	}

	response.Success(c, info)
}

// DeleteFile 删除文件
// @Summary 删除文件
// @Tags 上传管理
// @Produce json
// @Param file_id path string true "文件ID"
// @Success 200 {object} common.BasicResponse
// @Router /upload/{file_id} [delete]
func (h *UploadHandler) DeleteFile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	fileID := c.Param("file_id")
	if fileID == "" {
		response.BadRequest(c, "缺少文件ID")
		return
	}

	if err := h.service.DeleteFile(c.Request.Context(), userID.(int64), fileID); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, "删除成功")
}

// GetUserFiles 获取用户文件列表
// @Summary 获取用户文件列表
// @Tags 上传管理
// @Produce json
// @Param type query string false "文件类型"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} common.BasicResponse
// @Router /user/files [get]
func (h *UploadHandler) GetUserFiles(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	fileType := c.Query("type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	list, total, err := h.service.GetUserFiles(c.Request.Context(), userID.(int64), fileType, page, pageSize)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessPage(c, list, total, page, pageSize)
}

// ServeFile 提供文件访问（需要权限校验）
// @Summary 获取文件内容
// @Tags 上传管理
// @Produce image/jpeg,image/png,application/pdf
// @Param file_id path string true "文件ID"
// @Router /files/{file_id} [get]
func (h *UploadHandler) ServeFile(c *gin.Context) {
	fileID := c.Param("file_id")
	if fileID == "" {
		response.BadRequest(c, "缺少文件ID")
		return
	}

	info, err := h.service.GetFile(c.Request.Context(), fileID)
	if err != nil {
		response.NotFound(c, "文件不存在")
		return
	}

	// 权限检查：公开帖子/头像可公开访问，私密内容需要验证
	// TODO: 根据帖子权限判断

	// 设置响应头
	c.Header("Content-Type", info.MimeType)
	c.Header("Content-Disposition", "inline") // 预览而非下载

	c.File(info.StoredPath)
}
