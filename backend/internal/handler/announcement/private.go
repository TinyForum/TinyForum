package announcement

import (
	"strconv"

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
