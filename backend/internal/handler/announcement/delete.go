package announcement

import (
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// Delete 删除公告
func (h *AnnouncementHandler) Delete(c *gin.Context) {
	id, ok := parseAnnouncementID(c)
	if !ok {
		return
	}
	userID := c.GetUint("user_id")
	if err := h.service.Delete(c.Request.Context(), id, userID); err != nil {
		handleAnnouncementServiceError(c, err)
		return
	}
	response.Success(c, nil)
}
