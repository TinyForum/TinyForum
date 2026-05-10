package bot

import (
	"strconv"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/service/bot"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc bot.Service
}

func NewHandler(svc bot.Service) *Handler {
	return &Handler{svc: svc}
}

// Create 创建机器人
// @Summary 创建机器人
// @Description 在指定板块下创建新的机器人
// @Tags 机器人管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Param body body request.CreateBotRequest true "机器人创建信息"
// @Success 200 {object} common.BasicResponse{data=object{id integer}} "创建成功，返回机器人ID"
// @Failure 400 {object} common.BasicResponse "请求参数错误（如板块ID无效、机器人信息不合法）"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Failure 403 {object} common.BasicResponse "无权限（需要版主或管理员权限）"
// @Failure 404 {object} common.BasicResponse "不存在"
// @Router /bots [post]
func (h *Handler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.CreateBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	botObj, err := h.svc.Create(c.Request.Context(), userID, &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, gin.H{"id": botObj.ID})
}

// Update 更新机器人
// @Summary 更新机器人
// @Description 更新指定机器人的配置信息（只能操作自己的机器人）
// @Tags 机器人管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "机器人ID"
// @Param body body request.UpdateBotRequest true "机器人更新信息"
// @Success 200 {object} common.BasicResponse{data=object} "更新成功，返回空数据"
// @Failure 400 {object} common.BasicResponse "请求参数错误"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Failure 403 {object} common.BasicResponse "无权限（非本人机器人）"
// @Failure 404 {object} common.BasicResponse "机器人不存在"
// @Router /bots/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req request.UpdateBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	err := h.svc.Update(c.Request.Context(), userID, uint(id), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, nil)
}

// Delete 删除机器人
// @Summary 删除机器人
// @Description 删除指定机器人（只能删除自己的机器人）
// @Tags 机器人管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "机器人ID"
// @Success 200 {object} common.BasicResponse{data=object} "删除成功，返回空数据"
// @Failure 400 {object} common.BasicResponse "请求参数错误"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Failure 403 {object} common.BasicResponse "无权限（非本人机器人）"
// @Failure 404 {object} common.BasicResponse "机器人不存在"
// @Router /bots/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	err := h.svc.Delete(c.Request.Context(), userID, uint(id))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, nil)
}

// Get 获取机器人详情
// @Summary 获取机器人详情
// @Description 根据ID获取机器人详细信息（公开信息，不需要鉴权？注意代码中未使用userID，但路由挂在Auth中间件下，实际仍需要认证）
// @Tags 机器人管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "机器人ID"
// @Success 200 {object} common.BasicResponse{data=object} "返回机器人视图对象（botVO）"
// @Failure 400 {object} common.BasicResponse "请求参数错误"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Failure 404 {object} common.BasicResponse "机器人不存在"
// @Router /bots/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	botVO, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, botVO)
}

// List 获取当前用户的机器人列表
// @Summary 获取我的机器人列表
// @Description 分页查询当前用户创建的所有机器人
// @Tags 机器人管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1) example(1)
// @Param pageSize query int false "每页条数" default(20) example(20)
// @Success 200 {object} common.BasicResponse{data=object{list=array,total=integer,page=integer}} "返回机器人列表及总数"
// @Failure 400 {object} common.BasicResponse "请求参数错误"
// @Failure 401 {object} common.BasicResponse "未授权"
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
	response.Success(c, gin.H{
		"list":  bots,
		"total": total,
		"page":  page,
	})
}

// List 获取机器人列表
// @Summary 获取机器人列表
// @Description 分页查询所有机器人
// @Tags 机器人管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1) example(1)
// @Param pageSize query int false "每页条数" default(20) example(20)
// @Success 200 {object} common.BasicResponse{data=object{list=array,total=integer,page=integer}} "返回机器人列表及总数"
// @Failure 400 {object} common.BasicResponse "请求参数错误"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Router /bots [get]
func (h *Handler) List(c *gin.Context) {
	// userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	bots, total, err := h.svc.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, gin.H{
		"list":  bots,
		"total": total,
		"page":  page,
	})
}

// RunNow 立即执行机器人
// @Summary 手动触发机器人执行
// @Description 立即执行指定机器人脚本一次，可附带事件数据（可选）
// @Tags 机器人管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "机器人ID"
// @Param body body object false "事件数据（任意JSON）" example({})
// @Success 200 {object} common.BasicResponse{data=object{message=string}} "触发成功，返回message"
// @Failure 400 {object} common.BasicResponse "请求参数错误"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Failure 403 {object} common.BasicResponse "无权限（非本人机器人）"
// @Failure 404 {object} common.BasicResponse "机器人不存在"
// @Router /bots/{id}/run [post]
func (h *Handler) RunNow(c *gin.Context) {
	// userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	// 可选 event data
	var event map[string]any
	_ = c.ShouldBindJSON(&event)
	err := h.svc.RunNow(c.Request.Context(), uint(id), event)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, gin.H{"message": "triggered"})
}
