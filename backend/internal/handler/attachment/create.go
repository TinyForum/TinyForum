package attachment

import (
	"tiny-forum/internal/model/converter"
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

// Attachment 上传文件
// @Summary 上传文件
// @Tags 上传管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Param post_id formData int false "帖子ID"
// @Param type formData string true "文件类型 (avatar/post_image/comment_attachment)"
// @Success 200 {object} common.BasicResponse
// @Router /attachment/post_file [post]
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

// UploadPluginFile
// @Summary 上传插件
// @Tags 上传管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Param post_id formData int false "帖子ID"
// @Param type formData string true "文件类型 (avatar/post_image/comment_attachment)"
// @Success 200 {object} common.BasicResponse
// @Router /attachments/plugin [post]
func (h *UploadHandler) UploadPluginFile(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, apperrors.ErrUnauthorized.Error())
		return
	}
	logger.Infof("请求 Content-Type: %s", c.ContentType())

	// 绑定请求（仅包含文件字段）
	var req request.UploadPluginRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	// 转换为 BO，注入用户ID
	requestBO := converter.UploadPluginRequestToUploadPluginBo(req, userID.(uint))

	// 调用服务
	result, err := h.service.UploadPlugin(c.Request.Context(), requestBO)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, result)
}
