package announcement

import (
	"strconv"

	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// 解析公告ID，返回错误响应
func parseAnnouncementID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的公告ID")
		return 0, false
	}
	return uint(id), true
}

// 统一处理 service 层错误
func handleAnnouncementServiceError(c *gin.Context, err error) {
	switch err {
	case apperrors.ErrAnnouncementNotFound:
		response.NotFound(c, err.Error())
	case apperrors.ErrInvalidPublishTime:
		response.BadRequest(c, err.Error())
	default:
		response.InternalError(c, err.Error())
	}
}
