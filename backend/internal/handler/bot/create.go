package bot

import (
	"strconv"
	"tiny-forum/internal/model/request"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// Create 创建机器人（支持 Lua 脚本 / 零代码两种模式）
// @Summary 创建机器人
// @Tags 机器人管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body request.CreateBotRequest true "机器人创建信息"
// @Success 200 {object} common.BasicResponse{data=object{id=integer}}
// @Router /bots [post]
func (h *Handler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.CreateBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	bot, err := h.svc.Create(c.Request.Context(), userID, &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, gin.H{"id": bot.ID})
}

// Update 更新机器人
// @Summary 更新机器人
// @Tags 机器人管理
// @Security ApiKeyAuth
// @Param id path int true "机器人ID"
// @Param body body request.UpdateBotRequest true "更新信息"
// @Router /bots/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req request.UpdateBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	if err := h.svc.Update(c.Request.Context(), userID, uint(id), &req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, nil)
}

// Delete 删除机器人
// @Summary 删除机器人
// @Tags 机器人管理
// @Security ApiKeyAuth
// @Param id path int true "机器人ID"
// @Router /bots/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), userID, uint(id)); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, nil)
}

// Get 获取机器人详情
// @Summary 获取机器人详情
// @Tags 机器人管理
// @Security ApiKeyAuth
// @Param id path int true "机器人ID"
// @Router /bots/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	bot, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, bot)
}

// ListMyBot 获取当前用户创建的机器人列表
// @Summary 获取我的机器人列表
// @Tags 机器人管理
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页条数" default(20)
// @Router /bots/user/me [get]
func (h *Handler) ListMyBot(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	bots, total, err := h.svc.ListByUser(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, gin.H{"list": bots, "total": total, "page": page})
}

// List 获取所有机器人列表
// @Summary 获取机器人列表
// @Tags 机器人管理
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页条数" default(20)
// @Router /bots [get]
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	bots, total, err := h.svc.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, gin.H{"list": bots, "total": total, "page": page})
}

// RunNow 手动触发机器人执行
// @Summary 手动触发执行
// @Tags 机器人管理
// @Security ApiKeyAuth
// @Param id path int true "机器人ID"
// @Param body body object false "事件数据（任意JSON）"
// @Router /bots/{id}/run [post]
func (h *Handler) RunNow(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var event map[string]any
	_ = c.ShouldBindJSON(&event)
	if event == nil {
		event = map[string]any{}
	}
	if err := h.svc.RunNow(c.Request.Context(), uint(id), event); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, gin.H{"message": "triggered"})
}
