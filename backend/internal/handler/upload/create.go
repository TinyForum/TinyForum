package upload

import (
	"tiny-forum/internal/model/request"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/logger"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// UploadPostFile
// @Summary 上传帖子文件
// @Tags 上传管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Param post_id formData int false "帖子ID"
// @Param type formData string true "文件类型"
// @Success 200 {object} common.BasicResponse
// @Router /attachments/post/:post_id [post]
func (h *UploadHandler) UploadPostFile(c *gin.Context) {

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
	var req request.UploadPostFileRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Infof("上传参数错误: ", err)
		response.BadRequest(c, apperrors.ErrValidation.Error())
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

// Upload 上传文件
// @Summary 上传文件
// @Tags 上传管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Param post_id formData int false "帖子ID"
// @Param type formData string true "文件类型 (avatar/post_image/comment_attachment)"
// @Success 200 {object} common.BasicResponse
// @Router /upload/post_file [post]
func (h *UploadHandler) UploadCommentFile(c *gin.Context) {

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
	var req request.UploadPostFileRequest
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

// Upload 上传文件
// @Summary 上传文件
// @Tags 上传管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Param post_id formData int false "帖子ID"
// @Param type formData string true "文件类型 (avatar/post_image/comment_attachment)"
// @Success 200 {object} common.BasicResponse
// @Router /upload/post_file [post]
func (h *UploadHandler) UploadPluginFile(c *gin.Context) {

	// 获取用户ID（从认证中间件）
	userID, exists := c.Get("user_id")
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
	var req request.UploadPostFileRequest
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
