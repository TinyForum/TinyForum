package attachment

import (
	"strconv"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// ListMyFiles 获取当前用户的文件列表
// @Summary 获取用户文件列表
// @Description 分页获取当前登录用户上传的所有文件，可按文件类型过滤
// @Tags attachment
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页条数" default(20)
// @Param file_type query string false "文件类型过滤 (post_image, comment_file, avatar, plugin)"
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Router /attachments/user/me [get]
func (h *AttachmentHandler) ListMyFiles(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	fileType := c.Query("file_type")
	files, total, err := h.svc.GetUserFiles(c.Request.Context(), userID, fileType, page, pageSize)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, gin.H{"list": files, "total": total})
}
