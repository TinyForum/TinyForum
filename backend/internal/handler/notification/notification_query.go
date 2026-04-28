package notification

import (

	"tiny-forum/internal/dto"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// List 获取通知列表
// @Summary 获取通知列表
// @Description 分页获取当前用户的通知列表
// @Tags 通知管理
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.Notification}} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /notifications [get]
func (h *NotificationHandler) List(c *gin.Context) {
    var req dto.NotificationListRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        response.BadRequest(c, "参数错误")
        return
    }
    
    userID := c.GetUint("user_id")
    
    // 处理分页默认值
    page := req.Page
    if page < 1 {
        page = 1
    }
    pageSize := req.PageSize
    if pageSize < 1 {
        pageSize = 20
    }
    
    // 调用 Service，获取 BO
    result, err := h.notifSvc.List(userID, page, pageSize)
    if err != nil {
        response.InternalError(c, err.Error())
        return
    }
    
    // BO → DTO 转换
    resp := dto.NotificationListResponse{
        List:        make([]dto.NotificationResponse, 0, len(result.List)),
        Total:       result.Total,
        UnreadCount: result.UnreadCount,
        Page:        result.Page,
        PageSize:    result.PageSize,
    }
    
    for _, notif := range result.List {
        item := dto.NotificationResponse{
            ID:         notif.ID,
            Type:       notif.Type,
            Content:    notif.Content,
            IsRead:     notif.IsRead,
            CreatedAt:  notif.CreatedAt,
            TargetID:   notif.TargetID,
            TargetType: notif.TargetType,
        }
        
        if notif.Sender != nil {
            item.Sender = &dto.NotificationSenderResponse{
                ID:       notif.Sender.ID,
                Username: notif.Sender.Username,
                Avatar:   notif.Sender.Avatar,
            }
        }
        
        resp.List = append(resp.List, item)
    }
    
    response.Success(c, resp)
}

// UnreadCount 获取未读通知数量
// @Summary 获取未读通知数量
// @Description 获取当前用户未读通知的总数
// @Tags 通知管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=object} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /notifications/count/unread [get]
func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")
	count, err := h.notifSvc.UnreadCount(userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"count": count})
}