package upload

import (
	"strconv"

	"tiny-forum/internal/dto"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// Upload 上传文件
// @Summary 上传文件
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Param post_id formData int false "帖子ID"
// @Param type formData string true "文件类型 (avatar/post_image/comment_attachment)"
// @Success 200 {object} response.Response{data=dto.UploadResponse}
// @Router /upload [post]
func (h *UploadHandler) Upload(c *gin.Context) {
	// 获取用户ID（从认证中间件）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请选择要上传的文件")
		return
	}

	// 解析请求参数
	var req dto.UploadRequest
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 调用服务上传
	result, err := h.service.UploadFile(c.Request.Context(), userID.(int64), file, &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, result)
}

// GetFile 获取文件信息
// @Summary 获取文件信息
// @Tags Upload
// @Produce json
// @Param file_id path string true "文件ID"
// @Success 200 {object} response.Response{data=dto.FileInfo}
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
// @Tags Upload
// @Produce json
// @Param file_id path string true "文件ID"
// @Success 200 {object} response.Response
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
// @Tags Upload
// @Produce json
// @Param type query string false "文件类型"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response{data=response.PageData{list=[]dto.FileInfo}}
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
// @Tags Upload
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
